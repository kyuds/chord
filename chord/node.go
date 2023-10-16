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
	hashLen int
	hf func() hash.Hash
	ft fingertable

	// grpc
	pb.UnimplementedChordServer
	rpc rpc

	// chord ring configs
	successor string
	predecessor string
}

func newNode(conf *Config) (*node, error) {
	n := &node{
		id: getHash(conf.Hash, conf.Address),
		ip: conf.Address,
		hashLen: conf.Hash().Size(),
		hf: conf.Hash,
	}

	// set up finger table
	// set myself as successor & predecessor
	// start RPC server:
	// - register chord server
	// - start commLayer
	tmpRPC, err := newRPC(conf)
	if err != nil { return nil, err }
	n.rpc = tmpRPC
	pb.RegisterChordServer(tmpRPC.server, n)
	n.rpc.start()


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
	out, err := n.rpc.getHashFuncCheckSum(address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(out)

	return nil
}

// gRPC Server (chord_grpc.pb.go) Implementation
func (n *node) GetHashFuncCheckSum(ctx context.Context, e *pb.Empty) (*pb.HashFuncResponse, error) {
	hashValue := getFuncHash(n.hf)
	fmt.Printf("sending: %s\n", hashValue)
	return &pb.HashFuncResponse{HashVal: hashValue}, nil
}
