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
	checkSum := h().Sum([]byte(key))
	return hex.EncodeToString(checkSum[:])
}

// TODO: validate IP address
func (c *Config) validate() error {
	return nil
}

// big Power calculator (since big library is overcomplicated)
func bigPow(x, y int) *big.Int {
	return big.NewInt(0).Exp(big.NewInt(int64(x)), big.NewInt(int64(y)), big.NewInt(0))
}
