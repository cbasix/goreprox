// towards the router is "in" away from the router is "out"
//
//	connList
//
// -net.Conn-> packedReader -Frame-> Router <-Frame +UnpackedReader
//
//	packedWriter  <-Frame-       -> Frame UnpackedWriter
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	mode          string
	exposeAddress string
	proxyAddress  string
	logLevel      int
	//connectionCount     int
	//providerConnections chan net.Conn
)

var (
	newConnections = make(chan net.Conn)
	connections    = make(map[uint]net.Conn)
	mutex          = new(sync.RWMutex)
)

// Parse command line flags
func init() {
	flag.StringVar(&mode, "mode", "proxy", "the mode that is started 'proxy' or 'provider' (default: proxy)")
	//flag.IntVar(&connectionCount, "connectionCount", 5, "amount of connections keept ready")
	flag.StringVar(&exposeAddress, "expose", ":8080", "the interface:port that will we forwarded to the provider)")
	flag.StringVar(&proxyAddress, "proxy", ":9887", "the interface:port the provider connects to")
	flag.IntVar(&logLevel, "v", 1, "the log level to use 0=Nothing 1=Info 2=Warn 3=Debug")
}

/*
Run goreproxy.
Depending in the mode flag this either starts the Proxy or the Provider.

Note to devs: both use the same flags and vars but for different purposes.
*/
func main() {
	parseFlags()

	if mode == "proxy" {
		startProxy()
	} else {
		startProvider()
	}
}

func parseFlags() {
	flag.Parse()
	if !strings.Contains(exposeAddress, ":") {
		exposeAddress = ":" + exposeAddress
	}
	if !strings.Contains(proxyAddress, ":") {
		proxyAddress = ":" + proxyAddress
	}
}

// Start the Proxy that listens for incoming Connections from Providers and clients
func startProxy() {
	shared := CreateConn(50, 50)
	// TODO error handing. Idea: global error chan?

	fmt.Println("Starting reproxy on:", exposeAddress, " waiting for provider on:", proxyAddress)
	sharedConn, err := listenOnce("tcp", proxyAddress)
	if err != nil {
		panic(fmt.Sprintf("Error on socket listen: %s", err))
	}
	log.Printf("Provider connected: %s", sharedConn.RemoteAddr().String())

	router := CreateRouter(shared)

	go func() {
		for {
			router.join()
		}
	}() // TODO stop handling#
	go func() {
		for {
			router.route()
		}
	}()

	// init shared conn plumbings
	go packedWriter(sharedConn, shared.out)
	go packedReader(sharedConn, shared.in)

	// listen for clients
	connections := make(chan net.Conn)
	go listenLoop("tcp", exposeAddress, connections)
	log.Printf("Listening for clients")

	for conn := range connections {
		client := CreateConn(50, 50)
		client.conn = conn
		go unpackedWriter(conn, client.out)
		go unpackedReader(conn, client.in)

		router.createDest(client)
	}
}

// Start the Provider that connects to the Proxy and to the target system as soon as it'needed
func startProvider() {
	shared := CreateConn(50, 50)
	// TODO error handing. Idea: global error chan?

	log.Println("Starting provider. Connecting to proxy on:", proxyAddress)
	sharedConn, err := net.Dial("tcp", proxyAddress)
	if err != nil {
		panic(fmt.Sprintf("Error on connecting to proxy: %s", err))
	}
	log.Printf("Connected to proxy: %s", sharedConn.RemoteAddr().String())

	router := CreateRouter(shared)
	router.createConnection = func() *ChannelConn {
		log.Printf("Creating target conn to: %s", exposeAddress)
		client := CreateConn(5, 5)

		clientConn, err := net.Dial("tcp", exposeAddress)
		if err != nil {
			panic(fmt.Sprintf("Error on connecting to target: %s", err))
		}
		client.conn = clientConn

		go unpackedWriter(clientConn, client.out)
		go unpackedReader(clientConn, client.in)

		return client
	}

	// init shared conn plumbings
	go packedWriter(sharedConn, shared.out)
	go packedReader(sharedConn, shared.in)

	go func() {
		for {
			router.join()
		}
	}() // TODO stop handling#

	for {
		router.route()
	}
}

/* Listens for one connection and returns it */
func listenOnce(proto string, address string) (net.Conn, error) {
	listener, err := net.Listen(proto, address)
	defer listener.Close()

	if err != nil {
		return nil, err
	}
	sharedConn, err := listener.Accept()
	return sharedConn, err
}

/* Endlessly listens for new connections and puts them into the connections channel */
func listenLoop(proto string, address string, connections chan<- net.Conn) error {
	listener, err := net.Listen(proto, address)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error in connection %+v\n%+v", conn, err)
		} else {
			log.Printf("New client connection %v", conn.RemoteAddr().String())
			connections <- conn
		}
	}

}

func logD(msg string, args ...interface{}) {
	if logLevel >= 3 { // TODO prettify or remove
		log.Printf(msg, args...)
	}
}
