package middleware

import (
	"fmt"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/errors"
	"github.com/gin-gonic/gin"
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
