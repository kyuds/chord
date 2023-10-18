package main

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/kyuds/go-chord/chord"
)

func main() {
	var conf *chord.Config
	if os.Args[1] == "1" {
		conf = chord.DefaultConfigs("localhost:8000")
	} else {
		conf = chord.DefaultConfigs("localhost:8001")
		conf.SetJoinNode("localhost:8000")
		conf.Hash = sha256.New
	}
	c := chord.Initialize(conf)
	fmt.Println(c.Lookup("hai"))
	for {

	}
}
