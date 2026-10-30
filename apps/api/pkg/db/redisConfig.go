package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"forestwhisperGo/apps/api/pkg/utils/configutil"

	"github.com/alicebob/miniredis/v2"

	"github.com/redis/go-redis/v9"
)

var (
	ctx      = context.Background()
	rdbMu    sync.RWMutex
	rdbConn  *redis.Client
	realAddr string
	redisUP  atomic.Bool
)

func rdb() *redis.Client {
	rdbMu.RLock()
	defer rdbMu.RUnlock()
	return rdbConn
}

func newClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        addr,
		DialTimeout: 2 * time.Second,
		MaxRetries:  -1, // -1 才是关闭重试, 0 表示默认 3 次
	})
}

func pingAddr(addr string) error {
	c := newClient(addr)
	defer c.Close()
	pctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.Ping(pctx).Err()
}

func CheckRedis() error {
	return pingAddr(realAddr)
}

func useMemory() {
	m, err := miniredis.Run()
	if err != nil {
		log.Printf("[redis] memory Redis startup failed, continuing with real Redis: %v", err)
		return
	}
	rdbMu.Lock()
	rdbConn = newClient(m.Addr())
	rdbMu.Unlock()
}

func init() {
	loadConfig()
	host := configutil.GetSectionString(cfg, "redis", "REDIS_HOST")
	port := configutil.GetSectionInt(cfg, "redis", "REDIS_PORT")
	if host == "" {
		host = "127.0.0.1"
	}
	if port == 0 {
		port = 6379
	}

	realAddr = host + ":" + strconv.Itoa(port)

	rdbConn = newClient(realAddr)
	redisUP.Store(pingAddr(realAddr) == nil)
	if !redisUP.Load() {
		useMemory()
		log.Println("[redis] Redis is unavailable, using memory")
	}

	go func() {
		for range time.Tick(30 * time.Second) {
			ok := pingAddr(realAddr) == nil
			if ok == redisUP.Load() {
				continue
			}
			if ok {
				rdbMu.Lock()
				rdbConn = newClient(realAddr)
				rdbMu.Unlock()
				log.Println("[redis] Redis has been restored; switching back to Redis")
			} else {
				useMemory()
				log.Println("[redis] Redis is unavailable; switching to memory")
			}
			redisUP.Store(ok)
		}
	}()
}

func Set(key, value string, ttl ...time.Duration) error {
	if value == "" {
		return nil
	}
	if len(ttl) > 0 {
		return rdb().Set(ctx, key, value, ttl[0]).Err()
	}
	return rdb().Set(ctx, key, value, 0).Err()
}

func SetEX(key, value string, seconds int) error {
	return Set(key, value, time.Duration(seconds)*time.Second)
}

func Get(key string) (string, error) {
	return rdb().Get(ctx, key).Result()
}

func Del(keys ...string) error {
	return rdb().Del(ctx, keys...).Err()
}

func Exists(keys ...string) (int64, error) {
	return rdb().Exists(ctx, keys...).Result()
}

func Expire(ctx context.Context, key string, ttl time.Duration) error {
	return rdb().Expire(ctx, key, ttl).Err()
}

func TTL(key string) (time.Duration, error) {
	return rdb().TTL(ctx, key).Result()
}

func Keys(pattern string) ([]string, error) {
	return rdb().Keys(ctx, pattern).Result()
}

func Incr(key string) (int64, error) {
	return rdb().Incr(ctx, key).Result()
}

func Decr(key string) (int64, error) {
	return rdb().Decr(ctx, key).Result()
}

func HSet(key string, field string, value interface{}) error {
	return rdb().HSet(ctx, key, field, value).Err()
}

func HGet(key, field string) (string, error) {
	return rdb().HGet(ctx, key, field).Result()
}

func HGetAll(key string) (map[string]string, error) {
	return rdb().HGetAll(ctx, key).Result()
}

func HDel(key string, fields ...string) error {
	return rdb().HDel(ctx, key, fields...).Err()
}

func SAdd(ctx context.Context, key string, members ...interface{}) error {
	return rdb().SAdd(ctx, key, members...).Err()
}

func SMembers(key string) ([]string, error) {
	return rdb().SMembers(ctx, key).Result()
}

func SRem(key string, members ...interface{}) error {
	return rdb().SRem(ctx, key, members...).Err()
}

func ZAdd(key string, score float64, member string) error {
	return rdb().ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

func ZRange(key string, start, stop int64) ([]string, error) {
	return rdb().ZRange(ctx, key, start, stop).Result()
}

func ZRem(key string, members ...interface{}) error {
	return rdb().ZRem(ctx, key, members...).Err()
}

func Client() *redis.Client {
	return rdb()
}

func Close() error {
	return rdb().Close()
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
