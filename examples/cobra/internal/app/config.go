package app

import (
	"github.com/go-toho/toho/logger"

	"github.com/zbiljic/zuki-go/examples/cobra/internal/server"
)

type Config struct {
	Log  logger.Config `json:"log"`
	HTTP server.Config `json:"http"`
}
