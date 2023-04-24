package main

import "net"

type connHandler func(net.Conn)

package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"sync"
	"time"
)

var clientConns := make(chan net.Conn, 10)
var provider := make(chan net.Conn)

func startProxy() {
	fmt.Println("Listening for provider on: ", proxyAddress)
	go startListening("tcp", proxyAddress, provider)
	go startListening("tcp", exposeAddress, clientConns)

	for {
		p := <- provider
		fmt.Println("Provider connected from: ", p.RemoteAddr().String())

		go splitSharedConn(p, getConnection) // Todo if not found error instead
		
		for {
			c := <- clientConns
			fmt.Println("Client connected from: ", p.RemoteAddr().String())

			go joinToShared()
		}
	}

	
}


func startListening(protocol string, listenOn string, connHandler chan<- net.Conn) {
	ln, err := net.Listen(protocol, listenOn)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		connHandler <- conn
	}
}
