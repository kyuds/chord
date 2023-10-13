package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/kyuds/go-chord/chord"
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

	c := chord.Initialize(node_ip)

	if (joining) {
		fmt.Printf("Joining Chord ring through exisiting node IP: %s\n", join_ip)
		c.JoinRing(join_ip)
	}

	fmt.Printf("Current node IP: %s\n", node_ip)
	fmt.Println("\nTo lookup, type 'lookup' followed by a key.")
	fmt.Println("To put, type 'put' followed by a key and value.")
	fmt.Print("To quit, type 'quit' and the program will exit.\n\n")

	cli := bufio.NewReader(os.Stdin)
	loop: for {
		fmt.Print("CHORD: ")
		cmd, _ := cli.ReadString('\n')
		cmd = strings.Trim(cmd, " \n")
		args := strings.Split(cmd, " ")
		
		if (len(args) == 0) {
			fmt.Println("no args!")
			continue
		}

		switch args[0] {
		case "lookup":
			if (len(args) == 2) {
				c.Lookup(args[1])
			}
		case "put":
			if (len(args) == 3) {
				c.Put(args[1], args[2])
			}
		case "quit":
			break loop
		}
	}
	c.Exit()
}
