package thread

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"qzapp/hook"
	"qzapp/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gitlab.stagingvip.net/publicGroup/public/common"
	"gitlab.stagingvip.net/publicGroup/public/redis"
)

var f_date string = "2006-01-02 15:04:05"
var f_day string = "2006-01-02"

const (
	STATUS_SUCCESS = int(200)
	MSG_SUCCESS    = "success"
)

type Public struct{}

type User struct{}

type TaskThread struct{}
type Finance struct{}
type BlogThread struct{}

type FileMoreChan struct {
	Status       int
	Msg          string
	CompletePath string
	UrlPath      string
}

//配置数据的缓存key集合
var Conf_Redis_Key = "Conf_Redis_Key"

//数据缓存的key集合
var Data_Redis_Key = "Data_Redis_Key"

var push_header map[string]string

var log_path string

func init() {
	conf_byte, err := common.ReadFile("./conf/conf.json")

	if err != nil {
		panic(err)
	}
	var json_conf map[string]string
	//解析json格式r
	err = json.Unmarshal(conf_byte, &json_conf)
	if err != nil {
		panic(err)
	}
	log_path = json_conf["log_path"]
	push_header = make(map[string]string)
	push_header["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	push_header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36"
}

/**
*  处理分页
 */
func treadPage(page, page_size string) (int, int) {
	page_int, _ := strconv.Atoi(page)
	if page_int < 1 {
		page_int = 1
	}
	size_int, _ := strconv.Atoi(page_size)
	if size_int < 1 {
		size_int = 20
	} else if size_int > 100 {
		size_int = 100
	}
	return page_int, size_int
}

/**
*  删除session缓存
 */
func delSession(session_id string) {
	key_str := fmt.Sprintf("user_list:session_id:%s", session_id)
	redis.RediGo.KeyDel(key_str)
}

/**
*  用户日志
 */
func UserLog(status int, account, title, content, msg string, in_param map[string]string, ctx *gin.Context) (int, string) {
	t_status := 100
	t_msg := "参数错误"
	if len(account) < 1 {
		session := sessions.Default(ctx)
		account = fmt.Sprintf("%v", session.Get("account"))
		if account == "nil" || account == "" {
			t_msg = "缺少用户名"
			return t_status, t_msg
		}
	}

	u_map := map[string]string{}
	u_map["account"] = account
	u_map["title"] = title
	u_map["content"] = content
	u_map["res"] = msg
	u_map["status"] = fmt.Sprintf("%d", status)
	u_map["param"] = common.Interface2Json(in_param)
	u_map["id"] = model.GetKey(16)
	u_map["ip"] = ctx.ClientIP()
	u_map["url"] = ctx.Request.RequestURI
	u_map["create_time"] = time.Now().Format(f_date)
	u_map["host"] = hook.Host(ctx.Request.Host)
	table_name := "user_log"
	in_sql := common.InsertSql(table_name, u_map)
	err := model.Query(in_sql)
	if err != nil {
		t_msg = err.Error()
		return t_status, t_msg
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg
}

func sessInfo(sess_key string, ctx *gin.Context) string {
	sess_val := ""
	session := sessions.Default(ctx)
	sess := session.Get(sess_key)
	if sess != nil {
		sess_val = fmt.Sprintf("%v", sess)
	}

	return sess_val
}

/**
*  用户看视频的次数
* @vid 视频id
* @score 视频所需积分
 */
func (u *User) userViewNum_bak(account, vid, score string) int {
	v_num := 0

	//当前用户的积分
	var model_ui model.UserInfo
	u_info := model_ui.User(account)
	if len(u_info.Account) < 1 {
		return v_num
	}

	now_time := time.Now().Format(f_date)
	score_int, _ := strconv.Atoi(score)
	if score_int < 1 {
		code := "score"
		val := hook.WebConfVal(code)
		score_int, _ = strconv.Atoi(val)
	}
	if score_int == 0 {
		v_num = 1
	}

	sql_arr := []string{}

	// //更新视频观看次数
	// num_sql := fmt.Sprintf("update video_list set view_num=view_num+1 where id='%s';", vid)
	// sql_arr = append(sql_arr[0:], num_sql)

	//判断用户今天是否看过当前视频
	today_day := time.Now().Format(f_day) + " 00:00:00"
	his_sql := fmt.Sprintf("select id,account from video_history where video_id='%s' and account='%s' and create_time>='%s' limit 1;", vid, account, today_day)
	his_field := []string{"id", "account"}
	his_map, _ := model.SqlRow(his_sql, his_field)
	is_his := 0 //1=用户今天看过视频 不扣积分
	id := model.GetKey(16)
	if len(his_map) > 0 && len(his_map["account"]) > 3 {
		is_his = 1
		//存在观看记录，则更新时间
		up_sql := fmt.Sprintf("update video_history set create_time='%s' where id='%s';", now_time, his_map["id"])
		sql_arr = append(sql_arr[0:], up_sql)
	} else {
		//不存在观看记录，则新增观看记录
		t_name := "video_history"
		d_map := map[string]string{}

		d_map["id"] = id
		d_map["account"] = account
		d_map["video_id"] = vid
		d_map["create_time"] = now_time
		in_sql := common.InsertSql(t_name, d_map)
		sql_arr = append(sql_arr[0:], in_sql)
	}

	//如果视频免费免费或者今天已经看过视频，则免费，不扣除积分
	if score_int == 0 || is_his == 1 {

		v_num = u_info.Score / 1

		if score_int == 0 { //免费视频的剩余观看次数写死
			v_num = 9999
		}

		if is_his == 1 && score_int > 0 { //已经看过的视频,剩余观看次数正常显示
			v_num = u_info.Score / score_int
		}

		err := model.Trans(sql_arr)
		if err != nil {
			v_num = 0
			return v_num
		}
		return v_num
	}

	//积分不够
	if u_info.Score < score_int {
		v_num = 0
		return v_num
	}

	var model_ul model.UserList
	u_list := model_ul.UserRedis(account)
	v_num = u_info.Score / score_int

	//消费积分
	o_name := "order_list"
	o_map := map[string]string{}
	o_map["id"] = id
	o_map["account"] = account
	o_map["status"] = "3"
	o_map["order_number"] = id
	o_map["agent_name"] = u_list["agent_name"]
	o_map["order_type"] = "6"
	o_map["is_test"] = strconv.Itoa(u_info.Is_test)
	o_map["agent_path"] = u_list["agent_path"]
	o_map["create_time"] = now_time
	o_map["is_auto"] = "1"
	o_map["note"] = "观看视频:" + vid
	o_map["score"] = strconv.Itoa(score_int)
	o_map["creator"] = account
	o_map["operator"] = account
	o_map["pay_time"] = now_time
	o_sql := common.InsertSql(o_name, o_map)
	sql_arr = append(sql_arr[0:], o_sql)

	//账变记录
	after_score := u_info.Score - score_int
	a_name := "amount_list"
	a_map := map[string]string{}
	a_map["id"] = id
	a_map["account"] = account
	a_map["order_number"] = id
	a_map["agent_name"] = u_list["agent_name"]
	a_map["order_type"] = "6"
	a_map["is_test"] = strconv.Itoa(u_info.Is_test)
	a_map["agent_path"] = u_list["agent_path"]
	a_map["before_score"] = fmt.Sprintf("%d", u_info.Score)
	a_map["after_score"] = fmt.Sprintf("%d", after_score)
	a_map["create_time"] = now_time
	a_map["note"] = "观看视频:" + vid
	a_map["score"] = strconv.Itoa(score_int)
	a_sql := common.InsertSql(a_name, a_map)
	sql_arr = append(sql_arr[0:], a_sql)

	//更新积分
	user_sql := fmt.Sprintf("update user_info set score=score-%d where account='%s';", score_int, account)
	sql_arr = append(sql_arr[0:], user_sql)
	err := model.Trans(sql_arr)
	if err != nil {
		v_num = 0
	}
	return v_num
}

func (u *User) userViewNum(account, vid, score string) (int, int) {
	user_score := 0  //用户剩余积分
	video_score := 1 //视频需要的积分

	//当前用户的积分
	var model_ui model.UserInfo
	u_info := model_ui.User(account)
	if len(u_info.Account) < 1 {
		return user_score, video_score
	}
	user_score = u_info.Score

	now_time := time.Now().Format(f_date)
	video_score, _ = strconv.Atoi(score)
	if video_score < 1 {
		val := hook.WebConfVal("score")
		video_score, _ = strconv.Atoi(val)
	}

	sql_arr := []string{}

	// //更新视频观看次数
	// num_sql := fmt.Sprintf("update video_list set view_num=view_num+1 where id='%s';", vid)
	// sql_arr = append(sql_arr[0:], num_sql)

	//判断用户今天是否看过当前视频
	today_day := time.Now().Format(f_day) + " 00:00:00"
	his_sql := fmt.Sprintf("select id,account from video_history where video_id='%s' and account='%s' and create_time>='%s' limit 1;", vid, account, today_day)
	his_field := []string{"id", "account"}
	his_map, _ := model.SqlRow(his_sql, his_field)
	is_his := 0 //1=用户今天看过视频 不扣积分
	id := model.GetKey(16)
	if len(his_map) > 0 && len(his_map["account"]) > 3 {
		is_his = 1
		//存在观看记录，则更新时间
		up_sql := fmt.Sprintf("update video_history set create_time='%s' where id='%s';", now_time, his_map["id"])
		sql_arr = append(sql_arr[0:], up_sql)
	} else {
		//不存在观看记录，则新增观看记录
		t_name := "video_history"
		d_map := map[string]string{}

		d_map["id"] = id
		d_map["account"] = account
		d_map["video_id"] = vid
		d_map["create_time"] = now_time
		in_sql := common.InsertSql(t_name, d_map)
		sql_arr = append(sql_arr[0:], in_sql)
	}

	//如果视频不要积分 || 今天已经看过视频，不扣除积分
	if video_score == 0 || is_his == 1 {
		go model.Trans(sql_arr)
		video_score = 0
		return user_score, video_score
	}

	//积分不够
	if user_score < video_score {
		return user_score, video_score
	}

	var model_ul model.UserList
	u_list := model_ul.UserRedis(account)

	//消费积分
	o_name := "order_list"
	o_map := map[string]string{}
	o_map["id"] = id
	o_map["account"] = account
	o_map["status"] = "3"
	o_map["order_number"] = id
	o_map["agent_name"] = u_list["agent_name"]
	o_map["order_type"] = "6"
	o_map["is_test"] = strconv.Itoa(u_info.Is_test)
	o_map["agent_path"] = u_list["agent_path"]
	o_map["create_time"] = now_time
	o_map["is_auto"] = "1"
	o_map["note"] = "观看视频:" + vid
	o_map["score"] = strconv.Itoa(video_score)
	o_map["creator"] = account
	o_map["operator"] = account
	o_map["pay_time"] = now_time
	o_sql := common.InsertSql(o_name, o_map)
	sql_arr = append(sql_arr[0:], o_sql)

	//账变记录
	after_score := user_score - video_score
	a_name := "amount_list"
	a_map := map[string]string{}
	a_map["id"] = id
	a_map["account"] = account
	a_map["order_number"] = id
	a_map["agent_name"] = u_list["agent_name"]
	a_map["order_type"] = "6"
	a_map["is_test"] = strconv.Itoa(u_info.Is_test)
	a_map["agent_path"] = u_list["agent_path"]
	a_map["before_score"] = fmt.Sprintf("%d", user_score)
	a_map["after_score"] = fmt.Sprintf("%d", after_score)
	a_map["create_time"] = now_time
	a_map["note"] = "观看视频:" + vid
	a_map["score"] = strconv.Itoa(video_score)
	a_sql := common.InsertSql(a_name, a_map)
	sql_arr = append(sql_arr[0:], a_sql)

	//更新积分
	user_sql := fmt.Sprintf("update user_info set score=score-%d where account='%s';", video_score, account)
	sql_arr = append(sql_arr[0:], user_sql)
	err := model.Trans(sql_arr)
	if err != nil {
		return 0, 1
	}
	return user_score, video_score
}

/**
*  上传文件

func uploadFile(file_name, form_name, file_path string, filter_arr []string, ctx *gin.Context) (int, string, string) {
	t_status := 100
	t_msg := "error"

	file, header, err := ctx.Request.FormFile(form_name)
	if err != nil {
		t_msg = "获取文件失败"
		return t_status, t_msg, ""
	}
	tmp_file := header.Filename
	//获取文件类型
	tmp_arr := strings.Split(tmp_file, ".")
	tmp_len := len(tmp_arr)
	if tmp_len < 1 {
		t_msg = "文件名异常"
		return t_status, t_msg, ""
	}
	file_type := tmp_arr[tmp_len-1]
	if !common.Arr_In(filter_arr, file_type) {
		t_msg = "文件类型不合法"
		return t_status, t_msg, ""
	}
	//获取文件名
	if file_name == "" {
		file_name = model.GetKey(15) + "." + file_type
	} else {
		file_name = file_name + "." + file_type
	}

	file_path_str := file_path + file_name

	//写入文件
	out, err := os.Create(file_path_str)
	if err != nil {
		t_msg = "文件生成失败error->" + err.Error()
		return t_status, t_msg, file_name
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		t_msg = "文件上传失败error->" + err.Error()
		return t_status, t_msg, ""
	}
	t_status = 200
	t_msg = "success"
	return t_status, t_msg, file_name
}
*/

/**
*  上传文件
* @form_name		<input type="file" name="form_name">
* @folder_nginx		Nginx配置的目录, 不需要带到url里面
* @folder_url		相对目录, Nginx配置的链接 + folder_url + form_name = http地址
* @is_exist			1=强制覆盖,可空
* @filter_arr		文件类型
 */
func UploadFile(form_name, folder_nginx, folder_url, file_name, is_exist string, filter_arr []string, ctx *gin.Context) (int, string, string, string) {
	t_status := 100
	t_msg := "error"
	complete_path := "" //完整地址
	url_path := ""      //http访问地址

	if folder_url == "" {
		t_msg = "未命名文件上传主目录"
		return t_status, t_msg, complete_path, url_path
	}

	file, header, err := ctx.Request.FormFile(form_name)
	if err != nil {
		t_msg = "获取文件失败"
		return t_status, t_msg, complete_path, url_path
	}

	//获取文件类型
	tmp_file := header.Filename
	tmp_arr := strings.Split(tmp_file, ".")
	tmp_len := len(tmp_arr)
	if tmp_len < 1 {
		t_msg = "文件名异常"
		return t_status, t_msg, complete_path, url_path
	}
	file_type := tmp_arr[tmp_len-1]
	if !common.Arr_In(filter_arr, file_type) {
		t_msg = "文件类型不合法"
		return t_status, t_msg, complete_path, url_path
	}

	//获取文件名
	if file_name == "" {
		file_name = model.GetKey(15) + "." + file_type
	} else {
		file_name = file_name + "." + file_type
	}

	folder := folder_nginx + folder_url //完整目录
	os.MkdirAll(folder, os.ModePerm)    //创建目录

	complete_path = folder + "/" + file_name
	url_path = folder_url + "/" + file_name

	//判断文件是否存在
	if is_exist != "1" && common.IsExist(complete_path) {
		t_status = 300
		t_msg = "文件已存在"
		return t_status, t_msg, complete_path, url_path
	}

	//写入文件
	out, err := os.Create(complete_path)
	if err != nil {
		t_msg = "文件生成失败error->" + err.Error()
		return t_status, t_msg, complete_path, url_path
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		t_msg = "文件上传失败error->" + err.Error()
		return t_status, t_msg, complete_path, url_path
	}

	t_status = 200
	t_msg = "success"
	return t_status, t_msg, complete_path, url_path
}

/**
*  上传文件
* @form_name[]		<input type="file" name="form_name[]">
* @folder_nginx		Nginx配置的目录, 不需要带到url里面
* @folder_url		相对目录, Nginx配置的链接 + folder_url + form_name = http地址
* @is_exist			1=强制覆盖,可空
* @filter_arr		文件类型
 */
func UploadFileMore(form_name, folder_nginx, folder_url, file_name, is_exist string, filter_arr []string, ctx *gin.Context) (int, string, []string, []string) {

	t_status := 100
	t_msg := "error"
	complete_path_list := make([]string, 0) //完整地址
	url_path_list := make([]string, 0)      //http访问地址

	if folder_url == "" {
		t_msg = "未命名文件上传主目录"
		return t_status, t_msg, complete_path_list, url_path_list
	}

	up := make(chan FileMoreChan)
	isStart := make(chan bool)
	defer func() {
		close(up)
		close(isStart)
	}()
	form, _ := ctx.MultipartForm()
	files := form.File[form_name]

	//启动
	go func() {
		isStart <- true
	}()
	go func() {
		for _, file := range files {
			select {
			case b := <-isStart:
				if b {
					go func(file *multipart.FileHeader) {
						// t_status, t_msg, _, filepath := UploadFileOne(file, folder)
						t_status, t_msg, complete_path, url_path := UploadFileOne(file, folder_nginx, folder_url, file_name, is_exist, filter_arr)
						up <- FileMoreChan{Status: t_status, Msg: t_msg, CompletePath: complete_path, UrlPath: url_path}
					}(file)
				}
			}
		}
	}()
	//接收up结构
	for upC := range up {
		if upC.Status != STATUS_SUCCESS {

			isStart <- false

			t_status = upC.Status
			t_msg = upC.Msg
			return t_status, t_msg, complete_path_list, url_path_list

		} else {
			//所接收文件没有错误提示
			complete_path_list = append(complete_path_list, upC.CompletePath)
			url_path_list = append(url_path_list, upC.UrlPath)
			if len(complete_path_list) != len(files) {
				isStart <- true
			} else {
				break
			}
		}
	}

	t_status = STATUS_SUCCESS
	t_msg = MSG_SUCCESS
	return t_status, t_msg, complete_path_list, url_path_list
}

/**
*  上传文件
* @form_name		<input type="file" name="form_name">
* @folder_nginx		Nginx配置的目录, 不需要带到url里面
* @folder_url		相对目录, Nginx配置的链接 + folder_url + form_name = http地址
* @is_exist			1=强制覆盖,可空
* @filter_arr		文件类型
 */
func UploadFileOne(file_header *multipart.FileHeader, folder_nginx, folder_url, file_name, is_exist string, filter_arr []string) (int, string, string, string) {

	t_status := 100
	t_msg := "error"
	complete_path := "" //完整地址
	url_path := ""      //http访问地址

	if folder_url == "" {
		t_msg = "未命名文件上传主目录"
		return t_status, t_msg, complete_path, url_path
	}

	// file, header, err := ctx.Request.FormFile(form_name)
	file, err := file_header.Open()
	if err != nil {
		t_msg = "获取文件失败"
		return t_status, t_msg, complete_path, url_path
	}

	//获取文件类型
	// tmp_file := header.Filename
	tmp_file := file_header.Filename
	tmp_arr := strings.Split(tmp_file, ".")
	tmp_len := len(tmp_arr)
	if tmp_len < 1 {
		t_msg = "文件名异常"
		return t_status, t_msg, complete_path, url_path
	}
	file_type := tmp_arr[tmp_len-1]
	if !common.Arr_In(filter_arr, file_type) {
		t_msg = "文件类型不合法"
		return t_status, t_msg, complete_path, url_path
	}

	//获取文件名
	if file_name == "" {
		file_name = model.GetKey(15) + "." + file_type
	} else {
		file_name = file_name + "." + file_type
	}

	folder := folder_nginx + folder_url //完整目录
	os.MkdirAll(folder, os.ModePerm)    //创建目录

	complete_path = folder + "/" + file_name
	url_path = folder_url + "/" + file_name

	//判断文件是否存在
	if is_exist != "1" && common.IsExist(complete_path) {
		t_status = 300
		t_msg = "文件已存在"
		return t_status, t_msg, complete_path, url_path
	}

	//写入文件
	out, err := os.Create(complete_path)
	if err != nil {
		t_msg = "文件生成失败error->" + err.Error()
		return t_status, t_msg, complete_path, url_path
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		t_msg = "文件上传失败error->" + err.Error()
		return t_status, t_msg, complete_path, url_path
	}

	t_status = STATUS_SUCCESS
	t_msg = MSG_SUCCESS
	return t_status, t_msg, complete_path, url_path
}
