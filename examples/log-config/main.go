package main

import (
	"context"
	"log/slog"

	_ "github.com/go-toho/contrib/config/vipero/viperofx"
	"github.com/go-toho/toho/logger"
	"github.com/zbiljic/zuki-go"
	"go.uber.org/fx"
)

type Config struct {
	Log logger.Config `json:"log"`
}

func main() {
	cfg := Config{}

	err := zuki.Run(zuki.Options[Config]{
		Name:       "zuki",
		Version:    "dev",
		Config:     &cfg,
		ConfigMode: zuki.ConfigModeLoaded,
		FxOptions: []fx.Option{
			fx.Invoke(registerWorker),
		},
	})
	if err != nil {
		panic(err)
	}
}

func registerWorker(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, log *slog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Debug("debug logging enabled from Config.Log")
				log.Info("worker finished")
				_ = shutdowner.Shutdown()
			}()
			return nil
		},
	})
}
