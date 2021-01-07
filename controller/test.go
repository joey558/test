package controller

import (
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	c.HTML(200, "index.tpl", gin.H{
		"title": "app后端接口调试",
	})
}
