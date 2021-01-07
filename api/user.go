package api

/**
*  注册接口
 */
func Register() (int, string, string) {
	a_status := 100
	a_msg := "注册失败"
	token := ""
	return a_status, a_msg, token
}

/**
*  登录KOK接口
 */
func Login() (int, string, string) {
	a_status := 100
	a_msg := "登录失败"
	token := ""
	return a_status, a_msg, token
}

/**
*  查询用户名是否存在
 */
func IsExist() (int, string) {
	a_status := 100
	a_msg := "用户名存在"
	return a_status, a_msg
}

/**
*  查询拥有的平台类型
 */
func GameType() (int, string, []map[string]string) {
	a_status := 100
	a_msg := "平台类型为空"
	game_type := []map[string]string{}
	return a_status, a_msg, game_type
}

/**
*  查询拥有的平台列表
 */
func PlatList() (int, string, []map[string]interface{}) {
	a_status := 100
	a_msg := "平台为空"
	plat_list := []map[string]interface{}{}
	return a_status, a_msg, plat_list
}

/**
*  额度转入KOK
 */
func Transfer() (int, string) {
	a_status := 100
	a_msg := "转入失败"
	return a_status, a_msg
}

/**
*  额度转入KOK确认结果
 */
func TransferConfirm() (int, string) {
	a_status := 100
	a_msg := "转入确认失败"
	return a_status, a_msg
}

/**
*  登录游戏场馆
 */
func GameLogin() (int, string, string) {
	a_status := 100
	a_msg := "游戏登录失败"
	login_url := ""
	return a_status, a_msg, login_url
}

/**
*  查询KOK钱包余额
 */
func Balance() (int, string, float64) {
	a_status := 100
	a_msg := "游戏登录失败"
	balance := 0.00
	return a_status, a_msg, balance
}

/**
*  获取KOK的H5页面地址
 */
func H5Url() (int, string, string) {
	a_status := 100
	a_msg := "获取页面地址失败"
	h5_url := ""
	return a_status, a_msg, h5_url
}

/**
*  获取KOK钱包中心的页面
*  1、该页面可以直接将kok的钱包转入各个场馆
*  2、该页面可以直接进行充值操作
*  3、该页面可以查询交易记录(充值，提现等相关操作)
 */
func FinanceUrl() (int, string, string) {
	a_status := 100
	a_msg := "获取页面地址失败"
	fin_url := ""
	return a_status, a_msg, fin_url
}

/**
*  退出KOK内嵌页面
 */
func Logout() (int, string) {
	a_status := 100
	a_msg := "退出失败"
	return a_status, a_msg
}

/**
*  获取彩票页面
 */
func LotteryUrl() (int, string, string) {
	a_status := 100
	a_msg := "获取页面地址失败"
	lott_url := ""
	return a_status, a_msg, lott_url
}

/**
*  获取彩种类型
 */
func LotteryType() (int, string, []map[string]string) {
	a_status := 100
	a_msg := "获取页面地址失败"
	lott_type := []map[string]string{}
	return a_status, a_msg, lott_type
}

/**
*  获取彩票彩种列表
 */
func LotteryList() (int, string, []map[string]interface{}) {
	a_status := 100
	a_msg := "彩种为空"
	plat_list := []map[string]interface{}{}
	return a_status, a_msg, plat_list
}

/**
*  获取彩票每个彩种的最新开奖信息
 */
func LotteryOpen() (int, string, map[string]string) {
	a_status := 100
	a_msg := "获取页面地址失败"
	open_info := map[string]string{}
	return a_status, a_msg, open_info
}
