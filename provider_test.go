package main

import (
	"encoding/gob"
	"net"
	"sync"
	"testing"
)

func clearContext() {
	newConnections = make(chan net.Conn)
	connections = make(map[uint]net.Conn)
	mutex = new(sync.RWMutex)
}

func TestHandleProxyConnection(t *testing.T) {
	//clearContext()
	shared1, shared2 := net.Pipe()
	private1, private2 := net.Pipe()
	go func() { newConnections <- private1 }()
	go gob.NewEncoder(shared2).Encode(&Frame{ConnectionId: 0, DropConnection: true, Data: []byte{'o'}})
	go handleProxyConnection(shared1)

	data := make([]byte, 1)
	private2.Read(data)

	if data[0] != 'o' {
		t.Errorf("got %q, wanted %q", data[0], 'o')
	}
}

func TestJoinToShared(t *testing.T) {
	//clearContext()
	shared1, shared2 := net.Pipe()
	private1, private2 := net.Pipe()
	stop := make(chan bool)

	go private2.Write([]byte{'u', 'l'})
	go joinToShared(0, private1, shared1, stop)

	var frm Frame
	dec := gob.NewDecoder(shared2)
	err := dec.Decode(&frm)
	if err != nil {
		panic(err)
	}

	stop <- true

	if frm.Data[0] != 'u' {
		t.Errorf("got %q, wanted %q", frm.Data[0], 'u')
	}
}
