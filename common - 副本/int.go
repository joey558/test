package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

/**
* 整形转byte数组
 */
func IntToByte(data int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((data>>shift)&mask))
	}
	return result
}

/**
* byte数组转int64
 */
func ByteToInt(data []byte) int64 {
	b_buf := bytes.NewBuffer(data)
	var res int32
	binary.Read(b_buf, binary.BigEndian, &res)
	return int64(res)
}

func ByteToUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

//字符串转int64
func Str2Int64(str string) (int64, error) {
	i64, err := strconv.ParseInt(str, 10, 64)
	return i64, err
}

//字符串转int
func Str2Int(str string) (int, error) {
	// i64, err := Str2Int64(str)
	i, err := strconv.Atoi(str)
	return i, err
}

func Interface2Int(inter interface{}) (int, error) {
	str := fmt.Sprintf("%v", inter)
	return Str2Int(str)
}

/*
* 排序拼接
* mid 拼接符号
* url 是否是url模式
* conectnil 空值是否参与拼接 为true 参与拼接
 */
func MapCreatLinkSort(m map[string]string, mid string, url bool, conectnil bool) string {
	var result = ""

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, k := range keys {
		if conectnil {
			if url {
				result += k + "=" + fmt.Sprintf("%s", m[k]) + mid
			} else {
				result += k + fmt.Sprintf("%s", m[k]) + mid
			}
		} else if m[k] != "" {
			if url {
				result += k + "=" + fmt.Sprintf("%s", m[k]) + mid
			} else {
				result += k + fmt.Sprintf("%s", m[k]) + mid
			}
		}

	}

	if mid != "" {
		result = strings.TrimRight(result, mid)
	}

	return result
}

//根据想要的进行拼接
// mid 拼接符号
// dm 为要拼接的内容用，隔开
// url 0为url模式 ; 1 拼接：参数名值mid参数名值; 2 值拼接
func MapCreatLink(m map[string]string, dm string, mid string, url int) string {
	var result = ""
	Extract := strings.Split(dm, ",")
	for _, value := range Extract {
		if value != "" && m[value] != "" {
			switch url {
			case 0:
				result += value + "=" + fmt.Sprintf("%s", m[value]) + mid
			case 1:
				result += value + fmt.Sprintf("%s", m[value]) + mid
			case 2:
				result += fmt.Sprintf("%s", m[value]) + mid
			}
		}
	}

	if mid != "" {
		result = strings.TrimRight(result, mid)
	}

	return result
}
