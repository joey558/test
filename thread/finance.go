package thread

import (
	"encoding/json"
	"fmt"
	"qzapp/hook"
	"qzapp/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.stagingvip.net/publicGroup/public/common"
)

/**
*  获取支付渠道列表
 */
func (ft *Finance) PayList(ctx *gin.Context) (int, string, []map[string]interface{}) {
	t_status := 100
	t_msg := "用户信息异常"
	pay_list := []map[string]interface{}{}
	account := sessInfo("account", ctx)
	if account == "" {
		return t_status, t_msg, pay_list
	}
	t_status = 200
	t_msg = "success"
	p_where := map[string]interface{}{}
	p_where["status"] = 1
	table_name := "pay_config"
	fields := []string{"title", "id", "sort", "class_code", "min_amount", "max_amount", "is_fix", "fix_amount", "note"}
	order_by := "sort asc"
	page_size := 100
	offset := 0
	p_list, _ := model.PageList(table_name, order_by, page_size, offset, fields, p_where)
	if len(p_list) < 1 {
		return t_status, t_msg, pay_list
	}
	list_map := map[string][]map[string]string{}
	for _, p_val := range p_list {
		list_map[p_val["class_code"]] = append(list_map[p_val["class_code"]][0:], p_val)
	}

	var model_pc model.PayClass
	for l_k, l_v := range list_map {
		pc_info := model_pc.ClassRedis(l_k)
		if len(pc_info) < 1 || len(pc_info["code"]) < 1 {
			continue
		}
		p_map := map[string]interface{}{}
		p_map["title"] = pc_info["title"]
		p_map["code"] = pc_info["code"]
		p_map["icon"] = pc_info["icon"]
		p_map["sort"] = pc_info["sort"]
		p_map["list"] = l_v
		pay_list = append(pay_list[0:], p_map)
	}

	common.SortMapInterfaceInt(pay_list, "sort")

	return t_status, t_msg, pay_list
}

/**
*  支付接口
 */
func (ft *Finance) Pay(in_param map[string]string, ctx *gin.Context) (int, string, map[string]string) {
	t_status := 100
	t_msg := "参数不足"
	pay_map := map[string]string{}
	//判断额度是否正常
	if len(in_param["amount"]) < 1 || len(in_param["pay_id"]) < 1 {
		return t_status, t_msg, pay_map
	}
	amount, _ := strconv.ParseFloat(in_param["amount"], 64)
	if amount < 1.00 {
		t_msg = "额度错误"
		return t_status, t_msg, pay_map
	}
	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, pay_map
	}
	//查询配置信息
	var model_pc model.PayConfig
	p_conf := model_pc.PayConf(in_param["pay_id"])
	if p_conf.Id < 1 || p_conf.Status != 1 {
		t_msg = "支付信息错误"
		return t_status, t_msg, pay_map
	}
	if amount < p_conf.Min_amount && p_conf.Is_fix == 0 {
		t_msg = "付款额度不能小于" + fmt.Sprintf("%.2f", p_conf.Min_amount)
		return t_status, t_msg, pay_map
	}
	if amount > p_conf.Max_amount && p_conf.Is_fix == 0 {
		t_msg = "付款额度不能大于" + fmt.Sprintf("%.2f", p_conf.Max_amount)
		return t_status, t_msg, pay_map
	}

	var model_cl model.ConfList
	c_id := fmt.Sprintf("%d", p_conf.Conf_id)
	conf_info := model_cl.ConfRedis(c_id)
	if len(conf_info) < 1 || len(conf_info["id"]) < 1 {
		t_msg = "配置信息错误"
		return t_status, t_msg, pay_map
	}

	var p_c map[string]string
	err := json.Unmarshal([]byte(conf_info["conf"]), &p_c)
	if err != nil {
		t_msg = "配置信息异常"
		return t_status, t_msg, pay_map
	}

	var model_ul model.UserList
	u_info := model_ul.UserRedis(account)
	table_name := "order_list"
	order_number := model.GetKey(16)
	create_time := time.Now().Format(f_date)
	o_data := map[string]string{}
	o_data["id"] = order_number
	o_data["amount"] = in_param["amount"]
	o_data["account"] = account
	o_data["order_number"] = order_number
	o_data["agent_name"] = u_info["agent_name"]
	o_data["is_test"] = u_info["is_test"]
	o_data["agent_path"] = u_info["agent_path"]
	o_data["pay_id"] = in_param["pay_id"]
	o_data["create_time"] = create_time
	o_data["creator"] = account
	o_sql := common.InsertSql(table_name, o_data)
	err = model.Query(o_sql)
	if err != nil {
		t_msg = "订单创建失败"
		return t_status, t_msg, pay_map
	}

	t_status = 200
	t_msg = "success"

	ip := ctx.ClientIP()
	encode_str := fmt.Sprintf(`{"is_mobile":"1","amount":"%s","pay_id":"%s","order_number":"%s","class_code":"%s","bank_code":"%s","push_url":"%s","ip":"%s"}`, in_param["amount"], in_param["pay_id"], order_number, p_conf.Class_code, p_conf.Bank_code, p_c["push_url"], ip)
	aes := common.SetAES(p_c["private_key"], "", "", 16)
	aes_res := aes.AesEncryptString(encode_str)
	aes_res = strings.Replace(aes_res, "+", "%2B", -1)

	pay_map["mer_code"] = p_c["mer_code"]
	pay_map["pay_data"] = aes_res
	pay_map["pay_url"] = p_c["pay_url"]
	pay_map["method"] = p_c["method"]

	return t_status, t_msg, pay_map
}

