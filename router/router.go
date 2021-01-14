package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"qzapp/controller"
	"qzapp/hook"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"qzapp/common"
)

type JsonOut struct {
	Status int
	Msg    string
	Data   map[string]interface{}
}

/**
*  判断是否登录
 */
func LoginAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//判断是否登录
		r_status, r_msg := hook.AuthLogin(ctx)
		if r_status != 200 {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{"Status": r_status, "Msg": r_msg})
		}
		r_status, r_msg = hook.AppType(ctx)
		if r_status != 200 {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{"Status": r_status, "Msg": r_msg})
		}
	}
}

var Router *gin.Engine

func init() {
	Router = gin.New()

	file_name := "./conf/redis.json"

	conf_byte, err := common.ReadFile(file_name)
	if err != nil {
		panic(err)
	}
	var json_conf map[string]string
	//解析json格式
	err = json.Unmarshal(conf_byte, &json_conf)
	if err != nil {
		panic(err)
	}
	max_idle, _ := strconv.Atoi(json_conf["max_idle"])
	if max_idle < 1 {
		max_idle = 10
	}

	addr := fmt.Sprintf("%s:%s", json_conf["host"], json_conf["port"])

	store, err := redis.NewStore(max_idle, "tcp", addr, json_conf["auth"], []byte(json_conf["pre_key"]))

	Router.Use(sessions.Sessions(json_conf["pre_key"], store))

	//静态文件路径，一定需要
	Router.LoadHTMLGlob("view/*")
	Router.LoadHTMLFiles("./view/index.tpl")
	Router.GET("/test/test.do", controller.Test)

	//不需要登录的接口接口
	var pub controller.PublicController
	pub_group := Router.Group("/public")
	//注册的接口
	pub_group.POST("/reg.do", pub.Register)
	//登录的接口
	pub_group.POST("/login.do", pub.Login)
	//退出的接口
	pub_group.POST("/logout.do", pub.Logout)
	//判断用户名是否可用
	pub_group.POST("/account_allow.do", pub.AcountAllow)
	//判断密码是否合法
	pub_group.POST("/pwd_allow.do", pub.PwdAllow)
	//app启动页
	pub_group.POST("/loading.do", pub.Loading)
	//app轮播图
	pub_group.POST("/banner.do", pub.Banner)
	//弹窗
	pub_group.POST("/pop.do", pub.Pop)
	//视频分类
	pub_group.POST("/tag_type.do", pub.TagType)
	//热门推荐视频
	pub_group.POST("/hot_video.do", pub.HotVideo)
	//分类视频
	pub_group.POST("/type_video.do", pub.TypeVideo)
	//标签视频
	pub_group.POST("/tag_video.do", pub.TagVideo)
	//女优列表
	pub_group.POST("/actor_list.do", pub.ActorList)
	//女优资料;把女优按照字母排序返回
	pub_group.POST("/actor_data.do", pub.ActorData)
	//女优资料详情;通过女优id查询视频(根据观看次数排序)
	pub_group.POST("/actor_video.do", pub.ActorVideo)
	//平台类型列表
	pub_group.POST("/plat_type.do", pub.PlatType)
	//平台列表
	pub_group.POST("/plat_list.do", pub.PlatList)
	//关键字搜索
	pub_group.POST("/kw_search.do", pub.KwSearch)
	//优惠类型
	pub_group.POST("/active_type.do", pub.ActiveType)
	//优惠列表
	pub_group.POST("/active_list.do", pub.ActiveList)
	//客服列表
	pub_group.POST("/cs_info.do", pub.CsInfo)
	//app的版本信息
	pub_group.POST("/app_ver.do", pub.AppVer)
	//每次启动APP的时候，请求的接口
	pub_group.POST("/open_app.do", pub.OpenApp)
	//每次关闭APP的时候，请求的接口
	pub_group.POST("/close_app.do", pub.CloseApp)
	//标签列表
	pub_group.POST("/tag_list.do", pub.TagList)
	//购买积分的比例
	pub_group.POST("/rate.do", pub.Rate)
	//获取茄子主推视频
	pub_group.POST("/qz_recomm.do", pub.QzRecomm)
	//茄子主推短视频的观看次数+1
	pub_group.POST("/add_view.do", pub.AddView)
	//影片的观看次数+1,以及获取影片的总观看次数
	pub_group.POST("/add_see.do", pub.AddSee)

	pub_group.POST("/home_view.do", pub.HomeView) // 首页-标题栏,每个标题栏下面取4个视频

	//====================用户接口==================//
	//用户相关的接口
	var user controller.UserController
	user_group := Router.Group("/user", LoginAuth())
	//查询用户信息的接口
	user_group.POST("/user_info.do", user.UserInfo)
	//通过视频ID查询视频的详细信息
	user_group.POST("/video_info.do", user.VideoInfo)
	//视频详情页的推荐视频
	user_group.POST("/recomm_video.do", user.RecommVideo)
	//视频详情页的评论
	user_group.POST("/video_comm.do", user.VideoComm)
	//评论
	user_group.POST("/comm.do", user.Comm)
	//点赞/取消点赞视频,以及获取点赞总次数
	user_group.POST("/like_video.do", user.LikeVideo)
	//点赞/取消点赞评论,以及获取点赞总次数
	user_group.POST("/like_comm.do", user.LikeComm)
	//收藏视频
	user_group.POST("/star_video.do", user.StarVideo)
	//查看自己的评论
	user_group.POST("/user_comm.do", user.UserComm)
	//获取余额&积分
	user_group.POST("/user_money.do", user.UserMoney)
	//上传头像
	user_group.POST("/user_avatar.do", user.UserAvatar)
	//修改昵称
	user_group.POST("/user_nickname.do", user.UserNickName)
	//修改年龄
	user_group.POST("/user_age.do", user.UserAge)
	//修改性别
	user_group.POST("/user_sex.do", user.UserSex)
	//修改个性签名
	user_group.POST("/user_persign.do", user.UserPerSign)
	//修改密码
	user_group.POST("/user_editpwd.do", user.UserEditPwd)
	//获取茄子钱包存款、优惠、消费记录
	user_group.POST("/qz_wallet.do", user.QzWallet)
	//获取kok币的购买、优惠、消费记录
	user_group.POST("/kok_gold.do", user.KokGold)
	//查看自己的收藏
	user_group.POST("/user_collect.do", user.UserCollect)
	//删除指定的收藏视频
	user_group.POST("/del_collect.do", user.DelCollect)
	//清空收藏
	user_group.POST("/clear_collect.do", user.ClearCollect)
	//查看观看历史
	user_group.POST("/user_history.do", user.UserHistory)
	//删除指定的观看历史
	user_group.POST("/del_history.do", user.DelHistory)
	//清空观看历史
	user_group.POST("/clear_history.do", user.ClearHistory)
	//查看未读消息数量
	user_group.POST("/unread_num.do", user.UnreadNum)
	//查看消息通知
	user_group.POST("/sys_msg.do", user.SysMsg)
	//删除指定消息中心消息
	user_group.POST("/del_msg.do", user.DelMsg)
	//把此消息标记为已读
	user_group.POST("/view_msg.do", user.ViewMsg)
	//用户点赞/取消点赞茄子主推视频以及查看视频点赞总次数
	user_group.POST("/like_theme.do", user.LikeTheme)
	//查看自己是否点赞茄子主推视频以及查看视频点赞总次数
	user_group.POST("/see_like.do", user.SeeLike)

	user_group.POST("/phone_verify.do", user.PhoneVerify) //获取初次绑定手机验证码/获取已绑手机验证码
	user_group.POST("/phone_bind.do", user.PhoneBind)     //校验手机验证码
	user_group.POST("/phone_chg.do", user.PhoneChg)       //获取换绑手机验证码

	//====================财务相关接口==================//
	//用户相关的接口
	var finance controller.FinanceController
	finance_group := Router.Group("/finance", LoginAuth())
	//获取支付方式列表的接口
	finance_group.POST("/pay_list.do", finance.PayList)
	//获取支付方式列表的接口
	finance_group.POST("/pay.do", finance.Pay)
	//购买积分
	finance_group.POST("/score.do", finance.Score)

	//====================任务&活动==================//
	var c_task controller.TaskController
	task_group := Router.Group("/task", LoginAuth())
	task_group.POST("/checkin_log.do", c_task.CheckInLog)     //签到日志
	task_group.POST("/checkin_click.do", c_task.CheckinClick) //点击签到
	task_group.POST("/task_list.do", c_task.TaskList)         //任务列表
	task_group.POST("/task_click.do", c_task.TaskClick)       //完成任务

	//====================广场==================//
	var c_blog controller.BlogController
	blog_group := Router.Group("/blog", LoginAuth())
	blog_group.POST("/blog_list.do", c_blog.BlogList)            //获取博客列表
	blog_group.POST("/blog_comm.do", c_blog.BlogComm)            //获取博客列表的全部评论
	blog_group.POST("/blog_comm_reply.do", c_blog.BlogCommReply) //评论博客,或者回复评论
	blog_group.POST("/blog_like.do", c_blog.BlogLike)            //博客 点赞/取消
	blog_group.POST("/comm_like.do", c_blog.CommLike)            //评论 点赞/取消
	// blog_group.POST("/upload_img_more.do", c_blog.UploadImgMore) //批量上传图片
	blog_group.POST("/blog_publish.do", c_blog.BlogPublish)  //发帖
	blog_group.POST("/blog_one.do", c_blog.BlogOne)          //获取某个博客信息
	blog_group.POST("/blog_comm_one.do", c_blog.BlogCommOne) //获取某个评论的信息
	blog_group.POST("/del_blog.do", c_blog.DelBlog)          //删除博客

}
