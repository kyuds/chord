package chord

import (
	"fmt"
	"context"
	"hash"
	"github.com/kyuds/go-chord/rpc"
)

type node struct {
	// chord node settings
	id string
	ip string
	hashLen int
	hf func() hash.Hash
	ft fingertable

	// grpc
	rpc.UnimplementedChordServer
	cmm comm

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
	tmpCmm, err := newComm(conf)
	if err != nil { return nil, err }
	n.cmm = tmpCmm
	rpc.RegisterChordServer(tmpCmm.server, n)
	n.cmm.start()


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
	out, err := n.cmm.getHashFuncCheckSum(address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(out)

	return nil
}

// gRPC Server (chord_grpc.pb.go) Implementation
func (n *node) GetHashFuncCheckSum(ctx context.Context, e *rpc.Empty) (*rpc.HashFuncResponse, error) {
	hashValue := getFuncHash(n.hf)
	fmt.Printf("sending: %s\n", hashValue)
	return &rpc.HashFuncResponse{HashVal: hashValue}, nil
}
