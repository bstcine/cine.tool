package utils

import (
	"os"
)

/**
文件是否存在
 */
func Exists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

/**
获取JSON格式的文件信息
 */
func GetJsonFileInfo(url string) (string,error) {
	var cmd = "ffprobe -v quiet -print_format json -show_format " + url
	result,err := RunCMD(cmd)
	return result,err
}
