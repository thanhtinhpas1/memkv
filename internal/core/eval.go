package core

import (
	"errors"
	"io"

	"memkv/internal/constants"
)

const (
	CommandPing   = "PING"
	CommandSet    = "SET"
	CommandGet    = "GET"
	CommandTTL    = "TTL"
	CommandDel    = "DEL"
	CommandExists = "EXISTS"
	CommandKeys   = "KEYS"
)

func EvalAndResponse(cmd *MemkvCommand, c io.ReadWriter) error {
	var res []byte

	switch cmd.Cmd {
	case CommandPing:
		res = cmdPing(cmd, c)

		// Sorted set
	case "ZADD":
		res = cmdZADD(cmd.Args)
	case "ZRANK":
		res = cmdZRANK(cmd.Args)
	case "ZREM":
		res = cmdZREM(cmd.Args)
	case "ZSCORE":
		res = cmdZSCORE(cmd.Args)
	case "ZCARD":
		res = cmdZCARD(cmd.Args)
	}

	_, err := c.Write(res)
	return err
}

func cmdPing(cmd *MemkvCommand, c io.ReadWriter) []byte {
	var buf []byte
	if len(cmd.Args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}

	if len(cmd.Args) == 0 {
		buf = Encode(constants.ResponsePong, true)
	} else {
		buf = Encode(cmd.Args[0], false)
	}

	return buf
}
