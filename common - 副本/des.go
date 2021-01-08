package common

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

/**
* des加密
 */
func DesEncryptString(str, keyStr string) string {
	//key := make([]byte, 8)    //设置加密数组
	//copy(key, []byte(keyStr)) //合并数组补位
	result, err := DesDecrypt([]byte(str), []byte(keyStr))
	res := ""
	if err == nil {
		res = base64.StdEncoding.EncodeToString(result)
	}
	return res
}

/**
* des加密
 */
func HexDesEncryptString(str, key_str, iv_str string) string {
	if iv_str == "" {
		iv_str = key_str
	}
	result, err := HexDesEncrypt([]byte(str), []byte(key_str), []byte(iv_str))
	res := ""
	if err == nil {
		res = hex.EncodeToString(result)
	}
	return res
}

func HexDesEncrypt(orig_data, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	//src = DZeroPadding(src, bs)
	orig_data = PKCS5Padding(orig_data, bs)
	if len(orig_data)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	crypted := make([]byte, len(orig_data))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(crypted, orig_data)

	return crypted, nil
}

func DesEncrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	//src = DZeroPadding(src, bs)
	src = PKCS5Padding(src, bs)
	if len(src)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

/**
* des解密
 */
func DesDecryptString(str, key string) string {
	result, _ := base64.StdEncoding.DecodeString(str)
	res := ""
	if len(result) > 0 {
		origData, err := DesDecrypt(result, []byte(key))
		if err == nil {
			res = string(origData)
		}
	}
	return res
}

func DesDecrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	//out = ZeroUnPadding(out)
	out = PKCS5UnPadding(out)
	return out, nil
}

/**
* Zero补位算法
 */
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

/**
* PKCS5补位算法
 */
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
