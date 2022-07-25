package main

import (
	"os"

	"vmiklos.hu/go/turtle-cpm/commands"
)

func main() {
	// notest
	os.Exit(commands.Main(os.Stdout))
}
