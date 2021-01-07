package thread

import (
	"fmt"
	"qzapp/hook"
	"qzapp/model"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

func (t *Public) Register(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "账号密码不能为空"
	if len(in_param) < 1 || len(in_param["account"]) < 1 || len(in_param["pwd"]) < 1 {
		return t_status, t_msg
	}

	t_status, t_msg = t.AcountAllow(in_param)
	if t_status != 200 {
		return t_status, t_msg
	}
	t_status, t_msg = t.PwdAllow(in_param)
	if t_status != 200 {
		return t_status, t_msg
	}

	t_status, t_msg = hook.AppType(ctx)
	if t_status != 200 {
		t_msg = "请求异常"
		return t_status, t_msg
	}
	reg_app := t_msg
	agent_name := ""
	agent_path := ""
	session_id := common.Random("smallnumber", 30)
	t_status = 100
	if in_param["reg_code"] != "" {
		var model_ul model.UserList
		agent := model_ul.UserByCodeRedis(in_param["reg_code"])
		agent_name = agent["account"]
		agent_path = agent["agent_path"] + agent_name + "_"
	}
	u_id := model.GetKey(16)
	reg_time := time.Now().Format(f_date)
	reg_ip := ctx.ClientIP()

	list_table := "user_list"
	list_map := map[string]string{}
	list_map["id"] = u_id
	list_map["account"] = in_param["account"]
	list_map["pwd"] = hook.HookAesEncrypt(in_param["pwd"])
	list_map["reg_time"] = reg_time
	list_map["reg_ip"] = reg_ip
	list_map["reg_app"] = reg_app
	list_map["login_time"] = reg_time
	list_map["login_ip"] = reg_ip
	list_map["agent_name"] = agent_name
	list_map["agent_path"] = agent_path
	list_map["reg_code"] = in_param["reg_code"]
	list_map["is_line"] = "1"
	list_map["session_id"] = session_id
	list_sql := common.InsertSql(list_table, list_map)

	//网站的配置
	var model_wc model.WebConf
	avat := model_wc.ConfRedis("avatar")
	num := model_wc.ConfRedis("day_num")
	day_num := "0"
	if num["day_num"] != "" {
		day_num = num["day_num"]
	}

	info_table := "user_info"
	info_map := map[string]string{}
	info_map["id"] = u_id
	info_map["account"] = in_param["account"]
	info_map["avatar"] = avat["val"]
	info_map["day_num"] = day_num

	var model_us model.UserInfo
	count := 0
	nick_name := ""
	//随机获取昵称,如果昵称唯一,跳出循环
	for count = 0; count < 100; count++ {
		nick_name = GetNameAll()
		u_map := model_us.UserByNickRedis(nick_name)
		if len(u_map) < 1 || len(u_map["nick_name"]) < 1 {
			break
		}
	}
	if count == 100 {
		t_msg = "注册失败,请重试"
		return t_status, t_msg
	}

	info_map["nick_name"] = nick_name
	info_sql := common.InsertSql(info_table, info_map)

	sql_arr := []string{list_sql, info_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		t_msg = "注册失败"
		return t_status, t_msg
	}

	t_status = 200
	t_msg = "success"

	session := sessions.Default(ctx)
	session.Set("session_id", session_id)
	session.Set("account", in_param["account"])
	session.Set("agent_name", agent_name)
	session.Set("is_test", 0)
	session.Set("is_agent", 0)
	session.Save()

	return t_status, t_msg
}

/**
*  判断账号是否合法
 */
func (t *Public) AcountAllow(in_param map[string]string) (int, string) {
	t_status := 100
	t_msg := "账号已存在"
	var model_ul model.UserList
	u_info := model_ul.UserRedis(in_param["account"])
	if len(u_info) > 0 && len(u_info["account"]) > 0 {
		return t_status, t_msg
	}

	code := "account"
	t_status, t_msg = hook.AuthAllow(code, in_param["account"])

	return t_status, t_msg
}

/**
*  判断密码是否合法
 */
func (t *Public) PwdAllow(in_param map[string]string) (int, string) {
	code := "pwd"
	t_status, t_msg := hook.AuthAllow(code, in_param["pwd"])

	return t_status, t_msg
}

func (t *Public) Login(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "账号密码不能为空"
	if len(in_param) < 1 {
		return t_status, t_msg
	}
	for _, in_val := range in_param {
		if in_val == "" {
			return t_status, t_msg
		}
	}
	var model_ul model.UserList
	u_info := model_ul.User(in_param["account"])
	if len(u_info.Account) < 1 {
		t_msg = "账号或密码错误"
		return t_status, t_msg
	}
	if u_info.Status != 1 {
		t_msg = "账号已被锁定或删除,请联系管理员"
		return t_status, t_msg
	}
	if hook.HookAesDecrypt(u_info.Pwd) != in_param["pwd"] {
		t_msg = "账号或密码错误"
		return t_status, t_msg
	}

	session := sessions.Default(ctx)

	session_id := common.Random("smallnumber", 30)
	session.Set("session_id", session_id)
	session.Set("account", u_info.Account)
	session.Set("agent_name", u_info.Agent_name)
	session.Set("is_test", u_info.Is_test)
	session.Set("is_agent", u_info.Is_agent)
	session.Save()
	a_data := map[string]interface{}{}
	a_data["session_id"] = session_id
	a_data["login_time"] = time.Now().Format(f_date)
	a_data["login_ip"] = ctx.ClientIP()
	a_data["is_line"] = 1
	err := model_ul.Updates(u_info, a_data)
	if err != nil {
		t_status = 100
		t_msg = "登录失败"
		return t_status, t_msg
	}
	delSession(u_info.Session_id)

	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  退出
 */
func (t *Public) Logout(ctx *gin.Context) (int, string) {
	t_status := 200
	t_msg := "success"
	session := sessions.Default(ctx)
	session_id := session.Get("session_id")
	if session_id == nil {
		return t_status, t_msg
	}
	sess_id := fmt.Sprintf("%v", session_id)
	session.Clear()
	session.Save()
	delSession(sess_id)
	out_sql := fmt.Sprintf("update user_list set is_line=0,session_id='' where session_id='%s';", sess_id)
	model.Query(out_sql)
	return t_status, t_msg
}

/**
*  查询启动页
 */
func (t *Public) Loading() (int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	fields := []string{"img_url", "loading_time", "jump_url", "sort"}
	sql := "select img_url,loading_time,jump_url,sort from loading_list where `status`=1 order by sort asc;"
	load_list, _ := model.SqlRows(sql, fields)
	return t_status, t_msg, load_list
}

/**
*  查询轮播图
 */
func (t *Public) Banner() (int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	fields := []string{"h5_img", "title", "jump_url", "sort"}
	sql := "select h5_img,title,jump_url,sort from banner_list where `status`=1 order by sort asc;"
	banner_list, _ := model.SqlRows(sql, fields)
	return t_status, t_msg, banner_list
}

/**
*  查询轮播图
 */
func (t *Public) Pop() (int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	fields := []string{"content", "title", "jump_url", "sort", "butt_title"}
	sql := "select content,title,jump_url,sort,butt_title from pop_list where `status`=1 order by sort asc;"
	pop_list, _ := model.SqlRows(sql, fields)
	return t_status, t_msg, pop_list
}

func (t *Public) TagType(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	type_list := []map[string]string{}
	p_where := map[string]interface{}{}
	p_where["status"] = 1
	if in_param["is_top"] != "" {
		is_top, _ := strconv.Atoi(in_param["is_top"])
		p_where["is_top"] = is_top
	}
	table_name := "tag_type"
	fields := []string{"id", "title", "img_url", "is_top", "sort"}
	order_by := "sort asc"
	page_size := 100
	offset := 0
	count_field := "count(0) as num"
	total, _ := model.ListTotal(table_name, count_field, p_where)
	if total < 1 {
		return t_status, total, t_msg, type_list
	}
	type_list, _ = model.PageList(table_name, order_by, page_size, offset, fields, p_where)
	return t_status, total, t_msg, type_list
}

/**
*  标签列表
 */
func (t *Public) TagList(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	tag_list := []map[string]string{}
	p_where := map[string]interface{}{}

	table_name := "tag_list"
	fields := []string{"title", "id"}
	order_by := "id asc"
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size
	count_field := "count(0) as num"
	total, _ := model.ListTotal(table_name, count_field, p_where)
	if total <= offset {
		return t_status, total, t_msg, tag_list
	}
	tag_list, _ = model.PageList(table_name, order_by, page_size, offset, fields, p_where)
	return t_status, total, t_msg, tag_list
}

func (t *Public) HotVideo(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	video_list := []map[string]string{}
	p_where := map[string]interface{}{}
	p_where["status"] = 1
	if in_param["is_hot"] != "" {
		is_hot, _ := strconv.Atoi(in_param["is_hot"])
		p_where["is_hot"] = is_hot
	}
	if in_param["is_recomm"] != "" {
		is_recomm, _ := strconv.Atoi(in_param["is_recomm"])
		p_where["is_recomm"] = is_recomm
	}
	page, page_size := treadPage(in_param["page"], in_param["page_size"])

	table_name := "video_list"
	count_field := "count(0) as num"
	total, _ := model.ListTotal(table_name, count_field, p_where)
	offset := (page - 1) * page_size
	if total <= offset {
		return t_status, total, t_msg, video_list
	}

	fields := []string{"id", "title", "view_num", "h5_img"}
	order_by := "view_num DESC, id DESC"
	v_list, _ := model.PageList(table_name, order_by, page_size, offset, fields, p_where)
	if len(v_list) < 1 {
		return t_status, total, t_msg, video_list
	}
	for _, v_val := range v_list {
		v_map := map[string]string{}
		v_map["vid"] = v_val["id"]
		v_map["title"] = v_val["title"]
		v_map["view_num"] = v_val["view_num"]
		v_map["h5_img"] = v_val["h5_img"]
		video_list = append(video_list[0:], v_map)
	}

	return t_status, total, t_msg, video_list
}

/**
*  通过分类ID查询视频
 */
func (t *Public) TypeVideo(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 100
	t_msg := "请选择所属分类"
	video_list := []map[string]string{}
	total := 0
	if in_param["type_id"] == "" {
		return t_status, total, t_msg, video_list
	}
	t_status = 200
	t_msg = "success"
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size

	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM video_type vt LEFT JOIN video_list vl ON vt.video_id=vl.id where vl.status=1 and vt.type_id='%s' limit 1;", in_param["type_id"])
	total_map, _ := model.SqlRow(count_sql, count_field)
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		return t_status, total, t_msg, video_list
	}

	uni_sql := fmt.Sprintf("SELECT vl.id as vid,vl.title as title,vl.view_num as view_num,vl.h5_img as h5_img FROM video_type vt LEFT JOIN video_list vl ON vt.video_id=vl.id where vl.status=1 and vt.type_id='%s' order by vl.view_num desc limit %d,%d;", in_param["type_id"], offset, page_size)
	fields := []string{"vid", "title", "view_num", "h5_img"}

	video_list, _ = model.SqlRows(uni_sql, fields)
	return t_status, total, t_msg, video_list
}

/**
*  通过标签ID查询视频
 */
func (t *Public) TagVideo(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 100
	t_msg := "请选择所属标签"
	video_list := []map[string]string{}
	total := 0
	if in_param["tag_id"] == "" {
		return t_status, total, t_msg, video_list
	}
	t_status = 200
	t_msg = "success"
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size

	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM video_tag vt LEFT JOIN video_list vl ON vt.video_id=vl.id where vl.status=1 and vt.tag_id='%s' limit 1;", in_param["tag_id"])
	total_map, _ := model.SqlRow(count_sql, count_field)
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		return t_status, total, t_msg, video_list
	}

	uni_sql := fmt.Sprintf("SELECT vl.id as vid,vl.title as title,vl.view_num as view_num,vl.h5_img as h5_img FROM video_tag vt LEFT JOIN video_list vl ON vt.video_id=vl.id where vl.status=1 and vt.tag_id='%s' order by vl.view_num desc limit %d,%d;", in_param["tag_id"], offset, page_size)
	fields := []string{"vid", "title", "view_num", "h5_img"}

	video_list, _ = model.SqlRows(uni_sql, fields)
	return t_status, total, t_msg, video_list
}

