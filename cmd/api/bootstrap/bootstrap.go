package bootstrap

import (
	"context"
	"github.com/alexperezortuno/go-url-shortner/internal/config"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/server"
)

func Run() error {
	// Load configuration
	cfg := config.LoadConfig()
	ctx := context.Background()

	// Start server
	c, srv := server.New(ctx, cfg)
	return srv.Run(c)
}
