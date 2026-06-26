package zuki

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/go-toho/toho"
	tohoapp "github.com/go-toho/toho/app"
	"github.com/go-toho/toho/config/configfx"
	"github.com/go-toho/toho/contrib/log/slogo"
	"github.com/go-toho/toho/contrib/log/slogo/slogofx"
	"github.com/go-toho/toho/tohofx"
	"go.uber.org/fx"
)

// ConfigMode controls how the runtime makes the application config available.
type ConfigMode uint8

const (
	// ConfigModeSupplied supplies Options.Config as the resolved config value.
	ConfigModeSupplied ConfigMode = iota

	// ConfigModeLoaded lets a registered Toho config loader resolve Options.Config.
	ConfigModeLoaded
)

// ErrorSink reports background lifecycle errors back to the runtime.
type ErrorSink chan<- error

// Report sends a non-nil error to the runtime error sink.
func (s ErrorSink) Report(err error) {
	if err != nil && s != nil {
		s <- err
	}
}

// Options describes a Toho application started by this starter runtime.
type Options[C any] struct {
	ID        string
	Name      string
	Version   string
	Metadata  map[string]string
	Endpoints []*url.URL

	Context context.Context
	Config  *C

	ConfigMode  ConfigMode
	ConfigFiles []string

	Logger *slog.Logger

	// TohoOptions are applied before generated Toho wiring. Use the
	// dedicated Options fields for app info, context, config, logger, and Fx.
	TohoOptions []toho.Option

	AppOptions []tohoapp.Option
	FxOptions  []fx.Option
}

// Runtime holds a configured Toho application and its lifecycle error channel.
type Runtime[C any] struct {
	App    *toho.TohoApp[C, *slog.Logger]
	Config *C
	Errors <-chan error
}

// New builds a Toho application with the starter defaults.
func New[C any](opts Options[C]) *Runtime[C] {
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	cfg := opts.Config
	if cfg == nil {
		cfg = new(C)
	}

	errCh := make(chan error, 1)

	fxopts := starterFxOptions(opts, cfg, errCh)

	tohoopts := append([]toho.Option{}, opts.TohoOptions...)
	tohoopts = append(tohoopts,
		toho.AppInfo(appInfoOptions(opts)...),
		toho.AppCore(tohofx.NewCore()),
		toho.Context(ctx),
		toho.Config(cfg),
		toho.Options(fx.Options(fxopts...)),
	)
	if opts.Logger != nil {
		tohoopts = append(tohoopts, toho.Logger(opts.Logger))
	}

	return &Runtime[C]{
		App:    toho.NewC[C](tohoopts...),
		Config: cfg,
		Errors: errCh,
	}
}

// Run starts a configured Toho application, waits for shutdown, and stops it.
func Run[C any](opts Options[C]) error {
	return New(opts).Run()
}

// Run starts the runtime application, waits for a signal or lifecycle error,
// and then stops the application.
func (r *Runtime[C]) Run() error {
	if err := r.App.Start(); err != nil {
		return err
	}

	log := r.App.Logger()
	log.Info("started")

	var runErr error
	select {
	case err := <-r.App.Wait():
		log.Debug("exit", slogo.Err(err))
	case err, ok := <-r.Errors:
		if ok && err != nil {
			runErr = err
			log.Error("lifecycle", slogo.Err(err))
		}
	}

	log.Info("shutting down")

	if err := r.App.Stop(); err != nil {
		log.Error("shutdown", slogo.Err(err))
		if runErr != nil {
			return fmt.Errorf("%w; shutdown: %v", runErr, err)
		}
		return err
	}

	return runErr
}

func starterFxOptions[C any](opts Options[C], cfg *C, errCh chan error) []fx.Option {
	fxopts := []fx.Option{
		slogofx.FxPrinterLogger,
		slogofx.FxEventLogger,
		fx.Provide(func() ErrorSink { return ErrorSink(errCh) }),
	}

	if opts.ConfigMode == ConfigModeSupplied {
		fxopts = append(fxopts, configfx.SupplyConfig(cfg))
	}
	if len(opts.ConfigFiles) > 0 {
		fxopts = append(fxopts, configfx.SupplyConfigFiles(opts.ConfigFiles))
	}

	return append(fxopts, opts.FxOptions...)
}

func appInfoOptions[C any](opts Options[C]) []tohoapp.Option {
	id := opts.ID
	if id == "" {
		id, _ = os.Hostname()
	}

	appopts := make([]tohoapp.Option, 0, 5+len(opts.AppOptions))
	if id != "" {
		appopts = append(appopts, tohoapp.ID(id))
	}
	if opts.Name != "" {
		appopts = append(appopts, tohoapp.Name(opts.Name))
	}
	if opts.Version != "" {
		appopts = append(appopts, tohoapp.Version(opts.Version))
	}
	if len(opts.Metadata) > 0 {
		appopts = append(appopts, tohoapp.Metadata(opts.Metadata))
	}
	if len(opts.Endpoints) > 0 {
		appopts = append(appopts, tohoapp.Endpoint(opts.Endpoints...))
	}
	appopts = append(appopts, opts.AppOptions...)
	return appopts
}