func (t *Public) ActorData() (int, string, []map[string]interface{}) {

	t_status := 100
	t_msg := "查询错误"

	list := []map[string]interface{}{}

	table_name := "actor_list"
	// fields := []string{"id", "name", "nation", "birth", "sex", "cn_code",
	// 	"initials", "avatar", "is_hot", "height", "cup", "description",
	// 	"total", "bust", "waist", "hips", "status", "real_name", "nick_name"}

	fields := []string{"id", "name", "birth",
		"initials", "avatar", "height", "cup", "description",
		"bust", "waist", "hips"}

	actors_list, err := model.Rows(table_name, "cn_code asc", 9999, 1, fields, "status=?", 1)
	if err != nil {
		t_status = 501
		t_msg = "查询错误"
		return t_status, t_msg, list
	}

	img_url := hook.WebConfVal("img_url")
	if img_url == "" {
		t_status = 502
		t_msg = "图片地址未设置"
		return t_status, t_msg, list
	}

	index := ""
	index_arr := []map[string]string{}
	for key, actors_row := range actors_list {

		initials, ok := actors_row["initials"]
		if ok == false {
			continue
		}

		if IsZM(initials) == false {
			continue
		}

		avatar := img_url + actors_row["avatar"]

		if index == "" {
			index = initials
		}

		if index == initials { //叠加数据

			row := map[string]string{}
			row["id"] = actors_row["id"]
			row["name"] = actors_row["name"]
			row["avatar"] = avatar
			row["birth"] = actors_row["birth"]

			row["bust"] = actors_row["bust"]
			row["waist"] = actors_row["waist"]
			row["hips"] = actors_row["hips"]
			row["height"] = actors_row["height"]
			row["cup"] = actors_row["cup"]
			row["description"] = actors_row["description"]

			index_arr = append(index_arr, row)

		} else { //写入数据

			one := map[string]interface{}{}
			one["title"] = index
			one["items"] = index_arr
			list = append(list, one)

			//重置数据
			index_arr = []map[string]string{}
			index = initials

			//下一个数据
			row := map[string]string{}
			row["id"] = actors_row["id"]
			row["name"] = actors_row["name"]
			row["avatar"] = avatar
			row["birth"] = actors_row["birth"]

			row["bust"] = actors_row["bust"]
			row["waist"] = actors_row["waist"]
			row["hips"] = actors_row["hips"]
			row["height"] = actors_row["height"]
			row["cup"] = actors_row["cup"]
			row["description"] = actors_row["description"]

			index_arr = append(index_arr, row)
		}

		if key+1 == len(actors_list) {
			one := map[string]interface{}{}
			one["title"] = index
			one["items"] = index_arr
			list = append(list, one)
		}
	}

	t_status = STATUS_SUCCESS
	t_msg = MSG_SUCCESS
	return t_status, t_msg, list
}

