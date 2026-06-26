package zuki

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/go-toho/contrib/config/vipero/viperofx"
	tohoapp "github.com/go-toho/toho/app"
	"github.com/go-toho/toho/config/configfx"
	"github.com/go-toho/toho/logger"
	"github.com/go-toho/toho/tohofx"
	"go.uber.org/fx"
)

const registeredViperModule = "github.com/spf13/viper"

var registeredViperCreator tohofx.Creator

type testConfig struct {
	HTTP testHTTPConfig
}

type testHTTPConfig struct {
	Addr string
}

type testLoadedConfig struct {
	HTTP testLoadedHTTPConfig `json:"http"`
}

type testLoadedHTTPConfig struct {
	Addr string `default:"127.0.0.1:9090" json:"addr"`
}

type testConfigWithLogger struct {
	Log  logger.Config
	HTTP testHTTPConfig
}

func TestMain(m *testing.M) {
	// viperofx registers itself in Toho's global module registry from init().
	// Keep that package-wide side effect out of tests unless they opt into it.
	registeredViperCreator = tohofx.Options[registeredViperModule]
	delete(tohofx.Options, registeredViperModule)

	os.Exit(m.Run())
}

func useRegisteredViperModule(t *testing.T) {
	t.Helper()

	if registeredViperCreator == nil {
		t.Fatal("registered viper module is not available")
	}

	previous, hadPrevious := tohofx.Options[registeredViperModule]
	tohofx.Options[registeredViperModule] = registeredViperCreator
	t.Cleanup(func() {
		if hadPrevious {
			tohofx.Options[registeredViperModule] = previous
			return
		}
		delete(tohofx.Options, registeredViperModule)
	})
}

func TestRuntimeStartsWithSuppliedConfig(t *testing.T) {
	cfg := testConfig{
		HTTP: testHTTPConfig{Addr: "127.0.0.1:0"},
	}

	var got testHTTPConfig
	rt := New(Options[testConfig]{
		Name:   "zuki-test",
		Config: &cfg,
		FxOptions: []fx.Option{
			configfx.ProvideConfig[testHTTPConfig](),
			fx.Invoke(func(config testHTTPConfig) {
				got = config
			}),
		},
	})

	if err := rt.App.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.App.Stop(); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})

	if got.Addr != cfg.HTTP.Addr {
		t.Fatalf("got HTTP addr %q, want %q", got.Addr, cfg.HTTP.Addr)
	}
	if rt.Config.HTTP.Addr != cfg.HTTP.Addr {
		t.Fatalf("got runtime config HTTP addr %q, want %q", rt.Config.HTTP.Addr, cfg.HTTP.Addr)
	}
	if rt.App.Config().HTTP.Addr != cfg.HTTP.Addr {
		t.Fatalf("got app config HTTP addr %q, want %q", rt.App.Config().HTTP.Addr, cfg.HTTP.Addr)
	}
}

func TestRuntimeStartsWithLoadedConfigDefaults(t *testing.T) {
	useRegisteredViperModule(t)

	cfg := testLoadedConfig{}
	var got testLoadedHTTPConfig
	rt := New(Options[testLoadedConfig]{
		Name:       "zuki-test",
		Config:     &cfg,
		ConfigMode: ConfigModeLoaded,
		FxOptions: []fx.Option{
			configfx.ProvideConfig[testLoadedHTTPConfig](),
			fx.Invoke(func(config testLoadedHTTPConfig) {
				got = config
			}),
		},
	})

	if err := rt.App.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.App.Stop(); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})

	if got.Addr != "127.0.0.1:9090" {
		t.Fatalf("got HTTP addr %q, want 127.0.0.1:9090", got.Addr)
	}
	if cfg.HTTP.Addr != "127.0.0.1:9090" {
		t.Fatalf("got caller config HTTP addr %q, want 127.0.0.1:9090", cfg.HTTP.Addr)
	}
	if rt.Config.HTTP.Addr != "127.0.0.1:9090" {
		t.Fatalf("got runtime config HTTP addr %q, want 127.0.0.1:9090", rt.Config.HTTP.Addr)
	}
	if rt.App.Config().HTTP.Addr != "127.0.0.1:9090" {
		t.Fatalf("got app config HTTP addr %q, want 127.0.0.1:9090", rt.App.Config().HTTP.Addr)
	}
}

