package tools

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"../utils"
	"log"
	"io/ioutil"
	"net/http"
	"encoding/base64"
	"crypto/hmac"
	"hash"
	"crypto/sha1"
	"io"
	"sort"
	"bytes"
	"time"
	"errors"
	"encoding/json"
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
		bucket := tools.getBucket()

		_, err := os.Stat(serviceFilePath)
		isServiceRun := err == nil                    //是否在服务器运行
		isReplace := confMap["migrateReplace"] == "1" //是否覆盖上传

		jobs := make(chan OssInfo, rowCount)
		results := make(chan OssInfo, rowCount)

		for w := 1; w <= 10; w++ {
			go func(id int) {
				for ossObject := range jobs {
					objectKey := ossObject.ObjectKey
					migrateUrl := ossObject.MigrateUrl
					migratePath := ossObject.MigratePath

					isExist, err := bucket.IsObjectExist(objectKey)
					if err != nil {
						ossObject.Error = err
						results <- ossObject
						continue
					}

					//文件是否存在
					if isExist {
						if isReplace {
							bucket.DeleteObject(objectKey) //删除文件
						} else {
							ossObject.Error = errors.New("已经存在")
							results <- ossObject
							continue
						}
					}

					//上传
					if isServiceRun { //ECS
						err = bucket.PutObjectFromFile(objectKey, migratePath)
					} else { //本地
						/*localPath = localKjPath + objectKey
						utils.DownloadFile(objectUrl, localPath)
						err = bucket.PutObjectFromFile(objectKey, localPath)*/
						body, byteLen, err := utils.GetHttpFileBytes(migrateUrl)

						if err == nil && byteLen <= 1500 {
							err = errors.New("小于1.5K")
						}

						if err == nil {
							err = bucket.PutObject(objectKey, body)
						}
					}

					if err != nil {
						ossObject.Error = err
						results <- ossObject
					} else {
						results <- ossObject
					}
				}
			}(w)
		}

		for i := 0; i < rowCount; i++ {
			urls := strings.Split(rows[i].(string), ";")

			mediaUrl := urls[0]
			lessonId := urls[3]
			courseId := urls[4]

			var objectKey string
			var migrateUrl string
			var migratePath string

			if isCourseOrig {
				objectKey = "kj/" + mediaUrl
				migrateUrl = "http://www.bstcine.com/f/" + mediaUrl
				migratePath = serviceFilePath + mediaUrl
			} else {
				urlPrefix := urls[1]
				urlSuffix := urls[2]

				if strings.Contains(urlPrefix, "http://gcdn.bstcine.com/img") {
					objectKey = strings.Replace(urlPrefix, "http://gcdn.bstcine.com/", "", -1) + mediaUrl + urlSuffix
					objectKey = strings.Replace(objectKey, "/f/", "/", -1)
					objectKey = objectKey[0:strings.Index(objectKey, ".")] + ".jpg"
				} else {
					objectKey = "kj/" + mediaUrl
				}

				migrateUrl = urlPrefix + mediaUrl + urlSuffix
				migratePath = serviceKjFilePath + objectKey
			}

			jobs <- OssInfo{ObjectKey: objectKey, MigrateUrl: migrateUrl, MigratePath: migratePath, CourseId: courseId, LessonId: lessonId, Seq: strconv.Itoa(i + 1)}
		}
		close(jobs)

		migrateSuccessLogger := utils.GetLogger(workDir + "/log/migrate_success.log")
		migrateErrorLogger := utils.GetLogger(workDir + "/log/migrate_error.log")

		for a := 1; a <= rowCount; a++ {
			msg := <-results
			fmt.Printf("%v \n", msg)

			if msg.Error == nil {
				migrateSuccessLogger.Printf("%v", msg)
			} else {
				migrateErrorLogger.Printf("%v", msg)
			}
		}
	}
}

/**
资源权限设置
 */
