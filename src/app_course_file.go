package main

import (
	"fmt"
	"strings"
	"os"
	"./utils"
)

func main() {
	var debug = false

	var password string
	var fileType string
	var courseId string
	var outPutPath string

	if debug {
		password = ""
		fileType = "0"
		courseId = "42"

		outPutPath = "/Test/"
	} else {
		for len(password) <= 0 || len(fileType) <= 0 || len(courseId) <= 0 {
			password, fileType, courseId = getArags()
		}

		outPutPath = utils.GetCurPath() + string(os.PathSeparator)
	}

	_, rows := utils.GetFiles(password, "0", courseId)

	var kjOrigFiles, kjFiles []string
	for i := 0; i < len(rows); i++ {
		row := rows[i]

		if !strings.Contains(row.(string),"com") {
			kjFiles = append(kjFiles, regUrl(false, row.(string)))
			kjOrigFiles = append(kjOrigFiles, regUrl(true, row.(string)))
		}
	}
	utils.WriteLines(kjFiles, outPutPath+"kj_files.list")
	utils.WriteLines(kjOrigFiles, outPutPath+"kj_orig_files.list")

	fmt.Println(len(rows))
}

func getArags() (password, fileType, courseId string) {
	fmt.Println("请输入 Password,fileType,CourseId 中间用空格隔开! ")
	fmt.Println("注释：Type:(0:课件资源，1：其他资源) 、 LessonId：(-1:下载所有)")
	fmt.Println("例如：password 0 42")
	fmt.Scanln(&password, &fileType, &courseId)
	return password, fileType, courseId
}

func regUrl(isOrig bool, url string) string {
	mediaUrl := url[strings.Index(url, "/f/")+3:len(url)]

	if isOrig {
		url = "http://www.bstcine.com/ f/" + mediaUrl
	} else {
		urls := strings.Split(url, "/")
		srcType := urls[0]
		courseId := urls[1]

		if srcType == "img" {
			mediaUrl += ".jpg"
		}

		url = "http://gcdn.bstcine.com/" + srcType + "/ " + courseId + "/f/" + mediaUrl
	}

	fmt.Println(url)
	return url
}
