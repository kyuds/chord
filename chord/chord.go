package chord

import (
	"fmt"
)

type ChordCLI struct {
	IPAddress string
}

func Initialize(ip string) ChordCLI {
	fmt.Printf("Initialize! (%s)\n", ip)
	c := ChordCLI{IPAddress: ip}
	return c
}

func (c *ChordCLI) JoinRing(ip string) {
	fmt.Printf("Join! (%s)\n", ip)
}

func (c *ChordCLI) Put(key, val string) {
	fmt.Printf("Put! (%s %s)\n", key, val)
}

func (c *ChordCLI) Lookup(key string) {
	fmt.Printf("Lookup! (%s)\n", key)
}

func (c *ChordCLI) Exit() {
	fmt.Println("Exit!")
}
