package thread

import (
	"fmt"
	"qzapp/hook"

	"qzapp/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"qzapp/common"
	"qzapp/redis"
)

func (u *User) UserInfo(ctx *gin.Context) (int, string, map[string]string) {
	t_status := 100
	t_msg := "用户信息异常"
	u_info := map[string]string{}

	account := sessInfo("account", ctx)
	if account == "" {
		return t_status, t_msg, u_info
	}

	fields := []string{"id", "account", "nick_name", "avatar", "age", "sex", "per_sign", "phone"}
	s_sql := fmt.Sprintf("select id,account,nick_name,avatar,age,sex,per_sign,phone from user_info where account='%s' limit 1;", account)
	u_i, _ := model.SqlRow(s_sql, fields)

	if len(u_i) < 1 || len(u_i["account"]) < 1 {
		t_msg = "信息错误"
		return t_status, t_msg, u_info
	}

	u_info["id"] = u_i["id"]
	u_info["account"] = account
	u_info["nick_name"] = u_i["nick_name"]
	u_info["avatar"] = u_i["avatar"]
	u_info["age"] = u_i["age"]
	u_info["sex"] = "男"
	if u_i["sex"] == "1" {
		u_info["sex"] = "女"
	}
	u_info["note"] = u_i["per_sign"]
	u_info["phone"] = u_i["phone"]
	t_status = 200
	t_msg = "success"
	return t_status, t_msg, u_info
}

/**
*  视频详情
 */
func (u *User) VideoInfo(in_param map[string]string, ctx *gin.Context) (int, string, map[string]interface{}) {
	t_status := 100
	t_msg := "信息异常"
	v_info := map[string]interface{}{}
	account := sessInfo("account", ctx)
	if account == "" {
		return t_status, t_msg, v_info
	}
	if len(in_param) < 1 || len(in_param["vid"]) < 1 {
		t_msg = "视频ID不能为空"
		return t_status, t_msg, v_info
	}

	//查询视频信息
	v_field := []string{"id", "title", "content", "create_time", "h5_img", "url", "score", "like_num", "view_num", "star_num"}
	vl_sql := fmt.Sprintf("select id,title,content,create_time,h5_img,url,score,like_num,view_num,star_num from video_list where id='%s' and `status`='1' order by create_time desc limit 1;", in_param["vid"])
	v_map, _ := model.SqlRow(vl_sql, v_field)
	if len(v_map) < 1 || len(v_map["id"]) < 1 {
		t_msg = "视频信息异常"
		return t_status, t_msg, v_info
	}

	t_status = 200
	t_msg = "success"
	//查询自己是否点赞收藏
	v_info["is_like"] = 0
	v_info["is_star"] = 0
	var model_vs model.VideoStar
	var model_vl model.VideoLike
	is_like := model_vl.VideoLikeRedis(v_map["id"], account)
	if len(is_like) > 0 && len(is_like["id"]) > 5 {
		v_info["is_like"] = 1
	}
	is_star := model_vs.VideoStarRedis(v_map["id"], account)
	if len(is_star) > 0 && len(is_star["id"]) > 5 {
		v_info["is_star"] = 1
	}

	view_num := 0
	video_url := ""
	user_score, video_score := u.userViewNum(account, in_param["vid"], v_map["score"])
	if user_score >= video_score {

		if user_score == 0 && video_score == 0 {
			view_num = 0
		} else if user_score > 0 && video_score == 0 {
			view_num = 999
		} else {
			view_num = user_score / video_score
		}
		video_url = v_map["url"]
	}

	v_info["num"] = view_num
	v_map["url"] = video_url
	v_info["video_info"] = v_map

	//获取标签信息
	t_field := []string{"tid", "title"}
	tag_sql := fmt.Sprintf("SELECT tl.id as tid,tl.title as title FROM video_tag vt LEFT JOIN tag_list tl ON vt.tag_id = tl.id where vt.video_id='%s' order by tl.id asc;", in_param["vid"])
	tag_list, _ := model.SqlRows(tag_sql, t_field)
	v_info["tag_list"] = tag_list
	return t_status, t_msg, v_info
}

/**
*  获取推荐视频
 */
func (u *User) RecommVideo(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	video_list := []map[string]string{}
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size
	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM video_tag vt LEFT JOIN video_list vl ON vt.video_id=vl.id where vl.status=1 and vt.tag_id in (select tag_id from video_tag where video_id='%s') limit 1;", in_param["vid"])
	total_map, _ := model.SqlRow(count_sql, count_field)
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		return t_status, total, t_msg, video_list
	}
	uni_sql := fmt.Sprintf("SELECT vl.id as vid,vl.title as title,vl.view_num as view_num,vl.h5_img as h5_img FROM video_tag vt LEFT JOIN video_list vl ON vt.video_id=vl.id where vl.status=1 and vt.tag_id in (select tag_id from video_tag where video_id='%s') order by vl.view_num desc limit %d,%d;", in_param["vid"], offset, page_size)
	fields := []string{"vid", "title", "view_num", "h5_img"}
	video_list, _ = model.SqlRows(uni_sql, fields)
	return t_status, total, t_msg, video_list
}

/**
*  获取视频的评论
 */
