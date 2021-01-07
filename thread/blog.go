package thread

import (
	"encoding/json"
	"fmt"
	"qzapp/hook"
	"qzapp/model"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.stagingvip.net/publicGroup/public/common"
)

func (t_blog *BlogThread) getBlogCom(ctx *gin.Context, id, blog_id, p_id, page, page_size string) (int, []map[string]string) {

	table_name := "blog_comm"
	t_list := []map[string]string{}

	bc_w := "blog_id='" + blog_id + "' AND p_id = '" + p_id + "' AND status = 2"
	if id != "" {
		bc_w = "id='" + id + "' AND status = 2"
	}

	total := model.Count(table_name, bc_w)
	if total == 0 {
		return total, t_list
	}

	i_page, i_page_size := treadPage(page, page_size)

	bc_fields := []string{"id", "blog_id", "account", "to_account", "content", "p_id",
		"nick_name", "to_name", "like_num", "reply_num", "create_time", "is_up"}

	sql_order := "is_up DESC, like_num DESC, create_time ASC"
	if len(p_id) > 2 {
		sql_order = "create_time ASC"
	}

	bc_list, _ := model.Rows(table_name, sql_order, i_page_size, i_page, bc_fields, bc_w)

	account := sessInfo("account", ctx)
	var m_bcl model.BlogCommLike
	var m_user_info model.UserInfo

	for _, one := range bc_list {

		row := one

		//获取发帖人的头像
		user_info := m_user_info.UserRedis(one["account"])
		avatar, _ := user_info["avatar"]
		row["avatar"] = avatar
		fmt.Println("avatar->", avatar)

		comm_id := one["id"]

		//是否给博客点赞
		is_like := "0"
		is_like_map := m_bcl.RedisIsLike(comm_id, account)
		if r_is_like, _ := is_like_map["id"]; r_is_like != "" {
			is_like = "1"
		}
		row["is_like"] = is_like

		t_list = append(t_list, row)
	}

	return total, t_list

}

func (t_blog *BlogThread) BlogList(ctx *gin.Context, in_param map[string]string) (int, int, string, []map[string]interface{}) {

	t_status := STATUS_SUCCESS
	t_msg := MSG_SUCCESS
	t_list := []map[string]interface{}{}

	table_name := "blog_list"

	w_map := map[string]string{}
	w_map["status"] = "2"

	is_top, _ := in_param["is_top"]
	is_good, _ := in_param["is_good"]
	content, _ := in_param["content"]

	if is_top != "" {
		w_map["is_top"] = is_top
	}

	if is_good != "" {
		w_map["is_good"] = is_good
	}

	if content != "" {
		w_map["content LIKE"] = "%" + content + "%"
	}

	w_str := model.WhereMap2Str(w_map)

	total := model.Count(table_name, w_str)
	if total == 0 {
		return t_status, total, t_msg, t_list
	}

	page, page_size := treadPage(in_param["page"], in_param["page_size"])

	fields := []string{"id", "content", "img_url", "view_num", "account",
		"nick_name", "like_num", "avatar", "is_top", "is_good", "comm_num", "create_time"}

	sql_list, _ := model.Rows(table_name, "is_top DESC, is_good DESC, create_time DESC", page_size, page, fields, w_str)

	account := sessInfo("account", ctx)
	p_id := "0"
	var m_bll model.BlogListLike

	for _, one := range sql_list {

		row := map[string]interface{}{}
		for k, v := range one {
			if k == "img_url" {

				img_arr := []string{}
				json.Unmarshal([]byte(v), &img_arr)
				row["img_urls"] = img_arr

			} else {
				row[k] = v
			}
		}

		blog_id := one["id"]

		//是否给博客点赞
		is_like := "0"
		is_like_map := m_bll.RedisIsLike(blog_id, account)
		if r_is_like, _ := is_like_map["id"]; r_is_like != "" {
			is_like = "1"
		}
		row["is_like"] = is_like
		row["del_blog"] = "0"
		if account == one["account"] {
			row["del_blog"] = "1"
		}

		//每个博客下面显示3条评论
		bc_total, bc_list := t_blog.getBlogCom(ctx, "", blog_id, p_id, "1", "3")
		comment := map[string]interface{}{}
		comment["total"] = bc_total
		comment["list"] = bc_list

		row["comment"] = comment

		t_list = append(t_list, row)
	}

	return t_status, total, t_msg, t_list
}

