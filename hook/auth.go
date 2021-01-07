package hook

import (
	"qzapp/model"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.stagingvip.net/publicGroup/public/common"
)

func Host(host string) string {
	res := host
	if common.Substr(host, 0, 4) == "www." {
		res = strings.Replace(host, "www.", "", 1)
	}
	return res
}

func AppType(ctx *gin.Context) (int, string) {
	h_status := 500
	app_type := ""
	host := Host(ctx.Request.Host)
	var model_dl model.DomainList
	app_conf := model_dl.DomainRedis(host)
	if len(app_conf) < 1 || len(app_conf["app_type"]) < 1 {
		return h_status, app_type
	}
	h_status = 200
	app_type = app_conf["app_type"]
	return h_status, app_type
}
