package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"chord/chord"
)

// ./example create --address="localhost:8000"
// ./example join --address="localhost:8001" --join="localhost:8000"
// ./example join --address="localhost:8002" --join="localhost:8000"

func Cli() {
	jn := flag.NewFlagSet("join", flag.ExitOnError)
	jnAddr := jn.String("address", "", "IP address for Chord node start on.")
	jnRing := jn.String("join", "", "existing Chord node's IP to join to.")

	ct := flag.NewFlagSet("create", flag.ExitOnError)
	ctAddr := ct.String("address", "", "IP address for Chord node start on.")

	if len(os.Args) < 2 {
		fmt.Println("Expected more arguments.")
	}

	var join_ip string = "dummy"
	var node_ip string = "dummy"

	switch os.Args[1] {
	case "join":
		jn.Parse(os.Args[2:])
		join_ip = *jnRing
		node_ip = *jnAddr
	case "create":
		ct.Parse(os.Args[2:])
		node_ip = *ctAddr
		join_ip = ""
	default:
		fmt.Println("expected 'join' or 'create' commands.")
		os.Exit(1)
	}

	if node_ip == "" {
		fmt.Println("Please enter all arguments!")
		os.Exit(1)
	}

	fmt.Printf("Current node IP: %s\n", node_ip)
	fmt.Println("\nTo lookup, type 'lookup' followed by a key.")
	fmt.Print("To quit, type 'quit' and the program will exit.\n\n")

	conf := chord.DefaultConfig(node_ip, join_ip, nil)

	c, _ := chord.NewChord(conf)
	_ = c.Start()

	cli := bufio.NewReader(os.Stdin)

loop:
	for {
		fmt.Print("CHORD: ")
		cmd, _ := cli.ReadString('\n')
		cmd = strings.Trim(cmd, " \n")
		args := strings.Split(cmd, " ")

		if len(args) == 0 {
			fmt.Println("no args!")
			continue
		}

		switch args[0] {
		case "lookup":
			if len(args) == 2 {
				ret, err := c.LookUp(args[1])
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(ret)
			}
		case "quit":
			c.Exit()
			break loop
		}
	}

	for {

	}
}