func (t_blog *BlogThread) BlogComm(ctx *gin.Context, in_param map[string]string) (int, int, string, []map[string]interface{}) {

	t_status := STATUS_SUCCESS
	t_msg := MSG_SUCCESS
	t_list := []map[string]interface{}{}

	blog_id, _ := in_param["blog_id"]
	if blog_id == "" {
		t_status = 155
		t_msg = "博客id不能为空"
		return t_status, 0, t_msg, t_list
	}

	comment_id, _ := in_param["comment_id"]
	if comment_id == "" {
		comment_id = "0"
	}

	page, _ := in_param["page"]
	page_size, _ := in_param["page_size"]

	bc_total, bc_list := t_blog.getBlogCom(ctx, "", blog_id, comment_id, page, page_size)

	for _, one := range bc_list {

		row := map[string]interface{}{}
		for k, v := range one { // 没有办法 row = one; 只能一个个赋值
			row[k] = v
		}

		if len(one["p_id"]) < 2 {

			reply_total, reply_list := t_blog.getBlogCom(ctx, "", blog_id, one["id"], "1", "3")

			reply := map[string]interface{}{}
			reply["total"] = reply_total
			reply["list"] = reply_list

			row["reply"] = reply
		}

		t_list = append(t_list, row)
	}

	return t_status, bc_total, t_msg, t_list
}

func (t_blog *BlogThread) BlogCommReply(ctx *gin.Context, in_param map[string]string) (int, string, map[string]string) {

	t_status := STATUS_SUCCESS
	t_msg := MSG_SUCCESS
	t_list := map[string]string{}

	content := in_param["content"]
	content_len := ZhongWenChang(content)
	if content_len < 3 {
		t_status = 153
		t_msg = "字数不能小于3个字"
		return t_status, t_msg, t_list
	}

	if content_len > 500 {
		t_status = 154
		t_msg = "字数不能大于500个字"
		return t_status, t_msg, t_list
	}

	blog_id, _ := in_param["blog_id"]
	if blog_id == "" {
		t_status = 155
		t_msg = "博客id不能为空"
		return t_status, t_msg, t_list
	}

	var m_blog_list model.BlogList
	var m_blog_comm model.BlogComm
	var m_user_info model.UserInfo

	blog_list_one := m_blog_list.RedisGetOne("id", blog_id)
	if blog_list_id, _ := blog_list_one["id"]; blog_list_id == "" {
		t_status = 156
		t_msg = "博客id错误"
		return t_status, t_msg, t_list
	}

	account := sessInfo("account", ctx)
	user_info := m_user_info.UserRedis(account)
	if user_info_id, _ := user_info["id"]; user_info_id == "" {
		t_status = 157
		t_msg = "用户信息错误"
		return t_status, t_msg, t_list
	}
	nick_name := user_info["nick_name"]

	is_up := "0"
	if blog_list_one["account"] == account {
		is_up = "1"
	}

	sql_arr := []string{}

	to_account := ""
	to_name := ""
	sys_account := ""
	p_id, _ := in_param["comment_id"]
	if p_id == "" { //为空=评论博客
		p_id = "0"

	} else { //不为空=回复评论

		blog_comm_one := m_blog_comm.RedisGetOne("id", p_id)
		blog_comm_id, _ := blog_comm_one["id"]
		if blog_comm_id == "" {
			t_status = 158
			t_msg = "评论id错误"
			return t_status, t_msg, t_list
		}

		//回复谁就对谁发通知;如果是在祖先评论里面评论他人,不会对祖先评论发通知
		sys_account = blog_comm_one["account"]

		// 如果回复的是2级评论,所有这条线的祖先都是这个2级评论
		if blog_comm_one["p_id"] == "0" { //赋予祖先

			p_id = blog_comm_id

		} else { //继承祖先
			p_id = blog_comm_one["p_id"]

			// 回复祖先评论不需要 to_account
			to_account = sys_account
			to_name = blog_comm_one["nick_name"]
		}

		update_blog_comm := fmt.Sprintf("UPDATE blog_comm SET reply_num = reply_num + 1 WHERE id = '%s' ", p_id) //更新博客回复数,只+祖先
		sql_arr = append(sql_arr, update_blog_comm)
	}

	update_blog_list := fmt.Sprintf("UPDATE blog_list SET comm_num = comm_num + 1 WHERE id = '%s' ", blog_id) //更新博客评论数

	new_id := model.GetKey(16)
	create_time := time.Now().Format(f_date)

	t_list["id"] = new_id
	t_list["blog_id"] = blog_id
	t_list["account"] = account
	t_list["nick_name"] = nick_name
	t_list["to_account"] = to_account
	t_list["to_name"] = to_name
	t_list["content"] = content
	t_list["p_id"] = p_id
	t_list["like_num"] = "0"
	t_list["reply_num"] = "0"
	t_list["create_time"] = create_time
	t_list["status"] = "2"
	t_list["is_up"] = is_up
	inser_blog_comm := common.InsertSql("blog_comm", t_list) //插入评论

	if len(sys_account) > 1 { //评论回复发送通知,不会给发帖人发送通知
		blog_content := string([]rune(blog_list_one["content"])[:10])
		sub_title := "您在博客\"" + blog_content + "...\"中的评论有人回复了"

		sys_msg_map := map[string]string{}
		sys_msg_map["id"] = new_id
		sys_msg_map["title"] = "\"" + nick_name + "\"回复了您的评论"
		sys_msg_map["sub_title"] = sub_title
		sys_msg_map["content"] = content
		sys_msg_map["create_time"] = create_time
		sys_msg_map["to_user"] = sys_account
		sys_msg_map["msg_type"] = "1"
		//sys_msg_map["x_id"] = blog_id
		inser_sys_msg := common.InsertSql("sys_msg", sys_msg_map)
		sql_arr = append(sql_arr, inser_sys_msg)
	}

	//获取发帖人的头像
	avatar, _ := user_info["avatar"]
	t_list["is_like"] = "0"

	//保持结构一致
	t_list["avatar"] = avatar
	delete(t_list, "status")

	sql_arr = append(sql_arr, update_blog_list)
	sql_arr = append(sql_arr, inser_blog_comm)

	err := model.Trans(sql_arr)
	if err != nil {
		t_status = 159
		t_msg = "评论失败"
		return t_status, t_msg, t_list
	}

	return t_status, t_msg, t_list
}

