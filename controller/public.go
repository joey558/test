package controller

import (
	"qzapp/hook"
	"qzapp/thread"

	"github.com/gin-gonic/gin"
)

/**
* 会员注册
 */
func (pb *PublicController) Register(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}

	//接收值
	key_arr := []string{"account", "pwd", "reg_code"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = pub.Register(in_param, c)
	}

	account := in_param["account"]
	title := "用户注册"
	content := "用户注册的用户名:" + account
	go thread.UserLog(c_status, account, title, content, c_msg, in_param, c)
	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
* 判断用户名是否可用
 */
func (pb *PublicController) AcountAllow(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}

	//接收值
	key_arr := []string{"account"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = pub.AcountAllow(in_param)
	}

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
* 判断密码是否可用
 */
func (pb *PublicController) PwdAllow(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}

	//接收值
	key_arr := []string{"pwd"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = pub.PwdAllow(in_param)
	}

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
* 会员登录
 */
func (pb *PublicController) Login(c *gin.Context) {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	key_arr := []string{"account", "pwd"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = pub.Login(in_param, c)
	}

	account := in_param["account"]
	title := "用户登录"
	content := "登录的用户名:" + account
	go thread.UserLog(c_status, account, title, content, c_msg, in_param, c)

	//将数据装载到json返回值
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: d})
}

/**
*  退出
 */
func (pb *PublicController) Logout(c *gin.Context) {
	data := map[string]interface{}{}
	c_status, c_msg := pub.Logout(c)
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  获取启动页
 */
func (pb *PublicController) Loading(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	c_status, c_msg, data["load_list"] = pub.Loading()
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  获取轮播图
 */
func (pb *PublicController) Banner(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	c_status, c_msg, data["banner_list"] = pub.Banner()
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  弹窗公告
 */
func (pb *PublicController) Pop(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	c_status, c_msg, data["pop_list"] = pub.Pop()
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  视频分类
 */
func (pb *PublicController) TagType(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"is_top"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["type_list"] = pub.TagType(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  推荐/热门视频列表
 */
func (pb *PublicController) HotVideo(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"is_hot", "is_recomm", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["video_list"] = pub.HotVideo(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  分类视频列表
 */
func (pb *PublicController) TypeVideo(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"type_id", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["video_list"] = pub.TypeVideo(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  分类视频列表
 */
func (pb *PublicController) TagList(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["tag_list"] = pub.TagList(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  标签视频列表
 */
func (pb *PublicController) TagVideo(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"tag_id", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["video_list"] = pub.TagVideo(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  女优资料;把女优按照字母排序返回
 */
func (pb *PublicController) ActorData(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	c_status, c_msg, data["list"] = pub.ActorData()
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  女优视频列表
 */
func (pb *PublicController) ActorVideo(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"actor_id", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["video_list"] = pub.ActorVideo(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  女优列表
 */
func (pb *PublicController) ActorList(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"is_hot", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["actor_list"] = pub.ActorList(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  游戏平台类型
 */
func (pb *PublicController) PlatType(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	c_status, c_msg, data["plat_type"] = pub.PlatType()
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  游戏平台列表
 */
func (pb *PublicController) PlatList(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"type_code", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["plat_list"] = pub.PlatList(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  根据关键字搜索
 */
func (pb *PublicController) KwSearch(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"key_word", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["video_list"] = pub.KwSearch(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  获取优惠类型
 */
func (pb *PublicController) ActiveType(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	c_status, c_msg, data["type_list"] = pub.ActiveType()
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  根据关键字搜索
 */
func (pb *PublicController) ActiveList(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"type_id", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["active_list"] = pub.ActiveList(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  获取客服信息
 */
func (pb *PublicController) CsInfo(c *gin.Context) {
	data := map[string]interface{}{}
	c_status := 200
	c_msg := "success"
	data["cs_info"] = pub.CsInfo()
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  app的版本信息
 */
func (pb *PublicController) AppVer(c *gin.Context) {
	data := map[string]interface{}{}
	c_status := 100
	c_msg := "请求成功"
	c_status, c_msg, data["app_ver"] = pub.AppVer(c)
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  app的启动日志
 */
func (pb *PublicController) OpenApp(c *gin.Context) {
	data := map[string]interface{}{}
	c_status := 100
	c_msg := "请求成功"
	key_arr := []string{"phone_info", "phone_os", "uid", "phone_number", "account"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		go pub.OpenApp(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  app的关闭日志
 */
func (pb *PublicController) CloseApp(c *gin.Context) {
	data := map[string]interface{}{}
	c_status := 100
	c_msg := "请求成功"
	key_arr := []string{"uid"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		go pub.CloseApp(in_param, c)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  购买积分的比例
 */
func (pb *PublicController) Rate(c *gin.Context) {
	c_status := 200
	c_msg := "success"
	data := map[string]interface{}{}
	data["rate"] = pub.Rate()
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  获取茄子主推视频
 */
func (pb *PublicController) QzRecomm(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, data["total"], c_msg, data["recomm_list"] = pub.QzRecomm(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  茄子主推视频的观看次数+1
 */
func (pb *PublicController) AddView(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = pub.AddView(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  影片的观看次数+1,以及获取影片的总观看次数
 */
func (pb *PublicController) AddSee(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}
	key_arr := []string{"vid"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, data["view_num"] = pub.AddSee(in_param)
	}
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}

/**
*  首页-标题栏,每个标题栏下面取4个视频
 */
func (pb *PublicController) HomeView(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	data := map[string]interface{}{}

	c_status, c_msg, data["list"] = pub.HomeView()

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: data})
}
