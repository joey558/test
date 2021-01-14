package common

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
)

/**
* 数据加密，返回加密字符串
* @str	string	需要加密的字符串
* return	strintg
 */
func HexAesEncrypt(key_str, str string, cip int) string {
	key := make([]byte, 32)    //设置加密数组
	copy(key, []byte(key_str)) //合并数组补位
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
		return str
	}
	if str == "" {
		fmt.Println("plain content empty")
		return str
	}
	ecb := NewECBEncrypter(block)
	content := []byte(str)
	content = PKCS5Padding(content, block.BlockSize())

	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	res := hex.EncodeToString(crypted)
	return res
}

/**
* 数据解密，返回解密后的字符串
* @str	string	需要解密的字符串
* return string
 */
func HexAesDecrypt(key_str, str string) string {
	if len(key_str) < 1 {
		fmt.Println("key is null")
		return str
	}

	if str == "" {
		fmt.Println("plain content empty")
		return str
	}

	res := ""
	key := make([]byte, 32)    //设置加密数组
	copy(key, []byte(key_str)) //合并数组补位
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("err1 is->", err)
		return str
	}

	crypted, err := hex.DecodeString(str)
	if err != nil {
		fmt.Println("err is->", err)
		return str
	}
	//crypted, _ := base64.StdEncoding.DecodeString(str)
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	origData_pkcs5 := PKCS5UnPadding(origData)
	res = string(origData_pkcs5)

	return res
}
