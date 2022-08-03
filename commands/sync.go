package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newSyncCommand(ctx *Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "sync",
		Short: "copies a remote database to a local one",
		RunE: func(cmd *cobra.Command, args []string) error {
			databasePath, err := getDatabasePath()
			if err != nil {
				return fmt.Errorf("getDatabasePath() failed: %s", err)
			}

			command := Command("scp", "cpm:"+databasePath, databasePath)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			err = command.Run()
			if err != nil {
				return fmt.Errorf("command.Run() failed: %s", err)
			}

			ctx.NoWriteBack = true
			return nil
		},
	}

	return cmd
}
