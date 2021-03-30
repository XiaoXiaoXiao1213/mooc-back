package common

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func Client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func Set(key string, value string) error {
	err := Client().Set(ctx, key, value, time.Hour*24*15).Err()
	if err != nil {
		return err
	}
	return nil

}
func Get(key string) (string, error) {
	res := Client().Get(ctx, key)
	return res.Result()
}

func Delete(key string) error {
	err := Client().Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
