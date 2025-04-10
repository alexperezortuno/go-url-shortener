package middleware

import (
	"fmt"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/errors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"runtime/debug"
	"time"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[PANIC] %s panic recovered:\n%s | Stack: %s\n",
					time.Now().Format("2006/01/02 - 15:04:05"), err, debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError, errors.NewCustomError(errors.InternalServerError))
			}
		}()
		c.Next()
	}
}

func TracingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tracer := otel.Tracer("gin-server")
		spanCtx, span := tracer.Start(ctx.Request.Context(), ctx.FullPath())
		defer span.End()

		// Añadir atributos comunes
		span.SetAttributes(
			attribute.String("http.method", ctx.Request.Method),
			attribute.String("http.route", ctx.FullPath()),
			attribute.String("http.url", ctx.Request.URL.String()),
		)

		// Pasar el contexto con el span a la request
		ctx.Request = ctx.Request.WithContext(spanCtx)

		// Continuar con los demás handlers
		ctx.Next()

		// Registrar el status code después de que se complete el request
		span.SetAttributes(attribute.Int("http.status_code", ctx.Writer.Status()))
	}
}
