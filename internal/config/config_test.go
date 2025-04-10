package config

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func LoadConfigReturnsDefaultValuesWhenEnvVarsAreNotSet(t *testing.T) {
	config := LoadConfig()

	assert.Equal(t, "http", config.Protocol)
	assert.Equal(t, "0.0.0.0", config.Host)
	assert.Equal(t, "", config.Context)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, "UTC", config.TimeZone)
	assert.Equal(t, 10*time.Second, config.ShutdownTimeout)
	assert.Equal(t, "localhost:6379", config.RedisHost)
	assert.Equal(t, "", config.RedisPass)
	assert.Equal(t, 0, config.RedisDb)
	assert.Equal(t, "prod", config.Release)
	assert.Equal(t, []string{"*"}, config.CorsAllowsOrigin)
	assert.Equal(t, "otel-collector:4317", config.OtelExporterEndpoint)
	assert.Equal(t, "go-url-shortener", config.ServiceName)
	assert.False(t, config.TracingEnabled)
}

func LoadConfigReturnsOverriddenValuesWhenEnvVarsAreSet(t *testing.T) {
	setEnv("APP_PROTOCOL", "https")
	setEnv("APP_HOST", "127.0.0.1")
	setEnv("APP_CONTEXT", "api")
	setEnv("APP_PORT", "9090")
	setEnv("APP_TIME_ZONE", "PST")
	setEnv("REDIS_HOST", "redis:6379")
	setEnv("REDIS_PASSWORD", "password")
	setEnv("REDIS_DB", "1")
	setEnv("RELEASE", "dev")
	setEnv("CORS_ALLOW_ORIGIN", "http://example.com,http://test.com")
	setEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "custom-otel:4317")
	setEnv("SERVICE_NAME", "custom-service")
	setEnv("TRACING_ENABLED", "true")

	config := LoadConfig()

	assert.Equal(t, "https", config.Protocol)
	assert.Equal(t, "127.0.0.1", config.Host)
	assert.Equal(t, "/api", config.Context)
	assert.Equal(t, 9090, config.Port)
	assert.Equal(t, "PST", config.TimeZone)
	assert.Equal(t, "redis:6379", config.RedisHost)
	assert.Equal(t, "password", config.RedisPass)
	assert.Equal(t, 1, config.RedisDb)
	assert.Equal(t, "dev", config.Release)
	assert.Equal(t, []string{"http://example.com", "http://test.com"}, config.CorsAllowsOrigin)
	assert.Equal(t, "custom-otel:4317", config.OtelExporterEndpoint)
	assert.Equal(t, "custom-service", config.ServiceName)
	assert.True(t, config.TracingEnabled)
}

func SetGinModeSetsCorrectModeBasedOnRelease(t *testing.T) {
	config := &Config{}

	config.Release = "dev"
	config.SetGinMode()
	assert.Equal(t, gin.DebugMode, gin.Mode())

	config.Release = "test"
	config.SetGinMode()
	assert.Equal(t, gin.TestMode, gin.Mode())

	config.Release = "prod"
	config.SetGinMode()
	assert.Equal(t, gin.ReleaseMode, gin.Mode())
}

func SetGinModePanicsForInvalidRelease(t *testing.T) {
	config := &Config{Release: "invalid"}
	assert.Panics(t, func() {
		config.SetGinMode()
	})
}

func GetEnvStrReturnsFallbackWhenEnvVarIsNotSet(t *testing.T) {
	unsetEnv("UNSET_ENV_VAR")
	result := GetEnvStr("UNSET_ENV_VAR", "fallback")
	assert.Equal(t, "fallback", result)
}

func GetEnvStrReturnsValueWhenEnvVarIsSet(t *testing.T) {
	setEnv("SET_ENV_VAR", "value")
	result := GetEnvStr("SET_ENV_VAR", "fallback")
	assert.Equal(t, "value", result)
}

func GetEnvIntReturnsFallbackWhenEnvVarIsNotSet(t *testing.T) {
	unsetEnv("UNSET_INT_ENV_VAR")
	result := GetEnvInt("UNSET_INT_ENV_VAR", 42)
	assert.Equal(t, 42, result)
}

func GetEnvIntReturnsValueWhenEnvVarIsSet(t *testing.T) {
	setEnv("SET_INT_ENV_VAR", "100")
	result := GetEnvInt("SET_INT_ENV_VAR", 42)
	assert.Equal(t, 100, result)
}

func GetEnvIntReturnsFallbackForInvalidValue(t *testing.T) {
	setEnv("INVALID_INT_ENV_VAR", "invalid")
	result := GetEnvInt("INVALID_INT_ENV_VAR", 42)
	assert.Equal(t, 42, result)
}

func GetEnvBoolReturnsFallbackWhenEnvVarIsNotSet(t *testing.T) {
	unsetEnv("UNSET_BOOL_ENV_VAR")
	result := GetEnvBool("UNSET_BOOL_ENV_VAR", true)
	assert.True(t, result)
}

func GetEnvBoolReturnsValueWhenEnvVarIsSet(t *testing.T) {
	setEnv("SET_BOOL_ENV_VAR", "false")
	result := GetEnvBool("SET_BOOL_ENV_VAR", true)
	assert.False(t, result)
}

func GetEnvBoolReturnsFallbackForInvalidValue(t *testing.T) {
	setEnv("INVALID_BOOL_ENV_VAR", "invalid")
	result := GetEnvBool("INVALID_BOOL_ENV_VAR", true)
	assert.True(t, result)
}

func GetEnvStrArrayReturnsFallbackWhenEnvVarIsNotSet(t *testing.T) {
	unsetEnv("UNSET_ARRAY_ENV_VAR")
	result := GetEnvStrArray("UNSET_ARRAY_ENV_VAR", []string{"default"})
	assert.Equal(t, []string{"default"}, result)
}

func GetEnvStrArrayReturnsSplitValuesWhenEnvVarIsSet(t *testing.T) {
	setEnv("SET_ARRAY_ENV_VAR", "value1, value2, value3")
	result := GetEnvStrArray("SET_ARRAY_ENV_VAR", []string{"default"})
	assert.Equal(t, []string{"value1", "value2", "value3"}, result)
}

func setEnv(key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		panic(err)
	}
}

func unsetEnv(key string) {
	err := os.Unsetenv(key)
	if err != nil {
		panic(err)
	}
}
