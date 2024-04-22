package main

import "chord/cli"

/*
MockCli
./example localhost:8000 localhost:8001

Cli
./example create --address=localhost:8000
./example join --address="localhost:8001" --join="localhost:8000"
*/

func main() {
	// cli.MockCli()
	cli.Cli()
}
