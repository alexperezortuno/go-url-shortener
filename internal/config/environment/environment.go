package environment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type ServerValues struct {
	Protocol        string
	Host            string
	Port            int
	ShutdownTimeout time.Duration
	Context         string
	TimeZone        string
	RedisHost       string
	RedisPass       string
	RedisDb         int
	Release         string
	Cors            []string
}

func Server() ServerValues {
	return ServerValues{
		Protocol: GetEnvStr("APP_PROTOCOL", "http"),
		Host:     GetEnvStr("APP_HOST", "localhost"),
		Context: func() string {
			ctx := GetEnvStr("APP_CONTEXT", "")
			if ctx == "" {
				return ""
			}
			return fmt.Sprintf("/%s", ctx)
		}(),
		Port:            GetEnvInt("APP_PORT", 8080),
		TimeZone:        GetEnvStr("APP_TIME_ZONE", "UTC"),
		ShutdownTimeout: 10 * time.Second,
		RedisHost:       GetEnvStr("REDIS_HOST", "localhost:6379"),
		RedisPass:       GetEnvStr("REDIS_PASSWORD", ""),
		RedisDb:         GetEnvInt("REDIS_DB", 0),
		Release:         GetEnvStr("RELEASE", "prod"),
		Cors:            GetEnvStrArray("CORS", []string{"*"}),
	}
}

func (c *ServerValues) SetGinMode() {
	switch c.Release {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	default:
		log.Fatalf("Invalid environment: %s", c.Release)
	}
}
