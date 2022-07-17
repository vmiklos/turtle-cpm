package main

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
func OpenDatabaseForTesting(sqlDb *sql.DB) func() (*CpmDatabase, error) {
	return func() (*CpmDatabase, error) {
		var db CpmDatabase
		db.Database = sqlDb

		err := initDatabase(db.Database)
		if err != nil {
			return nil, fmt.Errorf("initDatabase() failed: %s", err)
		}

		return &db, nil
	}
}

// CloseDatabaseForTesting implements CloseDatabase and does nothing.
func CloseDatabaseForTesting(db *CpmDatabase) error {
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

func TestInsert(t *testing.T) {
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
	os.Args = []string{"", "create", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	results, err := readPasswords(db, "", "", "", "", false, []string{})
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

// Insert fails because the password is already inserted.
func TestInsertFail(t *testing.T) {
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
	os.Args = []string{"", "create", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword}
	buf := new(bytes.Buffer)

	// First run succeeds.
	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}

	// Second run fails as the machine/service/user already has a password.
	actualRet = Main(buf)

	expectedRet = 1
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
}

func TestSelect(t *testing.T) {
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
	err = createPassword(db, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine, service: myservice, user: myuser, password type: plain, password: mypassword\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
}

func TestSelectTotpCode(t *testing.T) {
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
	expectedPassword := "totppassword"
	expectedType := "totp"
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	err = createPassword(db, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "--totp", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine, service: myservice, user: myuser, password type: TOTP code, password: output-from-oathtool\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
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
	err = createPassword(db, expectedMachine, expectedService, expectedUser, "oldpassword", expectedType)
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
	results, err := readPasswords(db, "", "", "", "", false, []string{})
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
	err = createPassword(db, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
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
	results, err := readPasswords(db, "", "", "", "", false, []string{})
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
		} else if len(arg) == 1 && name == "gunzip" {
			compressedPath := arg[0]
			uncompressedPath := strings.ReplaceAll(compressedPath, ".gz", "")
			err := copyPath(compressedPath, uncompressedPath)
			if err != nil {
				t.Fatalf("copyPath() failed: %s", err)
			}
			return exec.Command("true")
		} else if len(arg) == 3 && name == "oathtool" && arg[0] == "-b" && arg[1] == "--totp" && arg[2] == "totppassword" {
			return exec.Command("echo", "output-from-oathtool")
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
	results, err := readPasswords(db, "", "", "", "", false, []string{})
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

func TestSelectMachineFilter(t *testing.T) {
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
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	err = createPassword(db, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	err = createPassword(db, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "-m", "mymachine1"}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	// mymachine1 is found, mymachine2 is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine1, service: myservice1, user: myuser1, password type: plain, password: mypassword1\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
}

func TestSelectServiceFilter(t *testing.T) {
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
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	err = createPassword(db, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	err = createPassword(db, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "-s", "myservice1"}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	// myservice1 is found, myservice2 is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine1, service: myservice1, user: myuser1, password type: plain, password: mypassword1\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
}

func TestSelectUserFilter(t *testing.T) {
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
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	err = createPassword(db, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	err = createPassword(db, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "-u", "myuser1"}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	// myuser1 is found, myuser2 is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine1, service: myservice1, user: myuser1, password type: plain, password: mypassword1\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
}

func TestSelectTypeFilter(t *testing.T) {
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
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	err = createPassword(db, "mymachine", "myservice", "myuser", "mypassword", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	err = createPassword(db, "mymachine", "myservice", "myuser", "mypassword", "totp")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "-t", "totp"}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	// totp is found, plain is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine, service: myservice, user: myuser, password type: TOTP shared secret, password: mypassword\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
}

func TestSelectImplicitFilter(t *testing.T) {
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
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	err = createPassword(db, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	err = createPassword(db, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	// Implicit search, also not telling that myservice1 is a service.
	os.Args = []string{"", "myservice1"}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	// myservice1 is found, myservice2 is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine1, service: myservice1, user: myuser1, password type: plain, password: mypassword1\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
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

func TestOpenCloseDatabase(t *testing.T) {
	// Intentionally not mocking OpenDatabase and CloseDatabase in this test.
	OldCommand := Command
	Command = CommandForTesting(t)
	defer func() { Command = OldCommand }()
	OldRemove := Remove
	Remove = RemoveForTesting
	defer func() { Remove = OldRemove }()
	OldStat := Stat
	Stat = StatForTesting
	defer func() { Stat = OldStat }()
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	os.Args = []string{"", "create", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword}
	buf := new(bytes.Buffer)
	os.Remove("qa/passwords.db")

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main(create) = %v, want %v, output is %q", actualRet, expectedRet, buf.String())
	}

	os.Args = []string{"", "search", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	buf = new(bytes.Buffer)

	actualRet = Main(buf)

	expectedRet = 0
	if actualRet != expectedRet {
		t.Fatalf("Main(search) = %q, want %q", actualRet, expectedRet)
	}

	expectedOutput := "machine: mymachine, service: myservice, user: myuser, password type: plain, password: mypassword\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
}
