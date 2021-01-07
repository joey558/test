package hook

import (
	"fmt"
	"qzapp/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/**
*  是否登录
 */
func AuthLogin(ctx *gin.Context) (int, string) {
	h_status := 600
	h_msg := "未登录"
	session := sessions.Default(ctx)
	session_id := session.Get("session_id")

	if session_id == nil {
		return h_status, h_msg
	}

	var model_ul model.UserList
	//判断唯一登录
	sess_id := fmt.Sprintf("%v", session_id)
	a_info := model_ul.UserBySessRedis(sess_id)

	if len(a_info["account"]) < 1 {
		h_msg = "账号已经在其他地方登录"
		return h_status, h_msg
	}
	h_status = 200
	h_msg = "success"

	return h_status, h_msg
}