/**
*  购买积分
 */
func (ft *Finance) Score(in_param map[string]string, ctx *gin.Context) (int, int, string) {
	t_status := 100
	t_msg := "购买的积分错误"
	score := 0

	buy_score, _ := strconv.Atoi(in_param["score"])
	if buy_score < 1 {
		return t_status, score, t_msg
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, score, t_msg
	}

	rate := hook.WebConfVal("rate")
	if rate == "" {
		rate = "10"
	}
	rate_int, _ := strconv.Atoi(rate)

	//购买的积分需要的额度
	amount := float64(buy_score) / float64(rate_int)
	if amount < 0.01 {
		t_msg = "购买的积分数量太少"
		return t_status, score, t_msg
	}

	var model_ui model.UserInfo
	u_info := model_ui.User(account)
	if amount > u_info.Amount {
		t_msg = "额度不足"
		return t_status, score, t_msg
	}

	var model_ul model.UserList
	u_list := model_ul.UserRedis(account)

	order_number := model.GetKey(16)
	create_time := time.Now().Format(f_date)

	//生成订单
	o_table := "order_list"
	o_data := map[string]string{}
	o_data["id"] = order_number
	o_data["amount"] = fmt.Sprintf("%.2f", amount)
	o_data["account"] = account
	o_data["status"] = "3"
	o_data["order_type"] = "5"
	o_data["order_number"] = order_number
	o_data["agent_name"] = u_list["agent_name"]
	o_data["is_test"] = u_list["is_test"]
	o_data["agent_path"] = u_list["agent_path"]
	o_data["score"] = in_param["score"]
	o_data["create_time"] = create_time
	o_data["creator"] = account
	o_data["pay_time"] = create_time
	o_sql := common.InsertSql(o_table, o_data)

	//生成账变记录
	after_amount := u_info.Amount - amount
	after_score := u_info.Score + buy_score
	a_table := "amount_list"
	a_data := map[string]string{}
	a_data["id"] = order_number
	a_data["amount"] = fmt.Sprintf("%.2f", amount)
	a_data["score"] = in_param["score"]
	a_data["account"] = account
	a_data["order_type"] = "5"
	a_data["order_number"] = order_number
	a_data["agent_name"] = u_list["agent_name"]
	a_data["is_test"] = u_list["is_test"]
	a_data["agent_path"] = u_list["agent_path"]
	a_data["before_amount"] = fmt.Sprintf("%.2f", u_info.Amount)
	a_data["create_time"] = create_time
	a_data["after_amount"] = fmt.Sprintf("%.2f", after_amount)
	a_data["before_score"] = fmt.Sprintf("%d", u_info.Score)
	a_data["after_score"] = fmt.Sprintf("%d", after_score)
	a_sql := common.InsertSql(a_table, a_data)

	//更新用户的积分信息
	up_sql := fmt.Sprintf("update user_info set amount=amount-%.2f,score=score+%d where account='%s';", amount, buy_score, account)

	sql_arr := []string{o_sql, a_sql, up_sql}

	err := model.Trans(sql_arr)
	if err != nil {
		t_msg = "积分购买失败"
		return t_status, score, t_msg
	}
	score = after_score
	t_status = 200
	t_msg = "success"

	return t_status, score, t_msg
}
