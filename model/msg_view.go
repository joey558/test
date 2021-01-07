package model

import (
	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (dl *MsgView) MsgView(msg_id, account string) MsgView {
	var d_info MsgView
	gdb.DB.Where("msg_id=? and account=?", msg_id, account).First(&d_info)
	return d_info
}

func (dl *MsgView) MsgRedis(msg_id, account string) map[string]string {
	redis_key := "msg_view:" + msg_id + ":" + account
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["msg_id"]) < 1 {
		d_info := dl.MsgView(msg_id, account)
		if len(d_info.Id) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
