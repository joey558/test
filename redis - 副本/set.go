package redis

import (
	"github.com/garyburd/redigo/redis"
)

/**
*  无序集合新增数据
*  SADD
 */
func (rp *RedisPool) Sadd(key_str, val string, ex_time int) int {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.Int(conn.Do("SADD", redis_key, val))
	if ex_time > 0 {
		rp.Expire(key_str, ex_time)
	}
	return res
}

/**
* 计算集合数量
* @key	string	集合的key
* @param int 返回集合中数据的数量
 */
func (rp *RedisPool) Scard(key_str string) int {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.Int(conn.Do("SCARD", redis_key))

	return res
}

/**
* 随机删除一条数据然后返回被删除的元素
* @return string  删除的元素
 */
func (rp *RedisPool) Spop(key_str string) string {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.String(conn.Do("SPOP", redis_key))

	return res
}
