package simple

import (
	"log"
	"net"
	"sync"
	"time"
)

type manager struct {
	shared  SharedConn
	private map[int]net.Conn
	lastConnId int
	lock    *sync.RWMutex
}

func CreateManager(shared SharedConn) manager {
	return manager{
		shared:  shared,
		private: make(map[int]net.Conn),
		lastConnId: 0,
		lock:    new(sync.RWMutex),
	}
}

func (m *manager) AddPrivateConn(conn net.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.lastConnId += 1
	m.private[m.lastConnId] = conn

	go handleIncoming(connId, conn)
}

func (m *manager) handleIncoming(connId int, conn net.Conn){
	buf := make([]byte, 1024)
	read, err = conn.Read(buf)
	if (err != nil) {
		log.Printf("Error in connId: %d %s in %s", err, conn.RemoteAddr().String())
		m.invalidate(connId)
	}

	buff = pack()
}

func (m *manager) invalidate(connId int){
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.private, connId)
	frm := Frame{
		ConnectionId: connId,
		DropConnection: true,
		Data: []byte{},
	}
	m.shared.Write()

}

func sendDisconnect()


func (m *manager) PrivateConn(conn net.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.private[]
}
