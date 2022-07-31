package commands

import (
	"bytes"
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

func ContainsString(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func TestUpdate(t *testing.T) {
	db, err := CreateDatabaseForTesting()
	defer db.Close()
	if err != nil {
		t.Fatalf("CreateDatabaseForTesting() err = %q, want nil", err)
	}
	OldOpenDatabase := OpenDatabase
	OpenDatabase = OpenDatabaseForTesting(db)
	defer func() { OpenDatabase = OldOpenDatabase }()
	OldCloseDatabase := CloseDatabase
	CloseDatabase = CloseDatabaseForTesting
	defer func() { CloseDatabase = OldCloseDatabase }()
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "newpassword"
	expectedType := "plain"
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, "oldpassword", expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "update", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := ""
	if buf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", buf.String(), expectedBuf)
	}
	results, err := readPasswords(db, "", "", "", "", false, false, []string{})
	if err != nil {
		t.Fatalf("readPasswords() err = %q, want nil", err)
	}
	actualLength := len(results)
	expectedLength := 1
	if actualLength != expectedLength {
		t.Fatalf("actualLength = %q, want %q", actualLength, expectedLength)
	}
	actualContains := ContainsString(results, fmt.Sprintf("machine: %s, service: %s, user: %s, password type: %s, password: %s", expectedMachine, expectedService, expectedUser, expectedType, expectedPassword))
	expectedContains := true
	if actualContains != expectedContains {
		t.Fatalf("actualContains = %v, want %v", actualContains, expectedContains)
	}
}

func TestPwgenUpdate(t *testing.T) {
	db, err := CreateDatabaseForTesting()
	defer db.Close()
	if err != nil {
		t.Fatalf("CreateDatabaseForTesting() err = %q, want nil", err)
	}
	OldOpenDatabase := OpenDatabase
	OpenDatabase = OpenDatabaseForTesting(db)
	defer func() { OpenDatabase = OldOpenDatabase }()
	OldCloseDatabase := CloseDatabase
	CloseDatabase = CloseDatabaseForTesting
	defer func() { CloseDatabase = OldCloseDatabase }()
	OldCommand := Command
	Command = CommandForTesting(t)
	defer func() { Command = OldCommand }()
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "output-from-pwgen"
	expectedType := "plain"
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, "oldpassword", expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "update", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := "Generated new password: output-from-pwgen\n"
	if buf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", buf.String(), expectedBuf)
	}
	results, err := readPasswords(db, "", "", "", "", false, false, []string{})
	if err != nil {
		t.Fatalf("readPasswords() err = %q, want nil", err)
	}
	actualLength := len(results)
	expectedLength := 1
	if actualLength != expectedLength {
		t.Fatalf("actualLength = %q, want %q", actualLength, expectedLength)
	}
	actualContains := ContainsString(results, fmt.Sprintf("machine: %s, service: %s, user: %s, password type: %s, password: %s", expectedMachine, expectedService, expectedUser, expectedType, expectedPassword))
	expectedContains := true
	if actualContains != expectedContains {
		t.Fatalf("actualContains = %v, want %v", actualContains, expectedContains)
	}
}

func TestDelete(t *testing.T) {
	db, err := CreateDatabaseForTesting()
	defer db.Close()
	if err != nil {
		t.Fatalf("CreateDatabaseForTesting() err = %q, want nil", err)
	}
	OldOpenDatabase := OpenDatabase
	OpenDatabase = OpenDatabaseForTesting(db)
	defer func() { OpenDatabase = OldOpenDatabase }()
	OldCloseDatabase := CloseDatabase
	CloseDatabase = CloseDatabaseForTesting
	defer func() { CloseDatabase = OldCloseDatabase }()
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	expectedType := "plain"
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "delete", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	results, err := readPasswords(db, "", "", "", "", false, false, []string{})
	if err != nil {
		t.Fatalf("readPasswords() err = %q, want nil", err)
	}
	actualLength := len(results)
	expectedLength := 0
	if actualLength != expectedLength {
		t.Fatalf("actualLength = %q, want %q", actualLength, expectedLength)
	}
}

func copyPath(inPath, outPath string) error {
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
			err := copyPath(encryptedQaPath, decryptedPath)
			if err != nil {
				t.Fatalf("copyPath() failed: %s", err)
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
			err := copyPath(decryptedPath, encryptedQaPath)
			if err != nil {
				t.Fatalf("copyPath() failed: %s", err)
			}
			return exec.Command("true")
		} else if len(arg) == 2 && name == "gunzip" && arg[0] == "--force" {
			compressedPath := arg[1]
			uncompressedPath := strings.ReplaceAll(compressedPath, ".gz", "")
			err := copyPath(compressedPath, uncompressedPath)
			if err != nil {
				t.Fatalf("copyPath() failed: %s", err)
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

func TestImport(t *testing.T) {
	db, err := CreateDatabaseForTesting()
	defer db.Close()
	if err != nil {
		t.Fatalf("CreateDatabaseForTesting() err = %q, want nil", err)
	}
	OldOpenDatabase := OpenDatabase
	OpenDatabase = OpenDatabaseForTesting(db)
	defer func() { OpenDatabase = OldOpenDatabase }()
	OldCloseDatabase := CloseDatabase
	CloseDatabase = CloseDatabaseForTesting
	defer func() { CloseDatabase = OldCloseDatabase }()
	OldCommand := Command
	Command = CommandForTesting(t)
	defer func() { Command = OldCommand }()
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	expectedType := "plain"
	os.Args = []string{"", "import"}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	results, err := readPasswords(db, "", "", "", "", false, false, []string{})
	if err != nil {
		t.Fatalf("readPasswords() err = %q, want nil", err)
	}
	actualLength := len(results)
	expectedLength := 2
	if actualLength != expectedLength {
		t.Fatalf("actualLength = %q, want %q", actualLength, expectedLength)
	}
	actualContains := ContainsString(results, fmt.Sprintf("machine: %s, service: %s, user: %s, password type: %s, password: %s", expectedMachine, expectedService, expectedUser, expectedType, expectedPassword))
	expectedContains := true
	if actualContains != expectedContains {
		t.Fatalf("actualContains = %v, want %v", actualContains, expectedContains)
	}
	expectedPassword = "totppassword"
	expectedType = "TOTP shared secret"
	actualContains = ContainsString(results, fmt.Sprintf("machine: %s, service: %s, user: %s, password type: %s, password: %s", expectedMachine, expectedService, expectedUser, expectedType, expectedPassword))
	if actualContains != expectedContains {
		t.Fatalf("actualContains = %v, want %v", actualContains, expectedContains)
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
