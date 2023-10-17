package chord

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"google.golang.org/grpc"
	"github.com/kyuds/go-chord/pb"
)

/*
Bindings for gRPC protocol in /rpc.

gRPC is fundamentally a client, server model. Therefore, typical gRPC programs implement
a client and a server to provide one-way communication (at least in most use cases) in
which a client sends a message to a server, the server performs some action, and then
responds with a response packet.

In this implementation of Chord, the connection logic for both the client and the server
is implemented in rpc.go, and the server side processing for gRPC requests are defined in
node.go (under "gRPC Server (chord_grpc.pb.go) Implementation" section). One can easily 
identify sections of the server (which is defined in rpcLayer) and corresponding functions
of start() and stop() to start and stop the server, and the client (which is defined as gConn)
which is instantiated with getClient(). The clients are stored in a pool map so that
repetitive queries can easily be established. 

Finally, the rpc interface provides an abstraction layer between bare metal gRPC invocations
and the processing of these responses to return user-familiar definitions to the code. These
definitions are under "gRPC Abstraction"
*/

// Connection/Communication Management
// Creating a pool server for gRPC connections. 
type rpc interface {
	// setup
	start() error
	stop() error
	getHashFuncCheckSum(string) (string, error)

	// chord
	getSuccessor(string) (string, error)

	// hash table
}

type rpcLayer struct {
	listener *net.TCPListener
	server *grpc.Server
	pool map[string]*gConn
	poolLock sync.RWMutex
	timeout time.Duration
	maxidle time.Duration
	dialOptions []grpc.DialOption
	shutdown int32
}

type gConn struct {
	address string
	client pb.ChordClient
	conn *grpc.ClientConn
	lastActive time.Time
}

func newRPC(conf *Config) (*rpcLayer, error) {
	lis, err := net.Listen("tcp", conf.Address)
	if err != nil { return nil, err }
	r := &rpcLayer {
		listener: lis.(*net.TCPListener),
		server: grpc.NewServer(conf.ServerOptions...),
		pool: make(map[string]*gConn),
		timeout: conf.Timeout,
		maxidle: conf.MaxIdle,
		dialOptions: conf.DialOptions,
		shutdown: 0,
	}
	return r, nil
}

func (r *rpcLayer) start() error {
	// start server to listen to port
	go func() {
		r.server.Serve(r.listener)
	}()

	// garbage collect timed out connections
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for {
			if atomic.LoadInt32(&r.shutdown) == 1 {
				return
			}
			select {
			case <-ticker.C:
				r.poolLock.Lock()
				for ip, conn := range r.pool {
					if time.Since(conn.lastActive) > r.maxidle {
						conn.conn.Close()
						delete(r.pool, ip)
					}
				}
				r.poolLock.Unlock()
			}
		}
	}()
	return nil
}

func (r *rpcLayer) stop() error {
	atomic.StoreInt32(&r.shutdown, 1)
	r.poolLock.Lock()
	defer r.poolLock.Unlock()
	r.server.Stop()
	for _, conn := range r.pool {
		conn.conn.Close()
	}
	r.pool = nil
	return nil
}

func (r *rpcLayer) getClient(address string) (pb.ChordClient, error) {
	r.poolLock.RLock()
	if atomic.LoadInt32(&r.shutdown) == 1 {
		return nil, fmt.Errorf("Server is shutdown.")
	}
	gconn, ok := r.pool[address]
	r.poolLock.RUnlock()

	if ok {
		return gconn.client, nil
	}

	conn, err := grpc.Dial(address, r.dialOptions...)
	if err != nil { return nil, err }

	client := pb.NewChordClient(conn)
	gconn = &gConn{
		address: address,
		client: client,
		conn: conn,
		lastActive: time.Now(),
	}
	r.poolLock.Lock()
	r.pool[address] = gconn
	r.poolLock.Unlock()
	return client, nil
}

// gRPC Abstraction
func (r *rpcLayer) getHashFuncCheckSum(address string) (string, error) {
	client, err := r.getClient(address)
	if err != nil { return "", err }
	
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := client.GetHashFuncCheckSum(ctx, &pb.Empty{})
	if err != nil { return "", err }

	return res.HashVal, nil
}

func (r *rpcLayer) getSuccessor(address string) (string, error) {
	client, err := r.getClient(address)
	if err != nil { return "", err }
	
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := client.GetSuccessor(ctx, &pb.Empty{})
	if err != nil { return "", err }

	return res.Address, nil
}
