package processor

const (
	MaxConnection = 1024
)

type Operation uint32

const (
	OperationRead  Operation = 0
	OperationWrite Operation = 1
)

type Event struct {
	Fd int
	Op Operation
}

type Multiplexer interface {
	Monitor(event Event) error
	Check() ([]Event, error)
	Close() error
}
