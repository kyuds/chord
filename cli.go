package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	//"crypto/sha256"
	"github.com/kyuds/go-chord/chord"
)

func main() {
	var conf *chord.Config

	if os.Args[1] == "1" {
		conf = chord.DefaultConfigs("localhost:8000")
	} else {
		conf = chord.DefaultConfigs("localhost:8001")
		conf.SetJoinNode("localhost:8000")
	}

	//conf.Hash = sha256.New
	_ = chord.Initialize(conf)
	for {

	}
}

func cli() {
	jn := flag.NewFlagSet("join", flag.ExitOnError)
	jnAddr := jn.String("address", "", "IP address for Chord node start on.")
	jnRing := jn.String("joinOn", "", "existing Chord node's IP to join to.")

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

	fmt.Printf("Current node IP: %s\n", node_ip)
	fmt.Println("\nTo lookup, type 'lookup' followed by a key.")
	fmt.Print("To quit, type 'quit' and the program will exit.\n\n")

	conf := chord.DefaultConfigs(node_ip)
	if joining {
		conf.SetJoinNode(join_ip)
	}
	c := chord.Initialize(conf)

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
				_, _ = c.Lookup(args[1])
			}
		case "quit":
			break loop
		}
	}
	c.Exit()
}
