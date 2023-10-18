package chord

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
API Settings for the Chord Implementation

All API for this implementation of Chord is present in this file and this file only. This
means that to use Chord without understanding the internals, only this file needs to be
referenced. As noted in the README, this implementation of Chord stays true to the definition
of the protocol given in the original paper (Stoica, et al). Therefore, the only APIs that
exist are "Initialize", "Lookup", and "Exit".

Initialization of a Chord node (and subsequently, a ring) takes in a Config struct that can
be modified based on the requirements of the given project. A DefaultConfig exists so that
developers can opt to use the default configurations when specified configurations are not
necessary. To make a Chord node join an existing ring (which is most likely what is going
to happen), calling SetJoinNode on the config file will allow the node to join an existing
Chord configuration.

When joining, note that the hash function used in the Chord ring and the hash function that
is given in the config struct needs to match. For instance, if the Chord ring uses SHA1 as
the hashing algorithm and the new node intending to join uses SHA256, then the new node will
not be able to join the ring.

Lookups are easy: simply call "Lookup" on any key, and the key will be HASHED and then will
locate the appropriate node that needs to handle the key. Then, the ip address of the node
will be returned to the user.

Exit is defined for a "planned" exit, as in the program intended for the specific Chord node
to terminate. This provides methods for the Chord ring to balance itself and respond to the
termination more quickly.
*/

type Config struct {
	Address       string
	Hash          func() hash.Hash
	Joining       bool
	JoinAddress   string
	MaxRepli      int
	ServerOptions []grpc.ServerOption
	DialOptions   []grpc.DialOption
	Timeout       time.Duration
	MaxIdle       time.Duration
	Stabilization time.Duration
	FingerFix     time.Duration
}

// Default configurations for Chord users.
func DefaultConfigs(address string) *Config {
	c := &Config{
		Address:       address,
		Hash:          sha1.New,
		Joining:       false,
		JoinAddress:   "",
		MaxRepli:      5,
		ServerOptions: nil,
		DialOptions:   make([]grpc.DialOption, 0, 2),
		Timeout:       10 * time.Millisecond,
		MaxIdle:       1000 * time.Millisecond,
		Stabilization: time.Second,
		FingerFix:     500 * time.Second,
	}
	c.DialOptions = append(
		c.DialOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	return c
}

// Set configuration file to tell the Chord node
// to join an already existing Chord ring in which
// address is located
func (c *Config) SetJoinNode(address string) {
	c.Joining = true
	c.JoinAddress = address
}

// Initializes the chord client for the user.
func Initialize(conf *Config) *chordcli {
	err := conf.validate()
	if err != nil {
		panic(err)
	}
	nd, err := initialize(conf)
	if err != nil {
		panic(err)
	}
	c := &chordcli{n: nd}
	return c
}

// Looks up the given key and returns the
// node's address responsible for the key.
func (c *chordcli) Lookup(key string) (string, error) {
	return c.n.findSuccessor(getHash(c.n.hf, key))
}

func (c *chordcli) Stat() {
	fmt.Println("Stats:")
}

type chordcli struct {
	n *node
}

// TODO: validate IP address
func (c *Config) validate() error {
	return nil
}
