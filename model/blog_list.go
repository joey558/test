package model

import (
	"qzapp/common"
	"qzapp/redis"
)

func (m *BlogList) TableName() string {
	return "blog_list"
}

func (m *BlogList) GetRedisKey(field, field_value string) string {
	//前缀:表名:字段:字段值
	return m.TableName() + ":" + field + ":" + field_value
}

func (m *BlogList) GetOne(where string) BlogList {
	var q BlogList
	gdb.DB.Where(where).First(&q)
	return q
}

func (m *BlogList) RedisGetOne(field, field_value string) map[string]string {

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

func (m *BlogList) RedisGetOneDel(field, field_value string) int {
	redis_key := m.GetRedisKey(field, field_value)
	return redis.RediGo.KeyDel(redis_key)
}
