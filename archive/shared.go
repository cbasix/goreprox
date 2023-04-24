package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"time"
)

type Frame struct {
	ConnectionId   uint
	DropConnection bool
	Data           []byte
}

type connectionProvider func(connId uint, shared net.Conn)

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

func splitSharedConn(shared net.Conn, getConn connectionProvider) {
	var frm Frame

	dec := gob.NewDecoder(shared)
	err := dec.Decode(&frm)
	if err != nil {
		panic(err)
	}

	// use existing connection or create a new one
	conn := getConn(frm.ConnectionId, shared)

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
