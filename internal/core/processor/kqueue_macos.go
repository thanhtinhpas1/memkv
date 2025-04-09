//go:build darwin

package processor

import "syscall"

type KqueueProcessor struct {
	fd            int
	kqEvents      []syscall.Kevent_t
	genericEvents []Event
}

var _ Multiplexer = (*KqueueProcessor)(nil)

func CreateIoMultiplexer() (*KqueueProcessor, error) {
	kqFd, err := syscall.Kqueue()
	if err != nil {
		return nil, err
	}

	return &KqueueProcessor{
		fd:            kqFd,
		kqEvents:      make([]syscall.Kevent_t, MaxConnection),
		genericEvents: make([]Event, MaxConnection),
	}, nil
}

// Monitor implements Multiplexer.
func (k *KqueueProcessor) Monitor(event Event) error {
	_, err := syscall.Kevent(k.fd, []syscall.Kevent_t{event.toNative(syscall.EV_ADD)}, nil, nil)
	return err
}

// Check implements Multiplexer.
func (k *KqueueProcessor) Check() ([]Event, error) {
	n, err := syscall.Kevent(k.fd, nil, k.kqEvents, nil)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		k.genericEvents[i] = createEvent(k.kqEvents[i])
	}

	return k.genericEvents[:n], nil
}

// Close implements Multiplexer.
func (k *KqueueProcessor) Close() error {
	return syscall.Close(k.fd)
}
