package main

import (
	"os"

	"github.com/vmiklos/turtle-cpm/commands"
)

func main() {
	// notest
	os.Exit(commands.Main(os.Stdout))
}
