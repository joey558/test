package model

import (
	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (dl *ThemeLike) ThemeLike(comm_id, account string) ThemeLike {
	var d_info ThemeLike
	gdb.DB.Where("comm_id=? and account=?", comm_id, account).First(&d_info)
	return d_info
}

func (dl *ThemeLike) ThemeLikeRedis(comm_id, account string) map[string]string {
	redis_key := "theme_like:" + comm_id + "_" + account
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["comm_id"]) < 1 {
		d_info := dl.ThemeLike(comm_id, account)
		if len(d_info.Comm_id) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
