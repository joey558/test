package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

/**
* 定义结构体
 */
type AESECB struct {
	Key_str string //加密的key
	Iv      string //加密的IV
	Ty      string //补位规则
	Mod     string //解码模式
	Cip     int
}

/**
*创建并实例化一个aes结构体
 */
func SetAESECB(key_str, iv, ty, mod string, cip int) *AESECB {
	c := new(AESECB)
	c.Key_str = key_str

	c.Ty = ty

	if len(c.Ty) < 1 {
		c.Ty = "pkcs5"
	}
	c.Cip = cip
	if cip%16 != 0 {
		c.Cip = 16
	}
	if iv == "" {
		iv = key_str
	}
	if len(iv) > 16 {
		iv = Substr(iv, 0, 16)
	}
	c.Mod = mod
	c.Iv = iv
	return c
}

/**
* 数据加密，返回加密字符串
* @str	string	需要加密的字符串
* return	strintg
 */
func (c *AESECB) AesEncryptString(str string) string {
	key := make([]byte, c.Cip)   //设置加密数组
	copy(key, []byte(c.Key_str)) //合并数组补位
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
	if c.Ty == "pkcs5" {
		content = c.PKCS5Padding(content, block.BlockSize())
	} else {
		content = c.ZeroPadding(content, block.BlockSize())
	}
	res := ""
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	if c.Mod == "hex" {
		res = hex.EncodeToString(crypted)
	} else {
		res = base64.StdEncoding.EncodeToString(crypted)
	}

	return res
}

/**
* 数据解密，返回解密后的字符串
* @str	string	需要解密的字符串
* return string
 */
func (c *AESECB) AesDecryptString(str string) string {
	if len(c.Key_str) < 1 {
		fmt.Println("key is null")
		return str
	}

	if str == "" {
		fmt.Println("plain content empty")
		return str
	}
	str = strings.Replace(str, "%2B", "+", -1)
	res := ""
	key := make([]byte, c.Cip)   //设置加密数组
	copy(key, []byte(c.Key_str)) //合并数组补位
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("err1 is->", err)
		return str
	}

	var crypted []byte
	if c.Mod == "hex" {
		crypted, err = hex.DecodeString(str)
	} else {
		crypted, err = base64.StdEncoding.DecodeString(str)
	}

	if err != nil {
		fmt.Println("err is->", err)
		return str
	}

	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	if c.Ty == "pkcs5" {
		origData_pkcs5 := c.PKCS5UnPadding(origData)
		res = string(origData_pkcs5)
	} else {
		origData = c.ZeroUnPadding(origData)
		res = string(origData)
	}

	return res
}

/**
* Zero补位算法
 */
func (c *AESECB) ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func (c *AESECB) ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

/**
* PKCS5补位算法
 */
func (c *AESECB) PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (c *AESECB) PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	data_index := length - unpadding
	res := origData
	if data_index < 0 || data_index > length {
		data_index = length - 1
	} else {
		res = origData[:(data_index)]
	}
	return res
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
