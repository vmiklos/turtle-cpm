package main

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func createPassword(db *sql.DB, machine, service, user, password, passwordType string) error {
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

func readPasswords(db *sql.DB, wantedMachine, wantedService, wantedUser, wantedType string, totp bool, args []string) ([]string, error) {
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
		var passwordType string
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

		results = append(results, fmt.Sprintf("machine: %s, service: %s, user: %s, password type: %s, password: %s", machine, service, user, passwordType, password))
	}

	return results, nil
}

func newCreateCommand(db *sql.DB) *cobra.Command {
	var machine string
	var service string
	var user string
	var password string
	var passwordType string
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "creates a new password",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := createPassword(db, machine, service, user, password, passwordType)
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
	cmd.Flags().StringVarP(&password, "password", "p", "", "password (required)")
	cmd.MarkFlagRequired("password")
	cmd.Flags().StringVarP(&passwordType, "type", "t", "plain", "password type ('plain' or 'totp', default: plain)")

	return cmd
}

func newUpdateCommand(db *sql.DB) *cobra.Command {
	var machine string
	var service string
	var user string
	var password string
	var passwordType string
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "updates an existing password",
		RunE: func(cmd *cobra.Command, args []string) error {
			query, err := db.Prepare("update passwords set password=? where machine=? and service=? and user=? and type=?")
			if err != nil {
				return fmt.Errorf("db.Prepare() failed: %s", err)
			}

			_, err = query.Exec(password, machine, service, user, passwordType)
			if err != nil {
				return fmt.Errorf("db.Exec() failed: %s", err)
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
	cmd.Flags().StringVarP(&password, "password", "p", "", "new password (required)")
	cmd.MarkFlagRequired("password")
	cmd.Flags().StringVarP(&passwordType, "type", "t", "plain", "password type ('plain' or 'totp', default: plain)")

	return cmd
}

func newDeleteCommand(db *sql.DB) *cobra.Command {
	var machine string
	var service string
	var user string
	var passwordType string
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "deletes an existing password",
		RunE: func(cmd *cobra.Command, args []string) error {
			query, err := db.Prepare("delete from passwords where machine=? and service=? and user=? and type=?")
			if err != nil {
				return fmt.Errorf("db.Prepare() failed: %s", err)
			}

			_, err = query.Exec(machine, service, user, passwordType)
			if err != nil {
				return fmt.Errorf("db.Exec() failed: %s", err)
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
	cmd.Flags().StringVarP(&passwordType, "type", "t", "plain", "password type ('plain' or 'totp', default: plain)")

	return cmd
}

// XMLPassword is the 4th <node> element from cpm's XML database.
type XMLPassword struct {
	XMLName xml.Name `xml:"node"`
	Label   string   `xml:"label,attr"`
	Totp    string   `xml:"totp,attr"`
}

// XMLUser is the 3rd <node> element from cpm's XML database.
type XMLUser struct {
	XMLName   xml.Name      `xml:"node"`
	Label     string        `xml:"label,attr"`
	Passwords []XMLPassword `xml:"node"`
}

// XMLService is the 2nd <node> element from cpm's XML database.
type XMLService struct {
	XMLName xml.Name  `xml:"node"`
	Label   string    `xml:"label,attr"`
	Users   []XMLUser `xml:"node"`
}

// XMLMachine is the 1st <node> element from cpm's XML database.
type XMLMachine struct {
	XMLName  xml.Name     `xml:"node"`
	Label    string       `xml:"label,attr"`
	Services []XMLService `xml:"node"`
}

// XMLMachines is the <root> element from cpm's XML database.
type XMLMachines struct {
	XMLName  xml.Name     `xml:"root"`
	Machines []XMLMachine `xml:"node"`
}

// Command returns the Cmd struct to execute the named program
var Command = exec.Command

// Remove removes the named file or (empty) directory.
var Remove = os.Remove

func newImportCommand(db *sql.DB) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "import",
		Short: "imports an old XML database",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Decrypt and uncompress ~/.cpmdb to a temp file.
			usr, err := user.Current()
			if err != nil {
				return fmt.Errorf("user.Current() failed: %s", err)
			}

			encryptedPath := usr.HomeDir + "/.cpmdb"
			decryptedFile, err := ioutil.TempFile("", "cpm")
			decryptedPath := decryptedFile.Name()
			defer Remove(decryptedPath)
			if err != nil {
				return fmt.Errorf("ioutil.TempFile() failed: %s", err)
			}

			Remove(decryptedPath)
			gpg := Command("gpg", "--decrypt", "-a", "-o", decryptedPath+".gz", encryptedPath)
			err = gpg.Start()
			if err != nil {
				return fmt.Errorf("cmd.Start() failed: %s", err)
			}
			err = gpg.Wait()
			if err != nil {
				return fmt.Errorf("cmd.Wait() failed: %s", err)
			}

			gunzip := Command("gunzip", decryptedPath+".gz")
			err = gunzip.Start()
			if err != nil {
				return fmt.Errorf("cmd.Start(gunzip) failed: %s", err)
			}
			err = gunzip.Wait()
			if err != nil {
				return fmt.Errorf("cmd.Wait(gunzip) failed: %s", err)
			}

			// Parse the XML.
			xmlFile, err := os.Open(decryptedPath)
			if err != nil {
				return fmt.Errorf("os.Open(decryptedPath) failed: %s", err)
			}
			defer xmlFile.Close()

			xmlBytes, err := ioutil.ReadAll(xmlFile)
			if err != nil {
				return fmt.Errorf("ioutil.ReadAll(xmlFile) failed: %s", err)
			}

			// Avoid 'encoding "ISO-8859-1" declared but Decoder.CharsetReader is nil'.
			xmlBytes = bytes.ReplaceAll(xmlBytes, []byte(`encoding="ISO-8859-1"`), []byte(`encoding="UTF-8"`))

			var machines XMLMachines
			err = xml.Unmarshal(xmlBytes, &machines)
			if err != nil {
				return fmt.Errorf("xml.Unmarshal() failed: %s", err)
			}

			// Import the parsed data.
			for _, machine := range machines.Machines {
				machineLabel := machine.Label
				for _, service := range machine.Services {
					serviceLabel := service.Label
					for _, user := range service.Users {
						userLabel := user.Label
						for _, password := range user.Passwords {
							passwordLabel := password.Label
							var passwordType string
							if password.Totp == "true" {
								passwordType = "totp"
							} else {
								passwordType = "plain"
							}

							err = createPassword(db, machineLabel, serviceLabel, userLabel, passwordLabel, passwordType)
							if err != nil {
								return fmt.Errorf("createPassword(machine='%s', service='%s', user='%s', type='%s') failed: %s", machineLabel, serviceLabel, userLabel, passwordType, err)
							}
						}
					}
				}
			}

			return nil
		},
	}

	return cmd
}

func newReadCommand(db *sql.DB) *cobra.Command {
	var machineFlag string
	var serviceFlag string
	var userFlag string
	var typeFlag string
	var totpFlag bool
	var cmd = &cobra.Command{
		Use:   "search",
		Short: "searches passwords",
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := readPasswords(db, machineFlag, serviceFlag, userFlag, typeFlag, totpFlag, args)
			if err != nil {
				return fmt.Errorf("readPasswords() failed: %s", err)
			}

			for _, result := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", result)
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&machineFlag, "machine", "m", "", "machine (required)")
	cmd.Flags().StringVarP(&serviceFlag, "service", "s", "", "service (required)")
	cmd.Flags().StringVarP(&userFlag, "user", "u", "", "user (required)")
	cmd.Flags().StringVarP(&typeFlag, "type", "t", "", "password type ('plain' or 'totp', default: '')")
	cmd.Flags().BoolVarP(&totpFlag, "totp", "T", false, "show current TOTP, not the TOTP key (default: false, implies '--type totp')")

	return cmd
}

func newRootCommand(db *sql.DB) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "cpm",
		Short: "turtle-cpm is a console password manager",
	}
	cmd.AddCommand(newCreateCommand(db))
	cmd.AddCommand(newReadCommand(db))
	cmd.AddCommand(newUpdateCommand(db))
	cmd.AddCommand(newDeleteCommand(db))
	cmd.AddCommand(newImportCommand(db))

	return cmd
}

