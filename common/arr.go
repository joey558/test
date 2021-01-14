package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

/*
创建数组
min 最小值
max 最大值
num 每个值的数量
*/
func LaborArr(min, max, num int) []int {
	a := []int{}

	for n := min; n <= max; n++ {

		for i := 0; i < num; i++ {
			a = append(a, n)
		}
	}
	return a
}

/**
*  合并数组
 */
func Merge(x []int, y []int) []int {
	z := x

	for _, y_v := range y {
		z = append(z, y_v)
	}

	return z
}

/**
* 将数组或者map组合成字符串
 */
func Implode(arr []string, sep string) string {
	res := ""
	if len(arr) > 0 {
		for _, v := range arr {
			//拼接sql语句
			if res == "" {
				res = v
			} else {
				res = res + sep + v
			}
		}
	}
	return res
}

/**
* 是否包含在数组中
 */
func Arr_In(arr []string, sep_str string) bool {
	res := false
	if len(arr) < 1 {
		return res
	}
	for _, arr_str := range arr {
		if arr_str == sep_str {
			res = true
			break
		}
	}
	return res
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Struct2MapStr(obj interface{}) map[string]string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	data := map[string]string{}
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = fmt.Sprintf("%v", v.Field(i).Interface())
	}
	return data
}

func Interface2Json(obj interface{}) string {
	res := ""
	b, err := json.Marshal(obj)
	if err != nil {
		return res
	}
	res = string(b)
	return res
}

func IntArr_Index(arr []int, sep_str int) int {
	res := -1
	if len(arr) > 0 {
		for arr_ind, arr_str := range arr {
			if arr_str == sep_str {
				res = arr_ind
				break
			}
		}
	}
	return res
}

func StrArr_Index(arr []string, sep_str string) int {
	res := -1
	if len(arr) > 0 {
		for arr_ind, arr_str := range arr {
			if arr_str == sep_str {
				res = arr_ind
				break
			}
		}
	}
	return res
}

/**
*  给Map数组排序
*  需要比较的值是字符串
 */
func SortMap(map_arr []map[string]string, map_filed string) {
	map_len := len(map_arr)
	if map_len < 1 {
		return
	}
	cop_arr := make([]map[string]string, map_len)
	for in_k, in_val := range map_arr {
		a := map[string]string{}
		for i_k, i_v := range in_val {
			a[i_k] = i_v
		}
		cop_arr[in_k] = a
	}

	ind_arr := make([]string, map_len)
	for m_k, m_val := range map_arr {
		r, ok := m_val[map_filed]
		if !ok {
			break
		}
		ind_arr[m_k] = r
	}

	if len(ind_arr) < 1 {
		return
	}
	cop_ind := make([]string, map_len)
	copy(cop_ind, ind_arr)
	sort.Strings(ind_arr)
	for in_k, in_val := range ind_arr {
		cop_in := StrArr_Index(cop_ind, in_val)
		cop_ind[cop_in] = "null"
		map_arr[in_k] = cop_arr[cop_in]
	}
}

/**
*  给Map数组排序
*  需要比较的值是int
*  只能是正数排序
 */
func SortMapInt(map_arr []map[string]string, map_filed string) {
	map_len := len(map_arr)
	if map_len < 1 {
		return
	}
	cop_arr := make([]map[string]string, map_len)
	for in_k, in_val := range map_arr {
		a := map[string]string{}
		for i_k, i_v := range in_val {
			a[i_k] = i_v
		}
		cop_arr[in_k] = a
	}

	ind_arr := make([]int, map_len)
	for m_k, m_val := range map_arr {
		r, ok := m_val[map_filed]
		if !ok {
			break
		}
		ind_arr[m_k], _ = strconv.Atoi(r)
	}

	if len(ind_arr) < 1 {
		return
	}
	cop_ind := make([]int, map_len)
	copy(cop_ind, ind_arr)
	sort.Ints(ind_arr)
	for in_k, in_val := range ind_arr {
		cop_in := IntArr_Index(cop_ind, in_val)
		cop_ind[cop_in] = -1
		map_arr[in_k] = cop_arr[cop_in]
	}
}

/**
*  给Map数组排序
 */
