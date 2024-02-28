package main

import (
	"chord/chord"
	"fmt"
)

func main() {
	c := chord.DefaultConfig("localhost:8000", "", nil)
	fmt.Println(c.HashName())
}
