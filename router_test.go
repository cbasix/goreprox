package main

import "testing"

func TestRouterRouting(t *testing.T) {
	shared := CreateConn(2, 2)
	client1 := CreateConn(1, 1)
	client2 := CreateConn(1, 1)
	router := CreateRouter(shared)
	router.createDest(client1)
	router.createDest(client2)

	shared.in <- Frame{ConnectionId: 1, Data: []byte{'1'}}
	shared.in <- Frame{ConnectionId: 2, Data: []byte{'2'}}

	for i := 0; i < 2; i++ {
		router.route()
	}

	if (<-client1.out).Data[0] != '1' {
		t.Errorf("Client 1 received invalid frame.")
	}

	if (<-client2.out).Data[0] != '2' {
		t.Errorf("Client 2 received invalid frame.")
	}
}

func TestRouterJoin(t *testing.T) {
	shared := CreateConn(2, 2)
	client1 := CreateConn(1, 1)
	client2 := CreateConn(1, 1)
	router := CreateRouter(shared)
	router.createDest(client1)
	router.createDest(client2)

	client1.in <- Frame{ConnectionId: 1, Data: []byte{1}}
	client2.in <- Frame{ConnectionId: 2, Data: []byte{2}}

	router.join()

	<-shared.out
	<-shared.out
}

func TestRouterRouteOrCreate(t *testing.T) {
	shared := CreateConn(2, 2)
	client := CreateConn(1, 1)
	router := CreateRouter(shared)
	router.createConnection = func() ChannelConn { return client }

	shared.in <- Frame{ConnectionId: 1, Data: []byte{'1'}}
	router.route()

	<-client.out
}

func TestMapFunctions(t *testing.T) {
	shared := CreateConn(0, 0)
	client := CreateConn(0, 0)
	router := CreateRouter(shared)

	connId := router.createDest(client)
	if connId != 1 {
		t.Errorf("Generated connId wrong. Expected %d got %d", 1, connId)
	}

	_, ok := router.getDest(1)
	if !ok {
		t.Errorf("Could not find connId 1 in conn map")
	}

	router.delDest(1)
	_, stillThere := router.getDest(1)
	if stillThere {
		t.Errorf("Deltion of connId 1 failed. It is still there")
	}

}
