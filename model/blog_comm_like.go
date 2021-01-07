package model

import (
	"fmt"

	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (m *BlogCommLike) TableName() string {
	return "blog_comm_like"
}

func (m *BlogCommLike) GetRedisKey(field, field_value string) string {
	//前缀:表名:字段:字段值
	return m.TableName() + ":" + field + ":" + field_value
}

func (m *BlogCommLike) GetOne(where string) BlogCommLike {
	var q BlogCommLike
	gdb.DB.Where(where).First(&q)
	return q
}

func (m *BlogCommLike) RedisGetOne(field, field_value string) map[string]string {

	redis_key := m.GetRedisKey(field, field_value)

	//优先查询redis
	a_map := redis.RediGo.HgetAll(redis_key)

	if _, ok := a_map["id"]; ok == false {

		where := field + "='" + field_value + "'"
		a_info := m.GetOne(where)

		if a_info.Id != "" {
			a_map = common.StructToMapSlow(a_info)
			redis.RediGo.Hmset(redis_key, a_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, redis_max_time)
		}
	}
	return a_map
}

func (m *BlogCommLike) RedisIsLike(blog_list_id, account string) map[string]string {

	redis_key := m.GetRedisKey(blog_list_id, account)

	//优先查询redis
	a_map := redis.RediGo.HgetAll(redis_key)

	if _, ok := a_map["id"]; ok == false {

		where := fmt.Sprintf("comm_id='%s' AND account='%s'", blog_list_id, account)
		a_info := m.GetOne(where)

		// if a_info.Id != "" { 没有数据也缓存查询结果
		a_map = common.StructToMapSlow(a_info)
		redis.RediGo.Hmset(redis_key, a_map, redis_data_time)
		redis.RediGo.Sadd(Data_Redis_Key, redis_key, redis_max_time)
		// }
	}
	return a_map
}

func (m *BlogCommLike) RedisIsLikeDel(blog_list_id, account string) int {

	redis_key := m.GetRedisKey(blog_list_id, account)
	return redis.RediGo.KeyDel(redis_key)
}
