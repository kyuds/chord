package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
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

	// Initialize Chord Node Here. 

	if (joining) {
		fmt.Printf("Joining Chord ring through exisiting node IP: %s\n", join_ip)
		// Join Chord Node here. 
	}

	fmt.Printf("Current node IP: %s\n", node_ip)
	fmt.Println("\nTo query, type 'query' followed by a key.")
	fmt.Println("To put, type 'put' followed by a key and value.")
	fmt.Print("To quit, type 'quit' and the program will exit.\n\n")

	cli := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("CHORD: ")
		cmd, _ := cli.ReadString('\n')
		cmd = strings.Trim(cmd, " \n")
		args := strings.Split(cmd, " ")
		
		if (len(args) == 0) {
			fmt.Println("no args!")
			continue
		}

		switch args[0] {
		case "query":
			fmt.Println(args[1])
			// Query logic here.
		case "put":
			fmt.Println(args[1])
			// Put logic here. 
		case "quit":
			break
		}
	}

	// Execute planned failures here. 
}
