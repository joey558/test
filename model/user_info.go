package model

import (
	"fmt"

	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (ui *UserInfo) User(account string) UserInfo {
	var u_info UserInfo
	gdb.DB.Where("account=?", account).First(&u_info)
	return u_info
}

func (ui *UserInfo) UserRedis(account string) map[string]string {
	redis_key := fmt.Sprintf("user_info:account:%s", account)
	//优先查询redis
	u_map := redis.RediGo.HgetAll(redis_key)
	if len(u_map["account"]) < 1 {
		u_info := ui.User(account)
		if len(u_info.Account) > 0 {
			u_map = common.StructToMapSlow(u_info)
			redis.RediGo.Hmset(redis_key, u_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return u_map
}

func (ui *UserInfo) UserByNick(nick_name string) UserInfo {
	var u_info UserInfo
	gdb.DB.Where("nick_name=?", nick_name).First(&u_info)
	return u_info
}

func (ui *UserInfo) UserByNickRedis(nick_name string) map[string]string {
	redis_key := fmt.Sprintf("user_info:nick_name:%s", nick_name)
	//优先查询redis
	u_map := redis.RediGo.HgetAll(redis_key)
	if len(u_map["nick_name"]) < 1 {
		u_info := ui.UserByNick(nick_name)
		if len(u_info.Nick_name) > 0 {
			u_map = common.StructToMapSlow(u_info)
			redis.RediGo.Hmset(redis_key, u_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return u_map
}
