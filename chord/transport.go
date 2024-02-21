package chord

import (
	"net"
	"sync"
	"time"

	"chord/pb"

	"google.golang.org/grpc"
)

// TODO: figure out a way to use channels to alert that grpc server failed for some reason.
// TODO: need to verify if using the connection pool will screw up any other client connection
// 		 in case the gRPC server is multiplexed.
// TODO: verify if gRPC multiplexing actually does work ok. When it does work, we can make
//		 multiplexing into an actual feature.

type transport struct {
	conf     *Config
	listener *net.TCPListener
	server   *grpc.Server
	pool     map[string]*connection
	poolLock sync.RWMutex
}

type connection struct {
	address    string
	client     pb.ChordClient
	conn       *grpc.ClientConn
	connLock   sync.RWMutex
	lastActive time.Time
}

// Constructor for the transport abstraction. There are two options
// for initializing the server. One might provide an exisiting gRPC
// server to which the Chord service will be registered, OR one might
// just simply have the gRPC service serve only for Chord.
func createTransport(conf *Config) (*transport, error) {
	listener, err := net.Listen("tcp", conf.Address)
	if err != nil {
		return nil, err
	}
	t := &transport{
		conf:     conf,
		listener: listener.(*net.TCPListener),
		server:   nil,
		pool:     make(map[string]*connection),
	}
	if conf.Server == nil {
		t.server = grpc.NewServer(conf.ServerOptions...)
	} else {
		t.server = conf.Server
		t.server.RegisterService(&pb.Chord_ServiceDesc, nil)
	}
	return t, nil
}

// might need to return channel here to keep track of failed gRPC listener

// Start gRPC server
func (t *transport) startServer() {

}

// Stop gRPC server
func (t *transport) stopServer() {

}

// Function to make sure that the gRPC server listening on a different go routine
// is still alive. This function will also be routinely called in a background
// go routine to ensure liveliness of the system.
func (t *transport) checkIfServerDead() bool {
	return false
}

// Function to be used to garbage collect unused gRPC connections. This will be
// routinely called in a background go routine.
func (t *transport) connectionGarbageCollect() {

}

// Get a client connection and cache it in a map.
func (t *transport) getClient(address string) (pb.ChordClient, error) {
	return nil, nil
}

// gRPC Client Implementation
