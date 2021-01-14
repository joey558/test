package common

/**
* 谷歌动态验证
 */

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

/**
* 定义结构体
 */
type GoogleAuth struct {
	IntervalLength int64 //验证码的有效时间 30
	PinLength      int   //验证码长度 6
}

/**
*创建并实例化一个GoogleAuth结构体
 */
func SetGoogleAuth(intervalLength int64, pinLength int) *GoogleAuth {
	c := new(GoogleAuth)
	c.IntervalLength = intervalLength
	c.PinLength = pinLength
	return c
}

/**
* 比较用户输入的谷歌验证码是否正确
* @Secret_char   int64  数据库保存的Secret_char
* @code_value    int64  用户输入的谷歌验证码
* @return int 状态值 200表示正确
   string   结果说明
*/
func (g *GoogleAuth) CheckGoogleCode(Secret_char string, code_value int64) (int, string) {
	c_status := 100
	c_msg := "验证码错误"

	secretUnix := time.Now().Unix()
	c_value := IntToByte(secretUnix / g.IntervalLength)

	c_key, err := base32.StdEncoding.DecodeString(Secret_char)

	if err != nil {
		return c_status, c_msg
	}
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, c_key)
	hmacSha1.Write(c_value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := ByteToUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits

	c_pwd := number % 1000000
	if int64(c_pwd) != code_value {
		return c_status, c_msg
	}
	c_status = 200
	c_msg = "success"

	return c_status, c_msg
}

/**
* 当前时间戳
* return int64
 */
func GetSecret() string {
	unix_str := strconv.FormatInt(time.Now().Unix(), 10)
	src := []byte(unix_str)
	// 解码后数据的最长长度
	maxLen := base32.StdEncoding.EncodedLen(len(src))
	// 解码后的缓存区
	dst := make([]byte, maxLen)
	// base32 解码
	base32.StdEncoding.Encode(dst, src)
	// 打印解码的数据

	secret := string(dst[:maxLen])
	return secret
}

/**
* 获取谷歌验证二维码的url
* @identifier string 谷歌的账号(test)
* @key 当前时间戳的[]byte
* @width int 图片的宽
* @height int 图片的高
 */
func GetImageUrl(identifier, secret string, width, height int) string {
	provision_url := fmt.Sprintf("otpauth://totp/%s?secret=%s", identifier, secret)
	provision_url = url.QueryEscape(provision_url)
	qr_url := fmt.Sprintf("http://chart.apis.google.com/chart?cht=qr&chs=%dx%d&chld=L|1&chl=%s", width, height, provision_url)
	return qr_url
}
