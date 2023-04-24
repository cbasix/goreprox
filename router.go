package main

import (
	"fmt"
	"sync"
)

//                          		 connList
// -net.Conn-> packedReader -Frame-> Router <-Frame +UnpackedReader
//             packedWriter  <-Frame-       -> Frame UnpackedWriter

type Router struct {
	shared        chan Frame
	destinations  map[uint]chan Frame
	destMutex     sync.RWMutex
	connIdCounter uint
}

func CreateRouter(shared chan Frame) Router {
	return Router{
		shared:        shared,
		destinations:  make(map[uint]chan Frame),
		destMutex:     sync.RWMutex{},
		connIdCounter: 0,
	}
}

func (r Router) route() {
	for {
		frame := <-r.shared
		target, found := r.destinations[frame.ConnectionId]

		if !found {
			panic(fmt.Sprintf("connection with id %d not found", frame.ConnectionId))
		}

		target <- frame
	}
}

func (r Router) routeOrCreate(createConnection func() chan Frame) {
	for {
		frame := <-r.shared
		target, found := r.getDest(frame.ConnectionId)

		if !found {
			target = createConnection()
			connId := r.createDest(target)
			frame.ConnectionId = connId
		}

		target <- frame
	}
}

/* Reads all destinations frames and forwards them to the shared connection. */
func (r Router) join() {
	for connId, frames := range r.destinations {
		frm := <-frames
		frm.ConnectionId = connId
		r.shared <- frm
	}
}

// make map get access threadsafe
func (r Router) getDest(connId uint) (chan Frame, bool) {
	r.destMutex.RLock()
	defer r.destMutex.RUnlock()

	s, f := r.destinations[connId] // todo improve
	return s, f
}

// make map put threadsafe
func (r Router) createDest(frames chan Frame) uint {
	r.destMutex.Lock()
	defer r.destMutex.Unlock()

	r.connIdCounter += 1 // TODO overflow / reuse
	r.destinations[r.connIdCounter] = frames

	return r.connIdCounter
}

// make map delete threadsafe
func (r Router) delDest(connId uint) {
	r.destMutex.Lock()
	defer r.destMutex.Unlock()

	delete(r.destinations, connId)
}
