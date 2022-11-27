package commands

import (
	"bytes"
	"os"
	"testing"
)

func TestSync(t *testing.T) {
	UseCommandForTesting(t)
	OldRemove := Remove
	Remove = RemoveForTesting
	defer func() { Remove = OldRemove }()
	OldStat := Stat
	Stat = StatForTesting
	defer func() { Stat = OldStat }()
	expectedMachine := "mymachine"
	expectedUser := "myuser"
	os.Args = []string{"", "create", "-m", expectedMachine, "-u", expectedUser}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	os.Remove("fixtures/passwords.db")
	actualRet := Main(inBuf, outBuf)
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main(create) = %v, want %v, output is %q", actualRet, expectedRet, outBuf.String())
	}
	os.Rename("fixtures/passwords.db", "fixtures/remote.db")
	os.Args = []string{"", "sync"}
	outBuf = new(bytes.Buffer)

	actualRet = Main(inBuf, outBuf)

	expectedRet = 0
	if actualRet != expectedRet {
		t.Fatalf("Main(sync) = %v, want %v, output is %q", actualRet, expectedRet, outBuf.String())
	}
	os.Args = []string{"", "search", "--noid", "-m", expectedMachine, "-u", expectedUser}
	outBuf = new(bytes.Buffer)
	actualRet = Main(inBuf, outBuf)
	expectedRet = 0
	if actualRet != expectedRet {
		t.Fatalf("Main(search) = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine, service: http, user: myuser, password type: plain, password: output-from-pwgen\n"
	actualOutput := outBuf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
	os.Remove("fixtures/remote.db")
	os.Remove("fixtures/passwords.db")
}
