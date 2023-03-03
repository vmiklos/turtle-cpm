// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	// register sqlite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

const (
	xdgStateHome = "XDG_STATE_HOME"
	// Version specifies the number for the version subcommand
	Version = "7.5"
)

// Command returns the Cmd struct to execute the named program
var Command = exec.Command

// Remove removes the named file or (empty) directory.
var Remove = os.Remove

// Stat returns a FileInfo describing the named file.
var Stat = os.Stat

// OpenDatabase opens the database before running a subcommand.
var OpenDatabase = openDatabase

// CloseDatabase opens the database before running a subcommand.
var CloseDatabase = closeDatabase

// NewRootCommand creates the parent of all subcommands.
func NewRootCommand(ctx *Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "cpm",
		Short: "turtle-cpm is a console password manager",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(os.Args) >= 2 && os.Args[1] == "version" {
				return nil
			}

			err := OpenDatabase(ctx)
			if err != nil {
				return fmt.Errorf("OpenDatabase() failed: %s", err)
			}

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			err := CloseDatabase(ctx)
			if err != nil {
				return fmt.Errorf("CloseDatabase() failed: %s", err)
			}

			return nil
		},
	}
	cmd.AddCommand(newCreateCommand(ctx))
	cmd.AddCommand(newReadCommand(ctx))
	cmd.AddCommand(newUpdateCommand(ctx))
	cmd.AddCommand(newDeleteCommand(ctx))
	cmd.AddCommand(newImportCommand(ctx))
	cmd.AddCommand(newPullCommand(ctx))
	cmd.AddCommand(newVersionCommand(ctx))

	return cmd
}

func getCommands() []string {
	return []string{
		"-h",
		"--help",
		"completion",
		"__complete",
		"create",
		"delete",
		"help",
		"import",
		"pull",
		"search",
		"update",
		"version",
	}
}

// Context is state that is preserved during PreRun / Run / PostRun.
type Context struct {
	TempFile         *os.File
	PermanentPath    string
	Database         *sql.DB
	NoWriteBack      bool
	DryRun           bool
	DatabaseMigrated bool
	OutOrStdout      *io.Writer
}

func pathExists(path string) bool {
	_, err := Stat(path)
	return err == nil
}

func getDatabasePath() (string, error) {
	var databaseDir string
	if a := os.Getenv(xdgStateHome); a != "" {
		databaseDir = filepath.Join(a, "cpm")
	} else {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("user.Current() failed: %s", err)
		}
		databaseDir = filepath.Join(usr.HomeDir, ".local", "state", "cpm")
	}

	databasePath := databaseDir + "/passwords.db"
	if !pathExists(databasePath) {
		err := os.MkdirAll(databaseDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("os.MkdirAll() failed: %s", err)
		}
	}

	return databasePath, nil
}

func openDatabase(ctx *Context) error {
	var err error
	ctx.TempFile, err = ioutil.TempFile("", "cpm")
	if err != nil {
		return fmt.Errorf("ioutil.TempFile() failed: %s", err)
	}

	ctx.PermanentPath, err = getDatabasePath()
	if err != nil {
		return fmt.Errorf("getDatabasePath() failed: %s", err)
	}
	createNew := true
	if pathExists(ctx.PermanentPath) {
		Remove(ctx.TempFile.Name())
		command := Command("gpg", "--decrypt", "-a", "-o", ctx.TempFile.Name(), ctx.PermanentPath)
		err = command.Run()
		if err != nil {
			return fmt.Errorf("Command() failed: %s", err)
		}
		createNew = false
	}

	ctx.Database, err = sql.Open("sqlite3", ctx.TempFile.Name())
	if err != nil {
		return fmt.Errorf("sql.Open() failed: %s", err)
	}

	err = initDatabase(ctx, createNew)
	if err != nil {
		return fmt.Errorf("initDatabase() failed: %s", err)
	}

	return nil
}

func initDatabaseWithVersion(ctx *Context, version int) error {
	var statement string
	if version == 0 {
		statement = `create table passwords (
		machine text not null,
		service text not null,
		user text not null,
		password text not null,
		type text not null,
		unique(machine, service, user, type)
		)`
	} else {
		statement = `create table passwords (
		id integer primary key autoincrement,
		machine text not null,
		service text not null,
		user text not null,
		password text not null,
		type text not null,
		unique(machine, service, user, type)
		)`
	}
	query, err := ctx.Database.Prepare(statement)
	if err != nil {
		return fmt.Errorf("db.Prepare() failed: %s", err)
	}
	_, err = query.Exec()
	if err != nil {
		return fmt.Errorf("db.Exec() failed: %s", err)
	}

	return nil
}