/**
*  通过女优ID查询视频
 */
func (t *Public) ActorVideo(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 100
	t_msg := "请选择所属标签"
	total := 0
	video_list := []map[string]string{}
	if in_param["actor_id"] == "" {
		return t_status, total, t_msg, video_list
	}
	t_status = 200
	t_msg = "success"
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size
	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM video_act va LEFT JOIN video_list vl ON va.video_id = vl.id where vl.status=1 and va.act_id='%s' limit 1;", in_param["actor_id"])
	total_map, _ := model.SqlRow(count_sql, count_field)
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		return t_status, total, t_msg, video_list
	}

	uni_sql := fmt.Sprintf("SELECT vl.id as vid,vl.title as title,vl.view_num as view_num,vl.h5_img as h5_img FROM video_act va LEFT JOIN video_list vl ON va.video_id = vl.id where vl.status=1 and va.act_id='%s' order by vl.view_num desc limit %d,%d;", in_param["actor_id"], offset, page_size)
	fields := []string{"vid", "title", "view_num", "h5_img"}

	video_list, _ = model.SqlRows(uni_sql, fields)
	return t_status, total, t_msg, video_list
}

func (t *Public) ActorList(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	actor_list := []map[string]string{}
	p_where := map[string]interface{}{}
	p_where["status"] = 1
	if in_param["is_hot"] != "" {
		is_hot, _ := strconv.Atoi(in_param["is_hot"])
		p_where["is_hot"] = is_hot
	}
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size
	table_name := "actor_list"
	count_field := "count(0) as num"
	total, _ := model.ListTotal(table_name, count_field, p_where)
	if total <= offset {
		return t_status, total, t_msg, actor_list
	}

	fields := []string{"cn_name", "nation", "total", "age", "id", "initials", "is_hot", "avatar"}
	order_by := "total desc"
	actor_list, _ = model.PageList(table_name, order_by, page_size, offset, fields, p_where)
	return t_status, total, t_msg, actor_list
}

