package main

import (
	"./utils"
	"fmt"
	"os"
)

func main() {
	///mnt/web/kj.bstcine.com/wwwroot/kj/d011502846526016MpDBrnsR8p/f/2017/11/03/135224745StPZy4Y.png.jpg
	//http://www.bstcine.com/img/
	//baseUrl = "http://www.bstcine.com/f/"
	isdebug := false
	var baseUrl string

	if isdebug {
		baseUrl = "/Volumes/Go/Test"
	} else {
		args := os.Args
		if len(args) < 2 {
			fmt.Println("Please input need check path...")
			return
		}

		baseUrl = args[1]
		if !utils.Exists(baseUrl) {
			fmt.Println("Input check path is not exists...")
			return
		}
	}

	files := utils.GetJsonFileList(baseUrl)
	fmt.Println(files)
}
