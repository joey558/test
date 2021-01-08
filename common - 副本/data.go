package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"reflect"
)

/**
*  结构体转成json字符串
 */
func StructToJson(obj interface{}) string {
	json_str := ""
	json_byte, err := json.Marshal(obj)
	if err != nil {
		return json_str
	}
	json_str = string(json_byte)
	return json_str
}

func StructToMap(obj interface{}) map[string]string {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = fmt.Sprintf("%v", obj2.Field(i).Interface())
	}
	return data
}

func StructToMapSlow(obj interface{}) map[string]string {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < obj1.NumField(); i++ {
		data[strings.ToLower(obj1.Field(i).Name)] = fmt.Sprintf("%v", obj2.Field(i).Interface())
	}
	return data
}