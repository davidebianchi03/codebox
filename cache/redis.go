package cache

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gomodule/redigo/redis"
	"gitlab.com/codebox4073715/codebox/config"
)

var redisPoolSingletonMutex sync.Mutex
var redisPool *redis.Pool

/*
Singleton function to get the redis pool for caching
*/
func GetRedisCachePool() *redis.Pool {
	if redisPool == nil {
		redisPoolSingletonMutex.Lock()
		if redisPool == nil {
			redisPool = &redis.Pool{
				MaxActive: 5,
				MaxIdle:   5,
				Wait:      true,
				Dial: func() (redis.Conn, error) {
					return redis.Dial(
						"tcp",
						fmt.Sprintf(
							"%s:%s",
							config.Environment.RedisHost,
							strconv.Itoa(config.Environment.RedisPort),
						),
					)
				},
			}
		}
		redisPoolSingletonMutex.Unlock()
	}
	return redisPool
}

/*
Set a key in cache, if expiration is set to a value that is equal or
less than zero the key won't expire
*/
func SetKeyToCache(key string, value []byte, ttlSeconds int) error {
	pool := GetRedisCachePool()

	conn := pool.Get()
	defer conn.Close()

	var err error
	if ttlSeconds <= 0 {
		_, err = conn.Do("SET", key, value)
	} else {
		_, err = conn.Do("SET", key, value, "EX", ttlSeconds)
	}

	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}
	return err
}

/*
Retrieve keys matching pattern from redis cache
*/
func GetKeysByPatternFromCache(pattern string) ([]string, error) {
	pool := GetRedisCachePool()

	conn := pool.Get()
	defer conn.Close()

	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

/*
Remove a key from db
*/
func DeleteKeyFromCache(key string) error {
	pool := GetRedisCachePool()

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}
