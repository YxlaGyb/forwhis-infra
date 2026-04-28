package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx    = context.Background()
	client *redis.Client
)

func init() {
	host := os.Getenv("REDIS_HOST")
	port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if host == "" {
		host = "127.0.0.1"
	}
	if port == 0 {
		port = 6379
	}

	client = redis.NewClient(&redis.Options{
		Addr: host + ":" + strconv.Itoa(port),
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("Redis Client Error: %v", err)
	} else {
		log.Printf("redis: %s:%d connection successful!", host, port)
	}
}

func Set(key, value string, ttl ...time.Duration) error {
	if value == "" {
		return nil
	}
	if len(ttl) > 0 {
		return client.Set(ctx, key, value, ttl[0]).Err()
	}
	return client.Set(ctx, key, value, 0).Err()
}

func SetEX(key, value string, seconds int) error {
	return Set(key, value, time.Duration(seconds)*time.Second)
}

func Get(key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func Del(keys ...string) error {
	return client.Del(ctx, keys...).Err()
}

func Exists(keys ...string) (int64, error) {
	return client.Exists(ctx, keys...).Result()
}

func Expire(ctx context.Context, key string, ttl time.Duration) error {
	return client.Expire(ctx, key, ttl).Err()
}

func TTL(key string) (time.Duration, error) {
	return client.TTL(ctx, key).Result()
}

func Keys(pattern string) ([]string, error) {
	return client.Keys(ctx, pattern).Result()
}

func Incr(key string) (int64, error) {
	return client.Incr(ctx, key).Result()
}

func Decr(key string) (int64, error) {
	return client.Decr(ctx, key).Result()
}

func HSet(key string, field string, value interface{}) error {
	return client.HSet(ctx, key, field, value).Err()
}

func HGet(key, field string) (string, error) {
	return client.HGet(ctx, key, field).Result()
}

func HGetAll(key string) (map[string]string, error) {
	return client.HGetAll(ctx, key).Result()
}

func HDel(key string, fields ...string) error {
	return client.HDel(ctx, key, fields...).Err()
}

func SAdd(ctx context.Context, key string, members ...interface{}) error {
	return client.SAdd(ctx, key, members...).Err()
}

func SMembers(key string) ([]string, error) {
	return client.SMembers(ctx, key).Result()
}

func SRem(key string, members ...interface{}) error {
	return client.SRem(ctx, key, members...).Err()
}

func ZAdd(key string, score float64, member string) error {
	return client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

func ZRange(key string, start, stop int64) ([]string, error) {
	return client.ZRange(ctx, key, start, stop).Result()
}

func ZRem(key string, members ...interface{}) error {
	return client.ZRem(ctx, key, members...).Err()
}

func Client() *redis.Client {
	return client
}

func Close() error {
	return client.Close()
}

func FormatKey(parts ...interface{}) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ":"
		}
		result += fmt.Sprintf("%v", p)
	}
	return result
}
