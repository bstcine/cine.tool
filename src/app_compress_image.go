package main

import (
	"./utils"
	"os"
	"fmt"
)

func main() {
	debug := false
	fmt.Println("欢迎使用图片压缩工具（800*450）....")

	localPath := "/Go/test/images/"
	outPutPath := "/Go/test/images/compress/"
	var imgWidth uint = 800
	var imgHeight uint = 450

	if !debug {
		localPath = utils.GetCurPath()
		outPutPath = utils.GetOutPath("compress")

		fmt.Println("请输入 width,height 中间用空格隔开! 例如：")
		fmt.Println("800 450")
		fmt.Scanln(&imgWidth, &imgHeight)
	}

	//创建输出目录
	os.MkdirAll(outPutPath, 0777)

	//获取图片和音频
	images, _ := utils.GetImageAudio(localPath)

	fmt.Println(">>>>>>>>>>   开始处理图片 ...")
	for i := 0; i < len(images); i++ {
		var name = images[i]
		fmt.Print(name + " 处理中...")
		utils.ResizeImgWH(imgWidth, imgHeight, localPath, outPutPath, name, name)
		fmt.Println(" 完成.")
	}
	fmt.Println("<<<<<<<<<<   图片处理成功...")

	var code string
	fmt.Printf("请输入任意键退出")
	fmt.Scanln(&code)
}

func GetImageArags() (width, height uint) {
	return width, height
}