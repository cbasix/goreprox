package main

type Frame struct {
	ConnectionId   uint
	DropConnection bool
	Data           []byte
}
