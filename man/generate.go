// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package main

import (
	"log"
	"time"

	"github.com/spf13/cobra/doc"
	"vmiklos.hu/go/cpm/commands"
)

func main() {
	var ctx commands.Context
	cmd := commands.NewRootCommand(&ctx)
	date := time.Date(2022, 7, 22, 12, 0, 0, 0, time.UTC)
	header := &doc.GenManHeader{
		Title:   "CPM",
		Section: "1",
		Date:    &date,
	}
	err := doc.GenManTree(cmd, header, "man")
	if err != nil {
		log.Fatal(err)
	}
}
