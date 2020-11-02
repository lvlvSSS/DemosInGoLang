package server

import (
	"bytes"
	"fmt"
	Log "github.com/sirupsen/logrus"
	"net"
	"strings"
	"sync"
	"time"
)

// The TCPServer is the object for TCP server.
type TCPServer struct {
	ln *net.TCPListener

	// The Clients : store all the clients that connected to this server.
	clients map[chan []byte]*net.TCPConn

	// The messages: the data in the message would be broadcast to all clients that connected to this server
	messages chan []byte

	// The addClient : add the client with thread safe way.
	addClients chan addClientUnion

	// the removeClients: remove the client with thread safe way.
	removeClients chan chan []byte

	// the duration is used for the clients, the heart beat from client to server.
	// use 60s as default.
	duration time.Duration

	// The callbacks store the callback function that will handle the msg from clients
	callbacks sync.Map
}

// define the callback after receive bytes from the specified client.
type Callback4Client func(msg []byte)

// The Register is used to store the callbacks
func (server *TCPServer) Register(ip string, callback Callback4Client) {
	if _, ok := server.callbacks.Load(ip); ok {
		Log.Warnf("[server.Register] already exist callback for client[%s], and will be replaced", ip)
	}
	server.callbacks.Store(ip, callback)
}

type addClientUnion struct {
	channel chan []byte
	client  *net.TCPConn
}

func (server *TCPServer) SetHeartBeat(duration time.Duration) {
	server.duration = duration
}

func (server *TCPServer) GetHeartBeat() (dur time.Duration) {
	dur = server.duration
	return dur
}

// The Close is used to close server manually.
func (server *TCPServer) Close() {
	defer func() {
		if err := recover(); err != nil {
			Log.Fatal(fmt.Sprintf("Close server[%[1]s] error: %[2]s", server.ln.Addr(), err))
		}
	}()
	close(server.messages)
	server.ln.Close()
}

// The Broadcast is used to send msg to all clients
func (server *TCPServer) Broadcast(msg []byte) {
	defer func() {
		if err := recover(); err != nil {
			Log.Fatalf("[server.Send] The server is already closed!!! : %s", err)
		}
	}()
	if server.messages == nil {
		panic("[server.Send] please start the server first!")
	}
	server.messages <- msg
}

// The Start is used to start listen the client to connect.
// It will not block the thread.
func (server *TCPServer) Start(address string) {
	addr, err := net.ResolveTCPAddr("", address)
	if err != nil {
		Log.Error(fmt.Sprintf("Resolve address[%[1]s] error: %[2]s", address, err))
	}

	server.ln, err = net.ListenTCP("tcp", addr)
	if err != nil {
		Log.Error(fmt.Sprintf("Listen the address[%[1]s] errors: %[2]s", addr, err))
	}

	//initialize the chan
	server.initChan()

	go server.listenForAccept()
}

func (server *TCPServer) initChan() {
	server.clients = make(map[chan []byte]*net.TCPConn)
	server.addClients = make(chan addClientUnion)
	server.removeClients = make(chan chan []byte)
	server.messages = make(chan []byte)
	server.duration = time.Second * 60

	go func() {
		for {
			select {
			case client := <-server.addClients:
				server.clients[client.channel] = client.client
				Log.Infof("[server.initChan] Client[%s] added", client.client.RemoteAddr())
			case client := <-server.removeClients:
				tmp, ok := server.clients[client]
				if !ok {
					Log.Warnf("[server.initChan] client[%s] need to be removed , but not exists", tmp.RemoteAddr())
					break
				}
				delete(server.clients, client)
				Log.Infof("[server.initChan] client[%s] removed", tmp.RemoteAddr())
			case msg, ok := <-server.messages:
				// if close the messages , then close all the client chan
				if !ok {
					for client, _ := range server.clients {
						close(client)
					}
					break
				}
				for client, _ := range server.clients {
					client <- msg
				}
			}
		}
	}()
}

func (server *TCPServer) listenForAccept() {
	defer func() {
		if err := recover(); err != nil {
			Log.Fatalf("[server.listenForAccept] Accept error: %s", err)
		}
	}()
	for {
		conn, err := server.ln.AcceptTCP()
		if err != nil {
			Log.Warnf("Accept failed: %s", err.Error())
			break
		}

		go server.handleRW(conn)
	}
}

func (server *TCPServer) handleRW(conn *net.TCPConn) {
	defer func() {
		if err := recover(); err != nil {
			Log.Errorf("[server.handleRW] errors: %s", err)
		}
	}()
	// Add the clients to server.clients
	channelFromServer := make(chan []byte)
	clientUnion := addClientUnion{
		channel: channelFromServer,
		client:  conn,
	}
	server.addClients <- clientUnion

	Log.Infof("[server.handleRW] Client[%s] connected to the server", conn.RemoteAddr().String())
	// set the heartbeat timer.
	timer := time.NewTimer(server.duration)
	// handle the read process
	go func() {
		defer func() {
			if err := recover(); err != nil {
				Log.Fatalf("[server.handleRW] read client[%s] errors: ", conn.RemoteAddr().String(), err)
			}
		}()
		clientAddr := strings.SplitN(conn.RemoteAddr().String(), ":", 2)[0]
		buf := make([]byte, 1024)
		var buffer bytes.Buffer
		for {
			total, err := conn.Read(buf)
			if err != nil {
				server.removeClients <- channelFromServer
				Log.Infof("[server.handleRW] client[%s] is disconnected", conn.RemoteAddr().String())
				conn.Close()
				return
			}
			callback, ok := server.callbacks.Load(clientAddr)
			if !ok {
				Log.Warnf("[server.listenForAccept] no callback fit for the client[%s], but received messages", clientAddr)
				continue
			}

			if total == 0 && callback != nil {
				if buffer.Len() > 0 {
					(callback.(Callback4Client))(buffer.Bytes())
					buffer.Reset()
				}
				continue
			}
			buffer.Write(buf[0:total])
		}
	}()
	for {
		select {
		case <-timer.C:
			server.removeClients <- channelFromServer
			conn.Close()
		case msg, ok := <-channelFromServer:
			if !ok {
				server.removeClients <- channelFromServer
				conn.Close()
				break
			}
			if _, err := conn.Write(msg); err != nil {
				server.removeClients <- channelFromServer
				Log.Errorf("[server.handleRW] write msg[%s] to client[%s] errors: %s", string(msg), conn.RemoteAddr(), err.Error())
				conn.Close()
			}
		}
	}
}
