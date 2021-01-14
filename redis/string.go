package redis

import (
	"github.com/garyburd/redigo/redis"
)

func (rp *RedisPool) StringWrite(key_str, val string, ex_time int) error {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	_, err := redis.String(conn.Do("SET", redis_key, val))
	if err == nil && ex_time > 0 {
		rp.Expire(key_str, ex_time)
	}
	return err
}

func (rp *RedisPool) StringRead(key_str string) string {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, _ := redis.String(conn.Do("GET", redis_key))

	return res
}

/**
* 向redis中写入一个key—val类型的字符串，如果存在，则不操作返回0
* @key 	string 	参数的主键
* @val 	string 	需要写入缓存的值
* @Ex	int		超时时间（秒）
				参数为0时，永远不过期

* return	返回影响行数
*/
func (rp *RedisPool) StringWriteNx(key_str string, val string, ex_time int) int {
	redis_key := rp.Pre_key + ":" + key_str
	conn := rp.Pool.Get()
	defer conn.Close()
	res, err := redis.Int(conn.Do("SETNX", redis_key, val))
	if err == nil && ex_time > 0 {
		rp.Expire(key_str, ex_time)
	}
	return res
}
