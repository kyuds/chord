package chord

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"chord/pb"

	"google.golang.org/grpc"
)

// Concurrency-safe transport wrapper for gRPC. Heavily uses
// gRPC's concurrency guarantees, but also enables for efficient
// garbage collection of idle connections, etc. Also establishes
// one connection per ip address and allows for multiple uses of
// such connections.

// TODO: need to verify if using the connection pool will screw up any other client connection
// 		 in case the gRPC server is multiplexed.
// TODO: verify if gRPC multiplexing actually does work ok. When it does work, we can make
//		 multiplexing into an actual feature.

// TCP Transport Abstraction. Backend of communication scheme
// is gRPC. The transport has a pool of persisted gRPC connections
// that can be used multiple times so that new connections can be
// avoided.
type transport struct {
	conf     *Config
	listener *net.TCPListener
	server   *grpc.Server
	alive    atomic.Bool
	pool     map[string]*connection
	poolLock sync.RWMutex
}

// Connection abstraction. This solely exists because old enough
// connections must be cleaned up.
type connection struct {
	address    string
	client     pb.ChordClient
	conn       *grpc.ClientConn
	lock       sync.RWMutex
	lastActive time.Time
}

// Constructor for the transport abstraction. There are two options
// for initializing the server. One might provide an exisiting gRPC
// server to which the Chord service will be registered, OR one might
// just simply have the gRPC service serve only for Chord.
func newTransport(conf *Config) (*transport, error) {
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
	t.alive.Store(false)
	if conf.Server == nil {
		t.server = grpc.NewServer(conf.ServerOptions...)
	} else {
		// TODO: implement feature
		// t.server = conf.Server
		// t.server.RegisterService(&pb.Chord_ServiceDesc, nil)
		return nil, fmt.Errorf("grpc multiplexing not supported yet")
	}
	return t, nil
}

// Start gRPC server. We use an atomic boolean and not a channel
// to check liveliness since liveliness may be checked more than once.
func (t *transport) startServer() {
	// TODO: channel to send error returned from grpc server.Serve()?
	go func() {
		t.alive.Store(true)
		t.server.Serve(t.listener)
		t.alive.Store(false)
	}()
}

// Stop gRPC server
func (t *transport) stopServer() {
	t.server.Stop()
}

// Function to make sure that the gRPC server listening on a different go routine
// is still alive. This function will also be routinely called in a background
// go routine to ensure liveliness of the system. Isn't too taxing on performance
// as atomic Load operations are relatively inexpensive.
func (t *transport) checkServerDead() bool {
	return !t.alive.Load()
}

// Function to be used to garbage collect unused gRPC connections. This will be
// routinely called in a background go routine.
func (t *transport) cleanIdleConnections() {
	t.poolLock.Lock()
	defer t.poolLock.Unlock()
	// we first write lock the entire pool so that Chord cannot request
	// for a new client, and then we iterate through the pool with the ip
	// and connections. We use "TryLock", not "Lock" to check if we are
	// ABLE to get write permissions to the connection (this serves to
	// check that the connection is not being held with Read Lock at the
	// time of checking, which also means lastActive will be updated anyways)
	// and then check for time constraints and close the connection if applicable.
	for ip, conn := range t.pool {
		if conn.lock.TryLock() {
			if time.Since(conn.lastActive) > t.conf.MaxIdle {
				conn.conn.Close()
				delete(t.pool, ip)
			}
			conn.lock.Unlock()
		}
	}
}

// Get a client connection and cache it in a map.
// Note that getClient may significantly block is grpc is used with
// the WithBlock dial option.
func (t *transport) getConnection(address string) (*connection, error) {
	// Phase 1: Check if client already exists.
	t.poolLock.RLock()
	conn, ok := t.pool[address]
	if ok {
		conn.lock.RLock()
		t.poolLock.Unlock()
		return conn, nil
	}
	t.poolLock.Unlock()

	// Phase 2: Create a new connection and store it to the pool.
	t.poolLock.Lock()
	defer t.poolLock.Unlock()

	// check one more time whether connection is established.
	conn, ok = t.pool[address]
	if ok {
		conn.lock.RLock()
		return conn, nil
	}
	// note that this is non blocking because we don't have
	// the WithBlock dial option in the configuration struct.
	dial, err := grpc.Dial(address, t.conf.DialOptions...)
	if err != nil {
		return nil, err
	}
	client := pb.NewChordClient(dial)
	conn = &connection{
		address:    address,
		client:     client,
		conn:       dial,
		lastActive: time.Now(),
	}
	// save to pool and read lock the connection.
	t.pool[address] = conn
	conn.lock.RLock()

	return conn, nil
}

