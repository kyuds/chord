package chord

import (
	"context"
	"encoding/hex"
	"fmt"
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

	// successors predecessors
	repli    int
	succ     *successors // basically a queue.
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
	go func() {
		// check predecessor existence --> clear predecessor if so
		// clear out all ft table whenever node is determined to be not active
		// normal stabilization process.
		// update successors accordingly --> always maintain correct queue.
		predchecker := time.NewTicker(4 * conf.Stabilization)
		stabilizer := time.NewTicker(1 * conf.Stabilization)
		for {
			select {
			case <-predchecker.C:
				n.checkPredecessor()
			case <-stabilizer.C:
				n.stabilize()
				/*
					case <-n.shutdownCh:
						ticker.Stop()
						return
				*/
			}
		}

	}()

	go func() {
		next := 0
		ticker := time.NewTicker(conf.FingerFix)
		for {
			select {
			case <-ticker.C:
				next = n.fixFinger(next)
				/*
					case <-node.shutdownCh:
						ticker.Stop()
						return
				*/
			}
		}
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

	succ, err := n.rpc.findSuccessor(address, getHash(n.hf, n.ip))
	if err != nil {
		return err
	}
	n.succ.push(succ)

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

	succ, err := n.rpc.getSuccessor(pred)
	if err != nil {
		// n.ft.invalidateAddress(pred)
		return "", err
	}
	return succ, nil
	// return "", nil
}

func (n *node) findPredecessor(hashed string) (string, error) {
	curr := n.ip
	succ := n.succ.getSuccessor()
	bigKey := bigify(hashed)
	var err error
	for {
		if curr == n.ip {
			succ = n.succ.getSuccessor()
		} else {
			succ, err = n.rpc.getSuccessor(curr)
			if err != nil {
				// n.ft.invalidateAddress(curr)
				return "", err
			}
		}
		if bigBetweenRightInclude(bigify(getHash(n.hf, curr)), bigify(getHash(n.hf, succ)), bigKey) {
			break
		}
		if curr == n.ip {
			curr = n.closestPrecedingFinger(hashed)
		} else {
			curr, err = n.rpc.closestPrecedingFinger(curr, hashed)
			if err != nil {
				// n.ft.invalidateAddress(curr)
				return "", err
			}
		}
	}
	return curr, nil
}

func (n *node) closestPrecedingFinger(hashed string) string {
	c1 := bigify(n.id)
	c2 := bigify(hashed)
	for i := n.ftLen - 1; i >= 0; i-- {
		f := n.ft.get(i)
		if !f.valid && i != 0 {
			continue
		}
		if i == 0 && bigBetween(c1, c2, f.id) {
			return n.succ.getSuccessor()
		}
		if i != 0 && bigBetween(c1, c2, f.id) {
			return f.ipaddr
		}
	}
	return n.ip
}

// finger fixing
func (n *node) fixFinger(i int) int {
	f := n.ft.get(i)
	succ, err := n.findSuccessor(hex.EncodeToString(f.id.Bytes()))
	if err != nil {
		return i
	}
	n.ft.lock.Lock()
	defer n.ft.lock.Unlock()
	f.ipaddr = succ
	f.iphash = getHash(n.hf, succ)
	f.valid = true

	return (i + 1) % n.ftLen
}

// stabilize
func (n *node) stabilize() {
	var succ string
	var pred string
	var err error

	for {
		succ = n.succ.getSuccessor()
		if succ == n.ip {
			pred = n.pred
			break
		} else {
			pred, err = n.rpc.getPredecessor(succ)
			if err != nil {
				n.succ.pop()
			} else {
				break
			}
		}
	}

	t1 := bigify(getHash(n.hf, n.ip))
	t2 := bigify(getHash(n.hf, succ))
	t3 := bigify(getHash(n.hf, pred))

	if pred != "" && bigBetween(t1, t2, t3) {
		succ = pred
	}
	n.succ.clear()
	if succ != n.ip {
		n.succ.push(succ)
		_ = n.rpc.notify(succ, n.ip)
		curr := succ
		for {
			if n.succ.length() >= n.repli {
				break
			}
			curr, err = n.rpc.getSuccessor(curr)
			if err != nil || curr == n.ip {
				break
			}
			n.succ.push(curr)
		}
	}
}

// check predecessor
func (n *node) checkPredecessor() {
	err := n.rpc.ping(n.pred)
	if err != nil {
		//n.ft.invalidateAddress(n.pred)
		n.predLock.Lock()
		n.pred = ""
		n.predLock.Unlock()
	}
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

func (n *node) Ping(ctx context.Context, r *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// successors struct
type successors struct {
	ip   string
	data []string
	lock sync.RWMutex
}

func createSuccessorQueue(ownIP string) *successors {
	return &successors{
		ip:   ownIP,
		data: make([]string, 0),
	}
}

func (q *successors) push(address string) {
	q.data = append(q.data, address)
}

func (q *successors) pop() string {
	d := q.data[0]
	q.data = q.data[1:]
	return d
}

func (q *successors) get(i int) string {
	return q.data[i]
}

func (q *successors) length() int {
	return len(q.data)
}

func (q *successors) empty() bool {
	return len(q.data) == 0
}

func (q *successors) getSuccessor() string {
	if q.empty() {
		return q.ip
	}
	return q.get(0)
}

func (q *successors) clear() {
	q.data = make([]string, 0)
}
