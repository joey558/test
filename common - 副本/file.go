package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

/**
*  文件内容读取,返回字符串
 */
func ReadFileString(file_pth string) (string, error) {
	str := ""
	f, err := os.Open(file_pth)
	if err != nil {
		return str, err
	}

	str_byte, err := ioutil.ReadAll(f)
	str = string(str_byte)
	return str, err
}

/**
*  文件内容读取,返回二进制数组
 */
func ReadFile(file_pth string) ([]byte, error) {
	f, err := os.Open(file_pth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

/**
* 新建文件并写入内容
* 如果文件已存在,则覆盖以前内容
 */
func WriteFile(file_pth, file_name, content string) (int, error) {
	_, err := CreateFile(file_pth)
	if err != nil {
		return 0, err
	}
	src := file_pth + "/" + file_name
	fs, e := os.Create(src)
	if e != nil {
		return 0, e
	}
	defer fs.Close()
	return fs.WriteString(content)
}

/**
* 创建文件
 */
func CreateFile(src string) (string, error) {
	if IsExist(src) {
		return src, nil
	}

	if err := os.MkdirAll(src, 0777); err != nil {
		if os.IsPermission(err) {
			fmt.Println("你不够权限创建文件")
		}
		return "", err
	}

	return src, nil
}

/**
* 判断文件或者目录是否存在
 */
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

/**
*  下载文件
 */
func DownloadFile(file_name string, con_read io.Reader) error {
	file, err := os.Create(file_name)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, con_read)
	return err
}

/**
* 获取文件的大小
 */
func FileSize(file_path string) int64 {
	var result int64
	is_exist := IsExist(file_path)
	if !is_exist {
		return result
	}
	filepath.Walk(file_path, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}
