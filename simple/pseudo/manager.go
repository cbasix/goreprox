package pseudo

import "net"

// listen for new "incomming connections"
// on server via "real server socket" on provider via "new connId listener?"

type manager interface {
	initialize(listener connectionListener)
	startListen()
}

type connectionListener interface {
	getConnectionChan() chan net.Conn
	// getErrorChan() chan error
	startListening()
}

// todo move to own?
type socketListener struct {
	listener net.Listener
}

// todo implement
func (s *socketListener) getConnectionChan() chan net.Conn {
	return nil
}

func (s *socketListener) startListening() {}

//-------------
type sharedListener struct {
	listener net.Conn
}

// todo implement
// waits for a frame and checks of
func (s *sharedListener) getConnectionChan() chan net.Conn {
	return nil
}

func (s *sharedListener) startListening() {}
