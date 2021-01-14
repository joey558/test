package model

import (
	"qzapp/common"
	"qzapp/redis"
)

func (cl *ConfList) ConfList(c_id string) ConfList {
	var d_info ConfList
	gdb.DB.Where("id=?", c_id).First(&d_info)
	return d_info
}

func (cl *ConfList) ConfRedis(c_id string) map[string]string {
	redis_key := "conf_list:" + c_id
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["id"]) < 1 {
		d_info := cl.ConfList(c_id)
		if d_info.Id > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Conf_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
