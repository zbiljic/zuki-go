package zuki

import (
	"context"
	"errors"
	"net/url"
	"testing"

	tohoapp "github.com/go-toho/toho/app"
	"github.com/go-toho/toho/config/configfx"
	"github.com/go-toho/toho/logger"
	"go.uber.org/fx"
)

type testConfig struct {
	HTTP testHTTPConfig
}

type testHTTPConfig struct {
	Addr string
}

type testConfigWithLogger struct {
	Log  logger.Config
	HTTP testHTTPConfig
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
