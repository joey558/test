package thread

import (
	"encoding/json"
	"fmt"
	"qzapp/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.stagingvip.net/publicGroup/public/common"
)

const (
	CHECKIN_TYPE       = "5"             //签到type
	CHECKIN_SCORE_CODE = "checkin_score" //签到积分配置key
)

//签到积分配置
type QDConf struct {
	Day   int
	Score int
}

func (t_my *TaskThread) CheckInLog(ctx *gin.Context) (int, string, string, int, int, []map[string]string) {
	t_status := 100
	t_msg := "错误"

	checkin_today := "0"                     //今日签到
	checkin_continuous := 0                  //info 连续签到天数
	checkin_total_score := 0                 //info 总积分
	checkin_history := []map[string]string{} //签到历史

	var m_tl model.TaskList
	tast_type := m_tl.RedisGetOne("tast_type", CHECKIN_TYPE)
	task_id, ok := tast_type["id"]
	if ok == false {
		t_msg = "没有签到任务"
		return t_status, t_msg, checkin_today, checkin_continuous, checkin_total_score, checkin_history
	}

	if tast_status, _ := tast_type["status"]; tast_status == "-1" {
		t_msg = "签到活动已下架"
		return t_status, t_msg, checkin_today, checkin_continuous, checkin_total_score, checkin_history
	}

	var m_webconf model.WebConf
	conf := m_webconf.ConfRedis(CHECKIN_SCORE_CODE)

	conf_arr := []QDConf{}
	err := json.Unmarshal([]byte(conf["val"]), &conf_arr)
	if err != nil {
		t_msg = "签到配置有误"
		return t_status, t_msg, checkin_today, checkin_continuous, checkin_total_score, checkin_history
	}

	account := sessInfo("account", ctx)
	tab_name := "task_log"

	time_now := time.Now().Format(f_day)
	monday := GetFirstDateOfWeek()

	fields := []string{"score", "create_time"}
	task_w := "task_id = ? AND account = ? AND create_time >= ?"

	task_logs, _ := model.Rows(tab_name, "create_time ASC", 100, 1, fields, task_w, task_id, account, monday)
	task_logs_arr := TaskLog2Map4Date(task_logs)

	week := common.DateDiff(monday, f_day, 6)           //得到一周日期
	w_int := common.DifferDays(monday, time_now, f_day) //假设今天星期1返回0,星期3返回2,星期7返回6

	week_day := 1 //一周中的第几天
	score_symbol := "+"
	for w_k, what_week := range week {

		history := map[string]string{}
		history["date"] = what_week

		wqd_score := "0" //未签到 展示签到可以领取的积分
		for _, row := range conf_arr {
			if row.Day == week_day {
				wqd_score = score_symbol + strconv.Itoa(row.Score)
			}
		}

		//默认未签到
		history["is_checkin"] = "0"
		history["score"] = wqd_score

		if score, ok := task_logs_arr[what_week]; ok { //已签到

			history["is_checkin"] = "1"
			history["score"] = score_symbol + score
			checkin_continuous++

			if what_week == time_now {
				checkin_today = "1"
			}

		} else if w_k < w_int { //循环到比当天小的日期没签到,当天没签到中断签到

			checkin_continuous = 0
		}

		week_day++
		checkin_history = append(checkin_history, history)
	}

	score_sql := fmt.Sprintf("SELECT SUM(score) as t_score FROM `task_log` WHERE task_id = '%s' AND account = '%s'", task_id, account)
	score_query, _ := model.SqlRow(score_sql, []string{"t_score"})
	if t_score, ok2 := score_query["t_score"]; ok2 {
		checkin_total_score, _ = common.Str2Int(t_score)
	}

	t_status = STATUS_SUCCESS
	t_msg = MSG_SUCCESS

	return t_status, t_msg, checkin_today, checkin_continuous, checkin_total_score, checkin_history
}

func (t_my *TaskThread) CheckinClick(ctx *gin.Context) (int, string, int) {

	t_status, t_msg, checkin_today, checkin_continuous, _, _ := t_my.CheckInLog(ctx)
	if t_status != STATUS_SUCCESS {
		return t_status, t_msg, 0
	}

	if checkin_today == "1" {
		return 144, "今日已签到", 0
	}

	score := -100
	qd_total := checkin_continuous + 1 //历史连续签到+今天
	// qd_total_str := fmt.Sprintf("%d", qd_total)

	var m_webconf model.WebConf
	conf := m_webconf.ConfRedis(CHECKIN_SCORE_CODE)

	conf_arr := []QDConf{}
	err := json.Unmarshal([]byte(conf["val"]), &conf_arr)
	if err != nil {
		return 145, "签到配置有误", 0
	}

	for _, row := range conf_arr {
		if row.Day == qd_total {
			score = row.Score
		}
	}

	if score == -100 {
		return 145, "签到配置有误", 0
	}

	var m_tl model.TaskList
	tast_type := m_tl.RedisGetOne("tast_type", CHECKIN_TYPE)
	task_id, ok := tast_type["id"]
	if ok == false {
		return 146, "没有签到任务", 0
	}

	is_task := "1"
	account := sessInfo("account", ctx)
	order_type := "3"
	note := "签到奖励"
	act_id, _ := strconv.Atoi((tast_type["act_id"]))
	c_ok := UpUserScore(score, act_id, account, order_type, note, is_task, task_id)
	if c_ok == false {
		return 146, "签到失败", 0
	}

	return STATUS_SUCCESS, MSG_SUCCESS, score
}

