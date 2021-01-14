package model

import (
	"qzapp/common"
	"qzapp/redis"
)

func (dl *VideoComm) Comm(cid string) VideoComm {
	var d_info VideoComm
	gdb.DB.Where("id=?", cid).First(&d_info)
	return d_info
}

func (dl *VideoComm) CommRedis(cid string) map[string]string {
	redis_key := "video_comm:" + cid
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["id"]) < 1 {
		d_info := dl.Comm(cid)
		if len(d_info.Id) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