func (u *User) VideoComm(in_param map[string]string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 100
	t_msg := "视频ID不能为空"
	total := 0
	comm_list := []map[string]string{}
	if len(in_param["vid"]) < 1 {
		return t_status, total, t_msg, comm_list
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, total, t_msg, comm_list
	}
	t_status = 200
	t_msg = "success"
	p_id := in_param["p_id"]
	if p_id == "" {
		p_id = "0"
	}
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size
	count_field := "count(0) as num"
	table_name := "video_comm"
	p_w := map[string]interface{}{}
	p_w["p_id"] = p_id
	p_w["video_id"] = in_param["vid"]
	total, _ = model.CommonTotal(table_name, count_field, p_w)
	if total <= offset {
		return t_status, total, t_msg, comm_list
	}

	fields := []string{"id", "account", "to_account", "content", "p_id", "nick_name", "to_name", "like_num", "reply_num", "create_time"}
	order_by := "like_num desc"
	c_list, _ := model.PageList(table_name, order_by, page_size, offset, fields, p_w)
	if len(c_list) < 1 {
		return t_status, total, t_msg, comm_list
	}
	var model_ui model.UserInfo
	var model_cl model.CommLike
	for _, c_val := range c_list {
		c_map := map[string]string{}
		u_info := model_ui.UserRedis(c_val["account"])
		c_map = c_val
		c_map["avatar"] = u_info["avatar"]
		c_map["is_like"] = "0"

		c_like := model_cl.CommLikeRedis(c_val["id"], account)
		if len(c_like) > 0 && len(c_like["id"]) > 1 {
			c_map["is_like"] = "1"
		}
		comm_list = append(comm_list[0:], c_map)
	}
	return t_status, total, t_msg, comm_list
}

/**
*  评论
 */
