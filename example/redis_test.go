package example

import (
	"github.com/go-redis/redis"
	"log"
	"testing"
	"time"
)

var redisClient *redis.Client

func TestRedisTtlResult(t *testing.T) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "foobared",
		DB:       0,
	})

	setRes, err := redisClient.Set("foo", "bar", 0).Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(setRes == "OK")

	ttlRes, err := redisClient.TTL("foo").Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(ttlRes.Seconds() == -1)

	ttlRes, err = redisClient.TTL("foo1").Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(ttlRes.Seconds() == -2)

	expRes, err := redisClient.Expire("foo", 30*24*60*60*time.Second).Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(expRes == true)

	ttlRes, err = redisClient.TTL("foo").Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(ttlRes.Seconds() == 30*24*60*60)
}
