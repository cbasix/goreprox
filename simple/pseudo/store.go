package pseudo

import (
	"net"
	"sync"
)

type connectionStore interface {
	addConnection(conn net.Conn) int
	removeConnnection(connId int)
	getConnnection(connId int) net.Conn
}

type providerStore interface {
	getProviderConnnection() net.Conn
	setProviderConnnection(net.Conn)
}

type managerState interface {
	providerStore
	connectionStore
}

// move to own file?
type managerStateMap struct {
	shared     net.Conn
	private    map[int]net.Conn
	lastConnId int
	lock       *sync.RWMutex
}

// todo implement
func (m *managerStateMap) addConnection(c net.Conn) int {
	return m.lastConnId
}

// todo implement
func (m *managerStateMap) removeConnnection(c net.Conn) {}

// todo implement
func (m *managerStateMap) getConnnection(connId int) {}

// todo implement
func (m *managerStateMap) getProviderConnnection() net.Conn {}

// todo implement
func (m *managerStateMap) setProviderConnnection(c net.Conn) {}
