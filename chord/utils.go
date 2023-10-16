package chord

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"reflect"
	"runtime"
)

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

func getHash(h func() hash.Hash, key string) string {
	checkSum := h().Sum([]byte(key))
	return hex.EncodeToString(checkSum[:])
}
