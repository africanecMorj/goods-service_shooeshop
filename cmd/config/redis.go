package config

import (
	"os"
	"log"
	"context"
	
	"github.com/redis/go-redis/v9"
)

func RedisInit () *redis.Client{
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	log.Println(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("Redis init error:"err)
	}

	client := redis.NewClient(opt)
	if err != nil {
		log.Fatal("Redis init error:", err)
	}

	if err := client.Ping(context.Background()); err != nil {
		log.Fatal("Redis ping error:", err)
	}

	log.Println("Redis connected")

	return client
}

