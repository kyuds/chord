package chord

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kyuds/go-chord/pb"
	"google.golang.org/grpc"
)

// Connection/Communication Management
// Creating a pool server for gRPC connections.
type rpc interface {
	// setup
	start() error
	stop() error
	getHashFuncCheckSum(string) (string, error)

	// chord
	getSuccessor(string) (string, error)
	closestPrecedingFinger(string, string) (string, error)

	// join
	findSuccessor(string, string) (string, error)
	getPredecessor(string) (string, error)
	notify(string, string) error

	// failure
	ping(string) error
}

type rpcLayer struct {
	listener    *net.TCPListener
	server      *grpc.Server
	pool        map[string]*gConn
	poolLock    sync.RWMutex
	timeout     time.Duration
	maxidle     time.Duration
	dialOptions []grpc.DialOption
	shutdown    int32
}

type gConn struct {
	address    string
	client     pb.ChordClient
	conn       *grpc.ClientConn
	lastActive time.Time
}

func newRPC(conf *Config) (*rpcLayer, error) {
	lis, err := net.Listen("tcp", conf.Address)
	if err != nil {
		return nil, err
	}
	r := &rpcLayer{
		listener:    lis.(*net.TCPListener),
		server:      grpc.NewServer(conf.ServerOptions...),
		pool:        make(map[string]*gConn),
		timeout:     conf.Timeout,
		maxidle:     conf.MaxIdle,
		dialOptions: conf.DialOptions,
		shutdown:    0,
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
	if err != nil {
		return nil, err
	}

	client := pb.NewChordClient(conn)
	gconn = &gConn{
		address:    address,
		client:     client,
		conn:       conn,
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
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := client.GetHashFuncCheckSum(ctx, &pb.Empty{})
	if err != nil {
		return "", err
	}

	return res.HexHashValue, nil
}

func (r *rpcLayer) getSuccessor(address string) (string, error) {
	client, err := r.getClient(address)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := client.GetSuccessor(ctx, &pb.Empty{})
	if err != nil {
		return "", err
	}

	return res.Address, nil
}

func (r *rpcLayer) closestPrecedingFinger(address, key string) (string, error) {
	client, err := r.getClient(address)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := client.ClosestPrecedingFinger(ctx, &pb.HashKeyRequest{HexHashValue: key})
	if err != nil {
		return "", err
	}

	return res.Address, nil
}

func (r *rpcLayer) findSuccessor(address, key string) (string, error) {
	client, err := r.getClient(address)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := client.FindSuccessor(ctx, &pb.HashKeyRequest{HexHashValue: key})
	if err != nil {
		return "", err
	}

	return res.Address, nil
}

func (r *rpcLayer) getPredecessor(address string) (string, error) {
	client, err := r.getClient(address)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := client.GetPredecessor(ctx, &pb.Empty{})
	if err != nil {
		return "", err
	}

	return res.Address, nil
}

func (r *rpcLayer) notify(address, key string) error {
	client, err := r.getClient(address)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err = client.Notify(ctx, &pb.AddressRequest{Address: key})
	if err != nil {
	}
	return err
}

func (r *rpcLayer) ping(address string) error {
	client, err := r.getClient(address)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err = client.Ping(ctx, &pb.Empty{})
	return err
}
