package tools

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"../utils"
	"log"
)

var serviceFilePath = "/mnt/web/app.bstcine.com/wwwroot/public/f/"
var serviceKjFilePath = "/mnt/web/kj.bstcine.com/wwwroot/"

/**
获取	阿里云清单路径
 */
func getHttpListUrl(isOrig bool, param string) (url string) {
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

/**
资源迁移
 */
func (tools Tools) MigrateObject() {
	workDir := tools.WorkPath
	confMap := tools.ConfMap

	if confMap["migrateType"] != "0" {
		fmt.Println("暂时只支持获取课件资源")
		return
	}

	//迁移课程资源的类型 是否为原始资源
	isCourseOrig := confMap["migrateCourseType"] == "orig"

	_, rows := utils.GetFiles(confMap["srcPassword"], confMap["migrateType"], confMap["migrateCourse"])
	rowCount := len(rows)

	if confMap["migrateModel"] == "list" { //获取资源清单
		var listFiles []string
		for i := 0; i < len(rows); i++ {
			row := rows[i]
			listFiles = append(listFiles, getHttpListUrl(isCourseOrig, row.(string)))
		}
		utils.WriteLines(listFiles, workDir+confMap["migrateListFileName"])

		fmt.Printf("%s 课程,共有 %d 个 %s 资源,已经生成到 %s", confMap["migrateCourse"], len(listFiles), confMap["migrateCourseType"], workDir+confMap["migrateListFileName"])
	} else if confMap["migrateModel"] == "local" { //本地资源上传
		client, err := oss.New(confMap["Endpoint"], confMap["AccessKeyId"], confMap["AccessKeySecret"])
		if err != nil {
			tools.HandleError(err)
			return
		}

		bucket, err := client.Bucket(confMap["Bucket"])
		if err != nil {
			tools.HandleError(err)
			return
		}

		//是否在服务器运行
		_, err = os.Stat(serviceFilePath)
		isServiceRun := err == nil

		jobs := make(chan []string, rowCount)
		results := make(chan string, rowCount)

		for w := 1; w <= 10; w++ {
			go func(id int) {
				for ossObject := range jobs {
					objectKey := ossObject[0]
					objectUrl := ossObject[1]
					localPath := ossObject[2]
					objectNo := ossObject[3]

					msg := "worker-" + strconv.Itoa(id) + "-" + objectNo + "/" + strconv.Itoa(rowCount) + ": "

					isExist, err := bucket.IsObjectExist(objectKey)
					if err != nil {
						results <- msg + objectKey + " 检查失败 => " + err.Error()
						continue
					}

					if isExist {
						results <- msg + objectKey + " 已经存在"
						continue
					}

					if isServiceRun { //ESC
						err = bucket.PutObjectFromFile(objectKey, localPath)
					} else { //本地
						/*localPath = localKjPath + objectKey
						utils.DownloadFile(objectUrl, localPath)
						err = bucket.PutObjectFromFile(objectKey, localPath)*/
						body, err := utils.GetHttpFileBytes(objectUrl)
						if err == nil {
							err = bucket.PutObject(objectKey, body)
						}
					}

					if err != nil {
						results <- msg + localPath + " => " + objectKey + " 上传失败 => " + err.Error()
					} else {
						results <- msg + localPath + " => " + objectKey + " 上传成功"
					}
				}
			}(w)
		}

		for i := 0; i < rowCount; i++ {
			urls := strings.Split(rows[i].(string), ";")

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
				if strings.Contains(urlPrefix, "http://gcdn.bstcine.com/img") {
					objectKey = strings.Replace(urlPrefix, "http://gcdn.bstcine.com/", "", -1) + mediaUrl + urlSuffix
					objectKey = strings.Replace(objectKey, "/f/", "/", -1)
					objectKey = objectKey[0:strings.Index(objectKey, ".")] + ".jpg"
				} else {
					objectKey = "kj/" + mediaUrl
				}

				objectUrl = urlPrefix + mediaUrl + urlSuffix
				localPath = serviceKjFilePath + objectKey
			}

			jobs <- []string{objectKey, objectUrl, localPath, strconv.Itoa(i + 1)}
		}
		close(jobs)

		for a := 1; a <= rowCount; a++ {
			msg := <-results
			fmt.Printf("%s \n", msg)
			tools.GetLogger().Printf("%s", msg)
		}
	}
}

