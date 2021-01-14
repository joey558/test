package model

import (
	"qzapp/common"
	"qzapp/redis"
)

func (dl *ThemeVideo) ThemeVideo(id string) ThemeVideo {
	var d_info ThemeVideo
	gdb.DB.Where("id=?", id).First(&d_info)
	return d_info
}

func (dl *ThemeVideo) VideoRedis(id string) map[string]string {
	redis_key := "theme_video:" + id
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["id"]) < 1 {
		d_info := dl.ThemeVideo(id)
		if len(d_info.Id) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
