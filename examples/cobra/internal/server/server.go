package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-toho/toho/config/configfx"
	"github.com/zbiljic/zuki-go"
	"go.uber.org/fx"
)

var Module = fx.Options(
	configfx.ProvideConfig[Config](),
	fx.Provide(New),
	fx.Invoke(Register),
)

type Config struct {
	Addr string `default:":8080" json:"addr"`
}

func New(config Config) *http.Server {
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

func Register(lifecycle fx.Lifecycle, server *http.Server, log *slog.Logger, sink zuki.ErrorSink) {
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
