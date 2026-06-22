package main

import (
	"context"
	"log/slog"

	_ "github.com/go-toho/contrib/config/vipero/viperofx"
	"github.com/go-toho/toho/config/configfx"
	"github.com/go-toho/toho/logger"
	"github.com/go-toho/toho/logger/loggerfx"
	"github.com/zbiljic/zuki-go"
	"go.uber.org/fx"
)

type Config struct {
	Log     logger.Config `json:"log"`
	Message MessageConfig `json:"message"`
}

type MessageConfig struct {
	Text string `default:"hello from Fx injection" json:"text"`
}

type Greeter struct {
	log    *slog.Logger
	config MessageConfig
}

func main() {
	cfg := Config{}

	err := zuki.Run(zuki.Options[Config]{
		Name:       "zuki",
		Version:    "dev",
		Config:     &cfg,
		ConfigMode: zuki.ConfigModeLoaded,
		FxOptions: []fx.Option{
			loggerfx.SupplyFxSetupConfig(logger.DebugTextConfig),
			configfx.ProvideConfig[MessageConfig](),
			fx.Provide(newGreeter),
			fx.Invoke(registerGreeter),
		},
	})
	if err != nil {
		panic(err)
	}
}

func newGreeter(log *slog.Logger, config MessageConfig) *Greeter {
	return &Greeter{
		log:    log,
		config: config,
	}
}

func registerGreeter(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, greeter *Greeter) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				greeter.log.Info("greeter started", slog.String("message", greeter.config.Text))
				_ = shutdowner.Shutdown()
			}()
			return nil
		},
	})
}
