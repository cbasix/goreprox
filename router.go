package main

import (
	"fmt"
	"log"
	"sync"
)

type Router struct {
	shared           *ChannelConn
	destinations     map[uint]*ChannelConn
	destMutex        *sync.RWMutex
	connIdCounter    uint
	createConnection func() *ChannelConn
}

func CreateRouter(shared *ChannelConn) *Router {
	return &Router{
		shared:        shared,
		destinations:  make(map[uint]*ChannelConn),
		destMutex:     new(sync.RWMutex),
		connIdCounter: 0,
	}
}

/*
	Reads one frame from the shared connection and forwards it to the correct connection.

If no matching connection exists it tries to create a new one via the createConnection function set on the router.
If createConnection is not set on the router it panics.
*/
func (r *Router) route() {
	frame := <-r.shared.in
	target, found := r.getDest(frame.ConnectionId)

	// forward fram only when it has data
	if len(frame.Data) > 0 {
		if !found {
			if r.createConnection == nil {
				panic(fmt.Sprintf("connection with id %d not found in destinations: %+v", frame.ConnectionId, r.destinations))
			} else {
				target = r.createConnection()
				connId := r.createDest(target)
				frame.ConnectionId = connId
			}
		}
		target.out <- frame
	}

	// close connection if ordered to do so
	if found && frame.DropConnection {
		close(target.out)
		target.conn.Close()
		r.delDest(frame.ConnectionId)
	}
}

/* Reads one frame from all destinations and forwards them to the shared connection. */
func (r *Router) join() {
	r.destMutex.RLock()
	defer r.destMutex.RUnlock()

	var wg sync.WaitGroup
	for connId, frames := range r.destinations {
		if frames.closed {
			continue
		}
		wg.Add(1)
		go func(connId uint, frames *ChannelConn) {
			defer wg.Done()

			select {
			case frm, more := <-frames.in:
				if frames.closed {
					return
				}
				if more {
					frm.ConnectionId = connId
					r.shared.out <- frm
				} else {
					r.destMutex.RUnlock()
					r.destMutex.Lock()
					if frames.closed {
						return
					}
					frames.closed = true

					log.Print("conn closed send")
					r.shared.out <- Frame{
						ConnectionId:   connId,
						DropConnection: true,
					}

					r.destMutex.Unlock()
					r.destMutex.RLock()
				}
			default:
			}
		}(connId, frames)
	}
	wg.Wait()
}

// make map get access threadsafe
func (r *Router) getDest(connId uint) (*ChannelConn, bool) {
	r.destMutex.RLock()
	defer r.destMutex.RUnlock()

	s, f := r.destinations[connId] // todo improve
	return s, f
}

// make map put threadsafe
func (r *Router) createDest(frames *ChannelConn) uint {
	r.destMutex.Lock()
	defer r.destMutex.Unlock()

	r.connIdCounter += 1 // TODO overflow / reuse
	r.destinations[r.connIdCounter] = frames

	return r.connIdCounter
}

// make map delete threadsafe
func (r *Router) delDest(connId uint) {
	r.destMutex.Lock()
	defer r.destMutex.Unlock()

	delete(r.destinations, connId)
}
