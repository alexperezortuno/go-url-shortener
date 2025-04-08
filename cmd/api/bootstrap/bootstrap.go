package bootstrap

import (
	"context"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/server"
	"github.com/alexperezortuno/go-url-shortner/tools/environment"
)

var params = environment.Server()

func Run() error {
	ctx, srv := server.New(context.Background(), params)
	return srv.Run(ctx)
}
