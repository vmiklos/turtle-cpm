package commands

import (
	"bufio"
	"database/sql"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newUpdateCommand(ctx *Context) *cobra.Command {
	var machine string
	var service string
	var user string
	var password string
	var passwordType PasswordType = "plain"
	var dryRun bool
	var id string
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
			if len(id) > 0 {
				var result sql.Result
				if len(machine) > 0 {
					query, err := transaction.Prepare("update passwords set machine=? where id=?")
					if err != nil {
						return fmt.Errorf("db.Prepare() failed: %s", err)
					}

					result, err = query.Exec(machine, id)
					if err != nil {
						return fmt.Errorf("db.Exec() failed: %s", err)
					}
				}

				affected, err = result.RowsAffected()
				if err != nil {
					return fmt.Errorf("result.RowsAffected() failed: %s", err)
				}
			} else {
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

				generatedPassword := password
				if len(password) == 0 {
					var err error
					generatedPassword, err = generatePassword()
					if err != nil {
						return fmt.Errorf("generatePassword() failed: %s", err)
					}
				}

				query, err := transaction.Prepare("update passwords set password=? where machine=? and service=? and user=? and type=?")
				if err != nil {
					return fmt.Errorf("db.Prepare() failed: %s", err)
				}

				result, err := query.Exec(generatedPassword, machine, service, user, passwordType)
				if err != nil {
					return fmt.Errorf("db.Exec() failed: %s", err)
				}

				if generatedPassword != password {
					fmt.Fprintf(cmd.OutOrStdout(), "Generated new password: %s\n", generatedPassword)
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
			return nil
		},
	}
	cmd.Flags().StringVarP(&machine, "machine", "m", "", "machine (default: ask)")
	cmd.Flags().StringVarP(&service, "service", "s", "http", "service")
	cmd.Flags().StringVarP(&user, "user", "u", "", "user (default: ask)")
	cmd.Flags().StringVarP(&password, "password", "p", "", "new password")
	cmd.Flags().VarP(&passwordType, "type", "t", `password type ("plain" or "totp")`)
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, `do everything except actually perform the database action (default: false)`)
	cmd.Flags().StringVarP(&id, "id", "i", "", `other parameters specify new values for this id (default: '')`)

	return cmd
}
