package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	_ "github.com/go-toho/contrib/config/vipero/viperofx"
	"github.com/go-toho/toho/config/configfx"
	"github.com/zbiljic/zuki-go"
	"go.uber.org/fx"
)

type Config struct {
	HTTP HTTPConfig `json:"http"`
}

type HTTPConfig struct {
	Addr string `default:":8080" json:"addr"`
}

func main() {
	cfg := Config{}
	err := zuki.Run(zuki.Options[Config]{
		Name:       "zuki",
		Version:    "dev",
		Config:     &cfg,
		ConfigMode: zuki.ConfigModeLoaded,
		FxOptions: []fx.Option{
			configfx.ProvideConfig[HTTPConfig](),
			fx.Provide(newServer),
			fx.Invoke(registerServer),
		},
	})
	if err != nil {
		panic(err)
	}
}

func newServer(config HTTPConfig) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	return &http.Server{
		Addr:    config.Addr,
		Handler: mux,
	}
}

func registerServer(lifecycle fx.Lifecycle, server *http.Server, log *slog.Logger, sink zuki.ErrorSink) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info("starting HTTP server", slog.String("addr", server.Addr))
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					sink.Report(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})
}
