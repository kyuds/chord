package chord

import (
	"fmt"
)

type chordcli struct {
	IPAddress string
}

// Initializes the chord client for the user.
// Will panic if the specified IP and PORT is
// already in use or has some other network
// error. 
func Initialize(ip string) chordcli {
	fmt.Printf("Initialize! (%s)\n", ip)
	c := chordcli{IPAddress: ip}
	return c
}

// Joins an already existing chord ring by
// attempting to contact an already existing
// chord node at address "ip". Will panic
// if such node is not found or there is
// any other network error. 
func (c *chordcli) JoinRing(ip string) {
	fmt.Printf("Join! (%s)\n", ip)
}

// Puts a key, value pair into the chord
// network. The key is hashed via SHA1. 
func (c *chordcli) Put(key, val string) {
	fmt.Printf("Put! (%s %s)\n", key, val)
}

// Looks up the given key and returns the
// IP address of the node that has the key. 
func (c *chordcli) Lookup(key string) {
	fmt.Printf("Lookup! (%s)\n", key)
}

// Performs a planned exit of the node. 
func (c *chordcli) Exit() {
	fmt.Println("Exit!")
}