/**
*  游戏平台类型
 */
func (t *Public) PlatType() (int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	fields := []string{"code", "title", "icon"}
	sql := "select code,title,icon from game_type order by id asc;"
	plat_type, _ := model.SqlRows(sql, fields)
	return t_status, t_msg, plat_type
}

/**
*  根据平台类型查询平台列表
 */
func (t *Public) PlatList(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 100
	t_msg := "请选择平台类型"
	total := 0
	plat_list := []map[string]string{}
	if in_param["type_code"] == "" {
		return t_status, total, t_msg, plat_list
	}
	t_status = 200
	t_msg = "success"
	p_where := map[string]interface{}{}
	p_where["type_code"] = in_param["type_code"]
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size
	table_name := "game_plat"
	count_field := "count(0) as num"
	total, _ = model.ListTotal(table_name, count_field, p_where)
	if total <= offset {
		return t_status, total, t_msg, plat_list
	}

	fields := []string{"plat_code", "title", "sort"}
	order_by := "sort asc"
	g_plat, _ := model.PageList(table_name, order_by, page_size, offset, fields, p_where)
	if len(g_plat) < 1 {
		return t_status, total, t_msg, plat_list
	}
	var model_pl model.PlatList
	for _, g_val := range g_plat {
		p_map := model_pl.PlatRedis(g_val["plat_code"])
		if len(p_map) < 1 || len(p_map["code"]) < 1 {
			continue
		}
		p_info := map[string]string{}
		p_info["code"] = g_val["plat_code"]
		p_info["title"] = g_val["title"]
		p_info["icon"] = g_val["icon"]
		plat_list = append(plat_list[0:], p_info)
	}
	return t_status, total, t_msg, plat_list
}

