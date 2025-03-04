package model

import (
	"qzapp/common"
	"qzapp/redis"
)

func (dl *DomainList) DomainList(domain string) DomainList {
	var d_info DomainList
	gdb.DB.Where("domain=?", domain).First(&d_info)
	return d_info
}

func (dl *DomainList) DomainRedis(domain string) map[string]string {
	redis_key := "domain_list:" + domain
	//优先查询redis
	d_map := redis.RediGo.HgetAll(redis_key)
	if len(d_map["domain"]) < 1 {
		d_info := dl.DomainList(domain)
		if len(d_info.Domain) > 0 {
			d_map = common.StructToMapSlow(d_info)
			redis.RediGo.Hmset(redis_key, d_map, redis_data_time)
			redis.RediGo.Sadd(Data_Redis_Key, redis_key, 0)
		}
	}
	return d_map
}