func SortMapInterface(map_arr []map[string]interface{}, map_filed string) {
	map_len := len(map_arr)
	if map_len < 1 {
		return
	}
	cop_arr := make([]map[string]interface{}, map_len)
	for in_k, in_val := range map_arr {
		a := map[string]interface{}{}
		for i_k, i_v := range in_val {
			a[i_k] = i_v
		}
		cop_arr[in_k] = a
	}

	ind_arr := make([]string, map_len)
	for m_k, m_val := range map_arr {
		r, ok := m_val[map_filed]
		if !ok {
			break
		}
		r_val := fmt.Sprintf("%v", r)
		if r_val == "" {
			break
		}
		ind_arr[m_k] = r_val
	}

	if len(ind_arr) < 1 {
		return
	}
	cop_ind := make([]string, map_len)
	copy(cop_ind, ind_arr)
	sort.Strings(ind_arr)
	for in_k, in_val := range ind_arr {
		cop_in := StrArr_Index(cop_ind, in_val)
		cop_ind[cop_in] = "null"
		map_arr[in_k] = cop_arr[cop_in]
	}
}

func SortMapInterfaceInt(map_arr []map[string]interface{}, map_filed string) {
	map_len := len(map_arr)
	if map_len < 1 {
		return
	}
	cop_arr := make([]map[string]interface{}, map_len)
	for in_k, in_val := range map_arr {
		a := map[string]interface{}{}
		for i_k, i_v := range in_val {
			a[i_k] = i_v
		}
		cop_arr[in_k] = a
	}

	ind_arr := make([]int, map_len)
	for m_k, m_val := range map_arr {
		r, ok := m_val[map_filed]
		if !ok {
			break
		}
		r_val, err := strconv.Atoi(fmt.Sprintf("%v", r))
		if err != nil {
			break
		}
		ind_arr[m_k] = r_val
	}

	if len(ind_arr) < 1 {
		return
	}
	cop_ind := make([]int, map_len)
	copy(cop_ind, ind_arr)
	sort.Ints(ind_arr)
	for in_k, in_val := range ind_arr {
		cop_in := IntArr_Index(cop_ind, in_val)
		cop_ind[cop_in] = -1
		map_arr[in_k] = cop_arr[cop_in]
	}
}

/**
*  从N个数组中,每个数组取出一个元素组成新的数组
*  is_mul  是否允许重复元素
 */
func ArrToNewArr(arr [][]string, len_arr, is_mul int) (int, [][]string) {
	num := 0
	res := [][]string{}
	if len_arr == 0 {
		return num, res
	}
	if len_arr != len(arr) {
		return num, res
	}

	if len_arr == 10 || len_arr == 1 {
		for i := 0; i < len_arr; i++ {
			if len(arr[i]) < 1 {
				num = 0
				return num, res
			}
			num = num + len(arr[i])
		}
		return num, res
	}
	a_arr := map[int][][]string{}
	for a := 0; a < len_arr; a++ {
		a_arr[a+1] = NewArr(a_arr[a], arr[a], is_mul)
	}
	num = len(a_arr[len_arr])

	return num, a_arr[len_arr]
}

func NewArr(arr [][]string, sel_arr []string, is_mul int) [][]string {
	arr_arr := [][]string{}
	if len(sel_arr) < 1 {
		return arr
	}
	arr_len := len(arr)
	if arr_len == 0 {
		for b := 0; b < len(sel_arr); b++ {
			for_arr := []string{}
			for_arr = append(for_arr[0:], sel_arr[b])
			arr_arr = append(arr_arr[0:], for_arr)
		}
	} else {
		for i := 0; i < arr_len; i++ {
			for b := 0; b < len(sel_arr); b++ {
				if Arr_In(arr[i], sel_arr[b]) && is_mul == 0 {
					continue
				}
				for_arr := make([]string, len(arr[i])+1)
				for t_k, t_val := range arr[i] {
					for_arr[t_k] = t_val
				}
				for_arr[len(arr[i])] = sel_arr[b]

				arr_arr = append(arr_arr[0:], for_arr)
			}
		}
	}

	return arr_arr
}

func UrlToMap(urlStr string) map[string]string {
	paramMap := map[string]string{}
	urlArr := strings.Split(urlStr, "&")
	if len(urlArr) < 1 {
		return paramMap
	}
	for _, val := range urlArr {
		if val == "" {
			continue
		}
		valArr := strings.Split(val, "=")
		if len(valArr) < 2 {
			continue
		}
		paramMap[valArr[0]] = valArr[1]
	}
	return paramMap
}
