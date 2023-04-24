package main

import "testing"

func TestRouterRouting(t *testing.T) {
	shared := make(chan Frame, 2)
	client1 := make(chan Frame, 1)
	client2 := make(chan Frame, 1)
	router := CreateRouter(shared)
	router.createDest(client1)
	router.createDest(client2)

	shared <- Frame{ConnectionId: 1, Data: []byte{'1'}}
	shared <- Frame{ConnectionId: 2, Data: []byte{'2'}}
	go router.route()

	if (<-client1).Data[0] != '1' {
		t.Errorf("Client 1 received invalid frame.")
	}

	if (<-client2).Data[0] != '2' {
		t.Errorf("Client 2 received invalid frame.")
	}
}

func TestRouterJoin(t *testing.T) {
	shared := make(chan Frame, 2)
	client1 := make(chan Frame, 1)
	client2 := make(chan Frame, 1)
	router := CreateRouter(shared)
	router.createDest(client1)
	router.createDest(client2)

	client1 <- Frame{ConnectionId: 1, Data: []byte{1}}
	client2 <- Frame{ConnectionId: 2, Data: []byte{2}}
	go router.join()

	<-shared
	<-shared
}

func TestRouterRouteOrCreate(t *testing.T) {
	shared := make(chan Frame, 2)
	client := make(chan Frame, 1)
	router := CreateRouter(shared)

	shared <- Frame{ConnectionId: 1, Data: []byte{'1'}}
	go router.routeOrCreate(func() chan Frame { return client })

	<-client
}

func TestMapFunctions(t *testing.T) {
	shared := make(chan Frame)
	client := make(chan Frame)
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
