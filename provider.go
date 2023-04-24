package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"sync"
	"time"
)

var (
	newConnections = make(chan net.Conn)
	connections    = make(map[uint]net.Conn)
	mutex          = new(sync.RWMutex)
)

/*func startProvider() {
	fmt.Println("Connecting to proxy: ", proxyAddress)
	proxy, err := net.Dial("tcp", proxyAddress)
	if err != nil {
		panic(err)
	}
	go createConnections()
	handleProxyConnection(proxy)
}*/

func createConnections() {
	conn, err := net.Dial("tcp", exposeAddress)
	if err != nil {
		panic(err)
	}
	newConnections <- conn
}

func getConnection(connId uint, sharedConn net.Conn) net.Conn {
	mutex.Lock()
	defer mutex.Unlock()

	conn, connectionExists := connections[connId]
	if !connectionExists {
		conn = <-newConnections
		connections[connId] = conn
		stop := make(chan bool)
		go joinToShared(connId, conn, sharedConn, stop)
	}

	return conn
}

func handleProxyConnection(proxy net.Conn) {
	var frm Frame

	dec := gob.NewDecoder(proxy)
	err := dec.Decode(&frm)
	if err != nil {
		panic(err)
	}

	// use existing connection or create a new one
	conn := getConnection(frm.ConnectionId, proxy)

	// forward the data
	if len(frm.Data) > 0 {
		_, err = conn.Write(frm.Data)
		if err != nil {
			panic(err)
		}
	}

	if frm.DropConnection {
		conn.Close()
	}

}

func joinToShared(connId uint, privateConn net.Conn, sharedConn net.Conn, stop <-chan bool) {
	bufConn := bufio.NewReader(privateConn)
	for {
		select {
		case <-stop:
			return
		default:
		}

		Data := make([]byte, 8192)
		privateConn.SetReadDeadline(time.Now().Add(time.Millisecond * 10))

		read, err := bufConn.Read(Data)
		if err != nil {
			log.Println("Client connection error:", err)
		}

		frm := &Frame{
			ConnectionId: connId,
			Data:         Data[:read],
		}
		enc := gob.NewEncoder(sharedConn)
		err = enc.Encode(&frm)
		if err != nil {
			panic(err)
		}
	}
}

// on provider: connect to proxy
// wait for packages,
// if one is received check if it belongs to an open connection then copy it into the given connection
// if not open a new connection to target and start a goroutine that copies all answers retrieved to the shared connection.

// if either side closes the connection a close message must be shared to via the shared channel to the other end.
