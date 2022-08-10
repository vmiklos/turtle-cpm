package main

import (
	"os"

	"vmiklos.hu/go/cpm/commands"
)

func main() {
	// notest
	os.Exit(commands.Main(os.Stdin, os.Stdout))
}
