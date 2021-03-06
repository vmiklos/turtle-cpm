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

// CreateDatabaseForTesting creates an in-memory database.
func CreateDatabaseForTesting() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("sql.Open() failed: %s", err)
	}

	return db, nil
}

// OpenDatabaseForTesting implements OpenDatabase and takes an already opened sql.DB.
func OpenDatabaseForTesting(sqlDb *sql.DB) func(*Context) error {
	return func(ctx *Context) error {
		ctx.Database = sqlDb

		err := initDatabase(ctx.Database)
		if err != nil {
			return fmt.Errorf("initDatabase() failed: %s", err)
		}

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
				encryptedQaPath = "qa/cpmdb.xml"
			} else if strings.HasSuffix(encryptedPath, "passwords.db") {
				encryptedQaPath = "qa/passwords.db"
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
				encryptedQaPath = "qa/passwords.db"
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
		} else if len(arg) == 3 && name == "oathtool" && arg[0] == "-b" && arg[1] == "--totp" && arg[2] == "totppassword" {
			return exec.Command("echo", "output-from-oathtool")
		} else if name == "pwgen" {
			return exec.Command("echo", "output-from-pwgen")
		}
		t.Fatalf("CommandForTesting: unhandled command: %v", arg)
		panic("unreachable")
	}
}

func RemoveForTesting(name string) error {
	if strings.HasSuffix(name, "passwords.db") {
		return os.Remove("qa/passwords.db")
	}

	return os.Remove(name)
}

func StatForTesting(name string) (os.FileInfo, error) {
	if strings.HasSuffix(name, "passwords.db") {
		return os.Stat("qa/passwords.db")
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
