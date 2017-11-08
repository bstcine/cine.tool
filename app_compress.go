package main

import (
	"os"
	"bufio"
	"fmt"
)

type InputArgs struct {
	OutputPath string /** 输出目录 */
	LocalPath  string /** 输入的目录或文件路径 */
	LogoPath   string /** 水印图片名称*/
}

var inputArgs InputArgs

func main() {
	debug := false
	if debug {
		inputArgs.LocalPath = "/Volumes/Go/test/"
		inputArgs.OutputPath = "/Volumes/Go/test/"
		inputArgs.LogoPath = inputArgs.OutputPath + ".logo.png"
	} else {
		inputArgs.LocalPath = GetCurPath()
		inputArgs.OutputPath = GetOutPath("doing")
		inputArgs.LogoPath = inputArgs.OutputPath + ".logo.png"
	}

	//创建输出目录
	os.MkdirAll(inputArgs.OutputPath, 0777)
	//检查创建水印图片
	checkLogoFile(inputArgs.LogoPath)

	//获取图片和音频
	images,audios := GetImageAudio(inputArgs.LocalPath)
	hasMagick := CheckHasMagick()

	fmt.Println(">>>>>>>>>>   开始处理图片（压缩加水印）...")
	for i := 0; i < len(images); i++ {
		var name = images[i]
		fmt.Print(name + " 处理中...")
		if debug {
			ResizeImgByMagick(inputArgs.LocalPath, inputArgs.OutputPath, name, "m-"+name)
			ResizeImg(inputArgs.LocalPath, inputArgs.OutputPath, name, "n-"+name)
		} else {
			ResizeImgByMagick(inputArgs.LocalPath, inputArgs.OutputPath, name, name)
			LogoImgByMagick(inputArgs.LogoPath, inputArgs.OutputPath, inputArgs.OutputPath, name+".jpg", name+".jpg")

			if !hasMagick {
				ResizeImg(inputArgs.LocalPath, inputArgs.OutputPath, name, name)
				LogoImg(inputArgs.LogoPath, inputArgs.OutputPath, inputArgs.OutputPath, name, name)
			}
		}
		fmt.Println(" 完成.")
	}
	fmt.Println("<<<<<<<<<<   图片处理成功...")

	hasffmpeg := CheckHasFFMPEG()
	if hasffmpeg {
		fmt.Println(">>>>>>>>>>   开始处理音频（压缩）...")
		for i := 0; i < len(audios); i++ {
			var name = audios[i]
			fmt.Print(name + " 处理中...")
			if debug {
				ResizeAudio(inputArgs.LocalPath, inputArgs.OutputPath, name, "f-"+name)
			}else {
				ResizeAudio(inputArgs.LocalPath, inputArgs.OutputPath, name, name)
			}

			fmt.Println(" 完成.")
		}
		fmt.Println("<<<<<<<<<<   音频处理成功...")
	}

	fmt.Println("请输入 end ,结束本程序...")
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		if input.Text() == "end" {
			break
		}
	}
}
