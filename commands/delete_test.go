// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"bytes"
	"os"
	"testing"
)

func TestDelete(t *testing.T) {
	ctx := CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	var expectedType PasswordType = "plain"
	secure := false
	_, err := createPassword(&ctx, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType, secure)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "delete", "--id", "1"}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	results, err := readPasswords(ctx.Database, searchOptions{})
	if err != nil {
		t.Fatalf("readPasswords() err = %q, want nil", err)
	}
	actualLength := len(results)
	expectedLength := 0
	if actualLength != expectedLength {
		t.Fatalf("actualLength = %q, want %q", actualLength, expectedLength)
	}
}

func TestInteractiveDelete(t *testing.T) {
	ctx := CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	var expectedType PasswordType = "plain"
	secure := false
	_, err := createPassword(&ctx, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType, secure)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "delete"}
	inBuf := new(bytes.Buffer)
	inBuf.Write([]byte("1\n"))
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := "Id: "
	expectedBuf += "Deleted 1 password\n"
	if outBuf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", outBuf.String(), expectedBuf)
	}
	results, err := readPasswords(ctx.Database, searchOptions{})
	if err != nil {
		t.Fatalf("readPasswords() err = %q, want nil", err)
	}
	actualLength := len(results)
	expectedLength := 0
	if actualLength != expectedLength {
		t.Fatalf("actualLength = %q, want %q", actualLength, expectedLength)
	}
}

func TestDryRunDelete(t *testing.T) {
	ctx := CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	var expectedType PasswordType = "plain"
	secure := false
	_, err := createPassword(&ctx, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType, secure)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "delete", "-n", "-i", "1"}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	results, err := readPasswords(ctx.Database, searchOptions{})
	if err != nil {
		t.Fatalf("readPasswords() err = %q, want nil", err)
	}
	actualLength := len(results)
	// dry run, so not 0
	expectedLength := 1
	if actualLength != expectedLength {
		t.Fatalf("actualLength = %q, want %q", actualLength, expectedLength)
	}
}
