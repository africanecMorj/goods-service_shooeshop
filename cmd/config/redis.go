package config

import (
	"os"
	"log"
	"context"
	
	"github.com/redis/go-redis/v9"
)

func RedisInit () *redis.Client{
	opt, err := redis.ParseURL("redis://default:kQeqoLlvJQCVcwuQmQeOzGaIIhxDlKul@redis.railway.internal:6379")
	log.Println(os.Getenv("REDIS_URL"))
	log.Println(os.Getenv("JWT-SECRET"))
	if err != nil {
		log.Fatal("Redis init error:",err)
	}

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Redis ping error:", err)
	}

	log.Println("Redis connected")

	return client
}

