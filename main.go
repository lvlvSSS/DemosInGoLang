package main

import (
	"bufio"
	"fmt"
	server "go_tools/tcp_server"
	_ "net/http/pprof"
	"os"
	"time"
)

func main() {

	myserver := server.TCPServer{}
	myserver.SetHeartBeat(time.Second * 3600)

	myserver.AddLast("192.168.50.251", func(msg []byte) {
		fmt.Printf("Server receive msg from 192.168.50.251 : \r\n%#X \r\n", msg)
	})
	myserver.AddLast("192.168.50.7", func(msg []byte) {
		fmt.Printf("Server receive msg from 192.168.50.7 : \r\n%#X \r\n", msg)
	})
	myserver.AddLast("192.168.50.17", func(msg []byte) {
		fmt.Printf("Server receive msg from 192.168.50.17 : \r\n%#X \r\n", msg)
	})
	myserver.AddLast("192.168.50.37", func(msg []byte) {
		fmt.Printf("Server receive msg from 192.168.50.37 : \r\n%#X \r\n", msg)
	})
	myserver.Start(":50000")

	input := bufio.NewReader(os.Stdin)

	line, _, _ := input.ReadLine()
	for string(line) != "quit" {
		mybytes := []byte{0x02,0x00,0x00,0x10,0x03,0x70,0x73,0x03}
		myserver.Broadcast(mybytes)
		line, _, _ = input.ReadLine()
	}

}