func (t_blog *BlogThread) BlogLike(ctx *gin.Context, in_param map[string]string) (int, string, int, int) {

	t_status := STATUS_SUCCESS
	t_msg := MSG_SUCCESS
	is_like := 0
	like_num := 0

	blog_id := in_param["blog_id"]

	if len(blog_id) < 1 {
		t_status = 100
		t_msg = "博客id不能为空"
		return t_status, t_msg, is_like, like_num
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_status = 131
		t_msg = "用户信息异常"
		return t_status, t_msg, is_like, like_num
	}

	var m_bloglist model.BlogList
	bl_map := m_bloglist.RedisGetOne("id", blog_id)
	if r_id, _ := bl_map["id"]; r_id == "" {
		t_status = 132
		t_msg = "博客id错误"
		return t_status, t_msg, is_like, like_num
	}

	var model_l model.BlogListLike
	c_info := model_l.RedisIsLike(blog_id, account)
	like_sql := ""
	blog_sql := ""
	if len(c_info) < 1 || len(c_info["id"]) < 1 {
		//没有点赞，则点赞，并且添加点赞次数
		is_like = 1
		id := model.GetKey(16)
		now_time := time.Now().Format(f_date)
		like_sql = fmt.Sprintf("insert into blog_list_like (id,blog_list_id,account,create_time) VALUES ('%s','%s','%s','%s');", id, blog_id, account, now_time)
		blog_sql = fmt.Sprintf("update blog_list set like_num=like_num+1 where id='%s';", blog_id)
	} else {
		like_sql = fmt.Sprintf("delete from blog_list_like where id='%s';", c_info["id"])
		blog_sql = fmt.Sprintf("update blog_list set like_num=like_num-1 where id='%s';", blog_id)
	}

	sql_arr := []string{like_sql, blog_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		is_like = 0
		t_status = 100
		t_msg = "点赞失败"
		return t_status, t_msg, is_like, like_num
	}

	//获取点赞数
	if bl_like_num, _ := bl_map["like_num"]; bl_like_num != "" {
		like_num, _ = common.Str2Int(bl_like_num)
		if is_like == 1 {
			like_num = like_num + 1
		} else {
			like_num = like_num - 1
			if like_num < 0 {
				update_sql := fmt.Sprintf("update blog_list set like_num=0 where id='%s';", blog_id)
				model.Query(update_sql)
				like_num = 0
			}
		}
	}

	//清除缓存
	model_l.RedisIsLikeDel(blog_id, account)
	m_bloglist.RedisGetOneDel("id", blog_id)

	return t_status, t_msg, is_like, like_num
}

