// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func newUpdateCommand(ctx *Context) *cobra.Command {
	var machine string
	var service string
	var user string
	var passwordType PasswordType
	var password string
	var dryRun bool
	var secure bool
	var id string
	var archived string
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "updates an existing password",
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
			if len(machine) > 0 {
				query, err := transaction.Prepare("update passwords set machine=? where id=?")
				if err != nil {
					return fmt.Errorf("db.Prepare() failed: %s", err)
				}

				result, err := query.Exec(machine, id)
				if err != nil {
					return fmt.Errorf("db.Exec() failed: %s", err)
				}

				affected, err = result.RowsAffected()
				if err != nil {
					return fmt.Errorf("result.RowsAffected() failed: %s", err)
				}
			}
			if len(service) > 0 {
				query, err := transaction.Prepare("update passwords set service=? where id=?")
				if err != nil {
					return fmt.Errorf("db.Prepare() failed: %s", err)
				}

				result, err := query.Exec(service, id)
				if err != nil {
					return fmt.Errorf("db.Exec() failed: %s", err)
				}

				affected, err = result.RowsAffected()
				if err != nil {
					return fmt.Errorf("result.RowsAffected() failed: %s", err)
				}
			}
			if len(user) > 0 {
				query, err := transaction.Prepare("update passwords set user=? where id=?")
				if err != nil {
					return fmt.Errorf("db.Prepare() failed: %s", err)
				}

				result, err := query.Exec(user, id)
				if err != nil {
					return fmt.Errorf("db.Exec() failed: %s", err)
				}

				affected, err = result.RowsAffected()
				if err != nil {
					return fmt.Errorf("result.RowsAffected() failed: %s", err)
				}
			}
			if len(passwordType) > 0 {
				query, err := transaction.Prepare("update passwords set type=? where id=?")
				if err != nil {
					return fmt.Errorf("db.Prepare() failed: %s", err)
				}

				result, err := query.Exec(passwordType, id)
				if err != nil {
					return fmt.Errorf("db.Exec() failed: %s", err)
				}

				affected, err = result.RowsAffected()
				if err != nil {
					return fmt.Errorf("result.RowsAffected() failed: %s", err)
				}
			}
			generatedPassword := false
			if len(password) > 0 {
				if password == "-" {
					password, err = generatePassword(secure)
					if err != nil {
						return fmt.Errorf("generatePassword() failed: %s", err)
					}
					generatedPassword = true
				}
				query, err := transaction.Prepare("update passwords set password=? where id=?")
				if err != nil {
					return fmt.Errorf("db.Prepare() failed: %s", err)
				}

				result, err := query.Exec(password, id)
				if err != nil {
					return fmt.Errorf("db.Exec() failed: %s", err)
				}

				affected, err = result.RowsAffected()
				if err != nil {
					return fmt.Errorf("result.RowsAffected() failed: %s", err)
				}
			}
			if len(archived) > 0 {
				query, err := transaction.Prepare("update passwords set archived=? where id=?")
				if err != nil {
					return fmt.Errorf("db.Prepare() failed: %s", err)
				}

				parsed, err := strconv.ParseBool(archived)
				if err != nil {
					return fmt.Errorf("ParseBool() failed: %s", err)
				}
				result, err := query.Exec(parsed, id)
				if err != nil {
					return fmt.Errorf("db.Exec() failed: %s", err)
				}

				affected, err = result.RowsAffected()
				if err != nil {
					return fmt.Errorf("result.RowsAffected() failed: %s", err)
				}
			}
			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Would update %v password\n", affected)
				ctx.NoWriteBack = true
			} else {
				transaction.Commit()
				fmt.Fprintf(cmd.OutOrStdout(), "Updated %v password\n", affected)
			}
			if generatedPassword {
				fmt.Fprintf(cmd.OutOrStdout(), "Generated password: %s\n", password)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, `do everything except actually perform the database action (default: false)`)
	cmd.Flags().BoolVarP(&secure, "secure", "y", false, `increase number of symbols from 0 to 3 (default: false)`)
	cmd.Flags().StringVarP(&id, "id", "i", "", `unique identifier (default: ask)`)
	cmd.Flags().StringVarP(&machine, "machine", "m", "", "new machine (default: keep unchanged)")
	cmd.Flags().StringVarP(&service, "service", "s", "", "new service (default: keep unchanged)")
	cmd.Flags().StringVarP(&user, "user", "u", "", "new user (default: keep unchanged)")
	cmd.Flags().VarP(&passwordType, "type", "t", `new password type ("plain" or "totp"; default: keep unchanged)`)
	cmd.Flags().StringVarP(&password, "password", "p", "", `new password ("-" generates a new one; default: keep unchanged)`)
	cmd.Flags().StringVarP(&archived, "archived", "a", "", `new archived value ("true" or "false"; default: keep unchanged)`)

	return cmd
}
