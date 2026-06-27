package cmd

import (
	"time"

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

	cmd.Flags().Duration("read-header-timeout", 2*time.Second, "maximum time to read request headers")
	_ = opts.Viper.BindPFlag("http.read_header_timeout", cmd.Flags().Lookup("read-header-timeout"))

	cmd.Flags().Duration("read-timeout", 10*time.Second, "maximum time to read the entire request")
	_ = opts.Viper.BindPFlag("http.read_timeout", cmd.Flags().Lookup("read-timeout"))

	cmd.Flags().Duration("write-timeout", 10*time.Second, "maximum time before timing out response writes")
	_ = opts.Viper.BindPFlag("http.write_timeout", cmd.Flags().Lookup("write-timeout"))

	cmd.Flags().Duration("idle-timeout", 120*time.Second, "maximum time to wait for the next request on keep-alive connections")
	_ = opts.Viper.BindPFlag("http.idle_timeout", cmd.Flags().Lookup("idle-timeout"))

	return cmd
}
