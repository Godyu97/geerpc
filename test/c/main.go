package main

import (
	"log"
	"fmt"
	"github.com/Godyu97/geerpc/codec"
	"encoding/json"
	"net"
	"github.com/Godyu97/geerpc/server"
)

func main() {
	// in fact, following code is like a simple geerpc client
	conn, err := net.Dial("tcp", "127.0.0.1:10928")
	if err != nil {
		panic(err)
	}
	defer func() { _ = conn.Close() }()

	log.Println("连接成功")
	// send options
	_ = json.NewEncoder(conn).Encode(server.Option{
		MagicNumber: server.MagicNumber,
		CodecType:   codec.JsonType,
	})
	cc := codec.NewJsonCodec(conn)
	// send request & receive response
	for i := 0; i < 5; i++ {
		log.Println("开始请求", i)
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}
