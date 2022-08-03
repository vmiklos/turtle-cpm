package commands

import (
	"bytes"
	"os"
	"testing"
)

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
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
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

func TestQuietSelect(t *testing.T) {
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
	os.Args = []string{"", "search", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-q"}
	buf := new(bytes.Buffer)

	actualRet := Main(buf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "mypassword\n"
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
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
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
	_, err = createPassword(db, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(db, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
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
	_, err = createPassword(db, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(db, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
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
	_, err = createPassword(db, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(db, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
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
	_, err = createPassword(db, "mymachine", "myservice", "myuser", "mypassword", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(db, "mymachine", "myservice", "myuser", "mypassword", "totp")
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
	_, err = createPassword(db, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(db, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
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
	os.Remove("qa/passwords.db")
}
