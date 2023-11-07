package main

import (
	"github.com/Godyu97/geerpc/server"
	"log"
	"net"
)

func main() {
	// pick a free port
	l, err := net.Listen("tcp", "127.0.0.1:10928")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	server.Accept(l)
}
