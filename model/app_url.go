package model

import (
	"fmt"

	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

const (
	IOS_APP_TYPE     = "ios"
	ANDROID_APP_TYPE = "android"
)

func (m *AppUrl) TableName() string {
	return "app_url"
}

func (m *AppUrl) GetRedisKey(field, field_value string) string {
	//前缀:表名:字段:字段值
	return m.TableName() + ":" + field + ":" + field_value
}

func (m *AppUrl) GetOne(where string) AppUrl {
	var q AppUrl
	gdb.DB.Where(where).First(&q)
	return q
}

func (m *AppUrl) RedisGetOne(field, field_value string) map[string]string {

	redis_key := m.GetRedisKey(field, field_value)

	//优先查询redis
	a_map := redis.RediGo.HgetAll(redis_key)

	if _, ok := a_map["id"]; ok == false {

		where := field + "='" + field_value + "'"
		a_info := m.GetOne(where)

		if a_info.Id > 0 {
			a_map = common.StructToMapSlow(a_info)
			redis.RediGo.Hmset(redis_key, a_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, redis_max_time)
		}
	}
	return a_map
}

func (m *AppUrl) RedisByCodeApptype(code, app_type string) map[string]string {

	redis_key := m.GetRedisKey(app_type, code)

	//优先查询redis
	a_map := redis.RediGo.HgetAll(redis_key)

	if _, ok := a_map["id"]; ok == false {

		where := fmt.Sprintf("app_type = '%s' AND code = '%s'", app_type, code)
		a_info := m.GetOne(where)

		if a_info.Id > 0 {
			a_map = common.StructToMapSlow(a_info)
			redis.RediGo.Hmset(redis_key, a_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, redis_max_time)
		}
	}
	return a_map
}
