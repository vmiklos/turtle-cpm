package commands

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/spf13/cobra"
)

// XMLPassword is the 4th <node> element from cpm's XML database.
type XMLPassword struct {
	XMLName xml.Name `xml:"node"`
	Label   string   `xml:"label,attr"`
	Totp    string   `xml:"totp,attr"`
}

// XMLUser is the 3rd <node> element from cpm's XML database.
type XMLUser struct {
	XMLName   xml.Name      `xml:"node"`
	Label     string        `xml:"label,attr"`
	Passwords []XMLPassword `xml:"node"`
}

// XMLService is the 2nd <node> element from cpm's XML database.
type XMLService struct {
	XMLName xml.Name  `xml:"node"`
	Label   string    `xml:"label,attr"`
	Users   []XMLUser `xml:"node"`
}

// XMLMachine is the 1st <node> element from cpm's XML database.
type XMLMachine struct {
	XMLName  xml.Name     `xml:"node"`
	Label    string       `xml:"label,attr"`
	Services []XMLService `xml:"node"`
}

// XMLMachines is the <root> element from cpm's XML database.
type XMLMachines struct {
	XMLName  xml.Name     `xml:"root"`
	Machines []XMLMachine `xml:"node"`
}

func newImportCommand(ctx *Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "import",
		Short: "imports an old XML database",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Decrypt and uncompress ~/.cpmdb to a temp file.
			usr, err := user.Current()
			if err != nil {
				return fmt.Errorf("user.Current() failed: %s", err)
			}

			encryptedPath := usr.HomeDir + "/.cpmdb"
			decryptedFile, err := ioutil.TempFile("", "cpm")
			if err != nil {
				return fmt.Errorf("ioutil.TempFile() failed: %s", err)
			}

			decryptedPath := decryptedFile.Name()
			defer Remove(decryptedPath)

			err = runCommand("gpg", "--decrypt", "-a", "-o", decryptedPath+".gz", encryptedPath)
			if err != nil {
				return fmt.Errorf("runCommand() failed: %s", err)
			}

			err = runCommand("gunzip", "--force", decryptedPath+".gz")
			if err != nil {
				return fmt.Errorf("runCommand() failed: %s", err)
			}

			// Parse the XML.
			xmlFile, err := os.Open(decryptedPath)
			if err != nil {
				return fmt.Errorf("os.Open(decryptedPath) failed: %s", err)
			}
			defer xmlFile.Close()

			xmlBytes, err := ioutil.ReadAll(xmlFile)
			if err != nil {
				return fmt.Errorf("ioutil.ReadAll(xmlFile) failed: %s", err)
			}

			// Avoid 'encoding "ISO-8859-1" declared but Decoder.CharsetReader is nil'.
			xmlBytes = bytes.ReplaceAll(xmlBytes, []byte(`encoding="ISO-8859-1"`), []byte(`encoding="UTF-8"`))

			var machines XMLMachines
			err = xml.Unmarshal(xmlBytes, &machines)
			if err != nil {
				return fmt.Errorf("xml.Unmarshal() failed: %s", err)
			}

			// Import the parsed data.
			for _, machine := range machines.Machines {
				machineLabel := machine.Label
				for _, service := range machine.Services {
					serviceLabel := service.Label
					for _, user := range service.Users {
						userLabel := user.Label
						for _, password := range user.Passwords {
							passwordLabel := password.Label
							var passwordType PasswordType
							if password.Totp == "true" {
								passwordType = "totp"
							} else {
								passwordType = "plain"
							}

							_, err = createPassword(ctx.Database, machineLabel, serviceLabel, userLabel, passwordLabel, passwordType)
							if err != nil {
								return fmt.Errorf("createPassword(machine='%s', service='%s', user='%s', type='%s') failed: %s", machineLabel, serviceLabel, userLabel, passwordType, err)
							}
						}
					}
				}
			}

			return nil
		},
	}

	return cmd
}
