package main

import (
	"log"

	"github.com/spf13/cobra/doc"
	"vmiklos.hu/go/turtle-cpm/commands"
)

func main() {
	var ctx commands.Context
	cmd := commands.NewRootCommand(&ctx)
	header := &doc.GenManHeader{
		Title:   "CPM",
		Section: "1",
	}
	err := doc.GenManTree(cmd, header, "man")
	if err != nil {
		log.Fatal(err)
	}
}
