package chord

import (
	"chord/pb"
	"context"
	"fmt"
	"time"
)

/*
For testing purposes only.
Testing transport struct
*/

type DummyNode struct {
	conf      *Config
	transport *transport
	cont      bool
	pb.UnimplementedChordServer
}

func NewDummy(conf *Config) *DummyNode {
	d := &DummyNode{
		conf: conf,
		cont: true,
	}
	t, err := newTransport(conf)
	if err != nil {
		panic("initializing transport failed")
	}
	d.transport = t
	pb.RegisterChordServer(t.server, d)
	return d
}

func (d *DummyNode) Start() {
	d.transport.startServer()
	go func(d *DummyNode) {
		cleanup := time.NewTicker(d.conf.CleanConnection)
		for {
			select {
			case <-cleanup.C:
				d.transport.cleanIdleConnections()
			}
			time.Sleep(10 * time.Millisecond)
		}
	}(d)
}

// to switch for "blocking" operations
func (d *DummyNode) Switch() {
	d.cont = !d.cont
}

// gRPC client handler for testing
func (d *DummyNode) CallRPC(address string) {
	ret, err := d.transport.getSuccessor(address)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ret)
	}
}

// mock gRPC server implementation
func (d *DummyNode) GetChordConfigs(ctx context.Context, r *pb.Empty) (*pb.ConfigResponse, error) {
	return &pb.ConfigResponse{
		HashFuncName:        d.conf.HashName(),
		SuccessorListLength: int32(d.conf.NumSuccessor),
	}, nil
}

func (d *DummyNode) GetSuccessor(ctx context.Context, r *pb.Empty) (*pb.AddressResponse, error) {
	return &pb.AddressResponse{
		Present: d.cont,
		Address: "GetSuccessor",
	}, nil
}

func (d *DummyNode) ClosestPrecedingFinger(ctx context.Context, r *pb.HashKeyRequest) (*pb.AddressResponse, error) {
	return &pb.AddressResponse{
		Present: true,
		Address: "ClosestPrecedingFinger",
	}, nil
}

func (d *DummyNode) FindPredecessor(ctx context.Context, r *pb.HashKeyRequest) (*pb.AddressResponse, error) {
	return &pb.AddressResponse{
		Present: true,
		Address: "FindPredecessor",
	}, nil
}

func (d *DummyNode) GetPredecessor(ctx context.Context, r *pb.Empty) (*pb.AddressResponse, error) {
	return &pb.AddressResponse{
		Present: true,
		Address: "GetPredecessor",
	}, nil
}

func (d *DummyNode) GetSuccessorList(ctx context.Context, r *pb.Empty) (*pb.AddressListResponse, error) {
	return &pb.AddressListResponse{
		Present:   true,
		Addresses: nil,
		Length:    5,
	}, nil
}

func (d *DummyNode) Notify(ctx context.Context, r *pb.AddressRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (d *DummyNode) Ping(ctx context.Context, r *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
