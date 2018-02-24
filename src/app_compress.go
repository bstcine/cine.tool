package main

import (
	"./utils"
	"./model"
	"bufio"
	"fmt"
	"os"
)

func main() {
	var debug = true
	var inputArgs model.InputArgs

	if debug {
		inputArgs.LocalPath = "/Volumes/Go/test/"
		inputArgs.OutputPath = "/Volumes/Go/test/compress/"
		inputArgs.LogoPath = inputArgs.OutputPath + ".logo.png"
	} else {
		inputArgs.LocalPath = utils.GetCurPath()
		inputArgs.OutputPath = utils.GetOutPath("compress")
		inputArgs.LogoPath = inputArgs.OutputPath + ".logo.png"
	}

	//创建输出目录
	os.MkdirAll(inputArgs.OutputPath, 0777)
	//检查创建水印图片
	utils.CheckLogoFile(inputArgs.LogoPath)

	//获取图片和音频
	images, audios := utils.GetImageAudio(inputArgs.LocalPath)

	if utils.CheckHasMagick() {
		fmt.Println(">>>>>>>>>>   开始处理图片（压缩加水印）...")
		for i := 0; i < len(images); i++ {
			var name = images[i]
			fmt.Print(name + " 处理中...")
			utils.ResizeImgByMagick(inputArgs.LocalPath, inputArgs.OutputPath, name, name)
			utils.LogoImgByMagick(inputArgs.LogoPath, inputArgs.OutputPath, inputArgs.OutputPath, name+".jpg", name+".jpg")
			fmt.Println(" 完成.")
		}
		fmt.Println("<<<<<<<<<<   图片处理成功...")
	}

	if utils.CheckHasFFMPEG() {
		fmt.Println(">>>>>>>>>>   开始处理音频（压缩）...")
		for i := 0; i < len(audios); i++ {
			var name = audios[i]
			fmt.Print(name + " 处理中...")
			if debug {
				utils.ResizeAudio(inputArgs.LocalPath, inputArgs.OutputPath, name, "f-"+name)
			} else {
				utils.ResizeAudio(inputArgs.LocalPath, inputArgs.OutputPath, name, name)
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
