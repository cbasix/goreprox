package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
)

var (
	mode                string
	exposeAddress       string
	proxyAddress        string
	connectionCount     int
	providerConnections chan net.Conn
)

var (
	newConnections = make(chan net.Conn)
	connections    = make(map[uint]net.Conn)
	mutex          = new(sync.RWMutex)
)

// Parse command line flags
func init() {
	flag.StringVar(&mode, "mode", "proxy", "the mode that is started 'proxy' or 'provider' (default: proxy)")
	flag.IntVar(&connectionCount, "connectionCount", 5, "amount of connections keept ready")
	flag.StringVar(&exposeAddress, "exposeAddress", ":8080", "the interface:port that will we forwarded to the provider)")
	flag.StringVar(&proxyAddress, "providerAddress", ":9887", "the interface:port the provider connects to")
}

/*
Run goreproxy.
Depending in the mode flag this either starts the Proxy or the Provider.

Note to devs: both use the same flags and vars but for different purposes.
*/
func main() {
	flag.Parse()
	providerConnections = make(chan net.Conn, connectionCount)

	if mode == "proxy" {
		startProxy()
	} else {
		//startProvider()
	}
}

// Start the Proxy that listens for incoming Connections from Providers and clients
func startProxy() {
	fmt.Println("Starting reproxy on:", exposeAddress, " waiting for provider on:", proxyAddress)
	var sharedConn net.Conn
	sharedIn := make(chan Frame, 20)
	sharedOut := make(chan Frame, 20)

	go packedWriter(sharedConn, sharedOut)
	go packedReader(sharedConn, sharedIn)

	router := CreateRouter(sharedIn, sharedOut)

	go router.join()
	//go startListening(handleExposed, exposeAddress)
	//startListening(handleProvider, proxyAddress)
}

/*
func handleProvider(conn net.Conn) {
	fmt.Println("New provider connection: ", conn.RemoteAddr().String())
	providerConnections <- conn
}

func handleExposed(conn net.Conn) {
	fmt.Println("New client: ", conn.RemoteAddr().String())
	provider := <-providerConnections
	go copyIO(conn, provider)
	go copyIO(provider, conn)
}

func startProvider() {
	for {
		fmt.Println("Connecting to proxy: ", proxyAddress)
		proxy, err := net.Dial("tcp", proxyAddress)
		if err != nil {
			panic(err)
		}
		providerConnections <- proxy

		fmt.Println("Connecting to exposed: ", exposeAddress)
		provider, err := net.Dial("tcp", exposeAddress)
		if err != nil {
			panic(err)
		}

		fmt.Println("Proxy <> provider connection established")
		go copyIO(provider, proxy)
		go copyIoWithCloseFree(proxy, provider)
	}
}

func copyIoWithCloseFree(proxy, provider net.Conn) {
	copyIO(proxy, provider)
	<-providerConnections // make place for a new one
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	defer fmt.Println("Conn closed:", dest.RemoteAddr().String())
	io.Copy(src, dest)
}
*/
