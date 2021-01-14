package thread

import (
	"fmt"
	"qzapp/model"
	"time"

	"qzapp/common"
)

var Loc, _ = time.LoadLocation("Local")

//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

/*
* 改变日期
* @years 增加或减少的年
* @months 增加或减少的月
* @days 增加或减少的日
 */
func ChangeDate(t time.Time, years int, months int, days int) time.Time {
	change_time := t.AddDate(years, months, days)
	return change_time
}

/**
获取本周周一的日期
*/
func GetFirstDateOfWeek() (weekMonday string) {
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekMonday = weekStartDate.Format("2006-01-02")
	return
}

/*
[]map[string]string 转 map[关键字]值
用于判断
*/
func Arr2MapKeyVal(arr []map[string]string, k_name, v_name string) map[string]string {

	m := map[string]string{}
	for _, row := range arr {

		v1, ok1 := row[k_name]
		v2, ok2 := row[v_name]

		if (ok1 == true) && (ok2 == true) {
			m[v1] = v2
		}

	}
	return m
}

/*
[]map[string]string 转 map[关键字]值
用于判断
*/
func TaskLog2Map4Date(task_logs []map[string]string) map[string]string {

	m := map[string]string{}
	for _, row := range task_logs {

		create_time := row["create_time"]
		if len(create_time) != len("2020-08-01 13:28:28") { //日期不对
			continue
		}

		date_ymd := create_time[0:10]
		m[date_ymd] = row["score"]
	}

	return m
}

/*
2020-01-01 10:10:10
转
time.Time
*/
func DateStr2Time(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, Loc)
}

/*
是否26个字母
*/
func IsZM(s string) bool {

	zm_arr := []string{"A", "B", "C", "D", "E", "F", "G",
		"H", "I", "G", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y", "Z"}

	for _, v := range zm_arr {
		if s == v {
			return true
		}
	}
	return false
}

/*
改变用户积分
更新 user_info
添加 amount_list
添加 order_list
添加 task_log

@score -10; 用户积分-10
@order_type 订单类型:1=支付,2=出款,3=优惠,4=额度调整,5=购买积分,6=消费积分
@note 备注
@is_task 是否任务积分;0=否,1=是
@task_id 任务id
*/
func UpUserScore(score, act_id int, account, order_type, note, is_task, task_id string) bool {

	var m_userlist model.UserList
	userlist := m_userlist.UserRedis(account)
	if _, ok1 := userlist["id"]; ok1 == false {
		return false
	}

	sql_uinfo := "SELECT id, score FROM `user_info` WHERE account = '" + account + "'"
	uinfo, _ := model.SqlRow(sql_uinfo, []string{"id", "score"})
	before_score, score_ok := uinfo["score"]
	if score_ok == false {
		return false
	}

	user_score, user_score_err := common.Str2Int(before_score)
	if user_score_err != nil {
		return false
	}

	after_score := user_score + score

	now_datetime := time.Now().Format(f_date)

	id := model.GetKey(16)
	amount := 0
	order_number := id
	create_time := now_datetime
	status := 3  //完成
	pay_id := 0  //支付配置的ID
	is_auto := 1 //是否自动创建:1=自动创建,0=手工创建
	//act_id := 0  //优惠ID
	// creator := "system" //订单创建人
	creator := account      //订单创建人
	operator := creator     //操作人
	pay_time := create_time //操作人

	agent_name := userlist["agent_name"]
	is_test := userlist["is_test"]
	agent_path := userlist["agent_path"]

	user_info := fmt.Sprintf("UPDATE `user_info` SET score = score + %d WHERE account = '%s'", score, account)

	amount_list := fmt.Sprintf("INSERT INTO `amount_list`(`id`, `amount`, `score`, `account`, `order_number`, `agent_name`, `order_type`, `is_test`, `agent_path`, `create_time`, `note`, `before_score`, `after_score`) VALUES ('%s', %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%d')",
		id, amount, score, account, order_number, agent_name, order_type, is_test, agent_path, create_time, note, before_score, after_score)

	order_list := fmt.Sprintf("INSERT INTO `order_list`(`id`, `amount`, `account`, `status`, `order_number`, `agent_name`, `order_type`, `is_test`, `agent_path`, `pay_id`, `create_time`, `is_auto`, `act_id`, `note`, `score`, `creator`, `operator`, `pay_time`) VALUES ('%s', %d, '%s', %d, '%s', '%s', '%s', '%s', '%s', %d, '%s', %d, %d, '%s', %d, '%s', '%s', '%s')",
		id, amount, account, status, order_number, agent_name, order_type, is_test, agent_path, pay_id, create_time, is_auto, act_id, note, score, creator, operator, pay_time)

	sql_arr := []string{user_info, amount_list, order_list}

	if is_task == "1" {

		task_log := fmt.Sprintf("INSERT INTO `task_log`(`id`, `task_id`, `note`, `score`, `account`, `create_time`) VALUES ('%s', '%s', '%s', %d, '%s', '%s')",
			id, task_id, note, score, account, create_time)
		sql_arr = append(sql_arr, task_log)
	}
	err := model.Trans(sql_arr)
	if err == nil {
		return true
	}

	return false
}

/*
中文字符串长度
*/
func ZhongWenChang(s string) int {
	return len(([]rune(s)))
}