/**
*  根据关键字搜索
 */
func (t *Public) KwSearch(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 100
	t_msg := "请输入关键字"
	total := 0
	video_list := []map[string]string{}
	if in_param["key_word"] == "" {
		return t_status, total, t_msg, video_list
	}
	t_status = 200
	t_msg = "success"
	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size
	count_field := []string{"num"}
	count_sql := "SELECT count(0) as num FROM video_list where title like '%" + in_param["key_word"] + "%' or content like '%" + in_param["key_word"] + "%' limit 1;"
	total_map, _ := model.SqlRow(count_sql, count_field)
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		return t_status, total, t_msg, video_list
	}

	uni_sql := "SELECT id as vid,title,view_num,h5_img FROM video_list where title like '%" + in_param["key_word"] + "%' or content like '%" + in_param["key_word"] + "%' order by view_num desc limit " + strconv.Itoa(offset) + "," + strconv.Itoa(page_size) + ";"
	fields := []string{"vid", "title", "view_num", "h5_img"}

	video_list, _ = model.SqlRows(uni_sql, fields)
	return t_status, total, t_msg, video_list
}

/**
*  优惠类型
 */
func (t *Public) ActiveType() (int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	fields := []string{"id", "title", "sort"}
	sql := "select id,title,sort from active_type order by sort asc;"
	plat_type, _ := model.SqlRows(sql, fields)
	return t_status, t_msg, plat_type
}

