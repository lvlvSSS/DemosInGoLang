package tcp_server

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	Log "github.com/sirupsen/logrus"
	util "go_tools/util"
	"net"
	"strings"
	"sync"
	"time"
)

/*
	The TCPServer is the object for TCP server.
*/
type TCPServer struct {
	ln *net.TCPListener

	/* 	The Clients : store all the clients that connected to this server. */
	clients map[chan []byte]*net.TCPConn

	/* 	The messages: the data in the message would be broadcast to all clients that connected to this server */
	messages chan []byte

	/* 	The addClient : add the client with thread safe way. */
	addClients chan addClientUnion

	/* 	the removeClients: remove the client with thread safe way. */
	removeClients chan chan []byte

	/*
		the duration is used for the clients, the heart beat from client to server.
		use 60s as default.
	*/
	duration time.Duration

	/*  The callbacksMap store the callback function that will handle the msg from clients */
	callbacksMap sync.Map
}

type clientsList struct {
	list.List
	/* the rwLock is used to lock the value of callbacks */
	rwLock sync.RWMutex
	/* indicates how many callback in the list*/
	count int

	ip string
}

/*
	The AddLast is used to store the callbacks
*/
func (server *TCPServer) AddLastWithName(ip string, name string, callback func(msg []byte)) error {
	name = strings.TrimSpace(name)
	callbackFunc := Callback4Client{callback: Callback4ClientFunc(callback), name: name}
	callbacks, ok := server.callbacksMap.Load(ip)
	if !ok {
		Log.Debugf(
			"[%s] callback list not exist for client[%s], and need be created",
			util.GetCallFuncName(1), ip)
		callbacks = &clientsList{count: 0, ip: ip}
		server.callbacksMap.Store(ip, callbacks)
	}
	callbacksList := callbacks.(*clientsList)
	if err := server.insertCallback(callbacksList.Len(), callbackFunc, callbacksList); err != nil {
		return err
	}
	return nil
}

func (server *TCPServer) AddLast(ip string, callback func(msg []byte)) error {
	return server.AddLastWithName(ip, "", callback)
}

