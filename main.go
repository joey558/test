package main

import (
	"encoding/json"
	"qzapp/router"

	"gitlab.stagingvip.net/publicGroup/public/common"
)

func main() {
	conf_byte, err := common.ReadFile("./conf/conf.json")

	if err != nil {
		panic(err)
	}
	var json_conf map[string]string
	//解析json格式
	err = json.Unmarshal(conf_byte, &json_conf)
	if err != nil {
		panic(err)
	}

	_ = router.Router.Run(json_conf["port"])
}
