package main

import (
	"bufio"
	"fmt"
	server "go_tools/tcp_server"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	myserver := server.TCPServer{}
	myserver.Start("127.0.0.1:50000")
	myserver.AddLast("127.0.0.1", func(msg []byte) {
		fmt.Printf("Server receive msg from 127.0.0.1 : \r\n%s \r\n", string(msg))
	})

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	input := bufio.NewReader(os.Stdin)

	line, _, _ := input.ReadLine()
	for string(line) != "quit" {
		myserver.Broadcast(line)
		line, _, _ = input.ReadLine()
	}
}