func (t_blog *BlogThread) CommLike(ctx *gin.Context, in_param map[string]string) (int, int, string, int) {
	t_status := 100
	t_msg := "参数不足"
	is_like := 0
	like_num := 0

	comm_id := in_param["comm_id"]

	if len(comm_id) < 1 {
		t_msg = "评论id不能为空"
		return t_status, is_like, t_msg, like_num
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, is_like, t_msg, like_num
	}

	var model_bc model.BlogComm
	c_info := model_bc.RedisGetOne("id", comm_id)
	if len(c_info) < 1 || len(c_info["id"]) < 1 {
		t_msg = "评论异常"
		return t_status, is_like, t_msg, like_num
	}

	//判断评论是否已经点赞
	var model_cl model.BlogCommLike
	like_info := model_cl.RedisIsLike(comm_id, account)
	like_sql := ""
	comm_sql := ""
	now_time := time.Now().Format(f_date)

	if len(like_info) < 1 || len(like_info["id"]) < 1 {
		//没有点赞，则点赞，并且添加点赞次数
		is_like = 1
		id := model.GetKey(16)
		like_sql = fmt.Sprintf("insert into blog_comm_like (id,comm_id,account,create_time) VALUES ('%s','%s','%s','%s');", id, comm_id, account, now_time)
		comm_sql = fmt.Sprintf("update blog_comm set like_num=like_num+1 where id='%s';", comm_id)
	} else {
		like_sql = fmt.Sprintf("delete from blog_comm_like where id='%s';", like_info["id"])
		comm_sql = fmt.Sprintf("update blog_comm set like_num=like_num-1 where id='%s';", comm_id)
	}

	sql_arr := []string{like_sql, comm_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		is_like = 0
		t_msg = "点赞失败"
		return t_status, is_like, t_msg, like_num
	}

	//获取点赞数
	if bl_like_num, _ := c_info["like_num"]; bl_like_num != "" {
		like_num, _ = common.Str2Int(bl_like_num)
		if is_like == 1 {
			like_num = like_num + 1
		} else {
			like_num = like_num - 1
		}
	}

	//清除缓存
	model_cl.RedisIsLikeDel(comm_id, account)
	model_bc.RedisGetOneDel("id", comm_id)

	t_status = 200
	t_msg = "success"

	return t_status, is_like, t_msg, like_num
}

func (t_blog *BlogThread) UploadImgMore(ctx *gin.Context) (int, string, []string) {

	t_status := STATUS_SUCCESS
	t_msg := MSG_SUCCESS
	t_list := []string{}

	dateymd := time.Now().Format("20060102")
	form_name := "img[]"
	folder_nginx := hook.WebConfVal("img_path")
	folder_url := "/blog/" + dateymd
	filter_arr := []string{"jpg", "jpeg", "png", "gif"}
	file_name := ""

	form, _ := ctx.MultipartForm()
	files := form.File[form_name]
	if len(files) > 9 {
		t_msg = "图片数量不能大于9张"
		return t_status, t_msg, t_list
	}

	t_status, t_msg, _, t_list = UploadFileMore(form_name, folder_nginx, folder_url, file_name, "0", filter_arr, ctx)

	if t_status == STATUS_SUCCESS {
		img_url := hook.WebConfVal("img_url")

		for k, v := range t_list {
			t_list[k] = img_url + v
		}
	}

	return t_status, t_msg, t_list
}

