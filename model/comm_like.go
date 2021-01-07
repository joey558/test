package model

import (
	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (cl *CommLike) CommLike(comm_id, account string) CommLike {
	var d_info CommLike
	gdb.DB.Where("comm_id=? and account=?", comm_id, account).First(&d_info)
	return d_info
}

func (cl *CommLike) CommLikeRedis(comm_id, account string) map[string]string {
	redis_key := "comm_like:" + comm_id + "_" + account
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["comm_id"]) < 1 {
		d_info := cl.CommLike(comm_id, account)
		if len(d_info.Comm_id) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
