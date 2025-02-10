// Copyright 2025 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func gcPasswords(context *Context) error {
	query, err := context.Database.Prepare("vacuum")
	if err != nil {
		return fmt.Errorf("db.Prepare() failed: %s", err)
	}

	_, err = query.Exec()
	if err != nil {
		return fmt.Errorf("query.Exec() failed: %s", err)
	}

	return nil
}

func newGcCommand(ctx *Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gc",
		Short: "rebuilds the database file, repacking it into a minimal amount of disk space",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := gcPasswords(ctx)
			if err != nil {
				return fmt.Errorf("gcPasswords() failed: %s", err)
			}
			return nil
		},
	}

	return cmd
}
