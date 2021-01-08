package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"gitlab.stagingvip.net/publicGroup/public/common"
)

//连接池的结构体
type RedisPool struct {
	Pool    *redis.Pool
	Pre_key string
}

var RediGo *RedisPool

func init() {
	ReloadConf("")
}

func ReloadConf(file_name string) {
	if file_name == "" {
		file_name = "./conf/redis.json"
	}
	conf_byte, err := common.ReadFile(file_name)
	if err != nil {
		panic(err)
	}

	var json_conf map[string]string
	//解析json格式
	err = json.Unmarshal(conf_byte, &json_conf)
	if err != nil {
		panic(err)
	}
	r_time, _ := time.ParseDuration(json_conf["read_time"])
	db, _ := strconv.Atoi(json_conf["db"])
	c_time, _ := time.ParseDuration(json_conf["connect_time"])

	life_time, _ := time.ParseDuration(json_conf["life_time"])
	max_open, _ := strconv.Atoi(json_conf["max_open"])
	if max_open < 1 {
		max_open = 40
	}
	max_idle, _ := strconv.Atoi(json_conf["max_idle"])
	if max_idle < 1 {
		max_idle = 10
	}

	opt_read := redis.DialReadTimeout(r_time)
	opt_db := redis.DialDatabase(db)
	opt_conn := redis.DialConnectTimeout(c_time)
	opt_pwd := redis.DialPassword(json_conf["auth"])
	addr := fmt.Sprintf("%s:%s", json_conf["host"], json_conf["port"])
	pool := &redis.Pool{
		IdleTimeout: life_time,
		MaxIdle:     max_idle,
		MaxActive:   max_open,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(json_conf["network"], addr,
				opt_read, opt_db, opt_conn, opt_pwd,
			)
			if err != nil {
				panic(err)
			}

			return c, nil
		},
	}

	RediGo = &RedisPool{pool, json_conf["pre_key"]}
}
