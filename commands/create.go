package commands

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

func createPassword(out io.Writer, db *sql.DB, machine, service, user, password, passwordType string) error {
	if len(password) == 0 {
		// Length of 15 and no symbols matches current Firefox.
		output, err := Command("pwgen", "--secure", "15", "1").Output()
		if err != nil {
			return fmt.Errorf("Command(pwgen) failed: %s", err)
		}
		password = strings.TrimSpace(string(output))
		fmt.Fprintf(out, "Generated password: %s\n", password)
	}

	query, err := db.Prepare("insert into passwords (machine, service, user, password, type) values(?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("db.Prepare() failed: %s", err)
	}

	_, err = query.Exec(machine, service, user, password, passwordType)
	if err != nil {
		return fmt.Errorf("query.Exec() failed: %s", err)
	}
	return nil
}

func newCreateCommand(ctx *Context) *cobra.Command {
	var machine string
	var service string
	var user string
	var password string
	var passwordType string
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "creates a new password",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := createPassword(cmd.OutOrStdout(), ctx.Database, machine, service, user, password, passwordType)
			if err != nil {
				return fmt.Errorf("createPassword() failed: %s", err)
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&machine, "machine", "m", "", "machine (required)")
	cmd.MarkFlagRequired("machine")
	cmd.Flags().StringVarP(&service, "service", "s", "", "service (required)")
	cmd.MarkFlagRequired("service")
	cmd.Flags().StringVarP(&user, "user", "u", "", "user (required)")
	cmd.MarkFlagRequired("user")
	cmd.Flags().StringVarP(&password, "password", "p", "", "password")
	cmd.Flags().StringVarP(&passwordType, "type", "t", "plain", "password type ('plain' or 'totp', default: plain)")

	return cmd
}