func (u *User) Comm(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "参数不足"

	if len(in_param["vid"]) < 1 || len(in_param["content"]) < 1 {
		return t_status, t_msg
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	var model_vl model.VideoList
	v_info := model_vl.VideoRedis(in_param["vid"])
	if len(v_info) < 1 || len(v_info["title"]) < 1 {
		t_msg = "视频信息异常"
		return t_status, t_msg
	}

	sql_arr := []string{}

	p_id := in_param["p_id"]
	if p_id == "" {
		p_id = "0"
	}

	var model_ui model.UserInfo
	user := model_ui.UserRedis(account)
	nick_name := user["nick_name"]

	to_account := ""
	to_name := ""
	id := model.GetKey(16)
	now_time := time.Now().Format(f_date)
	if p_id != "0" {
		var model_vc model.VideoComm
		c_info := model_vc.CommRedis(p_id)
		if len(c_info) < 1 || len(c_info["id"]) < 1 {
			t_msg = "评论异常"
			return t_status, t_msg
		}
		if c_info["p_id"] == "0" {
			//回复一级评论
			p_id = c_info["id"]
			to_name = ""
		} else {
			//回复二级以下评论
			p_id = c_info["p_id"]
			to_name = c_info["nick_name"]
		}

		to_account = c_info["account"]
		p_sql := fmt.Sprintf("update video_comm set reply_num=reply_num+1 where id='%s';", c_info["id"])

		//生成系统通知
		title := nick_name + "回复了您的评论"
		sub_title := "您在视频" + v_info["title"] + "中的评论有人回复了"
		sys_sql := fmt.Sprintf("insert into sys_msg (id,title,sub_title,content,create_time,to_user,msg_type) VALUES ('%s','%s','%s','%s','%s','%s','1');", id, title, sub_title, in_param["content"], now_time, to_account)
		sql_arr = append(sql_arr[0:], p_sql)
		sql_arr = append(sql_arr[0:], sys_sql)
	}

	table_name := "video_comm"
	c_data := map[string]string{}
	c_data["id"] = id
	c_data["video_id"] = in_param["vid"]
	c_data["account"] = account
	c_data["to_account"] = to_account
	c_data["content"] = in_param["content"]
	c_data["p_id"] = p_id
	c_data["nick_name"] = nick_name
	c_data["to_name"] = to_name
	c_data["create_time"] = now_time
	in_sql := common.InsertSql(table_name, c_data)
	sql_arr = append(sql_arr[0:], in_sql)
	err := model.Trans(sql_arr)
	if err != nil {
		t_msg = "评论失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  点赞/取消点赞评论,以及获取点赞总次数
 */
func (u *User) LikeComm(in_param map[string]string, ctx *gin.Context) (int, string, map[string]string) {
	t_status := 100
	t_msg := "参数不足"
	like_comm := map[string]string{}
	like_comm["is_like"] = "0"
	if len(in_param["comm_id"]) < 1 {
		return t_status, t_msg, like_comm
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, like_comm
	}

	var model_vc model.VideoComm
	c_info := model_vc.CommRedis(in_param["comm_id"])
	if len(c_info) < 1 || len(c_info["id"]) < 1 {
		t_msg = "评论异常"
		return t_status, t_msg, like_comm
	}

	//判断评论是否已经点赞
	var model_cl model.CommLike
	like_info := model_cl.CommLikeRedis(in_param["comm_id"], account)
	like_sql := ""
	comm_sql := ""
	sys_sql := ""

	now_time := time.Now().Format(f_date)

	if len(like_info) < 1 || len(like_info["id"]) < 1 {
		//没有点赞，则点赞，并且添加点赞次数
		like_comm["is_like"] = "1"
		id := model.GetKey(16)
		like_sql = fmt.Sprintf("insert into comm_like (id,comm_id,account,create_time) VALUES ('%s','%s','%s','%s');", id, in_param["comm_id"], account, now_time)
		comm_sql = fmt.Sprintf("update video_comm set like_num=like_num+1 where id='%s';", in_param["comm_id"])

		//查询视频信息
		var model_vl model.VideoList
		v_info := model_vl.VideoRedis(c_info["video_id"])
		var model_ui model.UserInfo
		u_info := model_ui.UserRedis(account)

		//生成系统通知
		title := u_info["nick_name"] + "赞了您的评论"
		sub_title := "您在视频" + v_info["title"] + "中的评论"
		sys_sql = fmt.Sprintf("insert into sys_msg (id,title,sub_title,content,create_time,to_user,msg_type) VALUES ('%s','%s','%s','%s','%s','%s','2');", id, title, sub_title, c_info["content"], now_time, c_info["account"])

	} else {
		like_sql = fmt.Sprintf("delete from comm_like where id='%s';", like_info["id"])
		comm_sql = fmt.Sprintf("update video_comm set like_num=like_num-1 where id='%s';", in_param["comm_id"])
		sys_sql = fmt.Sprintf("delete from sys_msg where id='%s';", like_info["id"])
	}

	sql_arr := []string{like_sql, comm_sql, sys_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		like_comm["is_like"] = "0"
		t_msg = "点赞失败"
		return t_status, t_msg, like_comm
	}

	//清除缓存
	redis_key1 := "comm_like:" + in_param["comm_id"] + "_" + account
	redis.RediGo.KeyDel(redis_key1)
	redis_key2 := "video_comm:" + in_param["comm_id"]
	redis.RediGo.KeyDel(redis_key2)

	c_info = model_vc.CommRedis(in_param["comm_id"])
	like_comm["like_num"] = c_info["like_num"]

	t_status = 200
	t_msg = "success"
	return t_status, t_msg, like_comm
}

/**
*  点赞视频,以及获取点赞总次数
 */
func (u *User) LikeVideo(in_param map[string]string, ctx *gin.Context) (int, string, map[string]string) {
	t_status := 100
	t_msg := "参数不足"
	like_video := map[string]string{}
	like_video["is_like"] = "0"

	if len(in_param["vid"]) < 1 {
		return t_status, t_msg, like_video
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, like_video
	}

	var model_v model.VideoList
	c_info := model_v.VideoRedis(in_param["vid"])
	if len(c_info) < 1 || len(c_info["id"]) < 1 {
		t_msg = "视频异常"
		return t_status, t_msg, like_video
	}

	//判断评论是否已经点赞
	var model_vl model.VideoLike
	like_info := model_vl.VideoLikeRedis(in_param["vid"], account)
	like_sql := ""
	video_sql := ""
	if len(like_info) < 1 || len(like_info["id"]) < 1 {
		//没有点赞，则点赞，并且添加点赞次数
		like_video["is_like"] = "1"
		id := model.GetKey(16)
		now_time := time.Now().Format(f_date)
		like_sql = fmt.Sprintf("insert into video_like (id,video_id,account,create_time) VALUES ('%s','%s','%s','%s');", id, in_param["vid"], account, now_time)
		video_sql = fmt.Sprintf("update video_list set like_num=like_num+1 where id='%s';", in_param["vid"])
	} else {
		like_sql = fmt.Sprintf("delete from video_like where id='%s';", like_info["id"])
		video_sql = fmt.Sprintf("update video_list set like_num=like_num-1 where id='%s';", in_param["vid"])
	}
	sql_arr := []string{like_sql, video_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		like_video["is_like"] = "0"
		t_msg = "点赞失败"
		return t_status, t_msg, like_video
	}

	//清除缓存
	redis_skey := "video_list:" + in_param["vid"]
	redis.RediGo.KeyDel(redis_skey)
	redis_key := "video_like:" + in_param["vid"] + "_" + account
	redis.RediGo.KeyDel(redis_key)

	v_info := model_v.VideoRedis(in_param["vid"])
	like_video["like_num"] = v_info["like_num"]

	t_status = 200
	t_msg = "success"
	return t_status, t_msg, like_video
}

/**
*  收藏视频
 */
func (u *User) StarVideo(in_param map[string]string, ctx *gin.Context) (int, int, string) {
	t_status := 100
	t_msg := "参数不足"
	is_star := 0
	if len(in_param["vid"]) < 1 {
		return t_status, is_star, t_msg
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, is_star, t_msg
	}

	var model_v model.VideoList
	c_info := model_v.VideoRedis(in_param["vid"])
	if len(c_info) < 1 || len(c_info["id"]) < 1 {
		t_msg = "视频异常"
		return t_status, is_star, t_msg
	}

	//判断视频是否已经收藏
	var model_cl model.VideoStar
	star_info := model_cl.VideoStarRedis(in_param["vid"], account)
	star_sql := ""
	video_sql := ""
	if len(star_info) < 1 || len(star_info["id"]) < 1 {
		//没有收藏，则收藏，并且添加收藏次数
		is_star = 1
		id := model.GetKey(16)
		now_time := time.Now().Format(f_date)
		star_sql = fmt.Sprintf("insert into video_star (id,video_id,account,create_time) VALUES ('%s','%s','%s','%s');", id, in_param["vid"], account, now_time)
		video_sql = fmt.Sprintf("update video_list set star_num=star_num+1 where id='%s';", in_param["vid"])
	} else {
		star_sql = fmt.Sprintf("delete from video_star where id='%s';", star_info["id"])
		video_sql = fmt.Sprintf("update video_list set star_num=star_num-1 where id='%s';", in_param["vid"])
		//清除缓存
		redis_key := "video_star:" + in_param["vid"] + "_" + account
		redis.RediGo.KeyDel(redis_key)
	}
	sql_arr := []string{star_sql, video_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		is_star = 0
		t_msg = "点赞失败"
		return t_status, is_star, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, is_star, t_msg
}

/**
*  查询自己的评论
 */
func (u *User) UserComm(in_param map[string]string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 100
	t_msg := "视频ID不能为空"
	total := 0
	comm_list := []map[string]string{}
	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, total, t_msg, comm_list
	}
	var model_ui model.UserInfo
	u_info := model_ui.UserRedis(account)
	t_status = 200
	t_msg = "success"
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size
	count_field := "count(0) as num"
	table_name := "video_comm"
	p_w := map[string]interface{}{}
	p_w["account"] = account
	total, _ = model.CommonTotal(table_name, count_field, p_w)
	if total <= offset {
		return t_status, total, t_msg, comm_list
	}

	fields := []string{"id", "account", "video_id", "to_account", "content", "p_id", "nick_name", "to_name", "like_num", "reply_num", "create_time"}
	order_by := "create_time desc"
	c_list, _ := model.PageList(table_name, order_by, page_size, offset, fields, p_w)
	if len(c_list) < 1 {
		return t_status, total, t_msg, comm_list
	}
	var model_vl model.VideoList
	for _, c_val := range c_list {
		v_info := model_vl.VideoRedis(c_val["video_id"])
		c_val["title"] = v_info["title"]
		c_val["avatar"] = u_info["avatar"]
		comm_list = append(comm_list[0:], c_val)
	}
	return t_status, total, t_msg, comm_list
}

/**
*  用户的额度和积分
 */
func (u *User) UserMoney(ctx *gin.Context) (int, string, float64, int) {
	t_status := 100
	t_msg := "错误"
	amount := 0.00
	score := 0

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, amount, score
	}

	t_status = 200
	t_msg = "success"

	var m_ui model.UserInfo
	u_info := m_ui.User(account)
	amount = u_info.Amount
	score = u_info.Score

	return t_status, t_msg, amount, score
}

/**
*  上传头像
 */
func (u *User) UserAvatar(in_param map[string]string, ctx *gin.Context) (int, string, string) {
	t_status := 100
	t_msg := "错误"
	avatar_url := ""
	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, avatar_url
	}

	dateymd := time.Now().Format("20060102")
	form_name := "avatar"
	img_path := hook.WebConfVal("img_path")
	url_path := "/avatar/" + dateymd
	filter_arr := []string{"jpg", "jpeg", "png", "gif"}
	file_name := ""
	file_path := ""

	// t_status, t_msg, file_name = uploadFile(file_name, form_name, img_path, filter_arr, ctx)
	t_status, t_msg, _, file_path = UploadFile(form_name, img_path, url_path, file_name, "0", filter_arr, ctx)
	if t_status != 200 {
		return t_status, t_msg, avatar_url
	}

	img_url := hook.WebConfVal("img_url")
	avatar_url = img_url + file_path

	t_msg = "success"

	update_sql := fmt.Sprintf("update user_info set avatar='%s' where account='%s';", avatar_url, account)
	err := model.Query(update_sql)
	if err != nil {
		t_status = 100
		t_msg = "更新失败"
		return t_status, t_msg, avatar_url
	}
	//清除缓存
	redis_key := "user_info:account:" + account
	redis.RediGo.KeyDel(redis_key)
	return t_status, t_msg, avatar_url
}

/**
*  修改昵称
 */
func (u *User) UserNickName(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "错误"

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}
	//判断昵称是否为空
	nick_name := strings.Replace(in_param["nick_name"], " ", "", -1)
	if len(nick_name) < 1 {
		t_msg = "昵称不能为空"
		return t_status, t_msg
	}

	if len(nick_name) < len(in_param["nick_name"]) { //含有字符串
		t_msg = "昵称不能有特殊字符"
		return t_status, t_msg
	}

	var model_us model.UserInfo
	u_map := model_us.UserByNickRedis(nick_name)
	if len(u_map) > 0 || len(u_map["nick_name"]) > 0 {
		t_msg = "该昵称已存在"
		return t_status, t_msg
	}

	t_status = 200
	t_msg = "success"

	update_sql := fmt.Sprintf("update user_info set nick_name='%s' where account='%s';", nick_name, account)
	err := model.Query(update_sql)
	if err != nil {
		t_status = 100
		t_msg = "更新失败"
		return t_status, t_msg
	}
	//清除缓存
	redis_key_a := "user_info:account:" + account
	redis.RediGo.KeyDel(redis_key_a)
	redis_key_n := "user_info:nick_name:" + nick_name
	redis.RediGo.KeyDel(redis_key_n)
	return t_status, t_msg

}

