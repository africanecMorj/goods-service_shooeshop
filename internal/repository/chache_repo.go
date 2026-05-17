package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ChacheRepo struct{
	RDB *redis.Client
}

func (r *ChacheRepo)StoreCode(ctx context.Context, email, code string) error {
    key := "reset:code:" + email

    return r.RDB.Set(ctx, key, code, 10*time.Minute).Err()
}

func (r *ChacheRepo)GetCode(ctx context.Context, email string) (string, error) {
    key := "reset:code:" + email

	code, err := r.RDB.Get(ctx, key).Result()
	if err != nil {
    	return "", err
	}


    return code, nil
}