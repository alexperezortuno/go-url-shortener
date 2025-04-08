package store

import (
	"context"
	"fmt"
	"github.com/alexperezortuno/go-url-shortner/internal/config/environment"
	"github.com/go-redis/redis/v9"
	"log"
	"time"
)

// StorageService is struct wrapper around raw Redis client
type StorageService struct {
	redisClient *redis.Client
}

// Top level declarations for the storeService and Redis context
var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

var params = environment.Server()

const CacheDuration = 6 * time.Hour

// InitializeStore is initializing the store service and return a store pointer
func InitializeStore() *StorageService {
	log.Println(fmt.Sprintf("Initializing Redis client... %s", params.RedisHost))
	rdb := redis.NewClient(&redis.Options{
		Addr:     params.RedisHost,
		Password: params.RedisPass, // no password set
		DB:       params.RedisDb,   // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Println(fmt.Sprintf("Error init Redis: %v", err))
	}

	log.Println(fmt.Sprintf("Redis started successfully: pong message = {%s}", pong))
	storeService.redisClient = rdb
	return storeService
}

func SaveURLInRedis(shortURL, originalURL string) {
	err := storeService.redisClient.Set(ctx, shortURL, originalURL, CacheDuration).Err()
	if err != nil {
		log.Println(fmt.Sprintf("Failed SaveURLInRedis | Error: %v - shortURL: %s - originalURL: %s\n",
			err, shortURL, originalURL))
	}
}

func RetrieveInitialURLFromRedis(shortURL string) string {
	result, err := storeService.redisClient.Get(ctx, shortURL).Result()
	if err != nil {
		log.Println(fmt.Sprintf("Failed RetrieveInitialURLFromRedis | Error: %v - shortURL: %s\n",
			err, shortURL))
	}
	return result
}