func (t_blog *BlogThread) BlogPublish(ctx *gin.Context, in_param map[string]string) (int, string, map[string]string) {

	t_status := 103
	t_msg := "error"
	t_list := map[string]string{}

	var m_user_info model.UserInfo

	account := sessInfo("account", ctx)
	user_info := m_user_info.UserRedis(account)
	if user_info_id, _ := user_info["id"]; user_info_id == "" {
		t_status = 157
		t_msg = "用户信息错误"
		return t_status, t_msg, t_list
	}
	nick_name := user_info["nick_name"]

	content := in_param["content"]
	content_len := ZhongWenChang(content)
	if content_len < 4 {
		t_status = 153
		t_msg = "字数不能小于3个字"
		return t_status, t_msg, t_list
	}

	if content_len > 500 {
		t_status = 154
		t_msg = "字数不能大于500个字"
		return t_status, t_msg, t_list
	}

	img_list := []string{}

	dateymd := time.Now().Format("20060102")
	form_name := "img[]"
	folder_nginx := hook.WebConfVal("img_path")
	folder_url := "/blog/" + dateymd
	filter_arr := []string{"jpg", "jpeg", "png", "gif"}
	file_name := ""

	form, _ := ctx.MultipartForm()
	files := form.File[form_name]
	if len(files) > 9 {
		t_msg = "图片数量不能大于9张"
		return t_status, t_msg, t_list
	}

	if len(files) > 0 {
		t_status, t_msg, _, img_list = UploadFileMore(form_name, folder_nginx, folder_url, file_name, "0", filter_arr, ctx)
		if t_status != STATUS_SUCCESS {
			t_msg = "上传图片失败"
			return t_status, t_msg, t_list
		}
	}

	blog_status := hook.WebConfVal("blog_status")
	if blog_status == "" {
		blog_status = "1" //默认审核
	}

	img_url := hook.WebConfVal("img_url")
	img_arr := []string{}
	for _, v := range img_list {
		img_http := img_url + v
		img_arr = append(img_arr, img_http)
	}

	t_list["id"] = model.GetKey(16)
	t_list["content"] = content
	t_list["img_url"] = common.StructToJson(img_arr)
	t_list["view_num"] = "0"
	t_list["account"] = account
	t_list["nick_name"] = nick_name
	t_list["like_num"] = "0"
	t_list["avatar"] = user_info["avatar"]
	t_list["is_top"] = "0"
	t_list["is_good"] = "0"
	t_list["comm_num"] = "0"
	t_list["create_time"] = time.Now().Format(f_date)
	t_list["status"] = blog_status

	insert_str := common.InsertSql("blog_list", t_list)
	sql_err := model.Query(insert_str)
	if sql_err != nil {
		t_msg = "发帖失败"
		return t_status, t_msg, t_list
	}

	t_status = STATUS_SUCCESS
	t_msg = MSG_SUCCESS

	return t_status, t_msg, t_list
}

func (t_blog *BlogThread) BlogOne(ctx *gin.Context, in_param map[string]string) (int, string, map[string]interface{}) {

	t_status := 330
	t_msg := "error"
	t_list := map[string]interface{}{}

	blog_id, _ := in_param["id"]
	if blog_id == "" {
		t_msg := "博客id不能为空"
		return t_status, t_msg, t_list
	}

	fields := []string{"id", "content", "img_url", "view_num", "account",
		"nick_name", "like_num", "avatar", "is_top", "is_good", "comm_num", "create_time"}

	sql_one := fmt.Sprintf("SELECT id, content, img_url, view_num, account, nick_name, like_num, avatar, is_top, is_good, comm_num, create_time FROM blog_list WHERE id='%s' AND status = 2 limit 1", blog_id)
	one, _ := model.SqlRow(sql_one, fields)
	if len(one) == 0 {
		t_msg := "错误的id"
		return t_status, t_msg, t_list
	}

	account := sessInfo("account", ctx)
	p_id := "0"
	var m_bll model.BlogListLike

	for k, v := range one {
		if k == "img_url" {

			img_arr := []string{}
			json.Unmarshal([]byte(v), &img_arr)
			t_list["img_urls"] = img_arr

		} else {
			t_list[k] = v
		}
	}

	//是否给博客点赞
	is_like := "0"
	is_like_map := m_bll.RedisIsLike(blog_id, account)
	if r_is_like, _ := is_like_map["id"]; r_is_like != "" {
		is_like = "1"
	}
	t_list["is_like"] = is_like

	//每个博客下面显示3条评论
	bc_total, bc_list := t_blog.getBlogCom(ctx, "", blog_id, p_id, "1", "3")
	comment := map[string]interface{}{}
	comment["total"] = bc_total
	comment["list"] = bc_list

	t_list["comment"] = comment

	t_status = STATUS_SUCCESS
	t_msg = MSG_SUCCESS

	return t_status, t_msg, t_list
}

