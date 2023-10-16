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

// bindings for gRPC protocol in /rpc.

// Connection/Communication Management
// Creating a pool server for gRPC connections. 
type comm interface {
	// setup
	start() error
	stop() error
	getHashFuncCheckSum(string) (string, error)

	// chord

	// hash table
}

type commLayer struct {
	listener *net.TCPListener
	server *grpc.Server
	pool map[string]*gConn
	poolLock sync.RWMutex
	timeout time.Duration
	maxidle time.Duration
	shutdown int32
}

type gConn struct {
	address string
	client pb.ChordClient
	conn *grpc.ClientConn
	lastActive time.Time
}

func newComm(conf *Config) (*commLayer, error) {
	lis, err := net.Listen("tcp", conf.Address)
	if err != nil { return nil, err }

	c := &commLayer {
		listener: lis.(*net.TCPListener),
		server: grpc.NewServer(/*conf.ServerOptions...*/),
		pool: make(map[string]*gConn),
		timeout: conf.Timeout,
		maxidle: conf.MaxIdle,
		shutdown: 0,
	}

	return c, nil
}

func (c *commLayer) start() error {
	// start server to listen to port
	go func() {
		c.server.Serve(c.listener)
	}()

	// garbage collect timed out connections
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for {
			if atomic.LoadInt32(&c.shutdown) == 1 {
				return
			}
			select {
			case <-ticker.C:
				c.poolLock.Lock()
				for ip, conn := range c.pool {
					if time.Since(conn.lastActive) > c.maxidle {
						conn.close()
						delete(c.pool, ip)
					}
				}
				c.poolLock.Unlock()
			}
		}
	}()
	return nil
}

func (c *commLayer) stop() error {
	atomic.StoreInt32(&c.shutdown, 1)
	c.poolLock.Lock()
	defer c.poolLock.Unlock()
	c.server.Stop()
	for _, conn := range c.pool {
		conn.close()
	}
	c.pool = nil
	return nil
}

func (c *commLayer) getClient(address string) (pb.ChordClient, error) {
	c.poolLock.RLock()
	if atomic.LoadInt32(&c.shutdown) == 1 {
		return nil, fmt.Errorf("Server is shutdown.")
	}
	gconn, ok := c.pool[address]
	c.poolLock.RUnlock()

	if ok {
		return gconn.client, nil
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil { return nil, err }

	client := pb.NewChordClient(conn)
	gconn = &gConn{
		address: address,
		client: client,
		conn: conn,
		lastActive: time.Now(),
	}
	c.poolLock.Lock()
	c.pool[address] = gconn
	c.poolLock.Unlock()
	return client, nil
}

func (c * gConn) close() {
	c.conn.Close()
}

// gRPC Abstraction
func (c *commLayer) getHashFuncCheckSum(address string) (string, error) {
	client, err := c.getClient(address)
	if err != nil { return "", err }
	
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	r, err := client.GetHashFuncCheckSum(ctx, &pb.Empty{})
	if err != nil { return "", err }

	return r.HashVal, nil
}
