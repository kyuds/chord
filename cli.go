package main

import (
	"fmt"
	"os"
	"flag"
)

func main() {
	jn := flag.NewFlagSet("join", flag.ExitOnError)
	jnAddr := jn.String("address", "", "IP address for Chord node start on.")
	jnRing := jn.String("", "", "existing Chord node's IP to join to.")

	ct := flag.NewFlagSet("create", flag.ExitOnError)
	ctAddr := ct.String("address", "", "IP address for Chord node start on.")

	if (len(os.Args) < 2) {
		fmt.Println("Expected more arguments.")
	}

	var joining bool = false
	var join_ip string = "dummy"
	var node_ip string = "dummy"

	switch os.Args[1] {
	case "join":
		jn.Parse(os.Args[2:])
		joining = true
		join_ip = *jnRing
		node_ip = *jnAddr
	case "create":
		ct.Parse(os.Args[2:])
		node_ip = *ctAddr
	default:
		fmt.Println("expected 'join' or 'create' commands.")
		os.Exit(1)
	}

	if (node_ip == "" || join_ip == "") {
		fmt.Println("Please enter all arguments!")
		os.Exit(1)
	}

	if (joining) {
		fmt.Println(join_ip)
	}

	fmt.Println(node_ip)
}
