package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

/** generate md5 checksum of URL in hex format **/
func HexMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	c := m.Sum(nil)
	return hex.EncodeToString(c)
}

func BaseHmacSha(source string, key_str string) string {
	key := []byte(key_str)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(source))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
