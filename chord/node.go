package chord

import (
	"chord/pb"
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
	stateLock sync.RWMutex
	// predecessor  string
	// numSuccessor int

	// chord optimization
	// ft fingerTable

	// gRPC server & transport layer
	transport *transport
	pb.UnimplementedChordServer
}

// constructor for new Chord Node instance. Initializes the
// TCP listener and registers gRPC server.
func NewChord(conf *Config) (*ChordNode, error) {
	c := &ChordNode{
		conf:   conf,
		ipHash: getHash(conf.Hash, conf.Address),
		// initialize predecessor, fingertable, successorqueue?
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
				c.stabilize()
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
	if address != "" {
		hashName, err := c.transport.getHashFuncName(address)
		if err != nil {
			return err
		}
		if c.conf.HashName() != hashName {
			return fmt.Errorf("hash function needs to be homogenous in chord ring")
		}
		h := getHash(c.conf.Hash, c.conf.Address)
		succ, err := c.transport.findSuccessor(address, bigToString(h))
		if err != nil {
			return err
		}
		// TODO: initialize succ value
		fmt.Println(succ)
	}
	return nil
}

// Helper function for stabilization
func (c *ChordNode) stabilize() {
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
}

// Helper function for rectifying
func (c *ChordNode) rectify(candidate string) {
	// should this be made TryLock?
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	// follow algorithm on Zave, et al
}

// Helper function for fixing fingers
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
func (c *ChordNode) LookUp(key string) (string, string, error) {
	c.stateLock.RLock()
	defer c.stateLock.Unlock()
	return "", "", nil
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

// gRPC Server Implementation
// TODO
