package commands

import (
	"fmt"

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

			err = runCommand("scp", "cpm:"+databasePath, databasePath)
			if err != nil {
				return fmt.Errorf("runCommand() failed: %s", err)
			}

			ctx.NoWriteBack = true
			return nil
		},
	}

	return cmd
}
