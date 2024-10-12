// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func generatePassword(secure bool) (string, error) {
	numSymbols := 0
	if secure {
		numSymbols = 3
	}
	// Length of 15 and no symbols matches current Firefox and `pwgen --secure 15 1`.
	output, err := GeneratePassword( /*length*/ 15 /*numDigits=*/, 3, numSymbols /*noUpper=*/, false /*allowRepeat=*/, false)
	if err != nil {
		return "", fmt.Errorf("password.Generate() failed: %s", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func createPassword(context *Context, machine, service, user, password string, passwordType PasswordType, secure bool) (string, error) {
	if len(password) == 0 {
		var err error
		password, err = generatePassword(secure)
		if err != nil {
			return "", fmt.Errorf("generatePassword() failed: %s", err)
		}
	}

	transaction, err := context.Database.Begin()
	if err != nil {
		return "", fmt.Errorf("db.Begin() failed: %s", err)
	}

	defer transaction.Rollback()
	query, err := transaction.Prepare("insert into passwords (machine, service, user, password, type, created, modified) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("db.Prepare() failed: %s", err)
	}

	now := Now().Format(time.RFC3339)
	result, err := query.Exec(machine, service, user, password, passwordType, now, now)
	if err != nil {
		return "", fmt.Errorf("query.Exec() failed: %s", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("result.RowsAffected() failed: %s", err)
	}

	if context.DryRun {
		if context.OutOrStdout != nil {
			fmt.Fprintf(*context.OutOrStdout, "Would create %v password\n", affected)
		}
		context.NoWriteBack = true
	} else {
		if context.OutOrStdout != nil {
			fmt.Fprintf(*context.OutOrStdout, "Created %v password\n", affected)
		}
		transaction.Commit()
	}
	return password, nil
}

func newCreateCommand(ctx *Context) *cobra.Command {
	var machine string
	var service string
	var user string
	var password string
	var passwordType PasswordType = "plain"
	var dryRun bool
	var secure bool
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "creates a new password",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(cmd.InOrStdin())
			if len(machine) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Machine: ")
				line, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("ReadString() failed: %s", err)
				}
				machine = strings.TrimSuffix(line, "\n")
			}
			if len(user) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "User: ")
				line, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("ReadString() failed: %s", err)
				}
				user = strings.TrimSuffix(line, "\n")
			}

			ctx.DryRun = dryRun
			writer := cmd.OutOrStdout()
			ctx.OutOrStdout = &writer
			defer func() { ctx.OutOrStdout = nil }()
			generatedPassword, err := createPassword(ctx, machine, service, user, password, passwordType, secure)
			if err != nil {
				return fmt.Errorf("createPassword() failed: %s", err)
			}

			if generatedPassword != password {
				fmt.Fprintf(cmd.OutOrStdout(), "Generated password: %s\n", generatedPassword)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&machine, "machine", "m", "", "machine (default: ask)")
	cmd.Flags().StringVarP(&service, "service", "s", "http", "service")
	cmd.Flags().StringVarP(&user, "user", "u", "", "user (default: ask)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "password (default: generate)")
	cmd.Flags().VarP(&passwordType, "type", "t", `password type ("plain" or "totp")`)
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, `do everything except actually perform the database action (default: false)`)
	cmd.Flags().BoolVarP(&secure, "secure", "y", false, `increase number of symbols from 0 to 3 (default: false)`)

	return cmd
}
