package hook

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

/**
* sql注入过滤判断
* @param	str	需要做判断的字符串
* return	bool	如果包含注入关键词输出false，否则true
 */
func isAllow(str string) bool {
	h_res := true
	str = strings.ToLower(str)
	danger_keys := `(?:')|(?:--)|(\b(select|update|delete|insert|trancate|char|into|substr|ascii|declare|exec|count|master|into|drop|execute|script)\b)`
	re, _ := regexp.Compile(danger_keys)
	if re.MatchString(str) {
		h_res = false
	}
	return h_res
}

/**
*  验证参数是否合法
 */
func AuthAllow(code, str string) (int, string) {
	h_status := 200
	h_msg := "success"
	if !isAllow(str) {
		h_status = 100
		h_msg = "参数不合法"
		return h_status, h_msg
	}
	val := WebConfVal(code)
	if val == "" {
		return h_status, h_msg
	}
	if !authReg(str, val) {
		h_status = 100
		h_msg = "参数不合法"
		return h_status, h_msg
	}
	return h_status, h_msg
}

/**
*  正则判断
*  reg  正则表达式
*  str  需要判断的字符串
 */
func authReg(str, reg string) bool {
	res, _ := regexp.MatchString(reg, str)
	return res
}

func AuthInput(key_arr []string, ctx *gin.Context) (int, string, map[string]string) {
	h_status := 100
	h_msg := "参数不合法"
	val_map := map[string]string{}
	if len(key_arr) < 1 {
		h_status = 200
		h_msg = "success"
		return h_status, h_msg, val_map
	}
	for _, k_val := range key_arr {
		val := ctx.PostForm(k_val)
		if val == "" {
			continue
		}
		if !isAllow(val) {
			return h_status, h_msg, val_map
		}
		val_map[k_val] = val
	}
	h_status = 200
	h_msg = "success"
	return h_status, h_msg, val_map
}

func AuthInputForArr(input_arr ...[]string) (int, string) {

	h_status := 200
	h_msg := "success"

	if len(input_arr) == 0 {
		return h_status, h_msg
	}

	for _, input_row := range input_arr {

		for _, input_row_v := range input_row {

			if input_row_v == "" {
				continue
			}

			if !isAllow(input_row_v) {
				h_status = 100
				h_msg = "参数不合法"
				return h_status, h_msg
			}
		}
	}

	return h_status, h_msg
}
