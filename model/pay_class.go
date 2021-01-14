package model

import (
	"qzapp/common"
	"qzapp/redis"
)

func (pc *PayClass) Class(code string) PayClass {
	var p_class PayClass
	gdb.DB.Where("code=?", code).First(&p_class)
	return p_class
}

func (pc *PayClass) ClassRedis(code string) map[string]string {
	redis_key := "pay_class:" + code
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["code"]) < 1 {
		d_info := pc.Class(code)
		if len(d_info.Code) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