/*
	the insertCallback is to insert callback to list. The first position starts with 0.
*/
func (server *TCPServer) insertCallback(
	position int, callback Callback4Client, callbacks *clientsList) error {
	if callbacks == nil {
		return errors.New("callbacks is nil")
	}

	/* write lock start */
	callbacks.rwLock.Lock()
	defer callbacks.rwLock.Unlock()

	/* set the default name if name is empty*/
	if strings.TrimSpace(callback.name) == "" {
		callback.name = fmt.Sprintf("[%s]-[%d]", callbacks.ip, callbacks.count)
	}
	callbacks.count++

	if position < 0 {
		position = callbacks.Len() + position
	}

	if callbacks.Len() <= 0 || position > (callbacks.Len()-1) || position < 0 {
		callbacks.PushBack(callback)
		return nil
	}
	at := callbacks.Front()
	for i := 0; i < position; i++ {
		at = at.Next()
	}
	if at != nil {
		callbacks.InsertBefore(callback, at)
	} else {
		callbacks.PushBack(callback)
	}
	return nil
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

/*
	The Close is used to close server manually.
*/
func (server *TCPServer) Close() {
	defer func() {
		if err := recover(); err != nil {
			Log.Fatalf("[%[1]s] close server[%[2]s], error: %[3]s", util.GetCallFuncName(1),
				server.ln.Addr(), err)
		}
	}()
	close(server.messages)
	server.clients = nil
	server.addClients = nil
	server.removeClients = nil
	_ = server.ln.Close()
	server.callbacksMap = sync.Map{}
}

/*
	The Broadcast is used to send msg to all clients
*/
func (server *TCPServer) Broadcast(msg []byte) {
	if server.messages == nil {
		panic(fmt.Sprintf("[%s] please start the server first!", util.GetCallFuncName(1)))
	}
	server.messages <- msg
}

/*
	The Start is used to start listen the client to connect.
 	It will not block the thread.
*/
func (server *TCPServer) Start(address string) {
	addr, err := net.ResolveTCPAddr("", address)
	if err != nil {
		Log.Errorf("[%[1]s] resolve address[%[2]s], error: %[3]s", util.GetCallFuncName(1),
			address, err)
	}

	server.ln, err = net.ListenTCP("tcp", addr)
	if err != nil {
		Log.Errorf("[%[1]s]Listen the address[%[2]s] errors: %[3]s", util.GetCallFuncName(1),
			addr, err)
	}

	/* initialize the chan */
	server.initChan()
	Log.Infof(
		"[%s] server start to listen[%s]", util.GetCallFuncName(1), addr.String())
	go server.listenForAccept()
}

func (server *TCPServer) initChan() {
	server.clients = make(map[chan []byte]*net.TCPConn)
	server.addClients = make(chan addClientUnion)
	server.removeClients = make(chan chan []byte)
	server.messages = make(chan []byte)
	if server.duration == 0 {
		server.duration = time.Second * 60
	}
	Log.SetLevel(Log.DebugLevel)
	Log.Debugf(
		"[%s] the fields of server initialized.", util.GetCallFuncName(1))
	go func() {
		for {
			select {
			case client := <-server.addClients:
				server.clients[client.channel] = client.client
				Log.Infof(
					"[%s] Client[%s] added", util.GetCallFuncName(2), client.client.RemoteAddr())
			case client := <-server.removeClients:
				tmp, ok := server.clients[client]
				if !ok {
					Log.Debugf(
						"[%s] client need to be removed , but not exists", util.GetCallFuncName(1))
					break
				}
				Log.Infof(
					"[%s] client[%s] removed", util.GetCallFuncName(2), tmp.RemoteAddr())

				if _, isOk := server.clients[client]; isOk {
					/* delete the callbacks attached to the address */
					//server.callbacksMap.LoadAndDelete(strings.SplitN(targetConn.RemoteAddr().String(), ":", 2)[0])
					delete(server.clients, client)
				}
				close(client)
			case msg, ok := <-server.messages:
				/* if close the messages chan , then close all the clients */
				if !ok {
					Log.Debugf("[%s] close all clients that connected to server",
						util.GetCallFuncName(1))
					for client := range server.clients {
						close(client)
					}
					return
				}
				Log.Debugf("[%s] ready to broadcast message[%v] to all clients",
					util.GetCallFuncName(1), msg)
				for client := range server.clients {
					client <- msg
				}
			}
		}
	}()
}

func (server *TCPServer) listenForAccept() {
	defer func() {
		if err := recover(); err != nil {
			Log.Fatalf(
				"[%s] Accept error: %s",
				util.GetCallFuncName(1), err)
		}
	}()
	for {
		conn, err := server.ln.AcceptTCP()
		Log.Infof(
			"[%s] accept a new connection[%s] ",
			util.GetCallFuncName(1),
			conn.RemoteAddr().String())
		if err != nil {
			Log.Warnf(
				"[%s] accept failed: %s",
				util.GetCallFuncName(1), err.Error())
			return
		}

		/* start to handle the read/write */
		go server.handleRW(conn)
	}
}

func (server *TCPServer) handleRW(conn *net.TCPConn) {
	defer func() {
		if err := recover(); err != nil {
			Log.Errorf(
				"[%s] errors: %s",
				util.GetCallFuncName(1), err)
		}
	}()
	/* Add the clients to server.clients */
	channelFromServer := make(chan []byte)
	clientUnion := addClientUnion{
		channel: channelFromServer,
		client:  conn,
	}
	server.addClients <- clientUnion

	Log.Infof(
		"[%s] Client[%s] connected to the server",
		util.GetCallFuncName(1),
		conn.RemoteAddr().String())

	/* set the heartbeat timer.*/
	timer := time.NewTimer(server.duration)

	/* handle the read process */
	go func() {
		defer func() {
			if err := recover(); err != nil {
				Log.Fatalf("[%s] read client[%s], errors: %v",
					util.GetCallFuncName(1),
					conn.RemoteAddr().String(), err)
			}
		}()
		clientAddr := strings.SplitN(conn.RemoteAddr().String(), ":", 2)[0]
		buf := make([]byte, 1024)
		var buffer bytes.Buffer
		for {
			total, err := conn.Read(buf)
			if err != nil || total <= 0 {
				Log.Infof(
					"[%s] client[%s] is disconnected, caused : %v",
					util.GetCallFuncName(1),
					conn.RemoteAddr().String(), err)

				server.removeClients <- channelFromServer
				conn.Close()
				timer.Stop()
				return
			}

			Log.Debugf(
				"[%s] receive msg from client[%s]",
				util.GetCallFuncName(1),
				conn.RemoteAddr())

			timer.Stop()
			timer.Reset(server.duration)

			callback, ok := server.callbacksMap.Load(clientAddr)
			if !ok {
				Log.Warnf(
					"[%s] no callback fit for the client[%s], but received messages",
					util.GetCallFuncName(1),
					clientAddr)
				continue
			}
			buffer.Write(buf[0:total])

			if callback != nil {
				if buffer.Len() > 0 {
					callbackLi := callback.(*clientsList)
					at := callbackLi.Front()
					for i := 0; i < callbackLi.Len(); i++ {
						if at == nil {
							break
						}
						(at.Value.(Callback4Client).callback)(buffer.Bytes())
						at = at.Next()
					}

					buffer.Reset()
				}
				continue
			}
		}
	}()
	for {
		select {
		case <-timer.C:
			/*
				heartbeat timeout, then close the connection ,
				close the channel between the server and the specified client.
			*/
			Log.Debugf(
				"[%s] client[%s] timeout - no heartbeat",
				util.GetCallFuncName(1),
				conn.RemoteAddr())

			server.removeClients <- channelFromServer
			conn.Close()
			return
		case msg, ok := <-channelFromServer:
			if !ok {
				Log.Debugf(
					"[%s] - someone channel closed by server",
					util.GetCallFuncName(1))
				conn.Close()
				return
			}
			Log.Debugf(
				"[%s] ready to send msg to client[%s]",
				util.GetCallFuncName(1),
				conn.RemoteAddr())

			if _, err := conn.Write(msg); err != nil {
				server.removeClients <- channelFromServer
				Log.Errorf(
					"[%s] write msg[%x] to client[%s] , errors : %s",
					util.GetCallFuncName(1),
					msg, conn.RemoteAddr(), err.Error())
				conn.Close()
			} else {
				timer.Stop()
				timer.Reset(server.duration)
			}
		}
	}
}
