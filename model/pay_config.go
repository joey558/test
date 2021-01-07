package model

import (
	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (pc *PayConfig) PayConf(p_id string) PayConfig {
	var d_info PayConfig
	gdb.DB.Where("id=?", p_id).First(&d_info)
	return d_info
}

func (pc *PayConfig) PayRedis(p_id string) map[string]string {
	redis_key := "pay_config:" + p_id
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["id"]) < 1 {
		d_info := pc.PayConf(p_id)
		if d_info.Id > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