func initDatabase(ctx *Context, createNew bool) error {
	// We need createNew because both an empty db and the first schema was user_version == 0.
	if createNew {
		initDatabaseWithVersion(ctx, 1)

		query, err := ctx.Database.Prepare("pragma user_version = 1")
		if err != nil {
			return fmt.Errorf("db.Prepare() failed: %s", err)
		}
		_, err = query.Exec()
		if err != nil {
			return fmt.Errorf("db.Exec() failed: %s", err)
		}
	}

	var version int
	rows, err := ctx.Database.Query("pragma user_version")
	if err != nil {
		return fmt.Errorf("db.Query(pragma) failed: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			return fmt.Errorf("rows.Scan() failed: %s", err)
		}
	}

	if version < 1 {
		statements := []string{`create table passwords_copy(
		id integer primary key autoincrement,
		machine text not null,
		service text not null,
		user text not null,
		password text not null,
		type text not null,
		unique(machine, service, user, type)
	)`,
			"insert into passwords_copy(machine, service, user, password, type) select machine, service, user, password, type from passwords",
			"drop table passwords",
			"alter table passwords_copy rename to passwords",
		}
		for _, statement := range statements {
			query, err := ctx.Database.Prepare(statement)
			if err != nil {
				return fmt.Errorf("db.Prepare() failed: %s", err)
			}
			_, err = query.Exec()
			if err != nil {
				return fmt.Errorf("db.Exec() failed: %s", err)
			}
		}

		query, err := ctx.Database.Prepare("pragma user_version = 1")
		if err != nil {
			return fmt.Errorf("db.Prepare() failed: %s", err)
		}
		_, err = query.Exec()
		if err != nil {
			return fmt.Errorf("db.Exec() failed: %s", err)
		}
		ctx.DatabaseMigrated = true
	}

	return nil
}

// The database is only closed in case of no errors.
func closeDatabase(ctx *Context) error {
	if ctx.NoWriteBack && !ctx.DatabaseMigrated {
		return nil
	}

	if ctx.Database != nil {
		err := ctx.Database.Close()
		if err != nil {
			return fmt.Errorf("db.Database.Close() failed: %s", err)
		}
	}

	Remove(ctx.PermanentPath)
	command := Command("gpg", "--encrypt", "--sign", "-a", "--default-recipient-self", "-o", ctx.PermanentPath, ctx.TempFile.Name())
	err := command.Run()
	if err != nil {
		return fmt.Errorf("Command() failed to run 'gpg --encrypt --sign -a --default-recipient-self -o %s %s': %s", ctx.PermanentPath, ctx.TempFile.Name(), err)
	}

	return nil
}

// The database is always cleaned to avoid decrypted data on disk (even in case of a failure).
func cleanDatabase(ctx *Context) {
	if ctx.TempFile != nil {
		Remove(ctx.TempFile.Name())
	}
}

// runCommand is a wrapper around Command() to invoke it in an interactive mode.
func runCommand(name string, arg ...string) error {
	cmd := Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd.Run(%s) failed: %s", name, err)
	}

	return nil
}

// PasswordType is an enum of possible password types.
type PasswordType string

const (
	// PasswordTypePlain is a password sent to a server as-is.
	PasswordTypePlain PasswordType = "plain"
	// PasswordTypeTotp is a TOTP shared secret.
	PasswordTypeTotp PasswordType = "totp"
)

func (t *PasswordType) String() string {
	return string(*t)
}

// Set sets the value of `t` from `v`.
func (t *PasswordType) Set(v string) error {
	switch v {
	case "plain", "totp":
		*t = PasswordType(v)
		return nil
	default:
		return errors.New(`must be one of "plain", or "totp"`)
	}
}

// Type returns the type of `t` as a string.
func (t *PasswordType) Type() string {
	return "PasswordType"
}

// Main is the commandline interface to this package.
func Main(input io.Reader, output io.Writer) int {
	var ctx Context
	defer cleanDatabase(&ctx)

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
	var cmd = NewRootCommand(&ctx)
	var args []string
	if commandFound {
		args = os.Args[1:]
	} else {
		// Default to the search subcommand.
		args = append([]string{"search"}, os.Args[1:]...)
	}
	cmd.SetArgs(args)
	cmd.SetIn(input)
	cmd.SetOut(output)

	err := cmd.Execute()
	if err != nil {
		// cobra reported its error already itself.
		return 1
	}

	return 0
}