func getCommands() []string {
	return []string{
		"--help",
		"completion",
		"create",
		"delete",
		"help",
		"import",
		"search",
		"update",
	}
}

// CpmDatabase is an opened tempfile, containing an sqlite database.
type CpmDatabase struct {
	TempFile      *os.File
	PermanentPath string
	Database      *sql.DB
}

// Stat returns a FileInfo describing the named file.
var Stat = os.Stat

func pathExists(path string) bool {
	_, err := Stat(path)
	return err == nil
}

func getDatabasePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("user.Current() failed: %s", err)
	}

	databaseDir := usr.HomeDir + "/.local/state/cpm"
	databasePath := databaseDir + "/passwords.db"
	if !pathExists(databasePath) {
		err = os.MkdirAll(databaseDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("os.MkdirAll() failed: %s", err)
		}
	}

	return databasePath, nil
}

func openDatabase() (*CpmDatabase, error) {
	var db CpmDatabase
	var err error
	db.TempFile, err = ioutil.TempFile("", "cpm")
	if err != nil {
		return nil, fmt.Errorf("ioutil.TempFile() failed: %s", err)
	}

	db.PermanentPath, err = getDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("getDatabasePath() failed: %s", err)
	}
	if pathExists(db.PermanentPath) {
		Remove(db.TempFile.Name())
		cmd := Command("gpg", "--decrypt", "-a", "-o", db.TempFile.Name(), db.PermanentPath)
		err := cmd.Start()
		if err != nil {
			return nil, fmt.Errorf("cmd.Start() failed: %s", err)
		}
		err = cmd.Wait()
		if err != nil {
			return nil, fmt.Errorf("cmd.Wait() failed: %s", err)
		}
	}

	db.Database, err = sql.Open("sqlite3", db.TempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("sql.Open() failed: %s", err)
	}

	err = initDatabase(db.Database)
	if err != nil {
		return nil, fmt.Errorf("initDatabase() failed: %s", err)
	}

	return &db, nil
}