/**
资源权限设置
 */
func (tools Tools) SetObjectACL() {
	confMap := tools.ConfMap

	_, rows := utils.GetFiles(confMap["srcPassword"], "0", confMap["aclCourse"])

	client, err := oss.New(confMap["Endpoint"], confMap["AccessKeyId"], confMap["AccessKeySecret"])
	if err != nil {
		tools.HandleError(err)
		return
	}

	bucket, err := client.Bucket(confMap["Bucket"])
	if err != nil {
		tools.HandleError(err)
		return
	}

	objectACL := oss.ACLDefault
	aclType := confMap["aclType"]

	switch aclType {
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
			tools.HandleError(err)
		} else {
			log.Printf("%s set acl: %s", objectKey, aclType)
		}
	}
}

/**
资源迁移校验
 */
func (tools Tools) MigrateCheck() {
	confMap := tools.ConfMap
	if confMap["migrateType"] != "0" {
		fmt.Println("暂时只支持获取课件资源")
		return
	}

	//迁移课程资源的类型 是否为原始资源
	isCourseOrig := confMap["migrateCourseType"] == "orig"

	_, rows := utils.GetFiles(confMap["srcPassword"], confMap["migrateType"], confMap["migrateCourse"])
	rowCount := len(rows)

	client, err := oss.New(confMap["Endpoint"], confMap["AccessKeyId"], confMap["AccessKeySecret"])
	if err != nil {
		tools.HandleError(err)
		return
	}

	bucket, err := client.Bucket(confMap["Bucket"])
	if err != nil {
		tools.HandleError(err)
		return
	}

	jobs := make(chan []string, rowCount)
	results := make(chan []string, rowCount)

	for w := 1; w <= 10; w++ {
		go func(id int) {
			for ossObject := range jobs {
				objectKey := ossObject[0]

				header, err := bucket.GetObjectDetailedMeta(objectKey)
				if err != nil {
					results <- append(ossObject, "0", err.Error())
					continue
				}

				length := header.Get("Content-Length")
				results <- append(ossObject, length)
			}
		}(w)
	}

	for i := 0; i < rowCount; i++ {
		urls := strings.Split(rows[i].(string), ";")

		mediaUrl := urls[0]
		urlPrefix := urls[1]
		urlSuffix := urls[2]

		var objectKey string
		var objectUrl string

		if isCourseOrig {
			objectKey = "kj/" + mediaUrl
			objectUrl = "http://www.bstcine.com/f/" + mediaUrl
		} else {
			if strings.Contains(urlPrefix, "http://gcdn.bstcine.com/img") {
				objectKey = strings.Replace(urlPrefix, "http://gcdn.bstcine.com/", "", -1) + mediaUrl + urlSuffix
				objectKey = strings.Replace(objectKey, "/f/", "/", -1)
				objectKey = objectKey[0:strings.Index(objectKey, ".")] + ".jpg"
			} else {
				objectKey = "kj/" + mediaUrl
			}

			objectUrl = urlPrefix + mediaUrl + urlSuffix
		}

		jobs <- []string{objectKey, objectUrl, strconv.Itoa(i + 1)}
	}
	close(jobs)

	for a := 1; a <= rowCount; a++ {
		msg := <-results
		objectKey := msg[0]
		objectUrl := msg[1]
		length := msg[3]

		if i, err := strconv.Atoi(length); i <= 162 || err != nil {
			if objectKey != "kj/" && len(objectKey) > 5 {
				bucket.DeleteObject(objectKey)
			}

			tools.GetLogger().Printf("%s", objectUrl+" 上传失败")
		}

		fmt.Printf("%s/%d %s \n", msg[2], rowCount, msg)
	}
}
