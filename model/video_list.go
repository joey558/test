package model

import (
	"qzapp/common"
	"qzapp/redis"
)

func (dl *VideoList) VideoList(vid string) VideoList {
	var d_info VideoList
	gdb.DB.Where("id=?", vid).First(&d_info)
	return d_info
}

func (dl *VideoList) VideoRedis(vid string) map[string]string {
	redis_key := "video_list:" + vid
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["id"]) < 1 {
		d_info := dl.VideoList(vid)
		if len(d_info.Id) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
