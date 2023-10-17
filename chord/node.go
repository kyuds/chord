package chord

import (
	"fmt"
	"context"
	"hash"
	"github.com/kyuds/go-chord/pb"
)

type node struct {
	// chord node settings
	id string
	ip string
	hf func() hash.Hash

	// fingertable
	ft fingerTable
	ftLen int

	// grpc
	pb.UnimplementedChordServer
	rpc rpc
}

func newNode(conf *Config) (*node, error) {
	n := &node{
		id: getHash(conf.Hash, conf.Address),
		ip: conf.Address,
		hf: conf.Hash,
		ftLen: conf.Hash().Size(),
	}
	n.ft = initFingerTable(n.id, n.ip, n.ftLen)
	// set myself as successor & predecessor

	// start RPC server: register chord server; start rpc server
	tmpRPC, err := newRPC(conf)
	if err != nil { return nil, err }
	n.rpc = tmpRPC
	pb.RegisterChordServer(tmpRPC.server, n)
	n.rpc.start()

	if (conf.Joining) {
		err = n.joinNode(conf.JoinIP)
		if err != nil {
			n.rpc.stop()
			return nil, err
		}
	}

	// run background:
	// - stabilize
	// - fix fingers
	// - check predecessor

	return n, nil
}

func (n *node) joinNode(address string) error {
	// communicate to address
	// check if hash function checksum align
	// set predecessors and successors. 
	h, err := n.rpc.getHashFuncCheckSum(address)
	if err != nil { return err }
	if h != getFuncHash(n.hf) {
		return fmt.Errorf(sameHash)
	}

	return nil
}

// predecessor and successor operations (Chord p.5)
func (n *node) findSuccessor(key string) (string, error) {
	pred, err := n.findPredecessor(key)
	if err != nil { return "", err }
	succ, err := n.rpc.getSuccessor(pred)
	if err != nil { return "", err }
	return succ, nil
}

func (n *node) findPredecessor(key string) (string, error) {

}

func (n *node) closestPrecedingFinger(key string) (string, error) {

}

// gRPC Server (chord_grpc.pb.go) Implementation
func (n *node) GetHashFuncCheckSum(ctx context.Context, r *pb.Empty) (*pb.HashFuncResponse, error) {
	hashValue := getFuncHash(n.hf)
	fmt.Printf("sending: %s\n", hashValue)
	return &pb.HashFuncResponse{HashVal: hashValue}, nil
}

func (n *node) GetSuccessor(ctx context.Context, r *pb.Empty) (*pb.AddressResponse, error) {
	return &pb.AddressResponse{Address: n.ft.getSuccessor()}, nil
}

func (n *node) FindPredecessor(ctx context.Context, r *pb.HashKeyRequest) (*pb.AddressResponse, error) {

	return &pb.AddressResponse{Address: ""}, nil
}

// Error Messages
var (
	sameHash = "Checksum of hash functions on the nodes do not match. Please check that the same hash functions are being used."
)
