// Copyright 2025 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type passwordRow struct {
	ID           int
	Machine      string
	Service      string
	User         string
	Password     string
	PasswordType PasswordType
	Archived     bool
	Created      string
	Modified     string
}

func exportPasswords(db *sql.DB, args []string) ([]byte, error) {
	var results []passwordRow
	rows, err := db.Query("select id, machine, service, user, password, type, archived, created, modified from passwords")
	if err != nil {
		return nil, fmt.Errorf("db.Query(select) failed: %s", err)
	}

	defer rows.Close()
	for rows.Next() {
		var row passwordRow
		err = rows.Scan(&row.ID, &row.Machine, &row.Service, &row.User, &row.Password, &row.PasswordType, &row.Archived, &row.Created, &row.Modified)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan() failed: %s", err)
		}

		if len(args) > 0 && !matchesQuery(args[0], row.ID, row.Machine, row.Service, row.User, row.PasswordType) {
			continue
		}

		results = append(results, row)
	}

	j, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal() failed: %s", err)
	}

	return j, nil
}

func newExportCommand(ctx *Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "export",
		Short: "exports passwords as JSON",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			j, err := exportPasswords(ctx.Database, args)
			if err != nil {
				return fmt.Errorf("readPasswords() failed: %s", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", j)

			ctx.NoWriteBack = true
			return nil
		},
	}

	return cmd
}
