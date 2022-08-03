package commands

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"

	// register sqlite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
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
	cmd.AddCommand(newSyncCommand(ctx))

	return cmd
}

func getCommands() []string {
	return []string{
		"-h",
		"--help",
		"completion",
		"create",
		"delete",
		"help",
		"import",
		"search",
		"update",
		"sync",
	}
}

// Context is state that is preserved during PreRun / Run / PostRun.
type Context struct {
	TempFile      *os.File
	PermanentPath string
	Database      *sql.DB
	NoWriteBack   bool
}

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
	if pathExists(ctx.PermanentPath) {
		Remove(ctx.TempFile.Name())
		command := Command("gpg", "--decrypt", "-a", "-o", ctx.TempFile.Name(), ctx.PermanentPath)
		err = command.Run()
		if err != nil {
			return fmt.Errorf("Command() failed: %s", err)
		}
	}

	ctx.Database, err = sql.Open("sqlite3", ctx.TempFile.Name())
	if err != nil {
		return fmt.Errorf("sql.Open() failed: %s", err)
	}

	err = initDatabase(ctx.Database)
	if err != nil {
		return fmt.Errorf("initDatabase() failed: %s", err)
	}

	return nil
}

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
func closeDatabase(ctx *Context) error {
	err := ctx.Database.Close()
	if err != nil {
		return fmt.Errorf("db.Database.Close() failed: %s", err)
	}

	if ctx.NoWriteBack {
		return nil
	}

	Remove(ctx.PermanentPath)
	command := Command("gpg", "--encrypt", "--sign", "-a", "--default-recipient-self", "-o", ctx.PermanentPath, ctx.TempFile.Name())
	err = command.Run()
	if err != nil {
		return fmt.Errorf("Command() failed: %s", err)
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

// Main is the commandline interface to this package.
func Main(stream io.Writer) int {
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
	cmd.SetOut(stream)
	cmd.SetErr(stream)

	err := cmd.Execute()
	if err != nil {
		// cobra reported its error already itself.
		return 1
	}

	return 0
}
