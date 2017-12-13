package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os"
	"./utils"
)

type Result struct {
	Status bool
	Data   []CineFile
}

type CineFile struct {
	Path string
	Course string
	Name string
}

func main() {
	fmt.Println("欢迎使用课程习题与单词下载工具....")

	isdebug := true

	var outPutPath string
	var getCourseFileApi string
	var courseId string

	if isdebug {
		outPutPath = "/Volumes/Go/test/习题单词库/"
		getCourseFileApi = "http://local.bstcine.com:9000/api/tool/content/course/exerciseWord"
		os.MkdirAll(outPutPath, 0777)
	} else {
		outPutPath = utils.GetCurPath() + string(os.PathSeparator) + "习题单词库" + string(os.PathSeparator)
		getCourseFileApi = "http://www.bstcine.com/api/tool/content/course/exerciseWord"
	}

	fmt.Println("请输入课程ID：(注：不填为所有课程)，下载文件到 '习题单词库' 文件夹下")
	fmt.Scanln(&courseId)

	if courseId != "" {
		getCourseFileApi += "?lesson_id=" + courseId
	}

	resp, err := http.Get(getCourseFileApi)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	var res Result
	json.Unmarshal(body, &res)

	if res.Status {
		files := res.Data

		fmt.Printf("一共%d个文件\n", len(files))

		for i := 0; i < len(files); i++ {
			file := files[i]
			fmt.Println(file)

			downPath := outPutPath + file.Course + string(os.PathSeparator)
			if _, err := os.Stat(downPath); err != nil {
				os.MkdirAll(downPath, 0777)
			}

			utils.DownloadFile(downPath+ file.Name,file.Path)
		}
	}
}
