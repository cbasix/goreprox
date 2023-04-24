package main

import "net"

type connHandler func(net.Conn)

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
