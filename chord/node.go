package chord

import (
	"chord/pb"
)

// Chord Node Abstraction. Because each
// hash function produces a hash of different
// size, we use a big Integer to store the
// integer versions of the node ip hashes.
// The successor list is integrated into the
// finger table.
type node struct {
	// // chord node settings
	// nodeId string
	// nodeIp string
	// ipHash *big.Int
	// conf   *Config

	// // successor, predecessor
	// predecessor  string
	// predLock     sync.RWMutex
	// numSuccessor int
	// // ft fingerTable
	// // ftLength int

	// gRPC server & transport layer
	transport *transport
	pb.UnimplementedChordServer
}

// func initialize(conf *Config) (*node, error) {
// 	return nil, nil
// }

// func (n *node) start() error {
// 	return nil
// }

// func (n *node) stop() {

// }

// gRPC Server Implementation
// TODO
