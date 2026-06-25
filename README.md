# zuki-go

`zuki-go` is a small starter layer for Go services built with
[`github.com/go-toho/toho`](https://github.com/go-toho/toho) and Uber Fx.

It keeps the common service bootstrap in one place:

- app info
- Toho Fx core
- config wiring
- `slog` logging
- start, wait, and graceful stop
- background error reporting with `zuki.ErrorSink`

## Install

```sh
go get github.com/zbiljic/zuki-go
```

## Quick Start

```go
package main

import (
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
		Name:   "my-service",
		Config: &cfg,
		FxOptions: []fx.Option{
			configfx.ProvideConfig[HTTPConfig](),
			myServiceModule,
		},
	})
	if err != nil {
		panic(err)
	}
}
```

By default, `zuki` uses the config value you pass in. Use `ConfigModeLoaded`
when a Toho config loader, such as `go-toho/contrib/config/vipero`, should fill
the config.

## Examples

```sh
go run ./examples/http
curl http://localhost:8080/healthz

go run ./examples/log-config
go run ./examples/fx-injection
```

`examples/cobra` is a separate module:

```sh
cd examples/cobra
go run . run
```

## Development

This runs the root module and standalone example modules.

```sh
make test-all
```
