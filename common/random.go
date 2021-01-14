package common

import (
	"math/rand"
	"time"
)

func Random(param string, length int) string {
	str := ""
	if length < 1 {
		return str
	}
	tmp := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	switch param {
	case "number":
		tmp = "1234567890"
	case "small":
		tmp = "abcdefghijklmnopqrstuvwxyz"
	case "big":
		tmp = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "smallnumber":
		tmp = "1234567890abcdefghijklmnopqrstuvwxyz"
	case "bignumber":
		tmp = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "bigsmall":
		tmp = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	leng := len(tmp)
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		s_ind := ran.Intn(leng)
		str = str + Substr(tmp, s_ind, 1)
	}

	return str
}

/**
* 生成两个数之间的随机数
* @min   int64  最小值
* @max   int64  最大值
* return  int64	 返回一个int64整形数字
 */
func RandomMaxAndMin(min, max int64) int64 {
	if min >= max {
		return max
	}
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := ran.Int63n(max-min) + min
	return res
}

/**
* 生成两个数之间的随机数
* @min   int  最小值
* @max   int  最大值
* return  int	 返回一个int整形数字
 */
func RandomMaxAndMinInt(min, max int) int {
	if min >= max {
		return max
	}
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := ran.Intn(max-min) + min
	return res
}