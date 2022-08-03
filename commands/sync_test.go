package commands

import (
	"bytes"
	"os"
	"testing"
)

func TestSync(t *testing.T) {
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
	expectedUser := "myuser"
	os.Args = []string{"", "create", "-m", expectedMachine, "-u", expectedUser}
	buf := new(bytes.Buffer)
	os.Remove("qa/passwords.db")
	actualRet := Main(buf)
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main(create) = %v, want %v, output is %q", actualRet, expectedRet, buf.String())
	}
	os.Rename("qa/passwords.db", "qa/remote.db")
	os.Args = []string{"", "sync"}
	buf = new(bytes.Buffer)

	actualRet = Main(buf)

	expectedRet = 0
	if actualRet != expectedRet {
		t.Fatalf("Main(sync) = %v, want %v, output is %q", actualRet, expectedRet, buf.String())
	}
	os.Args = []string{"", "search", "-m", expectedMachine, "-u", expectedUser}
	buf = new(bytes.Buffer)
	actualRet = Main(buf)
	expectedRet = 0
	if actualRet != expectedRet {
		t.Fatalf("Main(search) = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine, service: http, user: myuser, password type: plain, password: output-from-pwgen\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
	os.Remove("qa/remote.db")
	os.Remove("qa/passwords.db")
}