// OpenDatabase opens the database before running a subcommand.
var OpenDatabase = openDatabase

func initDatabase(db *sql.DB) error {
	query, err := db.Prepare(`create table if not exists passwords (
		machine text not null,
		service text not null,
		user text not null,
		password text not null,
		type text not null,
		unique(machine, service, user, type)
	)`)
	if err != nil {
		return err
	}
	query.Exec()

	return nil
}

// The database is only closed in case of no errors.
func closeDatabase(db *CpmDatabase) error {
	err := db.Database.Close()
	if err != nil {
		return fmt.Errorf("db.Database.Close() failed: %s", err)
	}

	Remove(db.PermanentPath)
	cmd := Command("gpg", "--encrypt", "--sign", "-a", "--default-recipient-self", "-o", db.PermanentPath, db.TempFile.Name())
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("cmd.Start(gpg encrypt) failed: %s", err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("cmd.Wait(gpg encrypt) failed: %s", err)
	}

	return nil
}

// CloseDatabase opens the database before running a subcommand.
var CloseDatabase = closeDatabase

// The database is always cleaned to avoid decrypted data on disk (even in case of a failure).
func cleanDatabase(db *CpmDatabase) {
	if db.TempFile != nil {
		Remove(db.TempFile.Name())
	}
}

// Main is the commandline interface to this package.
func Main(stream io.Writer) int {
	db, err := OpenDatabase()
	if err != nil {
		// notest
		fmt.Fprintf(stream, "OpenDatabase() failed: %s", err)
		return 1
	}
	defer cleanDatabase(db)

	var commandFound bool
	commands := getCommands()
	for _, a := range commands {
		for _, b := range os.Args[1:] {
			if a == b {
				commandFound = true
				break
			}
		}
	}
	var cmd = newRootCommand(db.Database)
	var args []string
	if commandFound {
		args = os.Args[1:]
	} else {
		// Default to the search subcommand.
		args = append([]string{"search"}, os.Args[1:]...)
	}
	cmd.SetArgs(args)
	cmd.SetOut(stream)
	cmd.SetErr(stream)

	err = cmd.Execute()
	if err != nil {
		// cobra reported its error already itself.
		return 1
	}

	err = CloseDatabase(db)
	if err != nil {
		// notest
		fmt.Fprintf(stream, "CloseDatabase() failed: %s", err)
		return 1
	}
	return 0
}
