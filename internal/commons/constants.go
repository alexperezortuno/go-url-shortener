package commons

import "time"

const (
	HealthPath    = "health"
	UrlPath       = "url"
	ShortenerPath = "r"
)

var (
	AllowMethods     = []string{"GET", "POST"}                             // Métodos permitidos
	AllowHeaders     = []string{"Origin", "Content-Type", "Authorization"} // Headers permitidos
	ExposeHeaders    = []string{"Content-Length"}                          // Headers expuestos
	AllowCredentials = true                                                // Permitir credenciales
	MaxAge           = 12 * time.Hour                                      // Tiempo de cacheo de preflight
)
