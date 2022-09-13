package commands

import (
	"bytes"
	"io"
	"os"
	"testing"

	"rsc.io/qr"
)

func TestSelect(t *testing.T) {
	ctx := CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	var expectedType PasswordType = "plain"
	err := initDatabase(ctx.Database)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
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
	ctx := CreateContextForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	var expectedType PasswordType = "plain"
	err := initDatabase(ctx.Database)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
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
	ctx := CreateContextForTesting(t)
	UseCommandForTesting(t)
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "totppassword"
	var expectedType PasswordType = "totp"
	err := initDatabase(ctx.Database)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
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

func GenerateQrCodeForTesting(text string, l qr.Level, w io.Writer) {
	w.Write([]byte("qrcode-output"))
}

func TestQrcodeSelect(t *testing.T) {
	ctx := CreateContextForTesting(t)
	OldGenerateQrCode := GenerateQrCode
	GenerateQrCode = GenerateQrCodeForTesting
	defer func() { GenerateQrCode = OldGenerateQrCode }()
	expectedMachine := "mymachine"
	expectedService := "myservice"
	expectedUser := "myuser"
	expectedPassword := "mypassword"
	var expectedType PasswordType = "totp"
	err := initDatabase(ctx.Database)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, expectedMachine, expectedService, expectedUser, expectedPassword, expectedType)
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	os.Args = []string{"", "search", "-m", expectedMachine, "-s", expectedService, "-u", expectedUser, "--qrcode"}
	inBuf := new(bytes.Buffer)
	outBuf := new(bytes.Buffer)

	actualRet := Main(inBuf, outBuf)

	expectedRet := 0
	if actualRet != expectedRet {
		t.Fatalf("Main() = %q, want %q", actualRet, expectedRet)
	}
	expectedOutput := "machine: mymachine, service: myservice, user: myuser, password type: TOTP shared secret, password:\n"
	expectedOutput += "qrcode-output\n"
	actualOutput := outBuf.String()
	if actualOutput != expectedOutput {
		t.Fatalf("actualOutput = %q, want %q", actualOutput, expectedOutput)
	}
}

func TestSelectMachineFilter(t *testing.T) {
	ctx := CreateContextForTesting(t)
	err := initDatabase(ctx.Database)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
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
	ctx := CreateContextForTesting(t)
	err := initDatabase(ctx.Database)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
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
	ctx := CreateContextForTesting(t)
	err := initDatabase(ctx.Database)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
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
	ctx := CreateContextForTesting(t)
	err := initDatabase(ctx.Database)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine", "myservice", "myuser", "mypassword", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine", "myservice", "myuser", "mypassword", "totp")
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
	ctx := CreateContextForTesting(t)
	err := initDatabase(ctx.Database)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
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
	ctx := CreateContextForTesting(t)
	err := initDatabase(ctx.Database)
	if err != nil {
		t.Fatalf("initDatabase() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine1", "myservice1", "myuser1", "mypassword1", "plain")
	if err != nil {
		t.Fatalf("createPassword() = %q, want nil", err)
	}
	_, err = createPassword(&ctx, "mymachine2", "myservice2", "myuser2", "mypassword2", "plain")
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
	UseCommandForTesting(t)
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
	os.Remove("fixtures/passwords.db")

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
	os.Remove("fixtures/passwords.db")
}

func TestParsePassword(t *testing.T) {
	s := "otpauth://totp/Myserver:myuser?secret=mysecret&digits=6&algorithm=SHA1&issuer=Myserver&period=30"
	expected := "mysecret"

	actual, err := parsePassword(s)
	if err != nil {
		t.Fatalf("err = %q, want nil", err)
	}

	if actual != expected {
		t.Fatalf("actual = %q, want %q", actual, expected)
	}
}

func TestParsePasswordBadURL(t *testing.T) {
	s := "otpauth://totp/Myserver:myuser?digits=6&algorithm=SHA1&issuer=Myserver&period=30"

	_, err := parsePassword(s)
	if err == nil {
		t.Fatalf("err = nil, want !nil")
	}
}
