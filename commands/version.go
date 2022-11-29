package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCommand(ctx *Context) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "shows version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "turtle-cpm version %s\n", Version)
			ctx.NoWriteBack = true
			return nil
		},
	}

	return cmd
}
