package chord

import (
	"crypto/sha1"
	"hash"
	"time"
	"google.golang.org/grpc"
)

type Config struct {
	Address string
	Hash func() hash.Hash
	Joining bool
	JoinIP string
	ServerOptions []grpc.ServerOption
	DialOptions []grpc.DialOption
	Timeout time.Duration
	MaxIdle time.Duration
}

func DefaultConfigs(address string) *Config {
	c := &Config {
		Address: address,
		Hash: sha1.New,
		Joining: false,
		JoinIP: "",
		// TODO: figure out how to deal with these later. 
		ServerOptions: nil,
		DialOptions: make([]grpc.DialOption, 0, 2),
		Timeout: 10 * time.Millisecond,
		MaxIdle: 1000 * time.Millisecond,
	}
	c.DialOptions = append(c.DialOptions, grpc.WithInsecure(), grpc.WithBlock())
	return c
}

func (c *Config) SetJoinNode(address string) {
	c.Joining = true
	c.JoinIP = address
}

// Initializes the chord client for the user.
// Will join an already existing chord
// ring if joining is true on the specified
// join_ip string. 
func Initialize(conf *Config) *chordcli {
	e := conf.validate()
	if e != nil { panic(e) }

	n, err := newNode(conf)
	if err != nil { panic(err) }

	c := &chordcli{chrd: n}
	return c
}

// Looks up the given key and returns the
// value corresponding to the key. 
func (c *chordcli) Lookup(key string) (string, error) {
	return "no ip address", nil
}

// Performs a planned exit of the node. 
func (c *chordcli) Exit() error {
	return nil
}

type chordcli struct {
	chrd *node
}
