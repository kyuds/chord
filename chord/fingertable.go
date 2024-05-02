package chord

import (
	"math/big"
	"sync"
)

// fingertable entry abstraction
type fentry struct {
	id    *big.Int
	ip    string
	valid bool
}

func (f *fentry) update(address string) {
	f.ip = address
	f.valid = true
}

func (f *fentry) invalidate() {
	f.valid = false
}

// fingertable struct
type ftable struct {
	tb    []fentry
	succ  *big.Int
	size  int
	_lock sync.RWMutex
}

// prebuild hash spaces
// note the size of the fingertable is size - 1 because the successor pointer
// is maintained by the successor list, and thus there is no need to store
// unnecessary data.
func initFingerTable(h *big.Int, address string, hashSize int) ftable {
	size := hashSize * 8 // bytes to bits
	mod := bigPowTwo(size)
	tb := make([]fentry, size-1)
	for i := 1; i < size; i++ {
		tb[i-1] = fentry{
			id:    computeFinger(h, i, mod),
			ip:    address,
			valid: false,
		}
	}
	return ftable{tb: tb, size: size - 1, succ: computeFinger(h, 0, mod)}
}

// from Chord paper: (hash + 2^i) mod 2^m
// mod is same across all fingers, so is precomputed.
func computeFinger(hashed *big.Int, i int, mod *big.Int) *big.Int {
	tmp := big.NewInt(0).Add(hashed, bigPowTwo(i))
	return big.NewInt(0).Mod(tmp, mod)
}

func (f *ftable) get(i int) *fentry {
	return &f.tb[i]
}

func (f *ftable) invalidate(i int) {
	f.tb[i].invalidate()
}

func (f *ftable) lock() {
	f._lock.Lock()
}

func (f *ftable) unlock() {
	f._lock.Lock()
}

func (f *ftable) rlock() {
	f._lock.RLock()
}

func (f *ftable) runlock() {
	f._lock.RUnlock()
}