func (tools Tools) SetObjectACL() {
	confMap := tools.ConfMap

	_, rows := utils.GetFiles(confMap["srcPassword"], "0", confMap["aclCourse"])

	bucket := tools.getBucket()

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
		err := bucket.SetObjectACL(objectKey, objectACL)
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

	bucket := tools.getBucket()

	jobs := make(chan OssInfo, rowCount)
	results := make(chan OssInfo, rowCount)

	for w := 1; w <= 10; w++ {
		go func(id int) {
			for ossObject := range jobs {
				header, err := bucket.GetObjectDetailedMeta(ossObject.ObjectKey)
				if err != nil {
					ossObject.Error = err
					results <- ossObject
					continue
				}

				length, err := strconv.Atoi(header.Get("Content-Length"))
				if err == nil {
					ossObject.Length = length
				} else {
					ossObject.Length = 0
				}

				headResp, err := http.Head(ossObject.MigrateUrl)
				if err == nil {
					ossObject.EcsLength = int(headResp.ContentLength)
				} else {
					ossObject.EcsLength = 0
				}

				results <- ossObject
			}
		}(w)
	}

	for i := 0; i < rowCount; i++ {
		urls := strings.Split(rows[i].(string), ";")

		mediaUrl := urls[0]
		lessonId := urls[3]
		courseId := urls[4]

		var objectKey string
		var migrateUrl string

		if isCourseOrig {
			objectKey = "kj/" + mediaUrl
			migrateUrl = "http://www.bstcine.com/f/" + mediaUrl
		} else {
			urlPrefix := urls[1]
			urlSuffix := urls[2]

			if strings.Contains(urlPrefix, "http://gcdn.bstcine.com/img") {
				objectKey = strings.Replace(urlPrefix, "http://gcdn.bstcine.com/", "", -1) + mediaUrl + urlSuffix
				objectKey = strings.Replace(objectKey, "/f/", "/", -1)
				objectKey = objectKey[0:strings.Index(objectKey, ".")] + ".jpg"
			} else {
				objectKey = "kj/" + mediaUrl
			}

			migrateUrl = urlPrefix + mediaUrl + urlSuffix
		}

		jobs <- OssInfo{ObjectKey: objectKey, MigrateUrl: migrateUrl, CourseId: courseId, LessonId: lessonId, Seq: strconv.Itoa(i + 1)}
	}
	close(jobs)

	migrateCheckOssLogger := utils.GetLogger(tools.WorkPath + "/log/migrate_check_oss.log")
	migrateCheckEcsLogger := utils.GetLogger(tools.WorkPath + "/log/migrate_check_ecs.log")
	migrateCheckEquallyLogger := utils.GetLogger(tools.WorkPath + "/log/migrate_check_equally.log")
	migrateCheckSmallLogger := utils.GetLogger(tools.WorkPath + "/log/migrate_check_small.log")

	for a := 1; a <= rowCount; a++ {
		msg := <-results
		ossLength := msg.Length
		ecsLength := msg.EcsLength

		if msg.Error != nil { //OSS 访问出错的资源
			migrateCheckOssLogger.Printf("CourseId: %s ; LessonId: %s ; OSS：%s ; ECS：%s ; ERROR: %+v \n", msg.CourseId, msg.LessonId, msg.ObjectKey, msg.MigrateUrl, msg.Error)
		} else {
			if ecsLength == 0 { //ECS 不存在的资源
				migrateCheckEcsLogger.Printf("CourseId: %s ; LessonId: %s ; OSS：%s ; ECS：%s ; OSS-SIZE: %+v; ECS-SIZE: %+v \n", msg.CourseId, msg.LessonId, msg.ObjectKey, msg.MigrateUrl, ossLength, ecsLength)
			} else if ossLength < ecsLength { //迁移失败的资源
				migrateCheckEquallyLogger.Printf("CourseId: %s ; LessonId: %s ; OSS：%s ; ECS：%s ; OSS-SIZE: %+v; ECS-SIZE: %+v \n", msg.CourseId, msg.LessonId, msg.ObjectKey, msg.MigrateUrl, ossLength, ecsLength)
			} else if ossLength <= 1500 && msg.ObjectKey != "kj/" && strings.Count(msg.ObjectKey, "/") >= 3 { //小于1.5KB的资源
				bucket.DeleteObject(msg.ObjectKey)
				migrateCheckSmallLogger.Printf("CourseId: %s ; LessonId: %s ; OSS：%s ; ECS：%s ; OSS-SIZE: %+v; ECS-SIZE: %+v \n", msg.CourseId, msg.LessonId, msg.ObjectKey, msg.MigrateUrl, ossLength, ecsLength)
			}
		}

		fmt.Printf("%s/%d %+v \n", msg.Seq, rowCount, msg)
	}
}

/**
资源(kj)中非 jpg 图片转 jpg
 */
