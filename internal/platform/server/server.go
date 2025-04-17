package server

import (
	"context"
	"fmt"
	"github.com/alexperezortuno/go-url-shortner/internal/commons"
	"github.com/alexperezortuno/go-url-shortner/internal/config"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/middleware"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/server/handler/health"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/server/handler/shortner"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/services/metrics"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/services/tracing"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/storage/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	httpAddr        string
	engine          *gin.Engine
	shutdownTimeout time.Duration
}

func New(ctx context.Context, cfg *config.Config) (context.Context, Server) {
	cfg.SetGinMode()

	if cfg.TracingEnabled {
		// Initialize tracing
		shutdownTracer, err := tracing.InitTracer(ctx, cfg)
		if err != nil {
			log.Fatalf("failed to initialize tracer: %v", err)
		}
		defer func() {
			if err := shutdownTracer(ctx); err != nil {
				log.Printf("failed to shutdown tracer: %v", err)
			}
		}()
	}

	if cfg.MetricsEnabled {
		// Initialize metrics
		metrics.InitializeMetricsService(cfg.MetricsAddress, cfg.MetricsTags)
		metricsCheck := metrics.GetInstance()
		if !metricsCheck.HealthCheck() {
			log.Fatal("Metrics service is not available")
		}
	}

	srv := Server{
		engine:          gin.New(),
		httpAddr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		shutdownTimeout: cfg.ShutdownTimeout,
	}
	store.InitializeStore(cfg)

	log.Printf("Check app in %s:%d%s/%s", cfg.Host, cfg.Port, cfg.Context, "health")
	srv.registerRoutes(cfg)
	return serverContext(ctx), srv
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server running on", s.httpAddr)
	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func (s *Server) registerRoutes(cfg *config.Config) {
	ctx := cfg.Context
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CorsAllowsOrigin,     // Dominios permitidos
		AllowMethods:     commons.AllowMethods,     // Métodos permitidos
		AllowHeaders:     commons.AllowHeaders,     // Headers permitidos
		ExposeHeaders:    commons.ExposeHeaders,    // Headers expuestos
		AllowCredentials: commons.AllowCredentials, // Permitir credenciales
		MaxAge:           commons.MaxAge,           // Tiempo de cacheo de preflight
	}))

	// Middlewares
	s.engine.Use(gin.Logger())
	s.engine.Use(middleware.Logging())
	s.engine.Use(middleware.Recovery())
	if cfg.TracingEnabled {
		s.engine.Use(middleware.TracingMiddleware())
	}

	// Routes
	s.engine.GET(fmt.Sprintf("%s/%s", ctx, commons.HealthPath), health.CheckHandler())
	s.engine.POST(fmt.Sprintf("%s/%s", ctx, commons.UrlPath), shortner.CreateShortURL(cfg))
	s.engine.GET(fmt.Sprintf("%s/%s", ctx, commons.UrlPath), shortner.ReturnLongURL())
	s.engine.GET(fmt.Sprintf("%s/%s/:s", ctx, commons.ShortenerPath), shortner.RedirectURL())
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}
