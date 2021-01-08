package common

import (
	"fmt"
	"os"
	"strings"
	"time"
)

/**
* 指定文件名写入文件
 */
func LogsWithFileName(log_path,FiLeName string, msg string) {
	if IsDirExists(log_path) != true {
		os.MkdirAll(log_path, 0777)
	}
	isEnd := strings.HasSuffix(log_path, "/")
	if isEnd != true {
		log_path = log_path + "/"
	}
	timeStr := time.Now().Format("2006-01-02")
	logFile := log_path + FiLeName + timeStr + ".log"

	fout, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		fmt.Println(log_path, err)
		return
	}

	fout.WriteString(time.Now().Format("2006-01-02 15:04:05") + "\r\n" + msg + "\r\n=====================\r\n")
	defer fout.Close()
}

/**
* 判断目录是否存在
 */
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}
