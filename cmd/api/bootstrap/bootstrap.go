package bootstrap

import (
    "context"
    "go-url-shortner/internal/platform/server"
    "go-url-shortner/tools/environment"
)

var params = environment.Server()

func Run() error {
    ctx, srv := server.New(context.Background(), params)
    return srv.Run(ctx)
}
