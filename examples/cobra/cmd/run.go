package cmd

import (
	"github.com/spf13/cobra"

	"github.com/zbiljic/zuki-go/examples/cobra/internal/app"
)

func newRunCommand(appName string, opts *app.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the application",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(cmd.Context(), appName, *opts)
		},
	}

	cmd.Flags().String("addr", ":8080", "HTTP listen address")
	_ = opts.Viper.BindPFlag("http.addr", cmd.Flags().Lookup("addr"))

	return cmd
}
