package model

import (
	"fmt"

	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (wc *WebConf) Config(code string) WebConf {
	var w_conf WebConf
	gdb.DB.Where("code=?", code).First(&w_conf)
	return w_conf
}

func (wc *WebConf) ConfRedis(code string) map[string]string {
	redis_key := fmt.Sprintf("web_conf:%s", code)
	//优先查询redis
	w_map := redis.RediGo.HgetAll(redis_key)
	if len(w_map["code"]) < 1 {
		w_conf := wc.Config(code)
		if len(w_conf.Code) > 0 {
			w_map = common.StructToMapSlow(w_conf)
			redis.RediGo.Hmset(redis_key, w_map, redis_max_time)
			redis.RediGo.Sadd(Conf_Redis_Key, redis_key, 0)
		}
	}
	return w_map
}
