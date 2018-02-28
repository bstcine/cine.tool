package utils

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
	"os"
)

/**
获取配置文件参数
 */
func GetConfArgs(path string) (argsMap map[string]string) {
	if _, err := os.Stat(path); err != nil {
		argsMap = nil
	} else {
		args, _ := ReadLines(path)
		argsMap = make(map[string]string)

		for i := 0; i < len(args); i++ {
			arg := args[i]
			if !strings.Contains(arg, "#") {
				argSplit := strings.Split(arg, "=")
				if len(argSplit) > 1 {
					argsMap[argSplit[0]] = argSplit[1]
				} else {
					argsMap[argSplit[0]] = ""
				}
			}
		}
	}
	return argsMap
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
