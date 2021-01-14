package redis

import (
	"github.com/garyburd/redigo/redis"
)

/**
*  有序集合新增数据
*  ZADD
 */
func (rp *RedisPool) Zadd(key_str, val string, score int64, ex_time int) int {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.Int(conn.Do("ZADD", redis_key, score, val))
	if ex_time > 0 {
		rp.Expire(key_str, ex_time)
	}
	return res
}

/**
*  移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员
*  ZREMRANGEBYSCORE
 */
func (rp *RedisPool) ZremRangeByScore(key_str string, min, max int64) error {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("ZREMRANGEBYSCORE", redis_key, min, max)

	return err
}

/**
*  返回有序集 key 中，所有 index 值介于 min 和 max 之间(包括等于 min 或 max )的成员
*  ZRANGE
 */
func (rp *RedisPool) Zrange(key_str string, min, max int) []string {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.Strings(conn.Do("ZRANGE", redis_key, min, max))

	return res
}

/**
*  返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员
*  ZRANGE
 */
func (rp *RedisPool) ZrangeByScore(key_str string, min, max int64, count int) []string {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.Strings(conn.Do("ZRANGEBYSCORE", redis_key, min, max, count))

	return res
}

/**
*  删除指定元素
*  ZRANGE
 */
func (rp *RedisPool) Zrem(key_str, member string) int {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.Int(conn.Do("ZREM", redis_key, member))

	return res
}
