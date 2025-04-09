//go:build linux

package processor

import "syscall"

func (e Event) toNative() syscall.EpollEvent {
	event := syscall.EPOLLIN
	if e.Op == OperationWrite {
		event = syscall.EPOLLOUT
	}

	return syscall.EpollEvent{
		Fd:     int32(e.Fd),
		Events: uint32(event),
	}
}

func createEvent(ep syscall.EpollEvent) Event {
	op := OperationRead
	if ep.Events == syscall.EPOLLOUT {
		op = OperationWrite
	}

	return Event{
		Fd: int(ep.Fd),
		Op: op,
	}
}
