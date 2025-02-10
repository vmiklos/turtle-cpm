// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"bytes"
	"os"
	"testing"
)

func TestGc(t *testing.T) {
	CreateContextForTesting(t)
	os.Args = []string{"", "gc"}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	// Just make sure we don't fail / don't report an error.
	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedBuf := ""
	if outBuf.String() != expectedBuf {
		t.Fatalf("Main() output is %q, want %q", outBuf.String(), expectedBuf)
	}
}
