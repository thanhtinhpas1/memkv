//go:build linux

package processor

import (
	"log"
	"syscall"
)

type EpollProcessor struct {
	fd            int
	epollEvents   []syscall.EpollEvent
	genericEvents []Event
}

func CreateIoMultiplexer() (*EpollProcessor, error) {
	epollFd, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &EpollProcessor{
		fd:            epollFd,
		epollEvents:   make([]syscall.EpollEvent, MaxConnection),
		genericEvents: make([]Event, MaxConnection),
	}, nil
}

var _ Multiplexer = (*EpollProcessor)(nil)

func (fd *EpollProcessor) Monitor(event Event) error {
	epollEvent := event.toNative()
	return syscall.EpollCtl(fd.fd, syscall.EPOLL_CTL_ADD, event.Fd, &epollEvent)
}

// Check implements Multiplexer.
func (fd *EpollProcessor) Check() ([]Event, error) {
	n, err := syscall.EpollWait(fd.fd, fd.epollEvents, -1)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		fd.genericEvents[i] = createEvent(fd.epollEvents[i])
	}
	return fd.genericEvents[:n], nil
}

// Close implements Multiplexer.
func (fd *EpollProcessor) Close() error {
	return syscall.Close(fd.fd)
}