func (tools Tools) ImgFormatJPG() {
	workDir := tools.WorkPath
	confMap := tools.ConfMap

	isFormatDel := confMap["imgFormatDel"] == "1"
	bucket := tools.getBucket()

	_, rows := utils.GetFiles(confMap["srcPassword"], "0", confMap["imgCourse"])
	rowCount := len(rows)

	jobs := make(chan OssInfo, rowCount)
	results := make(chan OssInfo, rowCount)

	for w := 1; w <= 25; w++ {
		go func(id int) {
			for ossObject := range jobs {
				objectKey := ossObject.ObjectKey

				if !strings.Contains(objectKey,".") {
					ossObject.Error = errors.New("非法文件")
					results <- ossObject
					continue
				}

				suf := objectKey[strings.LastIndex(objectKey, "."):]

				if suf != ".mp3" && suf != ".mp4" && suf != ".jpg" {
					msg, err := tools.imgProcessSave(objectKey, objectKey[0:strings.LastIndex(objectKey, ".")]+".jpg", "image/format,jpg")
					if err == nil {
						var res map[string]interface{}
						json.Unmarshal([]byte(msg), &res)

						status, ok := res["status"]
						if res != nil && ok && status == "OK" {
							ossObject.Remark = "格式化成功:" + msg

							if isFormatDel { //格式化并删除原文件
								_, err = bucket.CopyObject(objectKey, "del/"+objectKey)
								if err == nil {
									bucket.DeleteObject(objectKey)
									ossObject.Remark += " 迁移到 Del 文件夹成功，并删除原文件"
								} else {
									ossObject.Remark += " 迁移到 Del 文件夹失败"
								}
							}
						} else {
							ossObject.Error = errors.New("格式化失败：" + msg)
						}
					} else {
						ossObject.Error = err
					}
					results <- ossObject
				} else {
					ossObject.Remark = "无需格式化"
					results <- ossObject
				}

			}
		}(w)
	}

	for i := 0; i < rowCount; i++ {
		urls := strings.Split(rows[i].(string), ";")

		mediaUrl := urls[0]
		courseId := urls[4]
		lessonId := urls[3]

		jobs <- OssInfo{Seq: strconv.Itoa(i + 1), CourseId: courseId, LessonId: lessonId, ObjectKey: "kj/" + mediaUrl}
	}
	close(jobs)

	ossImgFormatLogger := utils.GetLogger(workDir + "/log/oss_img_format.log")

	for a := 1; a <= rowCount; a++ {
		ossObject := <-results

		ossImgFormatLogger.Println(ossObject)

		fmt.Println(ossObject)
	}
}

/**
课件资源图片加水印
 */
func (tools Tools) ImgWaterMark() {
	workDir := tools.WorkPath
	confMap := tools.ConfMap

	_, rows := utils.GetFiles(confMap["srcPassword"], "0", confMap["imgCourse"])
	rowCount := len(rows)

	jobs := make(chan OssInfo, rowCount)
	results := make(chan OssInfo, rowCount)

	for w := 1; w <= 25; w++ {
		go func(id int) {
			for ossObject := range jobs {
				courseId := ossObject.CourseId
				mediaUrl := ossObject.ObjectKey[3:]

				name := mediaUrl[0:strings.LastIndex(mediaUrl, ".")]
				suf := mediaUrl[strings.LastIndex(mediaUrl, "."):]

				if suf == ".jpg" {
					msg, err := tools.imgProcessSave("kj/"+name+".jpg", "img/"+courseId+"/"+name+".jpg", "style/"+confMap["imgStyle"])

					if err == nil {
						var res map[string]interface{}
						json.Unmarshal([]byte(msg), &res)

						status, ok := res["status"]
						if res != nil && ok && status == "OK" {
							ossObject.Remark = "原图生成水印图成功:" + msg
						} else {
							ossObject.Error = errors.New("原图生成水印图失败：" + msg)
						}
					} else {
						ossObject.Error = err
					}

					results <- ossObject
				} else {
					ossObject.Remark = "无需加水印"
					results <- ossObject
				}
			}
		}(w)
	}

	for i := 0; i < rowCount; i++ {
		urls := strings.Split(rows[i].(string), ";")

		courseId := urls[4]
		lessonId := urls[3]
		mediaUrl := urls[0]

		jobs <- OssInfo{Seq: strconv.Itoa(i + 1), CourseId: courseId, LessonId: lessonId, ObjectKey: "kj/" + mediaUrl}
	}
	close(jobs)

	ossImgWatermarkLogger := utils.GetLogger(workDir + "/log/oss_img_watermark.log")

	for a := 1; a <= rowCount; a++ {
		msg := <-results

		ossImgWatermarkLogger.Println(msg)

		fmt.Println(msg)
	}
}

/**
######################################################
###################  阿里提供的 OSS API  ##############
######################################################
 */

/**
资源信息类
*/
type OssInfo struct {
	ObjectKey   string //对象key
	ObjectUrl   string //对象路径
	MigrateUrl  string //迁移路径
	MigratePath string //迁移本地路径
	Length      int    //长度
	EcsLength   int    //ECS 长度
	CourseId    string //Course Id
	LessonId    string //Lesson Id
	Error       error  //error
	Seq         string //seq
	Remark      string // 备注
}

