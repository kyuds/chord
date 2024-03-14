package cli

import (
	"bufio"
	"chord/chord"
	"os"
	"strings"
)

func MockCli() {
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
		}
	}
}
