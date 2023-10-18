package chord

import (
	"fmt"
	"math/big"
	"sync"
)

type fingerEntry struct {
	id *big.Int
	iphash string
	ipaddr string
	valid bool
}

type fingerTable struct {
	tb []*fingerEntry
	lock sync.RWMutex
}

// prebuild hash spaces
func initFingerTable(h string, address string, size int) fingerTable {
	tb := make([]*fingerEntry, size)
	for i := 0; i < size; i++ {
		tb[i] = &fingerEntry {
			id: computeFingerId(h, i, size),
			iphash: h,
			ipaddr: address,
			valid: false,
		}
	}
	tb[0].valid = true
	return fingerTable { tb: tb }
}

// from Chord paper: (hash + 2^i) mod 2^m
func computeFingerId(hashed string, i int, size int) *big.Int {
	r := bigify(hashed)
	r.Add(r, bigPow(2, i))
	return big.NewInt(0).Mod(r, bigPow(2, size))
}

func (f *fingerTable) get(i int) *fingerEntry {
	return f.tb[i]
}

func (f *fingerTable) getSuccessor() string {
	return f.tb[0].ipaddr
}

func (f *fingerTable) printself() {
	for _, i := range f.tb {
		fmt.Print(i.id)
		fmt.Print(" ")
		fmt.Print(i.ipaddr)
		fmt.Print(" ")
		fmt.Println(i.valid)
	}
}

func (f *fingerTable) invalidateAddress(address string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	for _, finger := range f.tb {
		if finger.ipaddr == address {
			finger.valid = false
		}
	}
}
