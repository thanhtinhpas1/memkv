package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"syscall"

	"memkv/internal/constants"
	core "memkv/internal/core"
	"memkv/internal/core/processor"
)

var eStatus int32 = constants.EngineStatusWaiting

// Server represents our Redis-like server
type Server struct {
	host string
	port int
}

// NewServer creates a new server instance
func NewServer(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

// Start starts the server
func (s *Server) RunAsyncTCPServer(wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Println("starting memkv server...")

	var err error
	events := make([]processor.Event, processor.MaxConnection)
	clientNum := 0

	// Create fd socket server
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Println(err)
		return err
	}

	defer syscall.Close(serverFD)

	// Set the Socket operate in non-blocking mode
	// Default mode is blocking mode
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		log.Println(err)
		return err
	}

	ip4 := net.ParseIP(s.host)
	if ip4 == nil {
		log.Println("invalid ip addres")
		return errors.New("invalid ip address")
	}

	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: s.port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		log.Println(err)
		return err
	}

	// start listening
	if err = syscall.Listen(serverFD, processor.MaxConnection); err != nil {
		log.Println(err)
		return err
	}

	// multiplexer can monitor many server fds at the same time, whenever a fd ready to read or write, it will notify our server
	multiplexer, err := processor.CreateIoMultiplexer()
	if err != nil {
		log.Println(err)
		return err
	}

	defer multiplexer.Close()

	if err = multiplexer.Monitor(processor.Event{
		Fd: serverFD,
		Op: processor.OperationRead,
	}); err != nil {
		log.Println(err)
		return err
	}

	for atomic.LoadInt32(&eStatus) != constants.EngineStatusShuttingDown {
		events, err = multiplexer.Check()
		if err != nil {
			continue
		}

		if !atomic.CompareAndSwapInt32(&eStatus, constants.EngineStatusWaiting, constants.EngineStatusRunning) {
			if eStatus == constants.EngineStatusShuttingDown {
				return nil
			}
		}

		for _, event := range events {
			if event.Fd == serverFD {
				clientNum++
				log.Printf("new client connected: id=%d\n", clientNum)

				// accept new connection
				connFd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println(err)
					continue
				}

				if err = syscall.SetNonblock(connFd, true); err != nil {
					return err
				}

				if err = multiplexer.Monitor(processor.Event{
					Fd: connFd,
					Op: processor.OperationRead,
				}); err != nil {
					log.Println(err)
					return err
				}
			} else {
				// the client FD is ready for reading, means an existing client is sending us a message
				comm := core.FDCommand{
					Fd: event.Fd,
				}
				cmd, err := readCommandFD(event.Fd)
				if err != nil {
					syscall.Close(event.Fd)
					clientNum--
					log.Println("client quit")
					atomic.SwapInt32(&eStatus, constants.EngineStatusWaiting)
					continue
				}
				responseRw(cmd, comm)
			}
			atomic.SwapInt32(&eStatus, constants.EngineStatusWaiting)
		}
	}

	return nil
}

func responseRw(cmd *core.MemkvCommand, rw io.ReadWriter) {
	err := core.EvalAndResponse(cmd, rw)
	if err != nil {
		responseErrorRw(err, rw)
	}
}

func responseErrorRw(err error, rw io.ReadWriter) {
	rw.Write([]byte(fmt.Sprintf("-%s%s", err, core.CRLF)))
}

func readCommandFD(fd int) (*core.MemkvCommand, error) {
	var buf = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	return core.ParseCmd(buf[:n])
}

func WaitForSignal(wg *sync.WaitGroup, signals chan os.Signal) {
	defer wg.Done()
	<-signals

	for atomic.LoadInt32(&eStatus) == constants.EngineStatusRunning {
	}
	log.Println("Shutting down gracefully...")
	os.Exit(0)
}
