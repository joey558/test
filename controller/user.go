package controller

import (
	"qzapp/hook"

	"github.com/gin-gonic/gin"
)

/**
* 查询用户信息
 */
func (u *UserController) UserInfo(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	c_status, c_msg, d["user_info"] = user.UserInfo(c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
* 查询视频信息
 */
func (u *UserController) VideoInfo(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"vid"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["video_info"] = user.VideoInfo(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 视频推荐
 */
func (u *UserController) RecommVideo(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"vid", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["video_info"] = user.RecommVideo(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 视频评论
 */
func (u *UserController) VideoComm(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"vid", "p_id", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["comm_list"] = user.VideoComm(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 评论视频
 */
func (u *UserController) Comm(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"vid", "p_id", "content"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = user.Comm(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 点赞/取消点赞评论,以及获取点赞总次数
 */
func (u *UserController) LikeComm(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"comm_id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["like_comm"] = user.LikeComm(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 点赞视频,以及获取点赞总次数
 */
func (u *UserController) LikeVideo(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"vid"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["like_video"] = user.LikeVideo(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 收藏视频
 */
func (u *UserController) StarVideo(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"vid"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["is_star"], c_msg = user.StarVideo(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 查看自己的所有评论
 */
func (u *UserController) UserComm(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["comm_list"] = user.UserComm(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 用户的额度和积分
 */
func (u *UserController) UserMoney(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	c_status, c_msg, data["amount"], data["score"] = user.UserMoney(c)

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 上传图片
 */
func (u *UserController) UserAvatar(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"avatar"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["avatar"] = user.UserAvatar(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 修改昵称
 */
func (u *UserController) UserNickName(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"nick_name"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = user.UserNickName(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 修改年龄
 */
func (u *UserController) UserAge(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"age"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = user.UserAge(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 修改性别
 */
func (u *UserController) UserSex(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"sex"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = user.UserSex(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 修改个性签名
 */
func (u *UserController) UserPerSign(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"per_sign"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = user.UserPerSign(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 修改密码
 */
func (u *UserController) UserEditPwd(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"pwd", "old_pwd"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = user.UserEditPwd(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 获取茄子钱包存款、优惠、消费记录
 */
func (u *UserController) QzWallet(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"order_type", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)

	if c_status == 200 {
		c_status, data["total"], c_msg, data["wallet_info"] = user.QzWallet(in_param, c)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 获取kok币的购买、优惠、消费记录
 */
func (u *UserController) KokGold(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"order_type", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)

	if c_status == 200 {
		c_status, data["total"], c_msg, data["gold_info"] = user.KokGold(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 查看自己的收藏
 */
func (u *UserController) UserCollect(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["collect_list"] = user.UserCollect(in_param, c)

	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 清空收藏
 */
func (u *UserController) ClearCollect(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	c_status, c_msg = user.ClearCollect(c)

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 查看观看历史
 */
func (u *UserController) UserHistory(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["history_list"] = user.UserHistory(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 查看消息通知
 */
func (u *UserController) SysMsg(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"msg_type", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["msg_list"] = user.SysMsg(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 把此消息标记为已读
 */
func (u *UserController) ViewMsg(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"msg_id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = user.ViewMsg(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 点赞/取消点赞茄子主推视频
 */
func (u *UserController) LikeTheme(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["like_list"] = user.LikeTheme(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 查看自己是否点赞茄子主推视频
 */
func (u *UserController) SeeLike(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["like_info"] = user.SeeLike(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 删除指定的收藏视频
 */
func (u *UserController) DelCollect(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	video_id := c.PostFormArray("video_id[]")
	c_status, c_msg = hook.AuthInputForArr(video_id)
	if c_status == 200 {
		c_status, c_msg = user.DelCollect(video_id, c)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 删除指定的观看历史
 */
func (u *UserController) DelHistory(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	video_id := c.PostFormArray("video_id[]")
	c_status, c_msg = hook.AuthInputForArr(video_id)
	if c_status == 200 {
		c_status, c_msg = user.DelHistory(video_id, c)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 清空观看历史
 */
func (u *UserController) ClearHistory(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	c_status, c_msg = user.ClearHistory(c)

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 删除指定消息中心消息
 */
func (u *UserController) DelMsg(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	msg_id := c.PostFormArray("msg_id[]")
	c_status, c_msg = hook.AuthInputForArr(msg_id)
	if c_status == 200 {
		c_status, c_msg = user.DelMsg(msg_id, c)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 查看未读消息数量
 */
func (u *UserController) UnreadNum(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	c_status, c_msg, data["unread_info"] = user.UnreadNum(c)

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 绑定手机-发送验证码
 */
func (u *UserController) PhoneVerify(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	key_arr := []string{"phone"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["verify_code"] = user.PhoneVerify(in_param, c)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 绑定手机-绑定
 */
func (u *UserController) PhoneBind(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	key_arr := []string{"code"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = user.PhoneBind(in_param, c)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
* 换绑手机
 */
func (u *UserController) PhoneChg(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	key_arr := []string{"phone"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["verify_code"] = user.PhoneChg(in_param, c)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}
