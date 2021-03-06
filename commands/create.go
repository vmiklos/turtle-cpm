package commands

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func generatePassword() (string, error) {
	// Length of 15 and no symbols matches current Firefox.
	output, err := Command("pwgen", "--secure", "15", "1").Output()
	if err != nil {
		return "", fmt.Errorf("Command(pwgen) failed: %s", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func createPassword(db *sql.DB, machine, service, user, password, passwordType string) (string, error) {
	if len(password) == 0 {
		var err error
		password, err = generatePassword()
		if err != nil {
			return "", fmt.Errorf("generatePassword() failed: %s", err)
		}
	}

	query, err := db.Prepare("insert into passwords (machine, service, user, password, type) values(?, ?, ?, ?, ?)")
	if err != nil {
		return "", fmt.Errorf("db.Prepare() failed: %s", err)
	}

	_, err = query.Exec(machine, service, user, password, passwordType)
	if err != nil {
		return "", fmt.Errorf("query.Exec() failed: %s", err)
	}
	return password, nil
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
			generatedPassword, err := createPassword(ctx.Database, machine, service, user, password, passwordType)
			if err != nil {
				return fmt.Errorf("createPassword() failed: %s", err)
			}

			if generatedPassword != password {
				fmt.Fprintf(cmd.OutOrStdout(), "Generated password: %s\n", generatedPassword)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&machine, "machine", "m", "", "machine (required)")
	cmd.MarkFlagRequired("machine")
	cmd.Flags().StringVarP(&service, "service", "s", "http", "service")
	cmd.Flags().StringVarP(&user, "user", "u", "", "user (required)")
	cmd.MarkFlagRequired("user")
	cmd.Flags().StringVarP(&password, "password", "p", "", "password")
	cmd.Flags().StringVarP(&passwordType, "type", "t", "plain", `password type ("plain" or "totp")`)

	return cmd
}
