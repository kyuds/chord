package chord

import (
	"fmt"
	"hash"
	"sync"

	"github.com/kyuds/go-chord/pb"
)

type node struct {
	// chord node settings
	id       string
	ip       string
	hf       func() hash.Hash
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
	n.succ = createSuccessorQueue()

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

	}()

	// fix fingers
	go func() {

	}()

	// check predecessor
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
		if n.succ.empty() {
			return n.ip, nil
		}
		return n.succ.get(0), nil
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
	succ := n.ft.getSuccessor()
	bigKey := bigify(hashed)
	//var err error

	for {
		if curr == n.ip {
			succ = n.ft.getSuccessor()
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