/**
*  修改年龄
 */
func (u *User) UserAge(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "错误"

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	t_status = 200
	t_msg = "success"

	update_sql := fmt.Sprintf("update user_info set age='%s' where account='%s';", in_param["age"], account)
	err := model.Query(update_sql)
	if err != nil {
		t_status = 100
		t_msg = "更新失败"
		return t_status, t_msg
	}
	//清除缓存
	redis_key := "user_info:account:" + account
	redis.RediGo.KeyDel(redis_key)
	return t_status, t_msg

}

/**
*  修改性别
 */
func (u *User) UserSex(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "错误"

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	t_status = 200
	t_msg = "success"

	update_sql := fmt.Sprintf("update user_info set sex='%s' where account='%s';", in_param["sex"], account)
	err := model.Query(update_sql)
	if err != nil {
		t_status = 100
		t_msg = "更新失败"
		return t_status, t_msg
	}
	//清除缓存
	redis_key := "user_info:account:" + account
	redis.RediGo.KeyDel(redis_key)
	return t_status, t_msg

}

/**
*  修改个性签名
 */
func (u *User) UserPerSign(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "错误"

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	t_status = 200
	t_msg = "success"

	update_sql := fmt.Sprintf("update user_info set per_sign='%s' where account='%s';", in_param["per_sign"], account)
	err := model.Query(update_sql)
	if err != nil {
		t_status = 100
		t_msg = "更新失败"
		return t_status, t_msg
	}
	//清除缓存
	redis_key := "user_info:account:" + account
	redis.RediGo.KeyDel(redis_key)
	return t_status, t_msg
}

/**
*  修改密码
 */
