package controller

import (
	"qzapp/hook"

	"github.com/gin-gonic/gin"
)

/**
* 视频推荐
 */
// func (u *UserController) RecommVideo(c *gin.Context) {
// 	c_status := 100
// 	c_msg := "请求成功"
// 	data := map[string]interface{}{}
// 	key_arr := []string{"vid", "page", "page_size"}
// 	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
// 	if c_status == 200 {
// 		c_status, data["total"], c_msg, data["video_info"] = user.RecommVideo(in_param)
// 	}
// 	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
// }

/**
* 签到日志
 */
func (c_task *TaskController) CheckInLog(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	c_status, c_msg, c_data["checkin_today"], c_data["checkin_continuous"], c_data["checkin_total_score"], c_data["checkin_history"] = t_task.CheckInLog(c)

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 点击签到
 */
func (c_task *TaskController) CheckinClick(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	c_status, c_msg, c_data["score"] = t_task.CheckinClick(c)

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 任务列表
 */
func (c_task *TaskController) TaskList(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	c_status, c_msg, c_data["list"] = t_task.TaskList(c)

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 完成任务
 */
func (c_task *TaskController) TaskClick(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	key_arr := []string{"id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = t_task.TaskClick(c, in_param["id"])
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}
