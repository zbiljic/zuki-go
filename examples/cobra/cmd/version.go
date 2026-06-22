package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zbiljic/zuki-go/examples/cobra/internal/buildinfo"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), buildinfo.Version)
			return err
		},
	}
}
