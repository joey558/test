package thread

import (
	"encoding/json"
	"fmt"
	"qzapp/model"
	"regexp"
	"strconv"
	"time"

	"qzapp/common"
	"qzapp/redis"
)

const (
	SMS_URL          = "http://api.shunf.net/api/msg.do"
	SMS_KEY          = "EJO43Ddklsdfoln8" //AES加密key
	SMS_REDIS_KEY    = "sms:"             //redis key
	SMS_VALID_PERIOD = 900                //验证码有效期 15分钟
)

//加密
func AesEncrypt(private_key, encode_str string) string {
	mer_aes := common.SetAES(private_key, "", "", 16)
	aes_res := mer_aes.AesEncryptString(encode_str)
	return aes_res
}

// 发送短信
func MSMApi(phone_number, verify_code string) (int, string) {

	if !VerifyMobileFormat(phone_number) {
		return 321, "手机号码不规范"
	}

	header := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}

	str := fmt.Sprintf(`{"msg_id":"1","mobile":"%s","code":"%s","template":"524710"}`, phone_number, verify_code)
	code_str := AesEncrypt(SMS_KEY, str)

	param := "mer_id=1&post_data=" + code_str

	h_status, h_body := common.HttpBody(SMS_URL, "POST", param, header)
	if h_status != 200 {
		return h_status, "网络请求失败"
	}
	fmt.Println("h_body->", string(h_body))

	var json_res map[string]interface{}
	err := json.Unmarshal(h_body, &json_res)
	if err != nil {
		return 601, "返回数据解析失败"
	}

	api_status, _ := common.Interface2Int(json_res["Status"])
	api_msg := fmt.Sprintf("%s", json_res["Msg"])

	return api_status, api_msg
}

// 获取redis 今天 key
func VfctCodeKey(account string) string {

	day := time.Now().Format(f_day)
	k := SMS_REDIS_KEY + account + ":" + day
	return k
}

// 生成数字验证码
func VfctCode(red_key string, length int) string {

	rand_str := ""
	for i := 0; i < 50; i++ {
		rand_str = common.Random("number", length)
		red_res := redis.RediGo.Sadd(red_key, rand_str, (60 * 60 * 24 * 2)) //缓存2天
		if red_res > 0 {
			break
		}
	}

	return rand_str
}

// 发送验证码
func MSMSend(account, phone_number string) (int, string, string) {
	currentTime := time.Now()
	//今天的开始时间
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Format(f_date)
	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM captcha_list where account ='%s' and create_time >='%s' and create_time <='%s' ;", account, startTime, currentTime.Format(f_date))

	total_map, _ := model.SqlRow(count_sql, count_field) //获取今天发送的短信数量
	num_int, _ := strconv.Atoi(total_map["num"])
	if num_int >= 3 {
		return 322, "每个用户每天只能发送3条短信", ""
	}

	red_key := VfctCodeKey(account)
	vfct := VfctCode(red_key, 6) //6位验证码
	if vfct == "" {
		return 323, "生成验证码失败", ""
	}

	status, msg := MSMApi(phone_number, vfct)
	if status != 200 {
		//发送失败不计入今天的短信数量
		return status, msg, ""
	}
	id := model.GetKey(16)
	list_sql := fmt.Sprintf("insert into captcha_list (id,code,phone_number,account,create_time) VALUES ('%s','%s','%s','%s','%s');", id, vfct, phone_number, account, currentTime.Format(f_date))
	err := model.Query(list_sql)
	if err != nil {
		return 324, "保存验证码失败", ""
	}

	return 200, "发送成功", vfct
}

//mobile verify
func VerifyMobileFormat(mobileNum string) bool {
	// regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0-9])|(18[0-9])|166|198|199|(147))\\d{8}$"
	regular := "^(1)\\d{10}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

// 获取redis 今天 key
func MSMSendKey(account string) string {
	return SMS_REDIS_KEY + account + ":send"
}

func MSMSendCheck(account, code string) (int, string, string) {
	t_status := 100
	t_msg := "验证码错误"

	//查询验证码
	list_field := []string{"code", "phone_number", "create_time"}
	list_sql := fmt.Sprintf("SELECT code,phone_number,create_time FROM captcha_list where account='%s' order by create_time desc limit 1;", account)

	list_map, _ := model.SqlRow(list_sql, list_field)
	if len(list_map) < 1 {
		t_msg = "请先获取验证码"
		return t_status, t_msg, ""
	}

	phone_number, ok := list_map["phone_number"]
	if ok == false {
		t_status = 331
		return t_status, t_msg, ""
	}

	verify_code, ok := list_map["code"]
	if ok == false {
		t_status = 332
		return t_status, t_msg, ""
	}

	if verify_code != code {
		t_status = 333
		return t_status, t_msg, ""
	}

	now_str := time.Now().Format(f_date)
	now_time, _ := time.Parse("2006-01-02 15:04:05", now_str)
	create_time, _ := time.Parse("2006-01-02 15:04:05", list_map["create_time"])
	diff := now_time.Sub(create_time).Seconds()
	//判断验证码是否失效
	if diff > SMS_VALID_PERIOD {
		t_status = 334
		t_msg := "该验证码已过期"
		return t_status, t_msg, ""
	}

	t_status = 200
	t_msg = "验证码正确"
	return t_status, t_msg, phone_number
}

// 删除缓存
func MSMSendKeyDel(account string) {
	redis_key := SMS_REDIS_KEY + account + ":send"
	redis.RediGo.KeyDel(redis_key)
}
