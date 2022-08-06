package commands

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func readPasswords(db *sql.DB, wantedMachine, wantedService, wantedUser string, wantedType PasswordType, totp, quiet bool, args []string) ([]string, error) {
	var results []string
	if totp {
		wantedType = "totp"
	}
	rows, err := db.Query("select machine, service, user, password, type from passwords")
	if err != nil {
		return nil, fmt.Errorf("db.Query(insert) failed: %s", err)
	}

	defer rows.Close()
	for rows.Next() {
		var machine string
		var service string
		var user string
		var password string
		var passwordType PasswordType
		err = rows.Scan(&machine, &service, &user, &password, &passwordType)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan() failed: %s", err)
		}

		if len(wantedMachine) > 0 && machine != wantedMachine {
			continue
		}

		if len(wantedService) > 0 && service != wantedService {
			continue
		}

		if len(wantedUser) > 0 && user != wantedUser {
			continue
		}

		if len(wantedType) > 0 && passwordType != wantedType {
			continue
		}

		if len(args) > 0 {
			// Allow simply matching a sub-string: e.g. search for a service type or a part
			// of a machine without explicitly telling if the query is a service or a
			// machine.
			s := fmt.Sprintf("%s %s %s %s", machine, service, user, passwordType)
			if !strings.Contains(s, args[0]) {
				continue
			}
		}

		if passwordType == "totp" {
			if totp {
				// This is a TOTP password and the current value is required: invoke
				// oathtool to generate it.
				passwordType = "TOTP code"
				output, err := Command("oathtool", "-b", "--totp", password).Output()
				if err != nil {
					return nil, fmt.Errorf("exec.Command(oathtool) failed: %s", err)
				}
				password = strings.TrimSpace(string(output))
			} else {
				passwordType = "TOTP shared secret"
			}
		}

		var result string
		if quiet {
			result = password
		} else {
			result = fmt.Sprintf("machine: %s, service: %s, user: %s, password type: %s, password: %s", machine, service, user, passwordType, password)
		}
		results = append(results, result)
	}

	return results, nil
}

func newReadCommand(ctx *Context) *cobra.Command {
	var machineFlag string
	var serviceFlag string
	var userFlag string
	// show all types by default
	var typeFlag PasswordType
	var totpFlag bool
	var quietFlag bool
	var cmd = &cobra.Command{
		Use:   "search",
		Short: "searches passwords",
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := readPasswords(ctx.Database, machineFlag, serviceFlag, userFlag, typeFlag, totpFlag, quietFlag, args)
			if err != nil {
				return fmt.Errorf("readPasswords() failed: %s", err)
			}

			for _, result := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", result)
			}

			ctx.NoWriteBack = true
			return nil
		},
	}
	cmd.Flags().StringVarP(&machineFlag, "machine", "m", "", `machine (default: "")`)
	cmd.Flags().StringVarP(&serviceFlag, "service", "s", "", `service (default: "")`)
	cmd.Flags().StringVarP(&userFlag, "user", "u", "", `user (default: "")`)
	cmd.Flags().VarP(&typeFlag, "type", "t", `password type ("plain" or "totp", default: "")`)
	cmd.Flags().BoolVarP(&totpFlag, "totp", "T", false, `show current TOTP, not the TOTP key (default: false, implies "--type totp")`)
	cmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "quite mode: only print the password itself (default: false)")

	return cmd
}
