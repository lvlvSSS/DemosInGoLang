package tcp_server

import "net"

type TCPClient struct {
	net.TCPConn
	head *TCPClient
	tail *TCPClient
}

type TCPContext struct {
	next   *TCPContext
	prev   *TCPContext
	Client *TCPClient
}

/*
	The Callback4ClientFunc type define the callback after receive bytes from the specified client.
	It is an adapter to allow the use of ordinary functions as callback handlers.
*/
type Callback4ClientFunc func(msg []byte)

type Callback4Client struct {
	callback Callback4ClientFunc
	name     string
}

func (callback Callback4ClientFunc) HandleMsg(msg []byte) {
	callback(msg)
}
