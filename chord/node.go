package chord

import (
	"fmt"
	"context"
	"encoding/hex"
	"hash"
	"sync"
	"time"
	"github.com/kyuds/go-chord/pb"
)

type node struct {
	// chord node settings
	id string
	ip string
	hf func() hash.Hash
	predecessor string
	predLock sync.RWMutex

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
		ftLen: conf.Hash().Size() * 8,
	}
	n.ft = initFingerTable(n.id, n.ip, n.ftLen)
	n.predecessor = ""

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
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				n.stabilize()
			/*
			case <-node.shutdownCh:
				ticker.Stop()
				return
			*/
			}
		}
	}()
	
	go func() {
		//next := 0
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				//next = n.fixFinger(next)
			/*
			case <-node.shutdownCh:
				ticker.Stop()
				return
			*/
			}
		}
	}()
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

	succ, err := n.rpc.findSuccessor(address, n.ip)
	if err != nil { return err }

	n.ft.get(0).ipaddr = succ
	n.ft.get(0).iphash = getHash(n.hf, succ)

	return nil
}

// predecessor and successor operations (Chord p.5)
func (n *node) findSuccessor(hashed string) (string, error) {
	pred, err := n.findPredecessor(hashed)
	if err != nil { return "", err }

	if pred == n.ip {
		return n.ft.getSuccessor(), nil
	}

	succ, err := n.rpc.getSuccessor(pred)
	if err != nil { return "", err }
	return succ, nil
}

func (n *node) findPredecessor(hashed string) (string, error) {
	curr := n.ip
	succ := n.ft.getSuccessor()
	bigKey := bigify(hashed)
	var err error

	for {
		if curr == n.ip {
			succ = n.ft.getSuccessor()
		} else {
			succ, err = n.rpc.getSuccessor(curr)
			if err != nil { return "", err }
		}
		if bigBetweenRightInclude(bigify(getHash(n.hf, curr)), bigify(getHash(n.hf, succ)), bigKey) {
			break
		}
		if curr == n.ip {
			curr = n.closestPrecedingFinger(hashed)
		} else {
			curr, err = n.rpc.closestPrecedingFinger(curr, hashed)
			if err != nil { return "", err }
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

// chord extended concurrent logic
func (n *node) stabilize() {
	succ := n.ft.getSuccessor()
	var pred string
	var err error
	if succ == n.ip {
		pred = n.predecessor
	} else {
		pred, err = n.rpc.getPredecessor(succ)
		if err != nil {
			fmt.Println(err)
		}
	}
	if pred != "" && bigBetween(bigify(getHash(n.hf, n.ip)), bigify(getHash(n.hf, succ)), bigify(getHash(n.hf, pred))) {
		succ = pred
		n.ft.lock.Lock()
		defer n.ft.lock.Unlock()
		// we need to set our successor to the appropriate val. 
		n.ft.get(0).ipaddr = succ
		n.ft.get(0).iphash = getHash(n.hf, succ)
	}
	if succ != n.ip {
		n.rpc.notify(succ, n.ip)
	}
}

// not being called yet.
// taking the O(n) approach to find queries at the moment. 
func (n *node) fixFinger(i int) int{
	find := n.ft.get(i).id
	// TODO: need to fix this part!!
	succ, err := n.findSuccessor(hex.EncodeToString(find.Bytes()))
	if err != nil {
		fmt.Println("error in fixing fingertable")
		return i
	}
	n.ft.lock.Lock()
	defer n.ft.lock.Unlock()
	n.ft.get(i).ipaddr = succ
	n.ft.get(i).iphash = getHash(n.hf, succ)
	n.ft.get(i).valid = true

	return (i + 1) % n.ftLen
}

// gRPC Server (chord_grpc.pb.go) Implementation
func (n *node) GetHashFuncCheckSum(ctx context.Context, r *pb.Empty) (*pb.HashFuncResponse, error) {
	hashValue := getFuncHash(n.hf)
	return &pb.HashFuncResponse{HashVal: hashValue}, nil
}

func (n *node) GetSuccessor(ctx context.Context, r *pb.Empty) (*pb.AddressResponse, error) {
	return &pb.AddressResponse{Address: n.ft.getSuccessor()}, nil
}

func (n *node) ClosestPrecedingFinger(ctx context.Context, r *pb.KeyRequest) (*pb.AddressResponse, error) {
	a := n.closestPrecedingFinger(r.Key)
	return &pb.AddressResponse{Address: a}, nil
}

func (n *node) FindSuccessor(ctx context.Context, r *pb.KeyRequest) (*pb.AddressResponse, error) {
	a, err := n.findSuccessor(getHash(n.hf, r.Key))
	if err != nil { return &pb.AddressResponse{Address: ""}, nil }
	return &pb.AddressResponse{Address: a}, nil
}

func (n *node) GetPredecessor(ctx context.Context, r *pb.Empty) (*pb.AddressResponse, error) {
	return &pb.AddressResponse{Address: n.predecessor}, nil
}

func (n *node) Notify(ctx context.Context, r *pb.AddressRequest) (*pb.Empty, error) {
	t1 := bigify(getHash(n.hf, n.predecessor))
	t2 := bigify(getHash(n.hf, n.ip))
	t3 := bigify(getHash(n.hf, r.Address))

	if n.predecessor == "" || bigBetween(t1, t2, t3) {
		n.predecessor = r.Address
	}
	return &pb.Empty{}, nil
}

// Error Messages
var (
	sameHash = "Checksum of hash functions on the nodes do not match. Please check that the same hash functions are being used."
)
