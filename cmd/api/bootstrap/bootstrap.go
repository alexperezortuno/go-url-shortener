package bootstrap

import (
	"context"
	"github.com/alexperezortuno/go-url-shortner/internal/config"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/server"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/tracing"
	"log"
)

func Run() error {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize tracing
	ctx := context.Background()
	tp, err := tracing.InitTracer(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Start server
	c, srv := server.New(ctx, cfg)
	return srv.Run(c)
}
