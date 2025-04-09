package constants


const (
	ResponseOK                     = "OK"
	ResponsePong                   = "PONG"
	ResponseInvalidKey             = "ERR invalid key"
	ResponseInvalidValue           = "ERR invalid value"
	ResponseInvalidPattern         = "ERR invalid pattern"
	ResponseWrongNumberOfArguments = "ERR wrong number of arguments"
)

const (
	EngineStatusWaiting      = 1
	EngineStatusRunning      = 2
	EngineStatusShuttingDown = 3
)

var RespNil = []byte("$-1\r\n")
var RespOk = []byte("+OK\r\n")
var RespZero = []byte(":0\r\n")
var RespOne = []byte(":1\r\n")
var RespEmptyArray = []byte("*0\r\n")
var TtlKeyNotExist = []byte(":-2\r\n")
var TtlKeyExistNoExpire = []byte(":-1\r\n")
