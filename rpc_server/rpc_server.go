package main

import (
	"akv/akv"
	"fmt"
	"net"
	"net/rpc"
)

func main() {
	server := akv.CreateAndrewKeyValueStore()
	rpc.Register(server)
	l, err := net.Listen("tcp", "localhost:1234")

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Server started on port 1234")
	}

	defer l.Close()

	rpc.Accept(l)
}