func (u *User) UserEditPwd(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "错误"

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	var m_ui model.UserList
	u_info := m_ui.User(account)
	pwd := u_info.Pwd

	if in_param["old_pwd"] == "" {
		t_msg := "旧密码不能为空"
		return t_status, t_msg
	}

	if in_param["pwd"] == "" {
		t_msg := "新密码不能为空"
		return t_status, t_msg
	}

	if pwd != hook.HookAesEncrypt(in_param["old_pwd"]) {
		t_msg := "密码错误，请重新输入"
		return t_status, t_msg
	}

	code := "pwd"
	t_status, t_msg = hook.AuthAllow(code, in_param["pwd"])
	if t_status != 200 {
		t_msg := "密码格式错误"
		return t_status, t_msg
	}

	update_sql := fmt.Sprintf("update user_list set pwd='%s' where account='%s';", hook.HookAesEncrypt(in_param["pwd"]), account)
	err := model.Query(update_sql)
	if err != nil {
		t_status = 100
		t_msg = "更新失败"
		return t_status, t_msg
	}
	//清除缓存
	redis_key := "user_list:account:" + account
	redis.RediGo.KeyDel(redis_key)
	t_msg = "success"
	return t_status, t_msg
}

/**
*  获取茄子钱包存款、优惠、消费记录
 */
func (u *User) QzWallet(in_param map[string]string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 100
	t_msg := "参数不足"
	total := 0
	wallet_info := []map[string]string{}
	if len(in_param["order_type"]) < 1 {
		return t_status, total, t_msg, wallet_info
	}

	if in_param["order_type"] != "1" && in_param["order_type"] != "3" && in_param["order_type"] != "5" {
		t_msg = "参数错误"
		return t_status, total, t_msg, wallet_info
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, total, t_msg, wallet_info
	}
	t_status = 200
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size

	order_type := ""
	order_type = "order_type=" + in_param["order_type"]

	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM order_list where account ='%s' and amount >0.00 and %s limit 1;", account, order_type)

	total_map, _ := model.SqlRow(count_sql, count_field)
	if len(total_map) < 1 || len(total_map["num"]) < 1 {
		t_msg = "暂无记录"
		return t_status, total, t_msg, wallet_info
	}
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		t_msg = "暂无记录"
		return t_status, total, t_msg, wallet_info
	}
	fields := []string{"id", "amount", "status", "order_number", "create_time", "pay_id", "act_id"}
	list_sql := fmt.Sprintf("SELECT id,amount,status,order_number,create_time,pay_id,act_id FROM order_list where account ='%s' and amount >0.00 and %s order by create_time desc limit %d,%d;", account, order_type, offset, page_size)
	c_list, _ := model.SqlRows(list_sql, fields)

	var m_cfg model.PayConfig
	var m_bonus model.ActiveBonus
	for _, c_val := range c_list {
		pay_int, _ := strconv.Atoi(c_val["pay_id"])
		act_int, _ := strconv.Atoi(c_val["act_id"])
		if pay_int > 0 {
			p_info := m_cfg.PayRedis(c_val["pay_id"])
			c_val["title"] = p_info["title"] + "存入"
			wallet_info = append(wallet_info[0:], c_val)
		} else if act_int > 0 {
			p_info := m_bonus.BonusRedis(c_val["act_id"])
			c_val["title"] = p_info["title"]
			wallet_info = append(wallet_info[0:], c_val)
		} else {
			c_val["title"] = "购买积分"
			wallet_info = append(wallet_info[0:], c_val)
		}
	}

	t_msg = "success"
	return t_status, total, t_msg, wallet_info
}

/**
*   获取kok币的购买、优惠、消费记录
 */
func (u *User) KokGold(in_param map[string]string, ctx *gin.Context) (int, int, string, []map[string]string) {
	t_status := 100
	t_msg := "参数不足"
	total := 0
	gold_info := []map[string]string{}
	if len(in_param["order_type"]) < 1 {
		return t_status, total, t_msg, gold_info
	}

	if in_param["order_type"] != "3" && in_param["order_type"] != "5" && in_param["order_type"] != "6" {
		t_msg = "参数错误"
		return t_status, total, t_msg, gold_info
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, total, t_msg, gold_info
	}

	t_status = 200
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size

	order_type := ""
	order_type = "order_type=" + in_param["order_type"]

	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM order_list where account ='%s' and score >0 and %s limit 1;", account, order_type)

	total_map, _ := model.SqlRow(count_sql, count_field)
	if len(total_map) < 1 || len(total_map["num"]) < 1 {
		t_msg = "暂无记录"
		return t_status, total, t_msg, gold_info
	}
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		t_msg = "暂无记录"
		return t_status, total, t_msg, gold_info
	}
	fields := []string{"id", "score", "status", "order_number", "create_time", "act_id"}
	list_sql := fmt.Sprintf("SELECT id,score,status,order_number,create_time,act_id FROM order_list where account ='%s' and score >0 and %s order by create_time desc limit %d,%d;", account, order_type, offset, page_size)
	c_list, _ := model.SqlRows(list_sql, fields)
	for _, c_val := range c_list {
		if in_param["order_type"] == "5" {
			c_val["title"] = "购买积分"
			gold_info = append(gold_info[0:], c_val)
		} else if in_param["order_type"] == "6" {
			c_val["title"] = "观看视频消费积分"
			gold_info = append(gold_info[0:], c_val)
		} else {
			var m_bonus model.ActiveBonus
			p_info := m_bonus.BonusRedis(c_val["act_id"])
			c_val["title"] = p_info["title"]
			gold_info = append(gold_info[0:], c_val)
		}
	}

	t_msg = "success"
	return t_status, total, t_msg, gold_info
}

/**
*  查询自己的收藏
 */
