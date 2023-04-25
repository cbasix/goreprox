package main

import "net"

type Frame struct {
	ConnectionId   uint
	DropConnection bool
	Data           []byte
}

type ChannelConn struct {
	in     chan Frame
	out    chan Frame
	conn   net.Conn
	closed bool
}

func CreateConn(inSize int, outSize int) *ChannelConn {
	return &ChannelConn{
		in:     make(chan Frame, inSize),
		out:    make(chan Frame, outSize),
		closed: false,
	}
}
