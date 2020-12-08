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

	ttlRes := redisClient.TTL("foo").Val()
	log.Println(ttlRes.Seconds() == -1)

	ttlRes = redisClient.TTL("foo1").Val()
	log.Println(ttlRes.Seconds() == -2)

	expRes := redisClient.Expire("foo", 30*24*60*60*time.Second).Val()
	log.Println(expRes == true)

	ttlRes = redisClient.TTL("foo").Val()
	log.Println(ttlRes.Seconds() == 30*24*60*60)
}

func TestRedisExistsResult(t *testing.T) {
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

	extRes := redisClient.Exists("foo").Val()
	log.Println(extRes == 1)

	extRes = redisClient.Exists("foo1").Val()
	log.Println(extRes == 0)
}
