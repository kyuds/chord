package chord

import (
	"math/big"
	"sync"
)

// Successor Table maintains the correct pointers to a Chord Node's
// successor. The maintainence is done SOLELY in background operations,
// and while the node is rebalancing due to cluster membership changes,
// lookup operations can be expected to yield incorrect results. It will
// heal later on.
type successorTable struct {
	sentinel successorNode
	lock     sync.RWMutex
	length   int
}

// A linked list style is appropriate because we do not expect the number
// of successors kept track (replication factor) to be extremely high.
type successorNode struct {
	address string
	next    *successorNode
}

func newSuccessorTable(address string, length int) successorTable {
	return successorTable{
		sentinel: successorNode{
			address: address,
			next:    nil,
		},
		length: length,
	}
}

func (s *successorTable) getSuccessor() string {
	s.lock.RLock()
	defer s.lock.Unlock()
	if s.sentinel.next != nil {
		return s.sentinel.address
	}
	return s.sentinel.next.address
}

func (s *successorTable) clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sentinel.next = nil
}

// As proposed by the original Chord paper, the finger table
// stores information to speed up the search process.
type fingerTable struct {
	tb   []entry
	size int
	lock sync.RWMutex
}

type entry struct {
	id     *big.Int
	ipaddr string
	valid  bool
}

func newFingerTable(address string, size int) fingerTable {
	tb := make([]entry, size)
	for i := 0; i < size; i++ {
		tb[i] = entry{
			id:     big.NewInt(0),
			ipaddr: address,
			valid:  false,
		}
	}
	return fingerTable{
		tb:   tb,
		size: size,
	}
}
