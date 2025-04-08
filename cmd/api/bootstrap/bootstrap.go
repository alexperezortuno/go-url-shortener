package bootstrap

import (
	"context"
	"github.com/alexperezortuno/go-url-shortner/internal/config/environment"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/server"
)

var params = environment.Server()

func Run() error {
	ctx, srv := server.New(context.Background(), params)
	return srv.Run(ctx)
}
