package model

import (
	"fmt"

	"qzapp/common"
	"qzapp/redis"
)

func (m *TaskLog) TableName() string {
	return "task_log"
}

func (m *TaskLog) GetRedisKey(field, field_value string) string {
	//前缀:表名:字段:字段值
	return m.TableName() + ":" + field + ":" + field_value
}

func (m *TaskLog) GetOne(where string) TaskLog {
	var q TaskLog
	gdb.DB.Where(where).First(&q)
	return q
}

func (m *TaskLog) RedisByAccount(account, field, field_value string) map[string]string {

	redis_key := m.GetRedisKey(account, field+"_"+field_value)

	//优先查询redis
	a_map := redis.RediGo.HgetAll(redis_key)

	if _, ok := a_map["id"]; ok == false {

		where := fmt.Sprintf("account = '%s' AND %s = '%s'", account, field, field_value)
		a_info := m.GetOne(where)

		if a_info.Id != "" {
			a_map = common.StructToMapSlow(a_info)
			redis.RediGo.Hmset(redis_key, a_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, redis_max_time)
		}
	}
	return a_map
}

func (m *TaskLog) Create() (int64, error) {
	res := gdb.DB.Create(m)
	return res.RowsAffected, res.Error
}
