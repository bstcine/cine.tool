package main

import (
	"fmt"
	"strings"
	"os"
	"./utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"strconv"
)

var logName = "app_oss.log"
var debugLog *log.Logger

var cfgName = "app_oss.cfg"
var serviceFilePath = "/mnt/web/app.bstcine.com/wwwroot/public/f/"
var serviceKjFilePath = "/mnt/web/kj.bstcine.com/wwwroot/"
var localKjPath string

var curPath string
var outPutPath string
var argsMap map[string]string

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

func main() {
	var debug = false

	if debug {
		curPath = "/Go/Cine/cine.tool/assets/"
		outPutPath = "/Test/"
		localKjPath = "/Test/wwwroot/"
		cfgName = "app_oss_tmp.cfg"
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

		//是否在服务器运行
		_,err = os.Stat(serviceFilePath)
		isServiceRun := err == nil

		jobs := make(chan string, len(rows))
		results := make(chan string, len(rows))

		for w := 1; w <= 6; w++ {
			go func(id int) {
				for row:= range jobs {
					urls := strings.Split(row, ";")

					mediaUrl := urls[0]
					urlPrefix := urls[1]
					urlSuffix := urls[2]

					var objectKey string
					var objectUrl string
					var localPath string

					if isCourseOrig {
						objectKey = "kj/" + mediaUrl
						objectUrl = "http://www.bstcine.com/f/" + mediaUrl
						localPath = serviceFilePath + mediaUrl
					} else {
						if strings.Contains(urlPrefix,"http://gcdn.bstcine.com/img") {
							objectKey = strings.Replace(urlPrefix, "http://gcdn.bstcine.com/", "", -1) + mediaUrl + urlSuffix
							objectKey = strings.Replace(objectKey,"/f/","/",-1)
							objectKey = objectKey[0:strings.Index(objectKey,".")] + ".jpg"
						}else {
							objectKey = "kj/" + mediaUrl
						}

						objectUrl = urlPrefix + mediaUrl + urlSuffix
						localPath = serviceKjFilePath + objectKey
					}

					isExist, err := bucket.IsObjectExist(objectKey)
					if err != nil {
						handleError(err)
					}

					if isExist {
						results<- "worker " + strconv.Itoa(id) + ": " + objectKey + " 已经存在"
						continue
					}

					if !isServiceRun {//客户端下载
						localPath = localKjPath + objectKey
						utils.DownloadFile(objectUrl, localPath)
					}

					err = bucket.PutObjectFromFile(objectKey, localPath)
					if err != nil {
						handleError(err)
					}else {
						results<- "worker " + strconv.Itoa(id) + ": " + localPath + " => "+objectKey + " 上传成功"
					}
				}
			}(w)
		}

		for i := 0; i < len(rows); i++ {
			jobs <- rows[i].(string)
		}
		close(jobs)

		for a := 1; a <= len(rows); a++ {
			msg := <-results
			fmt.Printf("%d/%d: %s \n",a,len(rows),msg)
			debugLog.Printf("%d/%d: %s",a,len(rows),msg)
		}
	}

	fmt.Println("请输入任意键结束...")
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
