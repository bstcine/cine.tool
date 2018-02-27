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
		courseId = "d011502846526016MpDBrnsR8p"

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

		kjFiles = append(kjFiles, regUrl(false, row.(string)))
		kjOrigFiles = append(kjOrigFiles, regUrl(true, row.(string)))
	}
	utils.WriteLines(kjFiles, outPutPath+"kj.list")
	utils.WriteLines(kjOrigFiles, outPutPath+"kj_orig.list")

	fmt.Println(len(rows))
}

func getArags() (password, fileType, courseId string) {
	fmt.Println("请输入 Password,fileType,CourseId 中间用空格隔开! ")
	fmt.Println("注释：Type:(0:课件资源，1：其他资源) 、 LessonId：(-1:下载所有)")
	fmt.Println("例如：password 0 42")
	fmt.Scanln(&password, &fileType, &courseId)
	return password, fileType, courseId
}

func regUrl(isOrig bool, param string) (url string) {
	urls := strings.Split(param,";")

	mediaUrl := urls[0]
	urlPrefix := urls[1]
	urlSuffix := urls[2]

	if isOrig {
		url = "http://www.bstcine.com/ f/" + mediaUrl
	} else if strings.Contains(urlPrefix,"http://gcdn.bstcine.com") {
		urlPrefix = strings.Replace(urlPrefix,"/img/","/ img/",-1)
		urlPrefix = strings.Replace(urlPrefix,"/mp3/","/ mp3/",-1)

		url = urlPrefix + mediaUrl + urlSuffix
	}

	fmt.Println(url)
	return url
}
