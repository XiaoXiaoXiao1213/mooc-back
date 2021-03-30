package common

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func Client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func LSet(key string, values ...string) {

	for _, value := range values {
		err := Client().RPush(ctx, key, value).Err()
		if err != nil {
			panic(err)
		}
	}

}
//func LGet(key string, values ...string) {
//
//	for _, value := range values {
//		err := Client().R(ctx, key, value).Err()
//		if err != nil {
//			panic(err)
//		}
//	}
//
//}
