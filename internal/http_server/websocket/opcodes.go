package websocket

// The OpCode type is used to identify the different operations possible.
type OpCode string

const (
	GET   OpCode = "GET"
	SUB   OpCode = "SUB"
	UNSUB OpCode = "UNSUB"
	HIST  OpCode = "HIST"
	REPLY OpCode = "REPLY"
	ERR   OpCode = "ERR"
)
