package chord

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"math/big"
	"reflect"
	"runtime"
)

// Utility functions for Chord.

// Used to hash the name of the hash function
// used in Chord to ensure that all nodes share
// the same hash function.
// ex: hash "crypto/sha1.New":
//     getFuncHash(sha1.New) -> a6f0f8d9a226b2c8f385aed3583d14c3c0743629
func getFuncHash(i interface{}) string {
	funcName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	checkSum := sha1.Sum([]byte(funcName))
	return hex.EncodeToString(checkSum[:])
}

// use hash function to convert key to hashed string. 
func getHash(h func() hash.Hash, key string) string {
	hasher := h()
	hasher.Write([]byte(key))
	checkSum := hasher.Sum(nil)
	return hex.EncodeToString(checkSum[:])
}

// TODO: validate IP address
func (c *Config) validate() error {
	return nil
}

// math/big std library abstractions
func bigify(hashed string) *big.Int {
	return new(big.Int).SetBytes([]byte(hashed))
}

func bigPow(x, y int) *big.Int {
	return big.NewInt(0).Exp(big.NewInt(int64(x)), big.NewInt(int64(y)), big.NewInt(0))
}

func bigBetween(x, y, key *big.Int) bool {
	switch x.Cmp(y) {
	case -1:
		return x.Cmp(key) == -1 && y.Cmp(key) == 1
	case 0:
		return x.Cmp(key) != 0
	case 1:
		return x.Cmp(key) == -1 || y.Cmp(key) == 1
	}
	return true
}

func bigBetweenRightInclude(x, y, key *big.Int) bool {
	return bigBetween(x, y, key) || y.Cmp(key) == 0
}