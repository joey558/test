package hook

import (
	"qzapp/model"

	"qzapp/common"
)

var key_code = "private_key"

func WebConfVal(code string) string {
	var model_wc model.WebConf
	w_conf := model_wc.ConfRedis(code)
	res, _ := w_conf["val"]
	return res
}

func HookAesDecrypt(decode_str string) string {
	val := WebConfVal(key_code)
	aes_res := ""
	if val == "" {
		return aes_res
	}
	aes := common.SetAES(val, "", "", 16)
	aes_res = aes.AesDecryptString(decode_str)
	return aes_res
}

func HookAesEncrypt(encode_str string) string {
	val := WebConfVal(key_code)
	aes_res := ""
	if val == "" {
		return aes_res
	}
	aes := common.SetAES(val, "", "", 16)
	aes_res = aes.AesEncryptString(encode_str)
	return aes_res
}
