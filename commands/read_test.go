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
	var expectedType PasswordType = "plain"
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine, service: myservice, user: myuser, password type: plain, password: mypassword\n"
	actualOutput := outBuf.String()
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
	var expectedType PasswordType = "plain"
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-q"}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "mypassword\n"
	actualOutput := outBuf.String()
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
	var expectedType PasswordType = "totp"
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "--totp", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine, service: myservice, user: myuser, password type: TOTP code, password: output-from-oathtool\n"
	actualOutput := outBuf.String()
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
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	// mymachine1 is found, mymachine2 is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine1, service: myservice1, user: myuser1, password type: plain, password: mypassword1\n"
	actualOutput := outBuf.String()
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
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	// myservice1 is found, myservice2 is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine1, service: myservice1, user: myuser1, password type: plain, password: mypassword1\n"
	actualOutput := outBuf.String()
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
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	// myuser1 is found, myuser2 is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine1, service: myservice1, user: myuser1, password type: plain, password: mypassword1\n"
	actualOutput := outBuf.String()
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
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	// totp is found, plain is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine, service: myservice, user: myuser, password type: TOTP shared secret, password: mypassword\n"
	actualOutput := outBuf.String()
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
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	// myservice1 is found, myservice2 is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine1, service: myservice1, user: myuser1, password type: plain, password: mypassword1\n"
	actualOutput := outBuf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
}

func TestSelectInteractive(t *testing.T) {
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
	// Interactive search.
	os.Args = []string{""}
	inBuf := new(bytes.Buffer)
	inBuf.Write([]byte("mymachine1" + "\n"))
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	// myservice1 is found, myservice2 is not found.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "Search term: "
	expectedOutput += "machine: mymachine1, service: myservice1, user: myuser1, password type: plain, password: mypassword1\n"
	actualOutput := outBuf.String()
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
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	os.Remove("qa/passwords.db")

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main(create) = %v, want %v, output is %q", actualRet, expectedRet, outBuf.String())
	}

	os.Args = []string{"", "search", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	inBuf = new(bytes.Buffer)
	outBuf = new(bytes.Buffer)

	actualRet = Main(inBuf, outBuf)

	expectedRet = 0
	if actualRet != expectedRet {
		t.Fatalf("Main(search) = %q, want %q", actualRet, expectedRet)
	}

	expectedOutput := "machine: mymachine, service: myservice, user: myuser, password type: plain, password: mypassword\n"
	actualOutput := outBuf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
	os.Remove("qa/passwords.db")
}
