package common

import (
	"strings"
)

/**
*  生成插入的sql语句
 */
func InsertSql(table_name string, data map[string]string) string {
	sql := ""
	if table_name == "" || len(data) < 1 {
		return sql
	}

	key_str := ","
	val_str := ","
	for d_k, d_v := range data {
		key_str = key_str + ",`" + d_k + "`"
		val_str = val_str + `,'` + d_v + `'`
	}

	key_str = strings.Replace(key_str, ",,", "", 1)
	val_str = strings.Replace(val_str, ",,", "", 1)
	sql = "insert into `" + table_name + "` (" + key_str + ") VALUES (" + val_str + ");"
	return sql
}

/**
*  生成批量插入的sql语句
 */
func BatchInsertSql(table_name string, data []map[string]string) string {
	sql := ""
	if table_name == "" || len(data) < 1 {
		return sql
	}

	if len(data[0]) < 1 {
		return sql
	}

	key_str := ","
	key_arr := []string{}
	for d0_k, _ := range data[0] {
		key_arr = append(key_arr[0:], d0_k)
		key_str = key_str + ",`" + d0_k + "`"
	}
	val_str := ","
	for _, d_val := range data {
		v_str := ","
		for i := 0; i < len(key_arr); i++ {
			v_str = v_str + `,"` + d_val[key_arr[i]] + `"`
		}
		v_str = strings.Replace(v_str, ",,", "", 1)
		val_str = val_str + ",(" + v_str + ")"
	}

	key_str = strings.Replace(key_str, ",,", "", 1)
	val_str = strings.Replace(val_str, ",,", "", 1)
	sql = "insert ignore into `" + table_name + "` (" + key_str + ") VALUES " + val_str + ";"
	return sql
}
