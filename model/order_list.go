package model

import (
	"fmt"

	"qzapp/common"
	"qzapp/redis"
)

func (ui *OrderList) Order(account string) OrderList {
	var o_list OrderList
	gdb.DB.Where("account=?", account).First(&o_list)
	return o_list
}

func (ui *OrderList) OrderRedis(account string) map[string]string {
	redis_key := fmt.Sprintf("order_list:account:%s", account)
	//优先查询redis
	u_map := redis.RediGo.HgetAll(redis_key)
	if len(u_map["account"]) < 1 {
		o_list := ui.Order(account)
		if len(o_list.Account) > 0 {
			u_map = common.StructToMapSlow(o_list)
			redis.RediGo.Hmset(redis_key, u_map, redis_max_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return u_map
}
