package utils

import (
	"os/exec"
	"os"
	"path/filepath"
	"strings"
	"runtime"
	"io/ioutil"
	"fmt"
)

/**
获取当前路径
 */
func GetCurPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index] + string(os.PathSeparator)
	return ret
}

/**
获取输出路径
 */
func GetOutPath(dir string) string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index] + string(os.PathSeparator) + dir + string(os.PathSeparator)
	return ret
}

/**
获取图片和音频
 */
func GetImageAudio(path string) (images, audios []string) {
	//读取当前目录文件列表
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, nil
	}

	for i := 0; i < len(files); i++ {
		var info = files[i]
		var name = info.Name()

		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "resize-") || strings.HasPrefix(name, "logo-") ||
			strings.HasPrefix(name, "m-") || strings.HasPrefix(name, "n-") || strings.HasPrefix(name, "f-") {
			continue
		}

		//图片处理
		if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".png") {
			images = append(images, name)
		}

		//音频处理
		if strings.HasSuffix(name, ".mp3") {
			audios = append(audios, name)
		}
	}

	return images, audios
}

/**
运行命令
 */
func CineCMD(command string) bool {
	var result bool = true

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			result = false
			fmt.Print(err)
		}
		fmt.Print(string(out))
	} else {
		cmd := exec.Command("bash", "-c", command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			result = false
			fmt.Print(err)
		}
		fmt.Print(string(out))
	}
	return result
}

/**
运行命令
 */
func RunCMD(command string) (string, error) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", command)
		out, err := cmd.CombinedOutput()
		return string(out), err
	} else {
		cmd := exec.Command("bash", "-c", command)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
}
