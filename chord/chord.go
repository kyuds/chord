package chord

import (
	"google.golang.org/grpc"
	"crypto/sha1"
	"hash"
	"time"
)

// API for Chord Configurations
type Config struct {
	// node configs
	Address string
	Hash func() hash.Hash
	// Replication int (TODO)

	// gRPC configs
	ServerOptions []grpc.ServerOption
	DialOptions []grpc.DialOption
	Timeout time.Duration
	MaxIdle time.Duration
}

func DefaultConfigs(address string) *Config {
	return &Config {
		Address: address,
		Hash: sha1.New,
		// TODO: figure out how to deal with these later. 
		ServerOptions: nil,
		DialOptions: nil,
		Timeout: 10 * time.Millisecond,
		MaxIdle: 1000 * time.Millisecond,
	}
}

// TODO
// validate IP address, Hash func?
func (c *Config) validate() error {
	return nil
}

// API for Chord Process
type chordcli struct {
	chrd *node
}

// Initializes the chord client for the user.
// Will join an already existing chord
// ring if joining is true on the specified
// join_ip string. 
func Initialize(conf *Config, joining bool, join_ip string) *chordcli {
	e := conf.validate()
	if e != nil { panic(e) }

	n, err := newNode(conf)
	if err != nil { panic(err) }

	if joining {
		err = n.joinNode(join_ip)
		if err != nil { panic(err) }
	}

	c := &chordcli{chrd: n}
	return c
}

// Puts a key, value pair into the chord
// network. The key is hashed via SHA1. 
func (c *chordcli) Put(key, val string) error {
	return nil
}

// Looks up the given key and returns the
// value corresponding to the key. 
func (c *chordcli) Lookup(key string) (string, error) {
	return "empty val", nil
}

// Deletes the value corresponding to the
// given key. 
func (c *chordcli) Delete(key string) error {
	return nil
}

// Performs a planned exit of the node. 
func (c *chordcli) Exit() error {
	return nil
}
