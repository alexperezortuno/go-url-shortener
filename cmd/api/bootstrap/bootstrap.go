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
	shutdownTracer, err := tracing.InitTracer(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := shutdownTracer(ctx); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		}
	}()

	// Start server
	c, srv := server.New(ctx, cfg)
	return srv.Run(c)
}
