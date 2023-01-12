package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)
	os.Args = []string{"", "version"}

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main(version) = %v, want %v, output is %q", actualRet, expectedRet, outBuf.String())
	}
	expectedPrefix := "turtle-cpm "
	actualOutput := outBuf.String()
	if !strings.HasPrefix(actualOutput, expectedPrefix) {
		t.Fatalf("actualOutput = %q, want prefix %q", actualOutput, expectedPrefix)
	}
}
