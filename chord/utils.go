package chord

import (
	"encoding/hex"
	"hash"
	"math/big"
)

// hashing a string to big.Int representation
func getHash(h func() hash.Hash, key string) *big.Int {
	hasher := h()
	hasher.Write([]byte(key))
	return new(big.Int).SetBytes(hasher.Sum(nil))
}

// converter from big.Int to hex string (for rpc)
func bigToString(b *big.Int) string {
	return hex.EncodeToString(b.Bytes())
}

// converter hex from string to big.Int (for rpc)
func stringToBig(h string) *big.Int {
	b, _ := hex.DecodeString(h)
	return new(big.Int).SetBytes(b)
}

// b1 < key < b2
func bigInRange(b1, b2, key *big.Int) bool {
	switch b1.Cmp(b2) {
	case -1:
		return b1.Cmp(key) == -1 && b2.Cmp(key) == 1
	case 0:
		return b1.Cmp(key) != 0
	case 1:
		return b1.Cmp(key) == -1 || b2.Cmp(key) == 1
	}
	return true
}

// b1 < key <= b2
func bigInRangeRightInclude(b1, b2, key *big.Int) bool {
	return bigInRange(b1, b2, key) || b2.Cmp(key) == 0
}

func bigPowTwo(i int) *big.Int {
	return big.NewInt(0).Exp(big.NewInt(int64(2)), big.NewInt(int64(i)), big.NewInt(0))
}
