// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
)

// parsePassword parses a TOTP shared secret out of an otpauth:// URL, or just returns the input
// as-is.
func parsePassword(s string) (string, error) {
	if !strings.HasPrefix(s, "otpauth://") {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("url.Parse() failed: %s", err)
	}

	keyValues, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", fmt.Errorf("url.ParseQuery() failed: %s", err)
	}

	if !keyValues.Has("secret") {
		return "", fmt.Errorf("no 'secret' key in URL")
	}

	secrets := keyValues["secret"]
	return secrets[0], nil
}

type searchOptions struct {
	wantedMachine string
	wantedService string
	wantedUser    string
	wantedType    PasswordType
	totp          bool
	quiet         bool
	qrcode        bool
	noid          bool
	args          []string
}

// GenerateQrCode creates a QR Code and writes it out to io.Writer.
var GenerateQrCode = qrterminal.Generate

func readPasswords(db *sql.DB, opts searchOptions) ([]string, error) {
	var results []string
	if opts.totp {
		opts.wantedType = "totp"
	}
	rows, err := db.Query("select id, machine, service, user, password, type from passwords")
	if err != nil {
		return nil, fmt.Errorf("db.Query(insert) failed: %s", err)
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var machine string
		var service string
		var user string
		var password string
		var passwordType PasswordType
		err = rows.Scan(&id, &machine, &service, &user, &password, &passwordType)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan() failed: %s", err)
		}

		if len(opts.wantedMachine) > 0 && machine != opts.wantedMachine {
			continue
		}

		if len(opts.wantedService) > 0 && service != opts.wantedService {
			continue
		}

		if len(opts.wantedUser) > 0 && user != opts.wantedUser {
			continue
		}

		if len(opts.wantedType) > 0 && passwordType != opts.wantedType {
			continue
		}

		if len(opts.args) > 0 {
			// Allow simply matching a sub-string: e.g. search for a service type or a part
			// of a machine without explicitly telling if the query is a service or a
			// machine.
			s := fmt.Sprintf("%d %s %s %s %s", id, machine, service, user, passwordType)
			if !strings.Contains(s, opts.args[0]) {
				continue
			}
		}

		if passwordType == "totp" {
			if opts.totp {
				// This is a TOTP password and the current value is required: invoke
				// oathtool to generate it.
				passwordType = "TOTP code"
				sharedSecret, err := parsePassword(password)
				if err != nil {
					return nil, fmt.Errorf("parsePassword() failed: %s", err)
				}

				output, err := Command("oathtool", "-b", "--totp", sharedSecret).Output()
				if err != nil {
					return nil, fmt.Errorf("exec.Command(oathtool) failed: %s", err)
				}
				password = strings.TrimSpace(string(output))
			} else {
				passwordType = "TOTP shared secret"
			}
		}

		var result string
		if opts.quiet {
			result = password
		} else {
			if !opts.noid {
				result = fmt.Sprintf("id: %8d, ", id)
			}
			result += fmt.Sprintf("machine: %s, service: %s, user: %s, password type: %s, password:", machine, service, user, passwordType)
			if opts.qrcode {
				qrcode := new(bytes.Buffer)
				GenerateQrCode(password, qrterminal.L, qrcode)
				result += fmt.Sprintf("\n%s", qrcode)
			} else {
				result += fmt.Sprintf(" %s", password)
			}
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
	var qrcodeFlag bool
	var noidFlag bool
	var cmd = &cobra.Command{
		Use:   "search",
		Short: "searches passwords",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(machineFlag) == 0 && len(serviceFlag) == 0 && len(userFlag) == 0 && len(typeFlag) == 0 && len(args) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Search term: ")
				reader := bufio.NewReader(cmd.InOrStdin())
				term, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("ReadString() failed: %s", err)
				}

				args = append(args, strings.TrimSuffix(term, "\n"))
			}

			opts := searchOptions{}
			opts.wantedMachine = machineFlag
			opts.wantedService = serviceFlag
			opts.wantedUser = userFlag
			opts.wantedType = typeFlag
			opts.totp = totpFlag
			opts.quiet = quietFlag
			opts.qrcode = qrcodeFlag
			opts.noid = noidFlag
			opts.args = args
			results, err := readPasswords(ctx.Database, opts)
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
	cmd.Flags().BoolVarP(&totpFlag, "totp", "T", false, `show the current TOTP code, not the TOTP shared secret (default: false, implies "--type totp")`)
	cmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "quite mode: only print the password itself (default: false)")
	cmd.Flags().BoolVarP(&qrcodeFlag, "qrcode", "Q", false, "qrcode mode: print the TOTP shared secret as a QR code (default: false)")
	cmd.Flags().BoolVarP(&noidFlag, "noid", "I", false, "noid mode: omit password ID from the output (default: false)")

	return cmd
}
