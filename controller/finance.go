package controller

import (
	"qzapp/hook"

	"github.com/gin-gonic/gin"
)

/**
* 获取支付渠道列表
 */
func (fc *FinanceController) PayList(c *gin.Context) {
	c_status := 100
	c_msg := "请求完成"
	//定义需要输出的结果
	data := map[string]interface{}{}

	//接收值
	c_status, c_msg, data["pay_list"] = finance.PayList(c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 支付接口
 */
func (fc *FinanceController) Pay(c *gin.Context) {
	//定义需要输出的结果
	data := map[string]interface{}{}

	key_arr := []string{"amount", "pay_id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["pay_map"] = finance.Pay(in_param, c)
	}
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 支付接口
 */
func (fc *FinanceController) Score(c *gin.Context) {
	//定义需要输出的结果
	data := map[string]interface{}{}

	key_arr := []string{"score"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total_score"], c_msg = finance.Score(in_param, c)
	}
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}
