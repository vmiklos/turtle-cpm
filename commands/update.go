package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUpdateCommand(ctx *Context) *cobra.Command {
	var machine string
	var service string
	var user string
	var password string
	var passwordType PasswordType = "plain"
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "updates an existing password",
		RunE: func(cmd *cobra.Command, args []string) error {
			generatedPassword := password
			if len(password) == 0 {
				var err error
				generatedPassword, err = generatePassword()
				if err != nil {
					return fmt.Errorf("generatePassword() failed: %s", err)
				}
			}

			query, err := ctx.Database.Prepare("update passwords set password=? where machine=? and service=? and user=? and type=?")
			if err != nil {
				return fmt.Errorf("db.Prepare() failed: %s", err)
			}

			_, err = query.Exec(generatedPassword, machine, service, user, passwordType)
			if err != nil {
				return fmt.Errorf("db.Exec() failed: %s", err)
			}

			if generatedPassword != password {
				fmt.Fprintf(cmd.OutOrStdout(), "Generated new password: %s\n", generatedPassword)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&machine, "machine", "m", "", "machine (required)")
	cmd.MarkFlagRequired("machine")
	cmd.Flags().StringVarP(&service, "service", "s", "http", "service")
	cmd.Flags().StringVarP(&user, "user", "u", "", "user (required)")
	cmd.MarkFlagRequired("user")
	cmd.Flags().StringVarP(&password, "password", "p", "", "new password")
	cmd.Flags().VarP(&passwordType, "type", "t", `password type ("plain" or "totp")`)

	return cmd
}
