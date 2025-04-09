package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"memkv/internal/server"
)

var (
	host string
	port int
)

func init() {
	flag.StringVar(&host, "host", "0.0.0.0", "host")
	flag.IntVar(&port, "port", 6379, "port")
	flag.Parse()
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	s := server.NewServer(host, port)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go s.RunAsyncTCPServer(&wg)
	go server.WaitForSignal(&wg, signals)

	wg.Wait()
}