// Naive implementation for locking connections.
// THIS FUNCTION MUST BE CALLED WHEN CONNECTION HAS A LOCK HELD (RLOCK)
// RPCs hold a RLock on the connection, which must be upgraded to a
// write lock for the lastActive timestamp to be updated. TryLock is
// used such that the time is updated when the last RLock concurrently
// held is released.
func (c *connection) unlockAndTryUpdateTime() {
	// Keep in mind: between Unlock and TryLock, connection may be
	// deleted from the pool.
	c.lock.Unlock()
	if c.lock.TryLock() {
		c.lastActive = time.Now()
		c.lock.Unlock()
	}
}

// gRPC Client Implementation
// All rpcs should get connections using getConnection, which
// returns a connection with a RLock acquired. It is required
// by the rpc implementation to call unlockAndTryUpdateTime().

func (t *transport) getHashFuncName(address string) (string, error) {
	conn, err := t.getConnection(address)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.conf.Timeout)
	defer cancel()

	response, err := conn.client.GetHashFuncName(ctx, &pb.Empty{})
	if err != nil {
		conn.lock.Unlock()
		return "", err
	}
	conn.unlockAndTryUpdateTime()
	return response.GetHashFuncName(), nil
}

func (t *transport) getSuccessor(address string) (string, error) {
	conn, err := t.getConnection(address)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.conf.Timeout)
	defer cancel()

	response, err := conn.client.GetSuccessor(ctx, &pb.Empty{})
	if err != nil {
		conn.lock.Unlock()
		return "", err
	}
	conn.unlockAndTryUpdateTime()
	return response.GetAddress(), nil
}

func (t *transport) closestPrecedingFinger(address string, hash string) (string, error) {
	conn, err := t.getConnection(address)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.conf.Timeout)
	defer cancel()

	response, err := conn.client.ClosestPrecedingFinger(ctx, &pb.HashKeyRequest{HashValue: hash})
	if err != nil {
		conn.lock.Unlock()
		return "", err
	}
	conn.unlockAndTryUpdateTime()
	return response.GetAddress(), nil
}

func (t *transport) findSuccessor(address string, hash string) (string, error) {
	conn, err := t.getConnection(address)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.conf.Timeout)
	defer cancel()

	response, err := conn.client.FindSuccessor(ctx, &pb.HashKeyRequest{HashValue: hash})
	if err != nil {
		conn.lock.Unlock()
		return "", err
	}
	conn.unlockAndTryUpdateTime()
	return response.GetAddress(), nil
}

func (t *transport) getPredecessor(address string) (string, error) {
	conn, err := t.getConnection(address)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.conf.Timeout)
	defer cancel()

	response, err := conn.client.GetPredecessor(ctx, &pb.Empty{})
	if err != nil {
		conn.lock.Unlock()
		return "", err
	}
	conn.unlockAndTryUpdateTime()
	return response.GetAddress(), nil
}

func (t *transport) notify(address string, key string) error {
	conn, err := t.getConnection(address)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.conf.Timeout)
	defer cancel()

	_, err = conn.client.Notify(ctx, &pb.AddressRequest{Address: key})
	if err != nil {
		conn.lock.Unlock()
		return err
	}
	conn.unlockAndTryUpdateTime()
	return nil
}

func (t *transport) ping(address string) bool {
	conn, err := t.getConnection(address)
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.conf.Timeout)
	defer cancel()

	_, err = conn.client.Ping(ctx, &pb.Empty{})
	if err != nil {
		conn.lock.Unlock()
		return false
	}
	conn.unlockAndTryUpdateTime()
	return true
}
