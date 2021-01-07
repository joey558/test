package model

import (
	"fmt"

	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (b *SysBank) BankInfo(bank_code string) SysBank {
	var m_info SysBank
	gdb.DB.Where("code=?", bank_code).First(&m_info)
	return m_info
}

func (b *SysBank) BankInfoRedis(bank_code string) map[string]string {
	redis_key := fmt.Sprintf("sys_bank:%s", bank_code)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["code"]) < 1 {
		m_info := b.BankInfo(bank_code)
		if len(m_info.Code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Conf_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}

func (b *SysBank) BankByLine(hot_line string) SysBank {
	var m_info SysBank
	gdb.DB.Where("hot_line=?", hot_line).First(&m_info)
	return m_info
}

func (b *SysBank) BankByLineRedis(hot_line string) map[string]string {
	redis_key := fmt.Sprintf("sys_bank:%s", hot_line)
	//优先查询redis
	m_map := redis.RediGo.HgetAll(redis_key)
	if len(m_map["code"]) < 1 {
		m_info := b.BankByLine(hot_line)
		if len(m_info.Code) > 0 {
			m_map = common.StructToMapSlow(m_info)
			redis.RediGo.Hmset(redis_key, m_map, redis_max_time)
			redis.RediGo.Sadd(Conf_Redis_Key, redis_key, 0)
		}
	}
	return m_map
}
