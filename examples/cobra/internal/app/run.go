package app

import (
	"context"

	"github.com/go-toho/contrib/config/vipero/viperofx"
	"github.com/spf13/viper"
	"github.com/zbiljic/zuki-go"
	"go.uber.org/fx"

	"github.com/zbiljic/zuki-go/examples/cobra/internal/buildinfo"
	"github.com/zbiljic/zuki-go/examples/cobra/internal/server"
)

type Options struct {
	ConfigFiles []string
	Viper       *viper.Viper
}

func Run(ctx context.Context, appName string, opts Options) error {
	cfg := Config{}

	return zuki.Run(zuki.Options[Config]{
		Name:        appName,
		Version:     buildinfo.Version,
		Context:     ctx,
		Config:      &cfg,
		ConfigMode:  zuki.ConfigModeLoaded,
		ConfigFiles: opts.ConfigFiles,
		FxOptions: []fx.Option{
			viperofx.SupplyViper(opts.Viper),
			server.Module,
		},
	})
}
