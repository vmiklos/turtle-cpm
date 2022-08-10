package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDeleteCommand(ctx *Context) *cobra.Command {
	var machine string
	var service string
	var user string
	var passwordType PasswordType = "plain"
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes an existing password",
		RunE: func(cmd *cobra.Command, args []string) error {
			query, err := ctx.Database.Prepare("delete from passwords where machine=? and service=? and user=? and type=?")
			if err != nil {
				return fmt.Errorf("db.Prepare() failed: %s", err)
			}

			result, err := query.Exec(machine, service, user, passwordType)
			if err != nil {
				return fmt.Errorf("db.Exec() failed: %s", err)
			}

			affected, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("result.RowsAffected() failed: %s", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted %v password\n", affected)

			return nil
		},
	}
	cmd.Flags().StringVarP(&machine, "machine", "m", "", "machine (required)")
	cmd.MarkFlagRequired("machine")
	cmd.Flags().StringVarP(&service, "service", "s", "http", "service")
	cmd.Flags().StringVarP(&user, "user", "u", "", "user (required)")
	cmd.MarkFlagRequired("user")
	cmd.Flags().VarP(&passwordType, "type", "t", `password type ("plain" or "totp")`)

	return cmd
}
