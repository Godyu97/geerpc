package client

import (
	"sync"
	"github.com/Godyu97/geerpc/codec"
)

type Call struct {
	Seq           uint64
	ServiceMethod string     // format "<service>.<method>"
	Args          any        // arguments to the function
	Reply         any        // reply from the function
	Error         error      // if error occurs, it will be set
	Done          chan *Call // Strobes when call is complete.
}

//为了支持异步调用，Call 结构体中添加了一个字段 Done，Done 的类型是 chan *Call，当调用结束时，会调用 call.done() 通知调用方
func (c *Call) done() {
	c.Done <- c
}

// Client represents an RPC Client.
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
type Client struct {
	cc       codec.Codec
	opt      *codec.Option
	sending  sync.Mutex // protect following
	header   codec.Header
	mu       sync.Mutex // protect following
	seq      uint64
	pending  map[uint64]*Call
	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}
