package model

import (
	"qzapp/common"
	"qzapp/redis"
)

func (as *AppSetup) AppSet(uid string) AppSetup {
	var d_info AppSetup
	gdb.DB.Where("uid=?", uid).First(&d_info)
	return d_info
}

func (dl *AppSetup) AppSetRedis(uid string) map[string]string {
	redis_key := "app_setup:" + uid
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["uid"]) < 1 {
		d_info := dl.AppSet(uid)
		if len(d_info.Uid) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
