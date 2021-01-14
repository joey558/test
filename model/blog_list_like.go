package model

import (
	"qzapp/common"
	"qzapp/redis"
)

func (m *BlogListLike) TableName() string {
	return "blog_list_like"
}

func (m *BlogListLike) GetRedisKey(field, field_value string) string {
	//前缀:表名:字段:字段值
	return m.TableName() + ":" + field + ":" + field_value
}

func (m *BlogListLike) GetOne(where string) BlogListLike {
	var q BlogListLike
	gdb.DB.Where(where).First(&q)
	return q
}

func (m *BlogListLike) RedisGetOne(field, field_value string) map[string]string {

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

func (m *BlogListLike) BlogLike(blog_list_id, account string) BlogListLike {
	var d_info BlogListLike
	gdb.DB.Where("blog_list_id=? and account=?", blog_list_id, account).First(&d_info)
	return d_info
}

func (m *BlogListLike) RedisIsLikeKey(blog_list_id, account string) string {
	redis_key := "blog_list_like:" + blog_list_id + "_" + account
	return redis_key
}

func (m *BlogListLike) RedisIsLike(blog_list_id, account string) map[string]string {

	redis_key := m.RedisIsLikeKey(blog_list_id, account)

	//优先查询redis
	a_map := redis.RediGo.HgetAll(redis_key)

	if len(a_map["id"]) < 1 {
		d_info := m.BlogLike(blog_list_id, account)
		if len(d_info.Id) > 0 {
			a_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, a_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}

	return a_map
}

func (m *BlogListLike) RedisIsLikeDel(blog_list_id, account string) int {
	redis_key := m.RedisIsLikeKey(blog_list_id, account)
	return redis.RediGo.KeyDel(redis_key)
}