func (t_blog *BlogThread) BlogCommOne(ctx *gin.Context, in_param map[string]string) (int, int, string, map[string]interface{}) {

	t_status := STATUS_SUCCESS
	t_msg := MSG_SUCCESS
	t_one := map[string]interface{}{}
	t_total := 0

	comment_id, _ := in_param["comment_id"]
	if comment_id == "" {
		t_status = 155
		t_msg = "评论id不能为空"
		return t_status, t_total, t_msg, t_one
	}

	page, _ := in_param["page"]
	page_size, _ := in_param["page_size"]

	bc_total, bc_list := t_blog.getBlogCom(ctx, comment_id, "", "", "1", "1")
	if len(bc_list) == 0 {
		return t_status, bc_total, t_msg, t_one
	}

	one := bc_list[0]

	for k, v := range one { // 没有办法 t_one = one; 只能一个个赋值
		t_one[k] = v
	}

	if len(one["p_id"]) < 2 {

		bc_total, bc_list = t_blog.getBlogCom(ctx, "", one["blog_id"], one["id"], page, page_size)

		reply := map[string]interface{}{}
		reply["total"] = bc_total
		reply["list"] = bc_list

		t_one["reply"] = reply
	}

	return t_status, bc_total, t_msg, t_one
}

/**
*    删除博客
 */
func (t_blog *BlogThread) DelBlog(in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "博客ID异常"

	if len(in_param["blog_id"]) < 1 {
		t_msg = "博客ID不能为空"
		return t_status, t_msg
	}

	account := sessInfo("account", ctx)
	if account == "" {
		t_msg = "用户信息异常"
		return t_status, t_msg
	}
	// fields := []string{"id"}
	// select_sql := fmt.Sprintf("select id from blog_comm where blog_id='%s';", blog_id)
	// comm_map, _ := model.SqlRows(select_sql, fields)
	// if len(comm_map) > 0 {
	// 	for _, c_val := range comm_map {
	// 		bl_sql := fmt.Sprintf("delete from blog_comm_like where comm_id='%s';", c_val["id"])
	// 		err := model.Query(bl_sql)
	// 		if err != nil {
	// 			t_status = 100
	// 			t_msg = "删除评论点赞失败"
	// 			return t_status, t_msg
	// 		}
	// 	}
	// }

	list_sql := fmt.Sprintf("delete from blog_list where id='%s';", in_param["blog_id"])
	comm_sql := fmt.Sprintf("delete from blog_comm where blog_id='%s';", in_param["blog_id"])
	//msg_sql := fmt.Sprintf("update sys_msg set status=%d where x_id='%s';", 0, in_param["blog_id"])
	like_sql := fmt.Sprintf("delete from blog_list_like where blog_list_id='%s';", in_param["blog_id"])

	sql_arr := []string{list_sql, comm_sql, like_sql}
	err := model.Trans(sql_arr)
	if err != nil {
		t_msg = "删除失败"
		return t_status, t_msg
	}

	//清除缓存
	var m_bloglist model.BlogList
	m_bloglist.RedisGetOneDel("id", in_param["blog_id"])

	var m_blogcomm model.BlogComm
	m_blogcomm.RedisGetOneDel("blog_id", in_param["blog_id"])

	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}
