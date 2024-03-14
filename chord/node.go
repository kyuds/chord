package chord

import (
	"chord/pb"
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"
)

// ChordNode Node Abstraction. Because each
// hash function produces a hash of different
// size, we use a big Integer to store the
// integer versions of the node ip hashes.
// The successor list is integrated into the
// finger table.
type ChordNode struct {
	// chord node settings
	conf      *Config
	ipHash    *big.Int
	terminate atomic.Bool

	// chord variable state
	stateLock     sync.RWMutex
	predecessor   string
	successorList []string

	// chord fingertable
	// ft fingerTable

	// gRPC server & transport layer
	transport *transport
	pb.UnimplementedChordServer
}

// constructor for new Chord Node instance. Initializes the
// TCP listener and registers gRPC server.
func NewChord(conf *Config) (*ChordNode, error) {
	c := &ChordNode{
		conf:          conf,
		ipHash:        getHash(conf.Hash, conf.Address),
		predecessor:   "",
		successorList: nil,
	}
	c.successorList = make([]string, conf.NumSuccessor)
	for i := 0; i < conf.NumSuccessor; i++ {
		c.successorList[i] = conf.Address
	}
	c.terminate.Store(false)
	transport, err := newTransport(conf)
	if err != nil {
		return nil, err
	}
	c.transport = transport
	pb.RegisterChordServer(transport.server, c)
	return c, nil
}

// Starts the Chord Node. Does the actual connection process
// of the Chord Node to (if exists) a Chord ring. Starts
// background goroutines that ensures the liveliness of the
// system.
func (c *ChordNode) Start() error {
	c.transport.startServer()
	if err := c.join(); err != nil {
		c.transport.stopServer()
		return err
	}
	go func(c *ChordNode) {
		// timers
		stabilizer := time.NewTicker(c.conf.Stabilization)
		fingerfix := time.NewTicker(c.conf.FingerFix)
		liveliness := time.NewTicker(c.conf.CheckAlive)
		cleanup := time.NewTicker(c.conf.CleanConnection)

		nextFingerIndex := 0
		for {
			if c.terminate.Load() {
				break
			}
			select {
			case <-stabilizer.C:
				// c.stabilize()
			case <-fingerfix.C:
				nextFingerIndex = c.fixFingerTable(nextFingerIndex)
			case <-liveliness.C:
				if c.transport.checkServerDead() {
					c.terminate.Store(true)
				}
			case <-cleanup.C:
				c.transport.cleanIdleConnections()
			default:
				continue
			}
			time.Sleep(10 * time.Millisecond)
		}
	}(c)

	return nil
}

// Helper function to join existing chord ring provided that
// the address is valid.
func (c *ChordNode) join() error {
	address := c.conf.JoinAddr
	fmt.Println(address)
	if address == "" {
		return nil
	}
	hashFuncName, succLength, err := c.transport.getChordConfigs(address)
	if err != nil {
		return err
	} else if succLength != c.conf.NumSuccessor || hashFuncName != c.conf.HashName() {
		return fmt.Errorf("chord configurations (hash function, successor list length) doesn't match")
	}
	predecessor, err := c.transport.findPredecessor(address, bigToString(c.ipHash))
	if err != nil {
		return err
	}
	succList, err := c.transport.getSuccessorList(predecessor)
	if err != nil {
		return err
	}
	c.predecessor = predecessor
	c.successorList = succList
	return nil
}

// Helper function for stabilization
func (c *ChordNode) stabilize() {
	c.stateLock.Lock()
	succ := c.successorList[0]
	succList, err1 := c.transport.getSuccessorList(succ)
	succPred, err2 := c.transport.getPredecessor(succ)
	if err1 != nil || err2 != nil {
		if len(c.successorList) == 1 {
			panic("successorlist ran out")
		}
		c.successorList = c.successorList[:len(c.successorList)-1]
		c.stateLock.Unlock()
		c.stabilize()
		return
	}
	c.successorList = append(c.successorList[:1], succList[:len(succList)-1]...)
	if bigInRange(c.ipHash, getHash(c.conf.Hash, c.successorList[0]), getHash(c.conf.Hash, succPred)) {
		succList2, err := c.transport.getSuccessorList(succPred)
		if err == nil {
			c.successorList = append([]string{succPred}, succList2[:len(succList2)-1]...)
		}
	}
	c.stateLock.Unlock()
}

// Helper function for rectifying
func (c *ChordNode) rectify(predecessor string) {
	// TODO: should this be made TryLock?
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	b1 := getHash(c.conf.Hash, c.predecessor)
	key := getHash(c.conf.Hash, predecessor)
	if bigInRange(b1, c.ipHash, key) || !c.transport.ping(c.predecessor) {
		c.predecessor = predecessor
	}
}

// TODO: Helper function for fixing fingers
func (c *ChordNode) fixFingerTable(idx int) int {
	return idx
}

// Stops all background processes and the gRPC server. Effectively
// removes the Chord server from the Chord ring.
// TODO: implement graceful exit (key notification)
func (c *ChordNode) Stop() {
	c.terminate.Store(true)
	c.transport.stopServer()
}

