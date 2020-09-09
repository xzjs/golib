//go:generate mockgen -destination ../mock/cache.go -package mock -source cache.go

package golib

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

// ICache 接口列表
type ICache interface {
	Set(key string, val string, time int64) (err error)
	Get(key string) (val string, err error)
	Del(key string) (err error)
	Ins(key string) (err error)
	Zscore(key string, member string) (score int, err error)
	Zadd(key string, member string, score int) error
	// Zincrby 对有序集合中指定成员的分数加上增量 increment
	Zincrby(key string, increment int, member string) (score int, err error)
	// Zrevrange 命令返回有序集中，指定区间内的成员。其中成员的位置按分数值递减(从大到小)来排列。
	Zrevrange(key string, start int, stop int) ([]int, error)
	// Zcard 获取有序集合中元素的个数
	Zcard(key string) (int, error)
}

// MyCache 缓存单例
var cache ICache

// MyCache 缓存结构体
type MyCache struct {
	Pool redis.Pool
}

// Cache 获取单例
func Cache() ICache {
	if cache == nil {
		conf := Conf()
		addr := conf.GetConf("redis", "addr")
		pool := redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", addr)
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}

		cache = &MyCache{Pool: pool}
	}
	return cache
}

// SetCache 设置单例
func SetCache(c ICache) {
	cache = c
}

// Set set
func (c *MyCache) Set(key string, val string, time int64) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()
	if time > 0 {
		_, err = conn.Do("SETEX", key, time, val)
	} else {
		_, err = conn.Do("SET", key, val)
	}

	return
}

// Get get
func (c *MyCache) Get(key string) (val string, err error) {
	conn := c.Pool.Get()
	defer conn.Close()
	val, err = redis.String(conn.Do("GET", key))
	return
}

//Del 自减
func (c *MyCache) Del(key string) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("DECR", key)
	return
}

//Ins 自加
func (c *MyCache) Ins(key string) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("INCR", key)
	return
}

// Zscore 获取有序集合成员的分数
func (c *MyCache) Zscore(key string, member string) (score int, err error) {
	conn := c.Pool.Get()
	defer conn.Close()
	reply, err := conn.Do("ZSCORE", key, member)
	if reply == nil {
		return 0, errors.New("缓存中无数据")
	}
	score, err = redis.Int(reply, err)
	return
}

// Zadd 插入一个元素到有序集合
func (c *MyCache) Zadd(key string, member string, score int) (err error) {
	conn := c.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("ZADD", key, score, member)
	return
}

// Zincrby 对有序集合中指定成员的分数加上增量 increment
func (c *MyCache) Zincrby(key string, increment int, member string) (score int, err error) {
	conn := c.Pool.Get()
	defer conn.Close()
	reply, err := conn.Do("ZINCRBY", key, increment, member)
	if err != nil {
		return 0, err
	}
	score, err = redis.Int(reply, err)
	return
}

// Zrevrange 命令返回有序集中，指定区间内的成员。其中成员的位置按分数值递减(从大到小)来排列。
func (c *MyCache) Zrevrange(key string, start int, stop int) (result []int, err error) {
	conn := c.Pool.Get()
	defer conn.Close()
	return redis.Ints(conn.Do("ZREVRANGE", key, start, stop, "WITHSCORES"))
}

// Zcard 获取有序集合中元素的个数
func (c *MyCache) Zcard(key string) (count int, err error) {
	conn := c.Pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("ZCARD", key))
}
