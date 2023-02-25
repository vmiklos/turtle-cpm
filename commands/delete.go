// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newDeleteCommand(ctx *Context) *cobra.Command {
	var dryRun bool
	var id string
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes an existing password",
		RunE: func(cmd *cobra.Command, args []string) error {
			transaction, err := ctx.Database.Begin()
			if err != nil {
				return fmt.Errorf("db.Begin() failed: %s", err)
			}

			defer transaction.Rollback()

			var affected int64
			if len(id) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Id: ")
				reader := bufio.NewReader(cmd.InOrStdin())
				line, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("ReadString() failed: %s", err)
				}
				id = strings.TrimSuffix(line, "\n")
			}
			query, err := transaction.Prepare("delete from passwords where id=?")
			if err != nil {
				return fmt.Errorf("db.Prepare() failed: %s", err)
			}

			result, err := query.Exec(id)
			if err != nil {
				return fmt.Errorf("db.Exec() failed: %s", err)
			}

			affected, err = result.RowsAffected()
			if err != nil {
				return fmt.Errorf("result.RowsAffected() failed: %s", err)
			}
			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would delete %v password\n", affected)
				ctx.NoWriteBack = true
			} else {
				transaction.Commit()
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted %v password\n", affected)
			}

			return nil
		},
	}
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, `do everything except actually perform the database action (default: false)`)
	cmd.Flags().StringVarP(&id, "id", "i", "", `unique identifier (default: '')`)

	return cmd
}
