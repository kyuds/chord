package chord

import (
	"context"
	"fmt"
	"hash"
	"sync"

	"github.com/kyuds/go-chord/pb"
)

type node struct {
	// chord node settings
	id string
	ip string
	hf func() hash.Hash

	// successors predecessors
	repli    int
	succ     *successors
	succLock sync.RWMutex
	pred     string
	predLock sync.RWMutex
	ft       fingerTable
	ftLen    int

	// grpc
	pb.UnimplementedChordServer
	rpc rpc
}

func initialize(conf *Config) (*node, error) {
	n := &node{
		id:    getHash(conf.Hash, conf.Address),
		ip:    conf.Address,
		hf:    conf.Hash,
		repli: conf.MaxRepli,
		ftLen: conf.Hash().Size() * 8,
	}

	n.ft = initFingerTable(n.id, n.ip, n.ftLen)
	n.pred = ""
	n.succ = createSuccessorQueue(n.ip)

	// start RPC server: register chord server; start rpc server
	tmpRPC, err := newRPC(conf)
	if err != nil {
		return nil, err
	}
	n.rpc = tmpRPC
	pb.RegisterChordServer(tmpRPC.server, n)
	n.rpc.start()

	if conf.Joining {
		err = n.joinNode(conf.JoinAddress)
		if err != nil {
			n.rpc.stop()
			return nil, err
		}
	}

	// background processes
	// stabilize
	go func() {
		// check predecessor existence --> clear predecessor if so
		// clear out all ft table whenever node is determined to be not active
		// normal stabilization process.
		// update successors accordingly --> always maintain correct queue.

	}()

	// fix fingers
	go func() {

	}()

	return n, nil
}

func (n *node) joinNode(address string) error {
	// communicate to address
	// check if hash function checksum align
	// set predecessors and successors.
	h, err := n.rpc.getHashFuncCheckSum(address)
	if err != nil {
		return err
	}
	if h != getFuncHash(n.hf) {
		return fmt.Errorf("using different hash function from Chord ring.")
	}

	return nil
}

func (n *node) findSuccessor(hashedKey string) (string, error) {
	pred, err := n.findPredecessor(hashedKey)
	if err != nil {
		return "", err
	}

	if pred == n.ip {
		return n.succ.getSuccessor(), nil
	}

	// succ, err := n.rpc.getSuccessor(pred)
	// if err != nil {
	// 	n.ft.invalidateAddress(pred)
	// 	return "", err
	// }
	//return succ, nil
	return "", nil
}

func (n *node) findPredecessor(hashed string) (string, error) {
	curr := n.ip
	succ := n.succ.getSuccessor()
	bigKey := bigify(hashed)
	//var err error

	for {
		if curr == n.ip {
			succ = n.succ.getSuccessor()
		} else {
			// succ, err = n.rpc.getSuccessor(curr)
			// if err != nil {
			// 	n.ft.invalidateAddress(curr)
			// 	return "", err
			// }
		}
		if bigBetweenRightInclude(bigify(getHash(n.hf, curr)), bigify(getHash(n.hf, succ)), bigKey) {
			break
		}
		if curr == n.ip {
			curr = n.closestPrecedingFinger(hashed)
		} else {
			// curr, err = n.rpc.closestPrecedingFinger(curr, hashed)
			// if err != nil {
			// 	n.ft.invalidateAddress(curr)
			// 	return "", err
			// }
		}
	}
	return curr, nil
}

func (n *node) closestPrecedingFinger(hashed string) string {
	c1 := bigify(n.id)
	c2 := bigify(hashed)
	for i := n.ftLen - 1; i >= 0; i-- {
		f := n.ft.get(i)
		if !f.valid {
			continue
		}
		if bigBetween(c1, c2, f.id) {
			return f.ipaddr
		}
	}
	return n.ip
}

// gRPC Server (chord_grpc.pb.go) Implementation
func (n *node) GetHashFuncCheckSum(ctx context.Context, r *pb.Empty) (*pb.HashFuncResponse, error) {
	hashValue := getFuncHash(n.hf)
	return &pb.HashFuncResponse{HexHashValue: hashValue}, nil
}

func (n *node) GetSuccessor(ctx context.Context, r *pb.Empty) (*pb.AddressResponse, error) {
	return &pb.AddressResponse{Address: n.succ.getSuccessor()}, nil
}

func (n *node) ClosestPrecedingFinger(ctx context.Context, r *pb.HashKeyRequest) (*pb.AddressResponse, error) {
	a := n.closestPrecedingFinger(r.HexHashValue)
	return &pb.AddressResponse{Address: a}, nil
}

func (n *node) FindSuccessor(ctx context.Context, r *pb.HashKeyRequest) (*pb.AddressResponse, error) {
	a, err := n.findSuccessor(r.HexHashValue)
	if err != nil {
		return &pb.AddressResponse{Address: ""}, nil
	}
	return &pb.AddressResponse{Address: a}, nil
}

func (n *node) GetPredecessor(ctx context.Context, r *pb.Empty) (*pb.AddressResponse, error) {
	return &pb.AddressResponse{Address: n.pred}, nil
}

func (n *node) Notify(ctx context.Context, r *pb.AddressRequest) (*pb.Empty, error) {
	t1 := bigify(getHash(n.hf, n.pred))
	t2 := bigify(getHash(n.hf, n.ip))
	t3 := bigify(getHash(n.hf, r.Address))

	if n.pred == "" || bigBetween(t1, t2, t3) {
		n.predLock.Lock()
		n.pred = r.Address
		n.predLock.Unlock()
	}
	return &pb.Empty{}, nil
}

func (n *node) CheckPredecessor(ctx context.Context, r *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
