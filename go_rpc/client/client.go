package client

type Call struct {
	Seq          uint64
	ServerMethod string
	Args         interface{}
	Reply        interface{}
	Error        error
	Done         chan *Call
}

func (call *Call) done() {
	call.Done <- call
}
