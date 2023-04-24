package main

import (
	"fmt"
	"net"

	"golang.org/x/sync/errgroup"
)

var (
	packedData = make(chan Frame)
)

//                          		 connList
// -net.Conn-> packedReader -Frame-> Router <-Frame +UnpackedReader
//             packedWriter  <-Frame-       -> Frame UnpackedWriter

func runProvider() {
	proxy = getProxyConn
	errs := errgroup.Group()
	errs.Go()
	packager(connList, proxy)
	unpackager(proxy, connList)
}

func startProvider() {
	fmt.Println("Connecting to proxy: ", proxyAddress)
	proxy, err := net.Dial("tcp", proxyAddress)
	if err != nil {
		panic(err)
	}
	go createConnections()
	splitSharedConn(proxy, getConnection)
}

func createConnections() {
	conn, err := net.Dial("tcp", exposeAddress)
	if err != nil {
		panic(err)
	}
	newConnections <- conn
}

// on provider: connect to proxy
// wait for packages,
// if one is received check if it belongs to an open connection then copy it into the given connection
// if not open a new connection to target and start a goroutine that copies all answers retrieved to the shared connection.

// if either side closes the connection a close message must be shared to via the shared channel to the other end.
