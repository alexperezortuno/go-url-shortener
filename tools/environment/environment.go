package environment

import (
    "github.com/gin-gonic/gin"
    "os"
    "strconv"
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
}

func env() {
    env := os.Getenv("APP_ENV")
    
    if env == "" || env == "release" {
        gin.SetMode(gin.ReleaseMode)
    }
}

func Server() ServerValues {
    env()
    port, err := strconv.Atoi(os.Getenv("APP_PORT"))
    if err != nil {
        port = 8107
    }
    
    protocol := os.Getenv("APP_PROTOCOL")
    host := os.Getenv("APP_HOST")
    timeZone := os.Getenv("APP_TIME_ZONE")
    context := os.Getenv("APP_CONTEXT")
    redisHost := os.Getenv("REDIS_HOST")
    redisPass := os.Getenv("REDIS_PASS")
    redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
    
    if err != nil {
        redisDb = 0
    }
    
    if protocol == "" {
        protocol = "http"
    }
    
    if host == "" {
        host = "localhost"
    }
    
    if context == "" {
        context = "api"
    }
    
    if timeZone == "" {
        timeZone = "America/Santiago"
    }
    
    if redisHost == "" {
        redisHost = "localhost:6379"
    }
    
    return ServerValues{
        Protocol:        protocol,
        Host:            host,
        Context:         context,
        Port:            port,
        TimeZone:        timeZone,
        ShutdownTimeout: 10 * time.Second,
        RedisHost:       redisHost,
        RedisPass:       redisPass,
        RedisDb:         redisDb,
    }
}
