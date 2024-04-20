package cli

import (
	"bufio"
	"chord/chord"
	"os"
	"strings"
)

/*
Testing code to see if the transport layer is functioning
correctly.
Verified:
- gRPC returns correct results
- if it reports a not present, retries multiple times
- retry limited to 3 times
- caches connections
- garbage collects connections
- updates times so that garbage collection only happens
  for idle connections.
*/

func MockTransportTestCli() {
	conf := chord.DefaultConfig(os.Args[1], "", nil)
	dn := chord.NewDummy(conf)
	dn.Start()
	cli := bufio.NewReader(os.Stdin)
	for {
		cmd, _ := cli.ReadString('\n')
		cmd = strings.Trim(cmd, " \n")
		args := strings.Split(cmd, " ")

		switch args[0] {
		case "switch":
			dn.Switch()
		case "call":
			dn.CallRPC(os.Args[2])
		case "callself":
			dn.CallRPC(os.Args[1])
		}
	}
}
