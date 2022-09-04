package commands

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestUpdate(t *testing.T) {
	db, err := CreateDatabaseForTesting()
	defer db.Close()
	if err != nil {
		t.Fatalf("CreateDatabaseForTesting() err = %q, want nil", err)
	}
	UseDatabaseForTesting(t, db)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "newpassword"
	var expectedType PasswordType = "plain"
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, "oldpassword", expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "update", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := "Updated 1 password\n"
	if outBuf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", outBuf.String(), expectedBuf)
	}
	results, err := readPasswords(db, searchOptions{})
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
	UseDatabaseForTesting(t, db)
	OldCommand := Command
	Command = CommandForTesting(t)
	defer func() { Command = OldCommand }()
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "output-from-pwgen"
	var expectedType PasswordType = "plain"
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, "oldpassword", expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "update", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := "Generated new password: output-from-pwgen\nUpdated 1 password\n"
	if outBuf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", outBuf.String(), expectedBuf)
	}
	results, err := readPasswords(db, searchOptions{})
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

func TestInteractiveUpdate(t *testing.T) {
	db, err := CreateDatabaseForTesting()
	defer db.Close()
	if err != nil {
		t.Fatalf("CreateDatabaseForTesting() err = %q, want nil", err)
	}
	UseDatabaseForTesting(t, db)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "newpassword"
	var expectedType PasswordType = "plain"
	err = initDatabase(db)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(db, expectedMachine, expectedService, expectedUser, "oldpassword", expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "update", "-s", expectedService, "-p", expectedPassword}
	inBuf := new(bytes.Buffer)
	inBuf.Write([]byte(expectedMachine + "\n" + expectedUser + "\n"))
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := "Machine: User: "
	expectedBuf += "Updated 1 password\n"
	if outBuf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", outBuf.String(), expectedBuf)
	}
	results, err := readPasswords(db, searchOptions{})
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
