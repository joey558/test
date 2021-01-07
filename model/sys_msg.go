package model

import (
	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (dl *SysMsg) SysMsg(id string) SysMsg {
	var d_info SysMsg
	gdb.DB.Where("id=?", id).First(&d_info)
	return d_info
}

func (dl *SysMsg) MsgRedis(id string) map[string]string {
	redis_key := "sys_msg:" + id
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["id"]) < 1 {
		d_info := dl.SysMsg(id)
		if len(d_info.Id) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
