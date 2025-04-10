package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Protocol             string
	Host                 string
	Port                 int
	ShutdownTimeout      time.Duration
	Context              string
	TimeZone             string
	RedisHost            string
	RedisPass            string
	RedisDb              int
	Release              string
	CorsAllowsOrigin     []string
	OtelExporterEndpoint string
	ServiceName          string
}

func LoadConfig() *Config {
	return &Config{
		Protocol: GetEnvStr("APP_PROTOCOL", "http"),
		Host:     GetEnvStr("APP_HOST", "0.0.0.0"),
		Context: func() string {
			ctx := GetEnvStr("APP_CONTEXT", "")
			if ctx == "" {
				return ""
			}
			return fmt.Sprintf("/%s", ctx)
		}(),
		Port:                 GetEnvInt("APP_PORT", 8080),
		TimeZone:             GetEnvStr("APP_TIME_ZONE", "UTC"),
		ShutdownTimeout:      10 * time.Second,
		RedisHost:            GetEnvStr("REDIS_HOST", "localhost:6379"),
		RedisPass:            GetEnvStr("REDIS_PASSWORD", ""),
		RedisDb:              GetEnvInt("REDIS_DB", 0),
		Release:              GetEnvStr("RELEASE", "prod"),
		CorsAllowsOrigin:     GetEnvStrArray("CORS_ALLOW_ORIGIN", []string{"*"}),
		OtelExporterEndpoint: GetEnvStr("OTEL_EXPORTER_OTLP_ENDPOINT", "otel-collector:4317"),
		ServiceName:          GetEnvStr("SERVICE_NAME", "go-url-shortener"),
	}
}

func (c *Config) SetGinMode() {
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

func GetEnvStr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func GetEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}

func GetEnvStrArray(key string, fallback []string) []string {
	if value, ok := os.LookupEnv(key); ok {
		return splitString(value)
	}
	return fallback
}

func splitString(s string) []string {
	var result []string
	for _, str := range strings.Split(s, ",") {
		result = append(result, strings.TrimSpace(str))
	}
	return result
}
