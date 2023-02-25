// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestInsert(t *testing.T) {
	ctx := CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	expectedType := "plain"
	os.Args = []string{"", "create", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword, "-t", "plain"}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := "Created 1 password\n"
	if outBuf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", outBuf.String(), expectedBuf)
	}
	opts := searchOptions{}
	opts.noid = true
	results, err := readPasswords(ctx.Database, opts)
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

func TestNoServiceInsert(t *testing.T) {
	ctx := CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	expectedType := "plain"
	os.Args = []string{"", "create", "-m", expectedMachine, "-u", expectedUser, "-p", expectedPassword}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := "Created 1 password\n"
	if outBuf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", outBuf.String(), expectedBuf)
	}
	opts := searchOptions{}
	opts.noid = true
	results, err := readPasswords(ctx.Database, opts)
	if err != nil {
		t.Fatalf("readPasswords() err = %q, want nil", err)
	}
	actualLength := len(results)
	expectedLength := 1
	if actualLength != expectedLength {
		t.Fatalf("actualLength = %q, want %q", actualLength, expectedLength)
	}
	actualContains := ContainsString(results, fmt.Sprintf("machine: %s, service: http, user: %s, password type: %s, password: %s", expectedMachine, expectedUser, expectedType, expectedPassword))
	expectedContains := true
	if actualContains != expectedContains {
		t.Fatalf("actualContains = %v, want %v", actualContains, expectedContains)
	}
}

func GeneratePasswordForTesting(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error) {
	return "output-from-pwgen", nil
}

func TestPwgenInsert(t *testing.T) {
	ctx := CreateContextForTesting(t)
	UseCommandForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "output-from-pwgen"
	expectedType := "plain"
	os.Args = []string{"", "create", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := "Created 1 password\nGenerated password: output-from-pwgen\n"
	if outBuf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", outBuf.String(), expectedBuf)
	}
	opts := searchOptions{}
	opts.noid = true
	results, err := readPasswords(ctx.Database, opts)
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
	CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	os.Args = []string{"", "create", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	// First run succeeds.
	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}

	// Second run fails as the machine/service/user already has a password.
	actualRet = Main(inBuf, outBuf)

	expectedRet = 1
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedPrefix := "Error: createPassword() failed: query.Exec() failed: UNIQUE constraint failed\n"
	actualOutput := outBuf.String()
	if strings.HasPrefix(actualOutput, expectedPrefix) {
		t.Fatalf("actualOutput = %q, want prefix %q", actualOutput, expectedPrefix)
	}
}

// Insert fails because -t mytype is not a valid type.
func TestInsertFailBadType(t *testing.T) {
	CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	os.Args = []string{"", "create", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword, "-t", "mytype"}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 1
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
}

func TestInteractiveInsert(t *testing.T) {
	ctx := CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	expectedType := "plain"
	os.Args = []string{"", "create", "-s", expectedService, "-p", expectedPassword, "-t", "plain"}
	inBuf := new(bytes.Buffer)
	inBuf.Write([]byte(expectedMachine + "\n" + expectedUser + "\n"))
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedBuf := "Machine: User: Created 1 password\n"
	if outBuf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", outBuf.String(), expectedBuf)
	}
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	opts := searchOptions{}
	opts.noid = true
	results, err := readPasswords(ctx.Database, opts)
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

func TestDryRunInsert(t *testing.T) {
	ctx := CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	os.Args = []string{"", "create", "-n", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "-p", expectedPassword, "-t", "plain"}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := "Would create 1 password\n"
	if outBuf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", outBuf.String(), expectedBuf)
	}
	results, err := readPasswords(ctx.Database, searchOptions{})
	if err != nil {
		t.Fatalf("readPasswords() err = %q, want nil", err)
	}
	actualLength := len(results)
	// This is a dry run, so not 1.
	expectedLength := 0
	if actualLength != expectedLength {
		t.Fatalf("actualLength = %q, want %q", actualLength, expectedLength)
	}
}