/**
*  优惠列表
 */
func (t *Public) ActiveList(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	active_list := []map[string]string{}

	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size

	type_id := ""

	if len(in_param["type_id"]) > 0 {
		type_id = " and type_id=" + in_param["type_id"]
	}

	now_time := time.Now().Format(f_date)
	count_field := []string{"num"}

	count_sql := fmt.Sprintf("SELECT count(0) as num FROM active_list where `status`=1 and s_date <='%s' and e_date >='%s' %s limit 1;", now_time, now_time, type_id)

	total_map, _ := model.SqlRow(count_sql, count_field)
	if len(total_map) < 1 || len(total_map["num"]) < 1 {
		return t_status, total, t_msg, active_list
	}
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		return t_status, total, t_msg, active_list
	}
	fields := []string{"id", "title", "sub_title", "h5_img", "jump_url", "type_id", "sort"}
	list_sql := fmt.Sprintf("SELECT id,title,sub_title,h5_img,jump_url,type_id,sort FROM active_list where `status`=1 and s_date <='%s' and e_date >='%s' %s order by sort asc limit %d,%d;", now_time, now_time, type_id, offset, page_size)

	active_list, _ = model.SqlRows(list_sql, fields)
	return t_status, total, t_msg, active_list
}

/**
*  获取客服信息
 */
func (t *Public) CsInfo() []map[string]string {
	fields := []string{"code", "content", "note", "icon"}
	cs_sql := "select code,content,note,icon from contact_list where `status`=1;"
	cs_info, _ := model.SqlRows(cs_sql, fields)

	return cs_info
}

/**
*  app的版本信息
 */
func (t *Public) AppVer(ctx *gin.Context) (int, string, map[string]string) {
	t_status := 500
	t_msg := "请求异常"
	app_ver := map[string]string{}
	t_status, t_msg = hook.AppType(ctx)
	if t_status != 200 {
		t_msg = "请求异常"
		return t_status, t_msg, app_ver
	}
	reg_app := t_msg

	t_status = 200
	t_msg = "success"
	fields := []string{"version", "is_up", "content", "create_time", "up_url"}
	cs_sql := fmt.Sprintf("select version,is_up,content,create_time,up_url from app_version where `app_type`='%s' order by create_time desc limit 1;", reg_app)
	app_ver, _ = model.SqlRow(cs_sql, fields)
	return t_status, t_msg, app_ver
}

/**
*  app的安装记录
 */
