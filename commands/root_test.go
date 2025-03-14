// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// CreateContextForTesting creates an in-memory database and a context around it.
func CreateContextForTesting(t *testing.T) Context {
	// Create an in-memory database.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open() failed: %s", err)
	}

	t.Cleanup(func() { db.Close() })

	// Make sure that open/close database works with our test database.
	oldOpenDatabase := OpenDatabase
	OpenDatabase = OpenDatabaseForTesting(db)
	t.Cleanup(func() { OpenDatabase = oldOpenDatabase })
	oldCloseDatabase := CloseDatabase
	CloseDatabase = CloseDatabaseForTesting
	t.Cleanup(func() { CloseDatabase = oldCloseDatabase })

	oldGeneratePassword := GeneratePassword
	GeneratePassword = GeneratePasswordForTesting
	t.Cleanup(func() { GeneratePassword = oldGeneratePassword })

	oldGenerateTotpCode := GenerateTotpCode
	GenerateTotpCode = GenerateTotpCodeForTesting
	t.Cleanup(func() { GenerateTotpCode = oldGenerateTotpCode })

	oldNow := Now
	Now = NowForTesting
	t.Cleanup(func() { Now = oldNow })

	ctx := Context{Database: db}
	err = initDatabase(&ctx)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}

	return ctx
}

// OpenDatabaseForTesting implements OpenDatabase and takes an already opened sql.DB.
func OpenDatabaseForTesting(sqlDb *sql.DB) func(*Context) error {
	return func(ctx *Context) error {
		ctx.Database = sqlDb
		return nil
	}
}

// CloseDatabaseForTesting implements CloseDatabase and does nothing.
func CloseDatabaseForTesting(ctx *Context) error {
	return nil
}

func CommandForTesting(t *testing.T) func(name string, arg ...string) *exec.Cmd {
	return func(name string, arg ...string) *exec.Cmd {
		if len(arg) == 5 && name == "gpg" && arg[0] == "--decrypt" && arg[1] == "-a" && arg[2] == "-o" {
			decryptedPath := arg[3]
			encryptedPath := arg[4]
			var encryptedQaPath string
			if strings.HasSuffix(encryptedPath, ".cpmdb") {
				encryptedQaPath = "fixtures/cpmdb.xml"
			} else if strings.HasSuffix(encryptedPath, "passwords.db") {
				encryptedQaPath = "fixtures/passwords.db"
			} else {
				t.Fatalf("unexpected encryted path: %s", encryptedPath)
			}
			err := CopyPath(encryptedQaPath, decryptedPath)
			if err != nil {
				t.Fatalf("CopyPath() failed: %s", err)
			}
			return exec.Command("true")
		} else if len(arg) == 7 && name == "gpg" && arg[0] == "--encrypt" && arg[1] == "--sign" && arg[2] == "-a" && arg[3] == "--default-recipient-self" && arg[4] == "-o" {
			encryptedPath := arg[5]
			decryptedPath := arg[6]
			var encryptedQaPath string
			if strings.HasSuffix(encryptedPath, "passwords.db") {
				encryptedQaPath = "fixtures/passwords.db"
			} else {
				t.Fatalf("unexpected encryted path: %s", encryptedPath)
			}
			err := CopyPath(decryptedPath, encryptedQaPath)
			if err != nil {
				t.Fatalf("CopyPath() failed: %s", err)
			}
			return exec.Command("true")
		} else if len(arg) == 2 && name == "gunzip" && arg[0] == "--force" {
			compressedPath := arg[1]
			uncompressedPath := strings.ReplaceAll(compressedPath, ".gz", "")
			err := CopyPath(compressedPath, uncompressedPath)
			if err != nil {
				t.Fatalf("CopyPath() failed: %s", err)
			}
			return exec.Command("true")
		} else if name == "scp" && len(arg) == 2 && strings.HasPrefix(arg[0], "cpm:") && strings.HasSuffix(arg[0], "passwords.db") && strings.HasSuffix(arg[1], "passwords.db") {
			err := CopyPath("fixtures/remote.db", "fixtures/passwords.db")
			if err != nil {
				t.Fatalf("CopyPath() failed: %s", err)
			}
			return exec.Command("true")
		}
		t.Fatalf("CommandForTesting: unhandled command: %v", arg)
		panic("unreachable")
	}
}

func RemoveForTesting(name string) error {
	if strings.HasSuffix(name, "passwords.db") {
		return os.Remove("fixtures/passwords.db")
	}

	return os.Remove(name)
}

func StatForTesting(name string) (os.FileInfo, error) {
	if strings.HasSuffix(name, "passwords.db") {
		return os.Stat("fixtures/passwords.db")
	}

	return os.Stat(name)
}

// ContainsString checks if `items` contains `item`.
func ContainsString(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

// CopyPath copies from inPath to outPath, assuming they are file paths.
func CopyPath(inPath, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("os.Create() failed: %s", err)
	}

	inFile, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("os.Open() failed: %s", err)
	}

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		return fmt.Errorf("io.Copy() failed: %s", err)
	}

	return nil
}

func UseCommandForTesting(t *testing.T) {
	oldCommand := Command
	Command = CommandForTesting(t)
	t.Cleanup(func() { Command = oldCommand })
}

// TestGetDatabasePath checks if getDatabasePath() handles a custom XDG_STATE_HOME value.
func TestGetDatabasePath(t *testing.T) {
	oldEnv := os.Getenv(xdgStateHome)
	os.Setenv(xdgStateHome, "/tmp")
	defer func() { os.Setenv(xdgStateHome, oldEnv) }()

	actual, err := getDatabasePath()
	if err != nil {
		t.Fatalf("getDatabasePath() err = %q, want nil", err)
	}

	expected := "/tmp/cpm/passwords.db"
	if actual != expected {
		t.Fatalf("getDatabasePath() = %q, want %q", actual, expected)
	}
}
