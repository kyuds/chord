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
func bigToHex(b *big.Int) string {
	return hex.EncodeToString(b.Bytes())
}

// converter hex from string to big.Int (for rpc)
func hexToBig(h string) *big.Int {
	b, _ := hex.DecodeString(h)
	return new(big.Int).SetBytes(b)
}
