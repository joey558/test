package redis

import (
	"github.com/garyburd/redigo/redis"
)

/**
*  向HASH中写入hash
*  HSET
 */
func (rp *RedisPool) Hmset(key_str string, hash_val map[string]string, ex_time int) error {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("HMSET", redis.Args{}.Add(redis_key).AddFlat(hash_val)...)
	if err == nil && ex_time > 0 {
		rp.Expire(key_str, ex_time)
	}

	return err
}

/**
*  从HASH中读取其中一个字段的值
*  HGET
 */
func (rp *RedisPool) HgetField(key_str, field string) string {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.String(conn.Do("HGET", redis_key, field))

	return res
}

/**
*  读取HASH中的所有值
*  HGETALL
 */
func (rp *RedisPool) HgetAll(key_str string) map[string]string {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.StringMap(conn.Do("HGETALL", redis_key))

	return res
}
