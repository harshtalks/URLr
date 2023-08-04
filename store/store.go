package store

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	Ctx context.Context
}

// Sharing context among functions
var Context = Store{Ctx: context.Background()}

// Handler for creating a connection to the remote Redis Server
func CreateStore() *redis.Client {
	redisDb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DATABASE_URL"),
		Password: os.Getenv("PASSWORD"),
	})

	return redisDb
}
