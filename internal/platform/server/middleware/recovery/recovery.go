package recovery

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

func Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Recover from panic
        defer func() {
            if err := recover(); err != nil {
                fmt.Printf("[Middleware] %s panic recovered:\n%s\n",
                    time.Now().Format("2006/01/02 - 15:04:05"), err)
                
                c.Abort()
                c.AbortWithStatus(http.StatusInternalServerError)
            }
        }()
        
        // Process request
        c.Next()
    }
}
