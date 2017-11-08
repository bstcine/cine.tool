package main

import (
	"fmt"
	"runtime"
)

func CheckHasFFMPEG() bool {
	command := "ffmpeg -version"
	status := CineCMD(command)

	if !status {
		if runtime.GOOS == "windows" {
			fmt.Println("请安装音频处理工具 https://www.ffmpeg.org/download.html#build-windows")
		} else {
			fmt.Println("请安装音频处理工具 https://www.ffmpeg.org/download.html")
		}
	}

	return status
}

func ResizeAudio(oidPath, newPath, oidName, newName string)  {
	// -ar 设置音频采样频率
	// -ac 设置音频通道的数量
	// -aq 设置音频质量
	// -b:a 设置音频比特率
	command := "ffmpeg -y -i "+oidPath + oidName +" -ar 32k -ac 1 -b:a 128k "+newPath + newName
	status := CineCMD(command)
	if status {
		fmt.Print(" 压缩成功 ")
	} else {
		fmt.Print(" 压缩失败 ")
	}
}
