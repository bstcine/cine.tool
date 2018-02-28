package main

import (
	"fmt"
	"strings"
	"os"
	"./utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var cfgName = "app_oss.cfg"

var curPath string
var outPutPath string
var argsMap map[string]string

func main() {
	var debug = false

	if debug {
		curPath = "/Go/Cine/cine.tool/assets/"
		outPutPath = "/Test/"
	} else {
		curPath = utils.GetCurPath() + string(os.PathSeparator)
		outPutPath = utils.GetCurPath() + string(os.PathSeparator)
	}

	argsMap = utils.GetConfArgs(curPath + cfgName)
	if argsMap == nil || len(argsMap) <= 0 {
		fmt.Println("配置文件不存在")
		return
	}

	if argsMap["srcType"] == "list" {
		getObjectList()
	} else if argsMap["srcType"] == "acl" {
		setOssObjectACL()
	}

}

/**
获取资源清单
 */
func getObjectList() {
	if argsMap["listType"] != "0" {
		fmt.Println("暂时只支持获取课件资源")
		return
	}

	_, rows := utils.GetFiles(argsMap["srcPassword"], argsMap["listType"], argsMap["listCourse"])

	var kjFiles, kjCdnFiles []string
	for i := 0; i < len(rows); i++ {
		row := rows[i]

		kjCdnFiles = append(kjCdnFiles, regUrl(false, row.(string)))
		kjFiles = append(kjFiles, regUrl(true, row.(string)))
	}
	utils.WriteLines(kjFiles, outPutPath+argsMap["listOutFileName"])
	utils.WriteLines(kjCdnFiles, outPutPath+argsMap["listOutCdnFileName"])

	fmt.Printf("共有 %d 个资源",len(rows))
}

/**
设置资源权限
 */
func setOssObjectACL() {
	_, rows := utils.GetFiles(argsMap["srcPassword"], "0", argsMap["aclCourse"])

	client, err := oss.New(argsMap["Endpoint"], argsMap["AccessKeyId"], argsMap["AccessKeySecret"])
	if err != nil {
		fmt.Println("Oss 访问请求失败，请检查网络或密钥等配置项")
		return
	}

	bucket, err := client.Bucket(argsMap["Bucket"])
	if err != nil {
		fmt.Println("Oss Bucket 访问失败，请检查Bucket配置是否存在")
		return
	}

	objectACL := oss.ACLDefault

	switch argsMap["aclType"] {
	case "default":
		objectACL = oss.ACLDefault
	case "public-read-write":
		objectACL = oss.ACLPublicReadWrite
	case "public-read":
		objectACL = oss.ACLPublicRead
	case "private":
		objectACL = oss.ACLPrivate
	}

	for i := 0; i < len(rows); i++ {
		row := rows[i]
		urls := strings.Split(row.(string), ";")
		mediaUrl := urls[0]
		urlPrefix := urls[1]
		urlSuffix := urls[2]

		urlPrefix = strings.Replace(urlPrefix, "http://gcdn.bstcine.com/", "", -1)
		objectKey := urlPrefix + mediaUrl + urlSuffix

		// 设置Object的访问权限
		err = bucket.SetObjectACL(objectKey, objectACL)
		if err != nil {
			// HandleError(err)
			fmt.Println(err)
		}else {
			fmt.Println("http://oss.bstcine.com/"+objectKey+" ==> acl set :" + argsMap["aclType"])
		}

	}
}

func regUrl(isOrig bool, param string) (url string) {
	urls := strings.Split(param, ";")

	mediaUrl := urls[0]
	urlPrefix := urls[1]
	urlSuffix := urls[2]

	fmt.Println(mediaUrl)

	if isOrig {
		url = "http://www.bstcine.com/ f/" + mediaUrl
	} else if strings.Contains(urlPrefix, "http://gcdn.bstcine.com") {
		urlPrefix = strings.Replace(urlPrefix, "/img/", "/ img/", -1)
		urlPrefix = strings.Replace(urlPrefix, "/mp3/", "/ mp3/", -1)

		url = urlPrefix + mediaUrl + urlSuffix
	}
	return url
}

func downMedia(mediaUrl string) {
	downLessonPath := "/Test/f/" + mediaUrl[0:strings.LastIndex(mediaUrl, "/")+1]
	if _, err := os.Stat(downLessonPath); err != nil {
		os.MkdirAll(downLessonPath, 0777)
	}
	utils.DownloadFile("http://www.bstcine.com/f/"+mediaUrl, downLessonPath+mediaUrl[strings.LastIndex(mediaUrl, "/"):len(mediaUrl)])
}
