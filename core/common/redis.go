package common

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func Client() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
	//err := rdb.HSet(ctx, "key", "value", 0).Err()
	//if err != nil {
	//	panic(err)
	//}
	//
	//val, err := rdb.Get(ctx, "key").Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key", val)
	//
	//val2, err := rdb.Get(ctx, "key2").Result()
	//if err == redis.Nil {
	//	fmt.Println("key2 does not exist")
	//} else if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("key2", val2)
	//}
	//// Output: key value
	//// key2 does not exist
}

func Hset(key string, values ...string) {

	for _, value := range values {
		err := Client().RPush(ctx, key, value).Err()
		if err != nil {
			panic(err)
		}
	}

}