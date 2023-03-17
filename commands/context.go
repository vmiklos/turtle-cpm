// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package commands

import (
	"os"
	"os/exec"

	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp/totp"
	"github.com/sethvargo/go-password/password"
)

// Command returns the Cmd struct to execute the named program
var Command = exec.Command

// Remove removes the named file or (empty) directory.
var Remove = os.Remove

// Stat returns a FileInfo describing the named file.
var Stat = os.Stat

// GeneratePassword is the package shortcut for password.Generator.Generate.
var GeneratePassword = password.Generate

// GenerateQrCode creates a QR Code and writes it out to io.Writer.
var GenerateQrCode = qrterminal.Generate

// GenerateTotpCode creates a TOTP token using the current time.
var GenerateTotpCode = totp.GenerateCode

// OpenDatabase opens the database before running a subcommand.
var OpenDatabase = openDatabase

// CloseDatabase opens the database before running a subcommand.
var CloseDatabase = closeDatabase
