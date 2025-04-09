//go:build darwin

package processor

import "syscall"

func (e Event) toNative(flags uint16) syscall.Kevent_t {
	var filter int16 = syscall.EVFILT_WRITE
	if e.Op == OperationRead {
		filter = syscall.EVFILT_READ
	}

	return syscall.Kevent_t{
		Filter: filter,
		Ident:  uint64(e.Fd),
		Flags:  flags,
	}
}

func createEvent(kq syscall.Kevent_t) Event {
	op := OperationWrite
	if kq.Filter == syscall.EVFILT_READ {
		op = OperationRead
	}
	return Event{
		Fd: int(kq.Ident),
		Op: op,
	}
}
