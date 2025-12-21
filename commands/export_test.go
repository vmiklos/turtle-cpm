// Copyright 2025 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestExport(t *testing.T) {
	ctx := CreateContextForTesting(t)
	_, err := ctx.Database.Exec(`insert into passwords (machine, service, user, password, type) values('mymachine', 'myservice', 'myuser', 'mypassword', 'plain');`)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "export"}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	actualOutput := outBuf.String()
	var passwords []passwordRow
	err = json.NewDecoder(bytes.NewBufferString(actualOutput)).Decode(&passwords)
	if err != nil {
		t.Fatalf("json.Decode() = %q, want nil", err)
	}
	if len(passwords) != 1 {
		t.Fatalf("passwords len = %q, want %q", len(passwords), 1)
	}
	password := passwords[0]
	if password.ID != 1 {
		t.Fatalf("password.ID = %q, want %q", password.ID, 1)
	}
	if password.Machine != "mymachine" {
		t.Fatalf("password.Machine = %q, want %q", password.Machine, "mymachine")
	}
	if password.Service != "myservice" {
		t.Fatalf("password.Service = %q, want %q", password.Service, "myservice")
	}
	if password.User != "myuser" {
		t.Fatalf("password.User = %q, want %q", password.User, "myuser")
	}
	if password.Password != "mypassword" {
		t.Fatalf("password.Password = %q, want %q", password.Password, "mypassword")
	}
	if password.PasswordType != "plain" {
		t.Fatalf("password.PasswordType = %q, want %q", password.PasswordType, "plain")
	}
	if password.Archived != false {
		t.Fatalf("password.Archived = %t, want %t", password.Archived, false)
	}
	if password.Created != "" {
		t.Fatalf("password.Created = %q, want %q", password.Created, "")
	}
	if password.Modified != "" {
		t.Fatalf("password.Modified = %q, want %q", password.Modified, "")
	}
}