func (t_my *TaskThread) TaskList(ctx *gin.Context) (int, string, []map[string]interface{}) {

	t_status := STATUS_SUCCESS
	t_msg := MSG_SUCCESS
	t_data := []map[string]interface{}{}

	account := sessInfo("account", ctx)

	table_tali := "task_list"
	fields_tali := []string{"id", "title", "content", "score", "icon", "act_id", "tast_type", "app_url_code"}
	where_tali := "status = 1 AND tast_type != 5" //状态正常 & tast_type != 签到
	task_list, _ := model.Rows(table_tali, "", 100, 1, fields_tali, where_tali)

	var m_tasklog model.TaskLog
	var m_appurl model.AppUrl

	for _, row := range task_list {

		is_click := "1" //未完成任务
		task_log := m_tasklog.RedisByAccount(account, "task_id", row["id"])
		if _, ok := task_log["id"]; ok {
			is_click = "0" //已完成任务
		}

		ios_url := map[string]string{}
		android_url := map[string]string{}
		if app_url_code, ok1 := row["app_url_code"]; ok1 {
			ios_url = m_appurl.RedisByCodeApptype(app_url_code, model.IOS_APP_TYPE)
			android_url = m_appurl.RedisByCodeApptype(app_url_code, model.ANDROID_APP_TYPE)
		}
		// t_url := map[string]string{
		// 	"ios_url":     ios_url,
		// 	"android_url": android_url,
		// }
		t_url := make([]map[string]string, 2)
		t_url[0] = ios_url
		t_url[1] = android_url

		t_row := map[string]interface{}{}
		t_row["url"] = t_url
		t_row["is_click"] = is_click

		for k, v := range row {
			t_row[k] = v
		}

		t_data = append(t_data, t_row)
	}

	return t_status, t_msg, t_data
}

func (t_my *TaskThread) TaskClick(ctx *gin.Context, id string) (int, string) {

	t_status := 166
	t_msg := "错误"

	if id == "" {
		t_msg = "id不能为空"
		return t_status, t_msg
	}

	var m_tasklist model.TaskList
	var m_tasklog model.TaskLog

	task_one := m_tasklist.RedisGetOne("id", id)
	if r_id, _ := task_one["id"]; r_id == "" {
		t_msg = "错误的id"
		return t_status, t_msg
	}

	account := sessInfo("account", ctx)

	task_log := m_tasklog.RedisByAccount(account, "task_id", id)
	if _, ok := task_log["id"]; ok {
		t_msg = "已完成任务"
		return t_status, t_msg
	}

	app_url_code := task_one["app_url_code"]
	if app_url_code == "" {
		t_msg = "无法完成任务"
		return t_status, t_msg
	}

	if app_url_code == "bind_phone" { //绑定手机

		uinfo_sql := "SELECT id, phone FROM `user_info` WHERE account = '" + account + "' LIMIT 1;"
		uinfo_query, _ := model.SqlRow(uinfo_sql, []string{"id", "phone"})
		if phone, _ := uinfo_query["phone"]; phone == "" {
			t_msg = "用户未绑定手机号"
			return t_status, t_msg
		}
	}

	if app_url_code == "watch_video" { //观看视频

		video_history, _ := model.Rows("video_history", "", 1, 1, []string{"id"}, "account=?", account)
		if len(video_history) == 0 {
			t_msg = "用户未观看视频"
			return t_status, t_msg
		}
	}

	if app_url_code == "first_recharge" { //充值1元

		amo_w := "account=? AND order_type = 1 AND amount >= 1"
		amount_list, _ := model.Rows("amount_list", "", 1, 1, []string{"id"}, amo_w, account)
		if len(amount_list) == 0 {
			t_msg = "用户未充值"
			return t_status, t_msg
		}
	}

	score_int, _ := common.Str2Int(task_one["score"])

	is_task := "1"
	order_type := "3"
	note := task_one["title"]
	act_id, _ := strconv.Atoi((task_one["act_id"]))
	c_ok := UpUserScore(score_int, act_id, account, order_type, note, is_task, id)
	if c_ok == false {
		t_msg = "任务完成失败"
		return t_status, t_msg
	}

	t_status = STATUS_SUCCESS
	t_msg = MSG_SUCCESS
	return t_status, t_msg
}
