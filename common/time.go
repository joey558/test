package common

import (
	"time"
)

var Loc, _ = time.LoadLocation("Local")
var start_day string = "2006-01-02 00:00:00"

/**
*  日期格式转换
*  @date_str  转换前的日期字符串
*  @bef_format  转换前的日期格式
*  @aft_format  转换后的日期格式
 */
func DateFormat(date_str, bef_format, aft_format string) string {
	date_time, _ := time.ParseInLocation(bef_format, date_str, Loc)
	format_str := date_time.Format(aft_format)
	return format_str
}

/**
*  日期范围
*  @date_str  指定的日期字符串
*  @format  日期格式
*  @start  向前推几天,0:表示今天,负数:表示往前推,正数表示向后推
*  @end  向后推几天,0:表示今天,负数:表示往前推,正数表示向后推
 */
func DateExtent(date_str, format string, start, end int) []string {
	date_arr := []string{}
	if start >= end {
		return date_arr
	}
	date_time, _ := time.ParseInLocation(format, date_str, Loc)
	start_time := date_time.AddDate(0, 0, start)
	end_time := date_time.AddDate(0, 0, end)
	for end_time.After(start_time) {
		date_arr = append(date_arr[0:], start_time.Format(format))
		start_time = start_time.AddDate(0, 0, 1)
	}

	return date_arr
}

/*
* 改变日期
* @time_date 要改变的日期
* @format 转化所需模板 如 20060102
* @years 增加或减少的年
* @months 增加或减少的月
* @days 增加或减少的日
 */
func ChangeDate(time_date, format string, years int, months int, days int) string {
	time_type, _ := time.ParseInLocation(format, time_date, Loc)
	new_time_type := time_type.AddDate(years, months, days)
	new_time_date := new_time_type.Format(format)
	return new_time_date
}

/**
* 判断两个字符串日期相差的天数
 */
func DifferDays(start_date, end_date, time_format string) int {
	day_num := -1
	//将初始日期处理为0时0分0秒
	start_t, err_s := time.ParseInLocation(time_format, start_date, Loc)
	if err_s != nil {
		return day_num
	}
	start_str := start_t.Format(start_day)

	//将日期字符串转成时间格式
	start_str_t, _ := time.ParseInLocation(start_day, start_str, Loc)
	end_t, err_e := time.ParseInLocation(time_format, end_date, Loc)
	if err_e != nil {
		return day_num
	}
	//将时间格式转成时间戳
	start_unix := start_str_t.Unix()
	end_unix := end_t.Unix()
	//通过Unix时间戳获取两个日期相差的天数
	unix_count := int(end_unix - start_unix)

	day_num = unix_count / 86400
	return day_num
}

/**
* 获取两个时间间隔内的所有日期
* 单位是天
 */
func DateDiff(time_date, format string, days int) []string {
	date_arr := []string{}
	days_int := days
	day_i := 0
	date_str := ""
	if days == 0 {
		date_arr = []string{time_date}
		return date_arr
	}
	if days < 0 {
		days_int = days * (-1)
	}
	for i := 0; i <= days_int; i++ {
		day_i = days_int - days + (days_int/days)*i
		date_str = ChangeDate(time_date, format, 0, 0, day_i)
		date_arr = append(date_arr[0:], date_str)
	}
	return date_arr
}

/**
* 判断两个字符串日期相差的天数
 */
func DifferDate(start_date, end_date, time_format, date_diff string) int {
	diff_num := -1
	if date_diff == "day" {
		diff_num = DifferDays(start_date, end_date, time_format)
		return diff_num
	}

	//将日期字符串转成时间格式
	start_t, err_s := time.ParseInLocation(time_format, start_date, Loc)
	if err_s != nil {
		return diff_num
	}
	end_t, err_e := time.ParseInLocation(time_format, end_date, Loc)
	if err_e != nil {
		return diff_num
	}
	//将时间格式转成时间戳
	start_unix := start_t.Unix()
	end_unix := end_t.Unix()
	//通过Unix时间戳获取两个日期相差的天数
	unix_count := int(end_unix - start_unix)
	cout_num := 1
	switch date_diff {
	case "minute":
		cout_num = 60
	case "hour":
		cout_num = 3600
	}

	diff_num = unix_count / cout_num
	if unix_count%cout_num != 0 {
		diff_num = diff_num + 1
	}
	return diff_num
}