/**
获取 Oss Bucket
 */
func (tools Tools) getBucket() (*oss.Bucket) {
	confMap := tools.ConfMap

	client, err := oss.New(confMap["Endpoint"], confMap["AccessKeyId"], confMap["AccessKeySecret"])
	if err != nil {
		tools.HandleError(err)
	}

	bucket, err := client.Bucket(confMap["Bucket"])
	if err != nil {
		tools.HandleError(err)
	}

	return bucket
}

/**
OSS 图片处理并保存
 */
func (tools Tools) imgProcessSave(objKey, newObjKey, process string) (string, error) {
	var bucket = "static-bstcine"
	var region = "oss-cn-shanghai"
	var ossHost = "http://" + bucket + "." + region + ".aliyuncs.com/"

	newObjKey = base64.StdEncoding.EncodeToString([]byte(newObjKey))
	bucket = base64.StdEncoding.EncodeToString([]byte(bucket))

	client := &http.Client{}

	url := ossHost + objKey + "?x-oss-process"
	data := "x-oss-process=" + process + "|sys/saveas,o_" + newObjKey + ",b_" + bucket
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	ossDate := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set(HTTPHeaderDate, ossDate)

	tools.signHeader(req, "/static-bstcine/"+objKey+"?x-oss-process")

	resp, err := client.Do(req)
	if(err != nil) {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

const (
	HTTPHeaderAuthorization      = "Authorization"
	HTTPHeaderCacheControl       = "Cache-Control"
	HTTPHeaderContentDisposition = "Content-Disposition"
	HTTPHeaderContentEncoding    = "Content-Encoding"
	HTTPHeaderContentLength      = "Content-Length"
	HTTPHeaderContentMD5         = "Content-MD5"
	HTTPHeaderContentType        = "Content-Type"
	HTTPHeaderContentLanguage    = "Content-Language"
	HTTPHeaderDate               = "Date"
)

func (tools Tools) signHeader(req *http.Request, canonicalizedResource string) {
	// Get the final Authorization' string
	authorizationStr := "OSS " + tools.ConfMap["AccessKeyId"] + ":" + tools.getSignedStr(req, canonicalizedResource)

	// Give the parameter "Authorization" value
	req.Header.Set(HTTPHeaderAuthorization, authorizationStr)
}

func (tools Tools) getSignedStr(req *http.Request, canonicalizedResource string) string {
	// Find out the "x-oss-"'s address in this request'header
	temp := make(map[string]string)

	for k, v := range req.Header {
		if strings.HasPrefix(strings.ToLower(k), "x-oss-") {
			temp[strings.ToLower(k)] = v[0]
		}
	}
	hs := newHeaderSorter(temp)

	// Sort the temp by the Ascending Order
	hs.Sort()

	// Get the CanonicalizedOSSHeaders
	canonicalizedOSSHeaders := ""
	for i := range hs.Keys {
		canonicalizedOSSHeaders += hs.Keys[i] + ":" + hs.Vals[i] + "\n"
	}

	// Give other parameters values
	// when sign url, date is expires
	date := req.Header.Get(HTTPHeaderDate)
	contentType := req.Header.Get(HTTPHeaderContentType)
	contentMd5 := req.Header.Get(HTTPHeaderContentMD5)

	signStr := req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + canonicalizedResource
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte("XOssD3DnWffLiJaSgWjFdV0kHzJeIC"))
	io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signedStr
}

// 用于signHeader的字典排序存放容器。
type headerSorter struct {
	Keys []string
	Vals []string
}

// Additional function for function SignHeader.
func (hs *headerSorter) Sort() {
	sort.Sort(hs)
}

// Additional function for function SignHeader.
func (hs *headerSorter) Len() int {
	return len(hs.Vals)
}

// Additional function for function SignHeader.
func (hs *headerSorter) Less(i, j int) bool {
	return bytes.Compare([]byte(hs.Keys[i]), []byte(hs.Keys[j])) < 0
}

// Additional function for function SignHeader.
func (hs *headerSorter) Swap(i, j int) {
	hs.Vals[i], hs.Vals[j] = hs.Vals[j], hs.Vals[i]
	hs.Keys[i], hs.Keys[j] = hs.Keys[j], hs.Keys[i]
}

func newHeaderSorter(m map[string]string) *headerSorter {
	hs := &headerSorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]string, 0, len(m)),
	}

	for k, v := range m {
		hs.Keys = append(hs.Keys, k)
		hs.Vals = append(hs.Vals, v)
	}
	return hs
}
