package model

import (
	"qzapp/common"
	"qzapp/redis"
)

func (dl *VideoStar) VideoStar(video_id, account string) VideoStar {
	var d_info VideoStar
	gdb.DB.Where("video_id=? and account=?", video_id, account).First(&d_info)
	return d_info
}

func (dl *VideoStar) VideoStarRedis(video_id, account string) map[string]string {
	redis_key := "video_star:" + video_id + "_" + account
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["video_id"]) < 1 {
		d_info := dl.VideoStar(video_id, account)
		if len(d_info.Video_id) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
