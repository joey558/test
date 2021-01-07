package model

import (
	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (dl *PlatList) Plat(code string) PlatList {
	var p_info PlatList
	gdb.DB.Where("code=?", code).First(&p_info)
	return p_info
}

func (dl *PlatList) PlatRedis(code string) map[string]string {
	redis_key := "plat_list:" + code
	//优先查询redis
	p_map := redis.RediGo.HgetAll(redis_key)
	if len(p_map["code"]) < 1 {
		p_info := dl.Plat(code)
		if len(p_info.Code) > 0 {
			p_map = common.StructToMapSlow(p_info)
			redis.RediGo.Hmset(redis_key, p_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return p_map
}
