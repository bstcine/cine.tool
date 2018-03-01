package main

import (
	"fmt"
	"strings"
	"os"
	"./utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
)

var logName = "app_oss.log"
var debugLog *log.Logger

var cfgName = "app_oss.cfg"
var serviceAppPath = "/mnt/web/app.bstcine.com/wwwroot/public/f/"
var serviceKjPath = "/mnt/web/kj.bstcine.com/wwwroot/"
var localKjPath string

var curPath string
var outPutPath string
var argsMap map[string]string

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

func main() {
	var debug = true

	if debug {
		curPath = "/Go/Cine/cine.tool/assets/"
		outPutPath = "/Test/"
		localKjPath = "/Test/wwwroot/"
	} else {
		curPath = utils.GetCurPath() + string(os.PathSeparator)
		outPutPath = utils.GetCurPath() + string(os.PathSeparator)
		localKjPath = utils.GetCurPath() + string(os.PathSeparator) + "wwwroot" + string(os.PathSeparator)
	}

	logFile,err  := os.Create(outPutPath+logName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	debugLog = log.New(logFile,"[Info]",log.Llongfile)

	argsMap = utils.GetConfArgs(curPath + cfgName)
	if argsMap == nil || len(argsMap) <= 0 {
		fmt.Println("配置文件不存在")
		return
	}

	if argsMap["srcType"] == "migrate" {
		migrateObject()
	} else if argsMap["srcType"] == "acl" {
		setOssObjectACL()
	}

}

/**
资源迁移
 */
func migrateObject() {
	if argsMap["migrateType"] != "0" {
		fmt.Println("暂时只支持获取课件资源")
		return
	}

	//迁移课程资源的类型 是否为原始资源
	isCourseOrig := argsMap["migrateCourseType"] == "orig"

	_, rows := utils.GetFiles(argsMap["srcPassword"], argsMap["migrateType"], argsMap["migrateCourse"])

	if argsMap["migrateModel"] == "list" { //获取资源清单
		var listFiles []string
		for i := 0; i < len(rows); i++ {
			row := rows[i]
			listFiles = append(listFiles, regUrl(isCourseOrig, row.(string)))
		}
		utils.WriteLines(listFiles, outPutPath+argsMap["migrateListFileName"])

		fmt.Printf("%s 课程,共有 %d 个 %s 资源,已经生成到 %s", argsMap["migrateCourse"], len(listFiles), argsMap["migrateCourseType"], outPutPath+argsMap["migrateListFileName"])
	} else if argsMap["migrateModel"] == "local" { //本地资源上传
		client, err := oss.New(argsMap["Endpoint"], argsMap["AccessKeyId"], argsMap["AccessKeySecret"])
		if err != nil {
			handleError(err)
			return
		}

		bucket, err := client.Bucket(argsMap["Bucket"])
		if err != nil {
			handleError(err)
			return
		}

		for i := 0; i < len(rows); i++ {
			row := rows[i]
			urls := strings.Split(row.(string), ";")

			mediaUrl := urls[0]
			urlPrefix := urls[1]
			urlSuffix := urls[2]

			var objectKey string
			var objectPath string
			var objectUrl string

			if isCourseOrig {
				objectKey = "kj/" + mediaUrl
				objectUrl = "http://www.bstcine.com/f/" + mediaUrl
				objectPath = serviceAppPath + mediaUrl
			} else {
				objectKey = strings.Replace(urlPrefix, "http://gcdn.bstcine.com/", "", -1) + mediaUrl + urlSuffix
				objectUrl = urlPrefix + mediaUrl + urlSuffix
				objectPath = serviceKjPath + objectKey
			}

			isExist, err := bucket.IsObjectExist(objectKey)
			if err != nil {
				handleError(err)
			}

			if isExist {
				log.Printf("%d/%d: %s 已经存在",i+1,len(rows),objectKey)
				debugLog.Printf("%d/%d: %s 已经存在",i+1,len(rows),objectKey)
				continue
			}

			if _, err := os.Stat(serviceAppPath); err != nil { //客户端
				objectPath = downMedia(objectUrl)
			}

			err = bucket.PutObjectFromFile(objectKey, objectPath)
			if err != nil {
				handleError(err)
			}else {
				log.Printf("%d/%d: %s => %s 上传成功",i+1,len(rows),objectPath,objectKey)
				debugLog.Printf("%d/%d: %s => %s 上传成功",i+1,len(rows),objectPath,objectKey)
			}
		}
	}

	fmt.Println("请输入任意键结算进程...")
	fmt.Scanln()
}

/**
设置资源权限
 */
func setOssObjectACL() {
	_, rows := utils.GetFiles(argsMap["srcPassword"], "0", argsMap["aclCourse"])

	client, err := oss.New(argsMap["Endpoint"], argsMap["AccessKeyId"], argsMap["AccessKeySecret"])
	if err != nil {
		handleError(err)
		return
	}

	bucket, err := client.Bucket(argsMap["Bucket"])
	if err != nil {
		handleError(err)
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
			handleError(err)
		} else {
			log.Printf("%s set acl: %s",objectKey,argsMap["aclType"])
			debugLog.Printf("%s set acl: %s",objectKey,argsMap["aclType"])
		}

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
	return url
}

func downMedia(url string) (downPath string) {
	path := strings.Replace(url, "http://gcdn.bstcine.com/", "", -1)
	path = strings.Replace(path, "http://www.bstcine.com/f/", "", -1)

	downPrefix := localKjPath + path[0:strings.LastIndex(path, "/")+1]
	if _, err := os.Stat(downPrefix); err != nil {
		os.MkdirAll(downPrefix, 0777)
	}

	downPath = localKjPath + path

	utils.DownloadFile(url, downPath)

	return downPath
}
