package cmd

import (
	"github.com/go-toho/contrib/config/vipero"
	"github.com/spf13/cobra"

	"github.com/zbiljic/zuki-go/examples/cobra/internal/app"
)

const AppName = "zuki"

func New() *cobra.Command {
	opts := &app.Options{
		Viper: vipero.New(AppName),
	}

	cmd := &cobra.Command{
		Use:          "zuki-cli",
		Short:        "Zuki Cobra starter example",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringSliceVarP(
		&opts.ConfigFiles,
		"config",
		"c",
		nil,
		"config file to load",
	)

	cmd.AddCommand(
		newRunCommand(AppName, opts),
		newVersionCommand(),
	)

	return cmd
}