func (u *User) UserCollect(in_param map[string]string, ctx *gin.Context) (int, int, string, []map[string]interface{}) {
	t_status := 100
	t_msg := "错误"
	total := 0
	collect_list := []map[string]interface{}{}
	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, total, t_msg, collect_list
	}
	t_status = 200
	t_msg = "success"
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size

	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM video_star where account ='%s' limit 1;", account)

	total_map, _ := model.SqlRow(count_sql, count_field)
	if len(total_map) < 1 || len(total_map["num"]) < 1 {
		t_msg = "暂无收藏"
		return t_status, total, t_msg, collect_list
	}
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		t_msg = "暂无收藏"
		return t_status, total, t_msg, collect_list
	}

	fields := []string{"id", "video_id", "create_time"}
	list_sql := fmt.Sprintf("SELECT id,video_id,create_time FROM video_star where account ='%s' order by create_time desc limit %d,%d;", account, offset, page_size)
	c_list, _ := model.SqlRows(list_sql, fields)
	var model_vl model.VideoList
	for _, c_val := range c_list {

		video_id := c_val["video_id"]

		v_info := model_vl.VideoRedis(video_id)
		p_map := map[string]interface{}{}
		p_map["video_id"] = video_id
		p_map["title"] = v_info["title"]
		p_map["pc_img"] = v_info["pc_img"]
		p_map["h5_img"] = v_info["h5_img"]
		p_map["url"] = v_info["url"]
		p_map["create_time"] = c_val["create_time"]

		//获取标签信息
		t_field := []string{"tid", "title"}
		tag_sql := fmt.Sprintf("SELECT tl.id as tid,tl.title as title FROM video_tag vt LEFT JOIN tag_list tl ON vt.tag_id = tl.id where vt.video_id='%s' order by tl.id asc;", c_val["video_id"])
		tag_list, _ := model.SqlRows(tag_sql, t_field)
		p_map["tag_list"] = tag_list

		collect_list = append(collect_list[0:], p_map)
	}

	return t_status, total, t_msg, collect_list
}

/**
*  清空收藏
 */
func (u *User) ClearCollect(ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "参数不足"

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	collect_sql := ""
	collect_sql = fmt.Sprintf("delete from video_star where account='%s';", account)
	//清除缓存
	redis_key := "video_star:" + account
	redis.RediGo.KeyDel(redis_key)
	err := model.Query(collect_sql)
	if err != nil {
		t_msg = err.Error()
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  查询观看历史
 */
func (u *User) UserHistory(in_param map[string]string, ctx *gin.Context) (int, int, string, []map[string]interface{}) {
	t_status := 100
	t_msg := "错误"
	total := 0
	history_list := []map[string]interface{}{}
	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, total, t_msg, history_list
	}
	t_status = 200
	t_msg = "success"
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size

	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM video_history where account ='%s' limit 1;", account)

	total_map, _ := model.SqlRow(count_sql, count_field)
	if len(total_map) < 1 || len(total_map["num"]) < 1 {
		t_msg = "暂无观看历史"
		return t_status, total, t_msg, history_list
	}
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		t_msg = "暂无观看历史"
		return t_status, total, t_msg, history_list
	}

	fields := []string{"id", "video_id", "create_time"}
	list_sql := fmt.Sprintf("SELECT id,video_id,create_time FROM video_history where account ='%s' order by create_time desc limit %d,%d;", account, offset, page_size)
	c_list, _ := model.SqlRows(list_sql, fields)
	var model_vl model.VideoList
	for _, c_val := range c_list {

		video_id := c_val["video_id"]

		v_info := model_vl.VideoRedis(video_id)
		p_map := map[string]interface{}{}
		p_map["video_id"] = video_id
		p_map["title"] = v_info["title"]
		p_map["pc_img"] = v_info["pc_img"]
		p_map["h5_img"] = v_info["h5_img"]
		p_map["url"] = v_info["url"]
		p_map["create_time"] = c_val["create_time"]

		//获取标签信息
		t_field := []string{"tid", "title"}
		tag_sql := fmt.Sprintf("SELECT tl.id as tid,tl.title as title FROM video_tag vt LEFT JOIN tag_list tl ON vt.tag_id = tl.id where vt.video_id='%s' order by tl.id asc;", c_val["video_id"])
		tag_list, _ := model.SqlRows(tag_sql, t_field)
		p_map["tag_list"] = tag_list

		history_list = append(history_list[0:], p_map)
	}

	return t_status, total, t_msg, history_list
}

/**
*  消息中心
 */
func (u *User) SysMsg(in_param map[string]string, ctx *gin.Context) (int, string, []map[string]string) {
	t_status := 100
	t_msg := "参数不足"
	total := 0
	msg_list := []map[string]string{}
	if len(in_param["msg_type"]) < 1 {
		return t_status, t_msg, msg_list
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, msg_list
	}
	t_status = 200
	t_msg = "success"

	//查询总数
	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM sys_msg where msg_type='%s' and (to_user ='%s' or to_user='all') limit 1;", in_param["msg_type"], account)

	total_map, _ := model.SqlRow(count_sql, count_field)
	if len(total_map) < 1 || len(total_map["num"]) < 1 {
		return t_status, t_msg, msg_list
	}

	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size

	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		return t_status, t_msg, msg_list
	}

	//查询详情
	list_field := []string{"id", "title", "sub_title", "content", "create_time"}
	list_sql := fmt.Sprintf("SELECT id,title,sub_title,content,create_time FROM sys_msg where msg_type='%s' and (to_user ='%s' or to_user='all') order by create_time desc limit %d,%d;", in_param["msg_type"], account, offset, page_size)

	list_map, _ := model.SqlRows(list_sql, list_field)
	if len(list_map) < 1 {
		return t_status, t_msg, msg_list
	}
	var model_mv model.MsgView
	for _, l_val := range list_map {
		msg_map := map[string]string{}

		//查询是否被查看
		m_view := model_mv.MsgRedis(l_val["id"], account)

		if len(m_view) > 0 && len(m_view["id"]) > 0 {
			if m_view["status"] == "2" {
				msg_map = l_val
				msg_map["is_view"] = "1"
				msg_list = append(msg_list[0:], msg_map)
			}
		} else {
			msg_map = l_val
			msg_map["is_view"] = "0"
			msg_list = append(msg_list[0:], msg_map)
		}

	}
	return t_status, t_msg, msg_list
}

