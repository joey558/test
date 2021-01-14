package redis

import (
	"github.com/garyburd/redigo/redis"
)

func (rp *RedisPool) Expire(key_str string, ex_time int) int {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.Int(conn.Do("EXPIRE", redis_key, ex_time))
	return res
}

/**
* 删除一个指定的key
* @key string 需要删除的键值,多个用空格隔开
* return int  返回删除的个数
 */
func (rp *RedisPool) KeyDel(key_str string) int {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.Int(conn.Do("DEL", redis_key))
	return res
}
