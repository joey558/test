package model

import (
	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (ab *ActiveBonus) ActBonus(act_id string) ActiveBonus {
	var d_info ActiveBonus
	gdb.DB.Where("id=?", act_id).First(&d_info)
	return d_info
}

func (ab *ActiveBonus) BonusRedis(act_id string) map[string]string {
	redis_key := "active_bonus:" + act_id
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["id"]) < 1 {
		d_info := ab.ActBonus(act_id)
		if d_info.Id > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