func (u *User) ViewMsg(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "参数不足"
	if len(in_param["msg_id"]) < 1 {
		return t_status, t_msg
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	var model_sm model.SysMsg
	s_msg := model_sm.MsgRedis(in_param["msg_id"])
	if len(s_msg) < 1 || len(s_msg["id"]) < 1 {
		t_msg = "消息ID错误"
		return t_status, t_msg
	}

	var model_mv model.MsgView
	v_msg := model_mv.MsgRedis(in_param["msg_id"], account)
	if len(v_msg) > 0 && len(v_msg["id"]) > 1 {
		t_msg = "该消息已经查看"
		return t_status, t_msg
	}

	v_id := model.GetKey(16)
	create_time := time.Now().Format(f_date)
	//查询详情
	msg_sql := fmt.Sprintf("insert into msg_view (id,msg_id,account,create_time,status) VALUES ('%s','%s','%s','%s',%d); ", v_id, in_param["msg_id"], account, create_time, 2)

	err := model.Query(msg_sql)
	if err != nil {
		t_msg = "查看失败"
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  点赞/取消点赞茄子主推视频
 */
func (u *User) LikeTheme(in_param map[string]string, ctx *gin.Context) (int, string, map[string]string) {
	t_status := 100
	t_msg := "参数不足"
	like_list := map[string]string{}
	like_list["is_like"] = "0"

	if len(in_param["id"]) < 1 {
		return t_status, t_msg, like_list
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, like_list
	}

	var model_v model.ThemeVideo
	c_info := model_v.VideoRedis(in_param["id"])
	if len(c_info) < 1 || len(c_info["id"]) < 1 {
		t_msg = "视频异常"
		return t_status, t_msg, like_list
	}
	//判断视频是否已经点赞
	var model_vl model.ThemeLike
	like_info := model_vl.ThemeLikeRedis(in_param["id"], account)
	like_sql := ""
	video_sql := ""
	if len(like_info) < 1 || len(like_info["id"]) < 1 {
		//没有点赞，则点赞，并且添加点赞次数
		like_list["is_like"] = "1"
		id := model.GetKey(16)
		now_time := time.Now().Format(f_date)
		like_sql = fmt.Sprintf("insert into theme_like (id,comm_id,account,create_time) VALUES ('%s','%s','%s','%s');", id, in_param["id"], account, now_time)
		video_sql = fmt.Sprintf("update theme_video set like_num=like_num+1 where id='%s';", in_param["id"])
	} else {
		like_sql = fmt.Sprintf("delete from theme_like where id='%s';", like_info["id"])
		video_sql = fmt.Sprintf("update theme_video set like_num=like_num-1 where id='%s';", in_param["id"])
	}
	sql_arr := []string{like_sql, video_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		like_list["is_like"] = "0"
		t_msg = "点赞失败"
		return t_status, t_msg, like_list
	}

	//清除缓存
	redis_vkey := "theme_video:" + in_param["id"]
	redis.RediGo.KeyDel(redis_vkey)
	redis_lkey := "theme_like:" + in_param["id"] + "_" + account
	redis.RediGo.KeyDel(redis_lkey)

	v_info := model_v.VideoRedis(in_param["id"])
	like_list["like_num"] = v_info["like_num"]

	t_status = 200
	t_msg = "success"
	return t_status, t_msg, like_list
}

/**
*  查看自己是否点赞茄子主推视频
 */
func (u *User) SeeLike(in_param map[string]string, ctx *gin.Context) (int, string, map[string]string) {
	t_status := 100
	t_msg := "参数不足"
	like_info := map[string]string{}
	like_info["is_like"] = "0"

	if len(in_param["id"]) < 1 {
		return t_status, t_msg, like_info
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, like_info
	}

	var model_v model.ThemeVideo
	c_info := model_v.VideoRedis(in_param["id"])
	if len(c_info) < 1 || len(c_info["id"]) < 1 {
		t_msg = "视频异常"
		return t_status, t_msg, like_info
	}
	//获取点赞次数
	like_info["like_num"] = c_info["like_num"]

	t_status = 200
	t_msg = "success"

	//判断视频是否已经点赞
	var model_vl model.ThemeLike
	l_info := model_vl.ThemeLikeRedis(in_param["id"], account)

	if len(l_info) < 1 || len(l_info["id"]) < 1 {
		like_info["is_like"] = "0"
	} else {
		like_info["is_like"] = "1"
	}

	return t_status, t_msg, like_info
}

/**
*   删除指定的收藏视频
 */
func (u *User) DelCollect(video_id []string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "信息异常"

	if len(video_id) < 1 {
		t_msg = "视频ID不能为空"
		return t_status, t_msg
	}
	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}
	for _, val := range video_id {
		del_sql := fmt.Sprintf("delete from video_star where account='%s' and video_id='%s';", account, val)
		err := model.Query(del_sql)
		if err != nil {
			t_status = 100
			t_msg = "删除失败"
			return t_status, t_msg
		}

		//清除缓存
		redis_key := "video_star:" + val + "_" + account
		redis.RediGo.KeyDel(redis_key)
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*   删除指定的观看历史
 */
func (u *User) DelHistory(video_id []string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "信息异常"

	if len(video_id) < 1 {
		t_msg = "视频ID不能为空"
		return t_status, t_msg
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}
	for _, val := range video_id {
		del_sql := fmt.Sprintf("delete from video_history where account='%s' and video_id='%s';", account, val)
		err := model.Query(del_sql)
		if err != nil {
			t_status = 100
			t_msg = "删除失败"
			return t_status, t_msg
		}

		//清除缓存
		redis_key := "video_history:" + val + "_" + account
		redis.RediGo.KeyDel(redis_key)
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  清空观看历史
 */
func (u *User) ClearHistory(ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "参数不足"

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	del_sql := ""
	del_sql = fmt.Sprintf("delete from video_history where account='%s';", account)
	//清除缓存
	redis_key := "video_history:" + account
	redis.RediGo.KeyDel(redis_key)
	err := model.Query(del_sql)
	if err != nil {
		t_msg = err.Error()
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*   删除指定消息中心消息
 */
func (u *User) DelMsg(msg_id []string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "信息异常"

	if len(msg_id) < 1 {
		t_msg = "消息ID不能为空"
		return t_status, t_msg
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}
	for _, val := range msg_id {

		msg_sql := ""
		//判断消息是否已读
		var model_mv model.MsgView
		m_info := model_mv.MsgRedis(val, account)
		if len(m_info) < 1 || len(m_info["id"]) < 1 {
			v_id := model.GetKey(16)
			create_time := time.Now().Format(f_date)
			msg_sql = fmt.Sprintf("insert into msg_view (id,msg_id,account,create_time,status) VALUES ('%s','%s','%s','%s',%d); ", v_id, val, account, create_time, 3)
		} else {
			msg_sql = fmt.Sprintf("update msg_view set status=%d where id='%s';", 3, val)
		}

		err := model.Query(msg_sql)
		if err != nil {
			t_status = 100
			t_msg = "删除失败"
			return t_status, t_msg
		}

		//清除缓存
		redis_key := "msg_view:" + val + ":" + account
		redis.RediGo.KeyDel(redis_key)
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  查看未读消息数量
 */
func (u *User) UnreadNum(ctx *gin.Context) (int, string, map[string]int) {
	t_status := 100
	t_msg := "参数不足"
	unread_info := map[string]int{}
	unread_info["all_num"] = 0
	unread_info["sys_num"] = 0
	unread_info["reply_num"] = 0
	unread_info["like_num"] = 0

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, unread_info
	}
	t_status = 200
	t_msg = "success"

	msg_type := [...]int{0, 1, 2}

	for _, t_val := range msg_type {

		//查询详情
		list_field := []string{"id"}
		list_sql := fmt.Sprintf("SELECT id FROM sys_msg where msg_type=%d and (to_user ='%s' or to_user='all');", t_val, account)

		list_map, _ := model.SqlRows(list_sql, list_field)
		if len(list_map) < 1 {
			continue
		}
		var model_mv model.MsgView
		for _, l_val := range list_map {
			//查询是否已读或删除
			m_view := model_mv.MsgRedis(l_val["id"], account)
			//未读且未删除
			if len(m_view) < 1 && len(m_view["id"]) < 1 {
				if t_val == 0 {
					unread_info["sys_num"]++
				} else if t_val == 1 {
					unread_info["reply_num"]++
				} else {
					unread_info["like_num"]++
				}
			}
		}
	}
	unread_info["all_num"] = unread_info["sys_num"] + unread_info["reply_num"] + unread_info["like_num"]

	return t_status, t_msg, unread_info
}

func (u *User) PhoneVerify(in_param map[string]string, ctx *gin.Context) (int, string, string) {

	t_status := 100
	t_msg := "手机号不能为空"
	t_verify := ""

	if len(in_param["phone"]) < 1 {
		return t_status, t_msg, t_verify
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, t_verify
	}

	var model_u model.UserInfo
	u_map := model_u.UserRedis(account)
	if len(u_map["phone"]) > 0 && u_map["phone"] != in_param["phone"] {
		t_msg = "该账号已绑定手机,如需解绑，请输入绑定手机"
		return t_status, t_msg, t_verify
	}

	t_status, t_msg, t_verify = MSMSend(account, in_param["phone"])
	if t_status == 200 {
		return t_status, t_msg, t_verify
	}
	t_msg = "发送失败"
	return t_status, t_msg, ""
}

func (u *User) PhoneBind(in_param map[string]string, ctx *gin.Context) (int, string) {

	t_status := 100
	t_msg := "验证码不能为空"

	if len(in_param["code"]) < 1 {
		return t_status, t_msg
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}

	phone := ""
	t_status, t_msg, phone = MSMSendCheck(account, in_param["code"])
	if t_status == 200 { //更新手机号

		update_sql := fmt.Sprintf("update user_info set phone='%s' where account='%s';", phone, account)
		err := model.Query(update_sql)
		if err != nil {
			t_status = 100
			t_msg = "绑定失败"
			return t_status, t_msg
		}

		//清除缓存
		redis_key := fmt.Sprintf("user_info:account:%s", account)
		redis.RediGo.KeyDel(redis_key)
		MSMSendKeyDel(account)
	}

	return t_status, t_msg
}

func (u *User) PhoneChg(in_param map[string]string, ctx *gin.Context) (int, string, string) {

	t_status := 100
	t_msg := "手机号不能为空"
	t_verify := ""

	if len(in_param["phone"]) < 1 {
		return t_status, t_msg, t_verify
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg, t_verify
	}

	var model_u model.UserInfo
	u_map := model_u.UserRedis(account)
	if len(u_map["phone"]) > 0 && u_map["phone"] != in_param["phone"] {
		t_status, t_msg, t_verify = MSMSend(account, in_param["phone"])
		if t_status == 200 {
			return t_status, t_msg, t_verify
		}
	}
	t_msg = "该手机号与当前绑定的手机号相同"
	return t_status, t_msg, ""
}