func (t *Public) OpenApp(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "参数不足"
	if len(in_param) < 1 {
		return t_status, t_msg
	}
	for in_key, in_val := range in_param {
		if in_key != "phone_number" && in_key != "account" && in_val == "" {
			return t_status, t_msg
		}
	}
	t_status, t_msg = hook.AppType(ctx)
	if t_status != 200 {
		t_msg = "请求异常"
		return t_status, t_msg
	}
	table_name := "app_setup"
	now_time := time.Now().Format(f_date)
	app_type := t_msg
	a_map := map[string]string{}
	a_map["app_type"] = app_type
	a_map["phone_info"] = in_param["phone_info"]
	a_map["phone_os"] = in_param["phone_os"]
	a_map["uid"] = in_param["uid"]
	a_map["id"] = model.GetKey(16)
	a_map["ip"] = ctx.ClientIP()
	a_map["create_time"] = now_time
	a_map["phone_number"] = in_param["phone_number"]
	a_map["last_req"] = now_time

	var model_as model.AppSetup
	app_info := model_as.AppSetRedis(in_param["uid"])
	go t.appLog(in_param["account"], a_map)
	if len(app_info) > 1 && len(app_info["uid"]) > 0 {
		up_sql := fmt.Sprintf("update %s set last_req='%s' where id='%s';", table_name, now_time, app_info["id"])
		err := model.Query(up_sql)
		if err != nil {
			t_msg = err.Error()
			return t_status, t_msg
		}

		list_sql := fmt.Sprintf("update user_list set `is_line`=1 where account='%s';", in_param["account"])
		err = model.Query(list_sql)
		if err != nil {
			t_msg = "更新在线状态失败"
			return t_status, t_msg
		}
		//清除缓存
		redis_key := "user_list:account:" + in_param["account"]
		redis.RediGo.KeyDel(redis_key)

		t_status = 200
		t_msg = "success"
		return t_status, t_msg
	}

	in_sql := common.InsertSql(table_name, a_map)
	err := model.Query(in_sql)
	if err != nil {
		t_msg = err.Error()
		return t_status, t_msg
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  app的启动记录
 */
func (t *Public) appLog(account string, in_param map[string]string) {
	if account == "" {
		account = "游客"
	}
	if len(in_param) < 1 {
		return
	}

	a_map := map[string]string{}
	a_map["id"] = in_param["id"]
	a_map["app_type"] = in_param["app_type"]
	a_map["phone_info"] = in_param["phone_info"]
	a_map["phone_os"] = in_param["phone_os"]
	a_map["uid"] = in_param["uid"]
	a_map["ip"] = in_param["ip"]
	a_map["start_time"] = in_param["create_time"]
	a_map["phone_number"] = in_param["phone_number"]
	a_map["account"] = account
	a_map["day_date"] = common.Substr(in_param["create_time"], 0, 10)
	table_name := "app_log"
	in_sql := common.InsertSql(table_name, a_map)
	model.Query(in_sql)
	return
}

/**
*  app的关闭日志
 */
func (t *Public) CloseApp(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "参数不足"
	if len(in_param) < 1 || len(in_param["uid"]) < 1 {
		return t_status, t_msg
	}
	fields := []string{"id", "start_time"}
	sel_sql := fmt.Sprintf("select id,start_time from app_log where uid='%s' and app_minu=0 order by start_time desc limit 1;", in_param["uid"])
	log_info, _ := model.SqlRow(sel_sql, fields)
	if len(log_info) < 1 || len(log_info["id"]) < 1 {
		t_msg = "缺失日志"
		return t_status, t_msg
	}

	now_time := time.Now().Format(f_date)
	app_minu := common.DifferDate(log_info["start_time"], now_time, f_date, "minute")
	if app_minu < 1 {
		app_minu = 1
	}

	account := sessInfo("account", ctx)
	if account != "" {
		list_sql := fmt.Sprintf("update user_list set `is_line`=0 where account='%s';", account)
		err := model.Query(list_sql)
		if err != nil {
			t_msg = "更新在线状态失败"
			return t_status, t_msg
		}
	}
	//清除缓存
	redis_key := "user_list:account:" + account
	redis.RediGo.KeyDel(redis_key)

	up_sql := fmt.Sprintf("update app_log set close_time='%s',app_minu=%d where id='%s';", now_time, app_minu, log_info["id"])
	err := model.Query(up_sql)
	if err != nil {
		t_msg = err.Error()
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*  购买积分的比例
 */
func (t *Public) Rate() int {
	rate := hook.WebConfVal("rate")
	if rate == "" {
		rate = "10"
	}
	rate_int, _ := strconv.Atoi(rate)
	return rate_int
}

/**
*  获取茄子主推视频
 */
func (t *Public) QzRecomm(in_param map[string]string) (int, int, string, []map[string]string) {
	t_status := 200
	t_msg := "success"
	total := 0
	recomm_list := []map[string]string{}

	page, page_size := treadPage(in_param["page"], in_param["page_size"])
	offset := (page - 1) * page_size

	now_time := time.Now().Format(f_date)
	count_field := []string{"num"}
	count_sql := fmt.Sprintf("SELECT count(0) as num FROM theme_list tl LEFT JOIN theme_video tv ON tv.theme_id=tl.id where tv.status=1 and tl.start_date <='%s' and tl.end_date >='%s' limit 1;", now_time, now_time)
	total_map, _ := model.SqlRow(count_sql, count_field)
	total, _ = strconv.Atoi(total_map["num"])
	if total <= offset {
		return t_status, total, t_msg, recomm_list
	}

	fields := []string{"id", "title", "content", "view_num", "like_num", "h5_img", "url"}
	list_sql := fmt.Sprintf("SELECT tv.id as id,tv.title as title,tv.content as content,tv.view_num as view_num,tv.like_num as like_num,tv.h5_img as h5_img,tv.url as url FROM theme_list tl LEFT JOIN theme_video tv ON tv.theme_id=tl.id where tv.status=1 and tl.start_date <='%s' and tl.end_date >='%s' order by sort asc,id desc limit %d,%d;", now_time, now_time, offset, page_size)

	recomm_list, _ = model.SqlRows(list_sql, fields)
	return t_status, total, t_msg, recomm_list
}

/**
*   茄子主推视频的观看次数+1
 */
func (t *Public) AddView(in_param map[string]string) (int, string) {
	t_status := 100
	t_msg := "信息异常"

	if len(in_param) < 1 || len(in_param["id"]) < 1 {
		t_msg = "视频ID不能为空"
		return t_status, t_msg
	}

	//更新视频观看次数
	num_sql := fmt.Sprintf("update theme_video set view_num=view_num+1 where id='%s';", in_param["id"])
	err := model.Query(num_sql)
	if err != nil {
		t_status = 100
		t_msg = "更新失败"
		return t_status, t_msg
	}

	//清除缓存
	redis_key := "theme_video:" + in_param["id"]
	redis.RediGo.KeyDel(redis_key)

	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

/**
*   影片的观看次数+1,以及获取影片的总观看次数
 */
func (t *Public) AddSee(in_param map[string]string) (int, string, string) {
	t_status := 100
	t_msg := "信息异常"
	view_num := "0"

	if len(in_param) < 1 || len(in_param["vid"]) < 1 {
		t_msg = "视频ID不能为空"
		return t_status, t_msg, view_num
	}

	//更新视频观看次数
	num_sql := fmt.Sprintf("update video_list set view_num=view_num+1 where id='%s';", in_param["vid"])
	err := model.Query(num_sql)
	if err != nil {
		t_status = 100
		t_msg = "更新失败"
		return t_status, t_msg, view_num
	}

	//清除缓存
	redis_key := "video_list:" + in_param["vid"]
	redis.RediGo.KeyDel(redis_key)

	var model_v model.VideoList
	v_info := model_v.VideoRedis(in_param["vid"])
	view_num = v_info["view_num"]

	t_status = 200
	t_msg = "success"
	return t_status, t_msg, view_num
}

/**
*   影片的观看次数+1,以及获取影片的总观看次数
 */
func (t *Public) HomeView() (int, string, []map[string]interface{}) {
	t_status := 200
	t_msg := "success"
	t_list := []map[string]interface{}{}

	type_fields := []string{"id", "title", "img_url", "is_top", "sort"}
	type_list, _ := model.SqlRows("SELECT id, title, img_url, is_top, sort FROM `tag_type` WHERE `status` = 1 ORDER BY sort ASC LIMIT 1000", type_fields)

	video_fields := []string{"vid", "title", "view_num", "h5_img"}

	for _, type_row := range type_list {

		type_id, _ := type_row["id"]
		if type_id == "" {
			continue
		}

		uni_sql := fmt.Sprintf("SELECT vl.id as vid,vl.title as title,vl.view_num as view_num,vl.h5_img as h5_img FROM video_type vt LEFT JOIN video_list vl ON vt.video_id=vl.id where vl.status=1 and vt.type_id='%s' order by vl.view_num desc limit 0,4", type_id)

		video_list, _ := model.SqlRows(uni_sql, video_fields)

		t_row := map[string]interface{}{}
		t_row["id"] = type_id
		t_row["title"] = type_row["title"]
		t_row["img_url"] = type_row["img_url"]
		t_row["video_list"] = video_list

		t_list = append(t_list, t_row)
	}

	return t_status, t_msg, t_list
}
