package utils

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
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
获取JSON格式的目录下所有文件信息
 */
func GetJsonFileList(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var jsonFileList []string

	for _, f := range files {
		if f.IsDir() {
			jsonFileList = append(jsonFileList, GetJsonFileList(path+string(os.PathSeparator)+f.Name())...)
		} else if !strings.HasPrefix(f.Name(), ".") {
			jsonFileInfo := GetJsonFileInfo(path + string(os.PathSeparator) + f.Name())
			jsonFileList = append(jsonFileList, jsonFileInfo)
		}
	}

	return jsonFileList
}

/**
获取JSON格式的文件信息
 */
func GetJsonFileInfo(url string) string {
	var cmd = "ffprobe -v quiet -print_format json -show_format " + url
	result,_ := RunCMD(cmd)
	return result
}
