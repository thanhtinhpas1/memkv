//go:build linux

package processor

import (
	"log"
	"syscall"
)

// EpollProcessor is a structure that implements the Multiplexer interface
// using the epoll mechanism for I/O multiplexing on Linux systems.
//
// With Epoll, Linux register file descriptors once using epoll_ctl, and the kernel internally tracks them.
// When handling events, the kernel notifies the user space about the file descriptors that are ready for I/O.
// Because kernel already tracking them so kernel no need to scan entirely for file descriptors.
// This is more efficient than select or poll, especially for a large number of file descriptors.
//
// Internally, kernel using red black tree to store file descriptors and a readly list (linked list) to store ready ones.
// Lookup and updates are O(log n), and notifying events is often O(1).
//
// Check this source code: https://github.com/torvalds/linux/blob/master/fs/eventpoll.c
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