// O(log N) Lookup operation for a key. Returns the IP Address and Node ID
// of the chord node responsible for said key. Operation may error in the
// following scenarios: there was an error with gRPC (this can be tried
// again as the ring gets rebalanced), or there was an issue with the underlying
// goroutines or gRPC server, which is a fatal error.
func (c *ChordNode) LookUp(key string) (string, error) {
	pred, err := c.findPredecessor(getHash(c.conf.Hash, key))
	if err != nil {
		return "", err
	}
	// handle case when we are immediate predecessor
	if c.conf.Address == pred {
		return c.successorList[0], nil
	}
	return c.transport.getSuccessor(pred)
}

func (c *ChordNode) findPredecessor(key *big.Int) (string, error) {
	var curr, succ string
	var err error
	curr = c.conf.Address

	for {
		if curr == c.conf.Address {
			succ = c.successorList[0]
		} else {
			succ, err = c.transport.getSuccessor(curr)
			if err != nil {
				return "", err
			}
		}
		b1, b2 := getHash(c.conf.Hash, curr), getHash(c.conf.Hash, succ)
		if bigInRangeRightInclude(b1, b2, key) {
			break
		}
		if curr == c.conf.Address {
			curr = c.closestPrecedingFinger(key)
		} else {
			curr, err = c.transport.closestPrecedingFinger(curr, bigToString(key))
			if err != nil {
				return "", err
			}
		}
	}
	return "", nil
}

// TODO
// This is tentative for only immediate successor pointer
func (c *ChordNode) closestPrecedingFinger(key *big.Int) string {
	f := big.NewInt(1)
	f.Add(f, key)
	f.Mod(f, big.NewInt(2<<c.conf.Hash().Size()))
	if bigInRange(c.ipHash, key, f) {
		return c.successorList[0]
	}
	return c.conf.Address
}

// Returns status of Chord Node. gRPC errors might invalidate the entire
// node instance, in which all node activity terminates. In this case,
// it is possible to create a new chord instance and start it.
func (c *ChordNode) IsDead() bool {
	return c.terminate.Load()
}

// Return configuration struct of ChordNode.
func (c *ChordNode) GetConfig() *Config {
	return c.conf
}

func (c *ChordNode) Exit() {
	c.transport.stopServer()
	c.terminate.Store(true)
}

// gRPC Server Implementation

func (c *ChordNode) GetChordConfigs(ctx context.Context, r *pb.Empty) (*pb.ConfigResponse, error) {
	return &pb.ConfigResponse{
		HashFuncName:        c.conf.HashName(),
		SuccessorListLength: int32(c.conf.NumSuccessor),
	}, nil
}

func (c *ChordNode) GetSuccessor(ctx context.Context, r *pb.Empty) (*pb.AddressResponse, error) {
	if c.stateLock.TryRLock() {
		defer c.stateLock.Unlock()
		return &pb.AddressResponse{
			Present: true,
			Address: c.successorList[0],
		}, nil
	}
	return &pb.AddressResponse{
		Present: false,
		Address: "",
	}, nil
}

func (c *ChordNode) ClosestPrecedingFinger(ctx context.Context, r *pb.HashKeyRequest) (*pb.AddressResponse, error) {
	if c.stateLock.TryRLock() {
		defer c.stateLock.Unlock()
		return &pb.AddressResponse{
			Present: true,
			Address: c.closestPrecedingFinger(stringToBig(r.HashValue)),
		}, nil
	}
	return &pb.AddressResponse{
		Present: false,
		Address: "",
	}, nil
}

// for joining ONLY
func (c *ChordNode) FindPredecessor(ctx context.Context, r *pb.HashKeyRequest) (*pb.AddressResponse, error) {
	if c.stateLock.TryRLock() {
		defer c.stateLock.Unlock()
		found, err := c.findPredecessor(stringToBig(r.HashValue))
		if err != nil {
			return nil, err
		}
		return &pb.AddressResponse{
			Present: true,
			Address: found,
		}, nil
	}
	return &pb.AddressResponse{
		Present: false,
		Address: "",
	}, nil
}

func (c *ChordNode) GetPredecessor(ctx context.Context, r *pb.Empty) (*pb.AddressResponse, error) {
	if c.stateLock.TryRLock() {
		defer c.stateLock.Unlock()
		return &pb.AddressResponse{
			Present: true,
			Address: c.predecessor,
		}, nil
	}
	return &pb.AddressResponse{
		Present: false,
		Address: "",
	}, nil
}

func (c *ChordNode) GetSuccessorList(ctx context.Context, r *pb.Empty) (*pb.AddressListResponse, error) {
	if c.stateLock.TryRLock() {
		defer c.stateLock.Unlock()
		return &pb.AddressListResponse{
			Present:   true,
			Addresses: c.successorList,
			Length:    int32(len(c.successorList)),
		}, nil
	}
	return &pb.AddressListResponse{
		Present:   false,
		Addresses: nil,
		Length:    int32(len(c.successorList)),
	}, nil
}

func (c *ChordNode) Notify(ctx context.Context, r *pb.AddressRequest) (*pb.Empty, error) {
	c.rectify(r.Address)
	return &pb.Empty{}, nil
}

func (c *ChordNode) Ping(ctx context.Context, r *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
