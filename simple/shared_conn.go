package simple

import (
	"log"
	"net"
	"time"
)

type SharedConn struct {
	conn    net.Conn
	server  bool
	proto   string
	address string
}

func (s *SharedConn) Write(data []byte) {
	_, err := s.conn.Write(data)
	if err != nil {
		log.Printf("Error in sharedConn write: %s", err)

		s.reaquireLoop()

		_, _ = s.conn.Write(data) // TODO just reaquired what happens if this fails...? reaquire again?

	}
}

func (s *SharedConn) Read(data []byte) int {
	read, err := s.conn.Read(data)
	if err != nil {
		log.Printf("Error in sharedConn read: %s", err)

		// retry reaquiring
		s.reaquireLoop()
		read, _ = s.conn.Write(data) // TODO just reaquired what happens if this fails...? reaquire again?

	}

	return read
}

func (s *SharedConn) reaquireLoop() error {
	for {
		err := s.reaquire()
		if err != nil {
			log.Printf("Error during reaquiring the sharedConn: %s", err)
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (s *SharedConn) reaquire() error {
	if s.server {
		listener, err := net.Listen(s.proto, s.address)
		if err != nil {
			return err
		}
		defer listener.Close()

		s.conn, err = listener.Accept()
		if err != nil {
			return err
		}

	} else {
		s.conn, err := net.Dial(s.proto, s.address)
		return err
	}
}
