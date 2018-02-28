package main

import (
	"fmt"
	"strings"
	"os"
	"./utils"
)

var cfgName = "app_oss.cfg"

func main() {
	var debug = false

	var curPath string
	var outPutPath string

	if debug {
		curPath = "/Go/Cine/cine.tool/assets/"
		outPutPath = "/Test/"
	} else {
		curPath = utils.GetCurPath() + string(os.PathSeparator)
		outPutPath = utils.GetCurPath() + string(os.PathSeparator)
	}

	argsMap := utils.GetConfArgs(curPath+cfgName)
	if argsMap == nil || len(argsMap) <=0 {
		fmt.Println("配置文件不存在")
		return
	}

	if argsMap["srcType"] == "move" {
		_, rows := utils.GetFiles(argsMap["srcPassword"], "0", argsMap["moveCourse"])

		var kjFiles, kjCdnFiles []string
		for i := 0; i < len(rows); i++ {
			row := rows[i]

			kjCdnFiles = append(kjCdnFiles, regUrl(false, row.(string)))
			kjFiles = append(kjFiles, regUrl(true, row.(string)))
		}
		utils.WriteLines(kjFiles, outPutPath+argsMap["moveOutFileName"])
		utils.WriteLines(kjCdnFiles, outPutPath+argsMap["moveOutCdnFileName"])

		fmt.Println(len(rows))
	}

}

func regUrl(isOrig bool, param string) (url string) {
	urls := strings.Split(param, ";")

	mediaUrl := urls[0]
	urlPrefix := urls[1]
	urlSuffix := urls[2]

	if isOrig {
		url = "http://www.bstcine.com/ f/" + mediaUrl
	} else if strings.Contains(urlPrefix, "http://gcdn.bstcine.com") {
		urlPrefix = strings.Replace(urlPrefix, "/img/", "/ img/", -1)
		urlPrefix = strings.Replace(urlPrefix, "/mp3/", "/ mp3/", -1)

		url = urlPrefix + mediaUrl + urlSuffix
	}

	fmt.Println(url)
	return url
}

func downMedia(mediaUrl string) {
	downLessonPath := "/Test/f/" + mediaUrl[0:strings.LastIndex(mediaUrl, "/")+1]
	if _, err := os.Stat(downLessonPath); err != nil {
		os.MkdirAll(downLessonPath, 0777)
	}
	utils.DownloadFile("http://www.bstcine.com/f/"+mediaUrl, downLessonPath+mediaUrl[strings.LastIndex(mediaUrl, "/"):len(mediaUrl)])
}
