package commons

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAllowMethodsContainsExpectedMethods(t *testing.T) {
	assert.Contains(t, AllowMethods, "GET")
	assert.Contains(t, AllowMethods, "POST")
	assert.NotContains(t, AllowMethods, "PUT")
}

func TestAllowHeadersContainsExpectedHeaders(t *testing.T) {
	assert.Contains(t, AllowHeaders, "Origin")
	assert.Contains(t, AllowHeaders, "Content-Type")
	assert.Contains(t, AllowHeaders, "Authorization")
	assert.NotContains(t, AllowHeaders, "X-Custom-Header")
}

func TestExposeHeadersContainsExpectedHeaders(t *testing.T) {
	assert.Contains(t, ExposeHeaders, "Content-Length")
	assert.NotContains(t, ExposeHeaders, "X-Expose-Header")
}

func TestAllowCredentialsIsTrue(t *testing.T) {
	assert.True(t, AllowCredentials)
}

func TestMaxAgeIs12Hours(t *testing.T) {
	assert.Equal(t, 12*time.Hour, MaxAge)
}
