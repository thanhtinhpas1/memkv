package core

import "syscall"

const defaultBufferSize = 512

type MemkvCommand struct {
	Cmd  string
	Args []string
}

type FDCommand struct {
	Fd int
}

func (fd FDCommand) Read(data []byte) (int, error) {
	return syscall.Read(fd.Fd, data)
}

func (fd FDCommand) Write(data []byte) (int, error) {
	return syscall.Write(fd.Fd, data)
}