func TestRuntimeStartsWithLoadedConfigFile(t *testing.T) {
	useRegisteredViperModule(t)

	configFile := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configFile, []byte("http:\n  addr: 127.0.0.1:8181\n"), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg := testLoadedConfig{}
	var got testLoadedHTTPConfig
	rt := New(Options[testLoadedConfig]{
		Name:        "zuki-test",
		Config:      &cfg,
		ConfigMode:  ConfigModeLoaded,
		ConfigFiles: []string{configFile},
		FxOptions: []fx.Option{
			configfx.ProvideConfig[testLoadedHTTPConfig](),
			fx.Invoke(func(config testLoadedHTTPConfig) {
				got = config
			}),
		},
	})

	if err := rt.App.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.App.Stop(); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})

	if got.Addr != "127.0.0.1:8181" {
		t.Fatalf("got HTTP addr %q, want 127.0.0.1:8181", got.Addr)
	}
	if cfg.HTTP.Addr != "127.0.0.1:8181" {
		t.Fatalf("got caller config HTTP addr %q, want 127.0.0.1:8181", cfg.HTTP.Addr)
	}
	if rt.Config.HTTP.Addr != "127.0.0.1:8181" {
		t.Fatalf("got runtime config HTTP addr %q, want 127.0.0.1:8181", rt.Config.HTTP.Addr)
	}
	if rt.App.Config().HTTP.Addr != "127.0.0.1:8181" {
		t.Fatalf("got app config HTTP addr %q, want 127.0.0.1:8181", rt.App.Config().HTTP.Addr)
	}
}

func TestRuntimeStartsWithConfigLoggerField(t *testing.T) {
	cfg := testConfigWithLogger{
		Log: logger.Config{
			Level:  "info",
			Format: logger.TextFormat,
		},
		HTTP: testHTTPConfig{Addr: "127.0.0.1:0"},
	}

	rt := New(Options[testConfigWithLogger]{
		Name:   "zuki-test",
		Config: &cfg,
		FxOptions: []fx.Option{
			configfx.ProvideConfig[testHTTPConfig](),
			fx.Invoke(func(testHTTPConfig) {}),
		},
	})

	if err := rt.App.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.App.Stop(); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})
}

func TestRuntimeAppInfoIncludesMetadataAndEndpoints(t *testing.T) {
	endpoint, err := url.Parse("http://127.0.0.1:8080")
	if err != nil {
		t.Fatalf("parse endpoint: %v", err)
	}

	rt := New(Options[struct{}]{
		ID:        "node-1",
		Name:      "zuki-test",
		Version:   "v1.2.3",
		Metadata:  map[string]string{"region": "test"},
		Endpoints: []*url.URL{endpoint},
	})

	info := rt.App.AppInfo()
	if info.ID() != "node-1" {
		t.Fatalf("got id %q, want node-1", info.ID())
	}
	if info.Name() != "zuki-test" {
		t.Fatalf("got name %q, want zuki-test", info.Name())
	}
	if info.Version() != "v1.2.3" {
		t.Fatalf("got version %q, want v1.2.3", info.Version())
	}
	if info.Metadata()["region"] != "test" {
		t.Fatalf("got metadata %#v, want region=test", info.Metadata())
	}
	if len(info.Endpoint()) != 1 || info.Endpoint()[0] != "http://127.0.0.1:8080" {
		t.Fatalf("got endpoints %#v, want http://127.0.0.1:8080", info.Endpoint())
	}
}

func TestRuntimeAppOptionsCanOverrideFields(t *testing.T) {
	rt := New(Options[struct{}]{
		Name: "zuki-test",
		AppOptions: []tohoapp.Option{
			tohoapp.Name("override-name"),
		},
	})

	if got := rt.App.AppInfo().Name(); got != "override-name" {
		t.Fatalf("got name %q, want override-name", got)
	}
}

func TestRunReturnsLifecycleError(t *testing.T) {
	expected := errors.New("boom")

	err := Run(Options[struct{}]{
		Name: "zuki-test",
		FxOptions: []fx.Option{
			fx.Invoke(func(lifecycle fx.Lifecycle, sink ErrorSink) {
				lifecycle.Append(fx.Hook{
					OnStart: func(context.Context) error {
						sink.Report(expected)
						return nil
					},
				})
			}),
		},
	})

	if !errors.Is(err, expected) {
		t.Fatalf("got error %v, want %v", err, expected)
	}
}

func TestRuntimeExposesReportedErrors(t *testing.T) {
	expected := errors.New("boom")

	rt := New(Options[struct{}]{
		Name: "zuki-test",
		FxOptions: []fx.Option{
			fx.Invoke(func(lifecycle fx.Lifecycle, sink ErrorSink) {
				lifecycle.Append(fx.Hook{
					OnStart: func(context.Context) error {
						sink.Report(expected)
						return nil
					},
				})
			}),
		},
	})

	if err := rt.App.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		if err := rt.App.Stop(); err != nil {
			t.Fatalf("stop: %v", err)
		}
	})

	select {
	case err := <-rt.Errors:
		if !errors.Is(err, expected) {
			t.Fatalf("got error %v, want %v", err, expected)
		}
	default:
		t.Fatal("expected reported error")
	}
}
