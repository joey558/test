package controller

import (
	"qzapp/hook"

	"github.com/gin-gonic/gin"
)

/**
* 获取博客列表
 */
func (c_blog *BlogController) BlogList(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	key_arr := []string{"is_top", "is_good", "content", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_data["total"], c_msg, c_data["list"] = t_blog.BlogList(c, in_param)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 获取博客列表的全部评论
 */
func (c_blog *BlogController) BlogComm(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	key_arr := []string{"blog_id", "comment_id", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_data["total"], c_msg, c_data["list"] = t_blog.BlogComm(c, in_param)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 评论博客,或者回复评论
 */
func (c_blog *BlogController) BlogCommReply(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	key_arr := []string{"blog_id", "comment_id", "content"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, c_data["info"] = t_blog.BlogCommReply(c, in_param)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 博客 点赞/取消
 */
func (c_blog *BlogController) BlogLike(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	key_arr := []string{"blog_id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, c_data["is_like"], c_data["like_num"] = t_blog.BlogLike(c, in_param)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 评论 点赞/取消
 */
func (c_blog *BlogController) CommLike(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	key_arr := []string{"comm_id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_data["is_like"], c_msg, c_data["like_num"] = t_blog.CommLike(c, in_param)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 批量上传图片
 */
func (c_blog *BlogController) UploadImgMore(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	c_status, c_msg, c_data["img_list"] = t_blog.UploadImgMore(c)
	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 发帖
 */
func (c_blog *BlogController) BlogPublish(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	key_arr := []string{"content"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, c_data["info"] = t_blog.BlogPublish(c, in_param)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 获取某个博客信息
 */
func (c_blog *BlogController) BlogOne(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	key_arr := []string{"id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg, c_data["info"] = t_blog.BlogOne(c, in_param)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 获取某个评论的信息
 */
func (c_blog *BlogController) BlogCommOne(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	key_arr := []string{"comment_id", "page", "page_size"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_data["total"], c_msg, c_data["list"] = t_blog.BlogCommOne(c, in_param)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}

/**
* 删除博客
 */
func (c_blog *BlogController) DelBlog(c *gin.Context) {
	c_status := 100
	c_msg := "请求成功"
	c_data := map[string]interface{}{}

	key_arr := []string{"blog_id"}
	c_status, c_msg, in_param := hook.AuthInput(key_arr, c)
	if c_status == 200 {
		c_status, c_msg = t_blog.DelBlog(in_param, c)
	}

	c.JSON(http_status, &JsonOut{Status: c_status, Msg: c_msg, Data: c_data})
}
