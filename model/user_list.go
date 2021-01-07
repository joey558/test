package model

import (
	"fmt"

	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (ul *UserList) User(account string) UserList {
	var u_info UserList
	gdb.DB.Where("account=?", account).First(&u_info)
	return u_info
}

func (ul *UserList) UserRedis(account string) map[string]string {
	redis_key := fmt.Sprintf("user_list:account:%s", account)
	//优先查询redis
	u_map := redis.RediGo.HgetAll(redis_key)
	if len(u_map["account"]) < 1 {
		u_info := ul.User(account)
		if len(u_info.Account) > 0 {
			u_map = common.StructToMapSlow(u_info)
			redis.RediGo.Hmset(redis_key, u_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return u_map
}

func (ul *UserList) UserByCode(link_code string) UserList {
	var u_info UserList
	gdb.DB.Where("link_code=?", link_code).First(&u_info)
	return u_info
}

func (ul *UserList) UserByCodeRedis(link_code string) map[string]string {
	redis_key := fmt.Sprintf("user_list:link_code:%s", link_code)
	//优先查询redis
	u_map := redis.RediGo.HgetAll(redis_key)
	if len(u_map["account"]) < 1 {
		u_info := ul.UserByCode(link_code)
		if len(u_info.Account) > 0 {
			u_map = common.StructToMapSlow(u_info)
			redis.RediGo.Hmset(redis_key, u_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return u_map
}

func (ul *UserList) UserBySession(session_id string) UserList {
	var u_info UserList
	gdb.DB.Where("session_id=?", session_id).First(&u_info)
	return u_info
}

func (ul *UserList) UserBySessRedis(session_id string) map[string]string {
	redis_key := fmt.Sprintf("user_list:session_id:%s", session_id)
	//优先查询redis
	u_map := redis.RediGo.HgetAll(redis_key)
	if len(u_map["account"]) < 1 {
		u_info := ul.UserBySession(session_id)
		if len(u_info.Account) > 0 {
			u_map = common.StructToMapSlow(u_info)
			redis.RediGo.Hmset(redis_key, u_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return u_map
}

func (ul *UserList) Updates(u_info UserList, u_data map[string]interface{}) error {
	res := gdb.DB.Model(&u_info).UpdateColumns(u_data)
	return res.Error
}
