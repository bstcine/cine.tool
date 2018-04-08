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

		_, err = os.Stat(serviceFilePath)
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
						body, err := utils.GetHttpFileBytes(migrateUrl)
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

				length := header.Get("Content-Length")
				ossObject.Length = length

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

	for a := 1; a <= rowCount; a++ {
		msg := <-results
		length := msg.Length

		if i, err := strconv.Atoi(length); i <= 10000 || err != nil || msg.Error != nil {
			if err == nil && i <= 10000 && msg.ObjectKey != "kj/" && len(msg.ObjectKey) > 5 {
				//bucket.DeleteObject(msg.ObjectKey)
				tools.GetLogger().Printf("CourseId: %s ; LessonId: %s ; OSS：%s ; ECS：%s ;SIZE: %sB ; ERROR: %+v ; DEL\n", msg.CourseId, msg.LessonId, msg.ObjectKey, msg.MigrateUrl, msg.Length, msg.Error)
			} else {
				tools.GetLogger().Printf("CourseId: %s ; LessonId: %s ; OSS：%s ; ECS：%s ;SIZE: %sB ; ERROR: %+v \n", msg.CourseId, msg.LessonId, msg.ObjectKey, msg.MigrateUrl, msg.Length, msg.Error)
			}
		}

		fmt.Printf("%s/%d %s \n", msg.Seq, rowCount, msg)
	}
}

/**
资源(kj)中非 jpg 图片转 jpg
 */
func (tools Tools) ImgFormatJPG() {
	confMap := tools.ConfMap

	_, rows := utils.GetFiles(confMap["srcPassword"], "0", confMap["imgCourse"])
	rowCount := len(rows)

	jobs := make(chan []string, rowCount)
	results := make(chan []string, rowCount)

	for w := 1; w <= 25; w++ {
		go func(id int) {
			for ossObject := range jobs {
				objectKey := ossObject[2]

				suf := objectKey[strings.LastIndex(objectKey, "."):len(objectKey)]

				if suf != ".mp3" && suf != ".mp4" && suf != ".jpg" {
					msg := tools.imgProcessSave(objectKey, objectKey[0:strings.LastIndex(objectKey, ".")]+".jpg", "image/format,jpg")
					results <- append(ossObject, "格式化成功:"+msg)
				} else {
					results <- append(ossObject, "无需格式化")
				}

			}
		}(w)
	}

	for i := 0; i < rowCount; i++ {
		urls := strings.Split(rows[i].(string), ";")

		mediaUrl := urls[0]
		lessonId := urls[3]

		var objectKey string
		var objectUrl string

		objectKey = "kj/" + mediaUrl
		objectUrl = "http://oss.bstcine.com/" + objectKey

		jobs <- []string{strconv.Itoa(i+1) + "/" + strconv.Itoa(rowCount), lessonId, objectKey, objectUrl}
	}
	close(jobs)

	for a := 1; a <= rowCount; a++ {
		msg := <-results
		tools.GetLogger().Println(msg)
		fmt.Println(msg)
	}
}

/**
课件资源图片加水印
 */
func (tools Tools) ImgWaterMark() {
	confMap := tools.ConfMap

	_, rows := utils.GetFiles(confMap["srcPassword"], "0", confMap["imgCourse"])
	rowCount := len(rows)

	jobs := make(chan []string, rowCount)
	results := make(chan []string, rowCount)

	for w := 1; w <= 25; w++ {
		go func(id int) {
			for ossObject := range jobs {
				courseId := ossObject[1]
				mediaUrl := ossObject[3]

				name := mediaUrl[0:strings.LastIndex(mediaUrl, ".")]
				suf := mediaUrl[strings.LastIndex(mediaUrl, "."):len(mediaUrl)]

				if suf == ".jpg" {
					msg := tools.imgProcessSave("kj/"+name+".jpg", "img/"+courseId+"/"+name+".jpg", "style/"+confMap["imgStyle"])
					results <- append(ossObject, "图片加水印-"+msg)
				} else {
					results <- append(ossObject, "无需加水印")
				}

			}
		}(w)
	}

	for i := 0; i < rowCount; i++ {
		urls := strings.Split(rows[i].(string), ";")

		courseId := urls[4]
		lessonId := urls[3]
		mediaUrl := urls[0]

		jobs <- []string{strconv.Itoa(i+1) + "/" + strconv.Itoa(rowCount), courseId, lessonId, mediaUrl}
	}
	close(jobs)

	for a := 1; a <= rowCount; a++ {
		msg := <-results
		tools.GetLogger().Println(msg)
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
	MigrateUrl  string //迁移路径
	MigratePath string //迁移本地路径
	Length      string //长度
	CourseId    string //Course Id
	LessonId    string //Lesson Id
	Error       error  //error
	Seq         string //seq
	Remark      string // 备注
}

func (tools Tools) getClient() (*oss.Client, error) {
	if tools.OSSClient != nil {
		return tools.OSSClient, nil
	}

	confMap := tools.ConfMap

	client, err := oss.New(confMap["Endpoint"], confMap["AccessKeyId"], confMap["AccessKeySecret"])
	if err != nil {
		tools.HandleError(err)
	}

	tools.OSSClient = client

	return client, err
}

func (tools Tools) getBucket() (*oss.Bucket, error) {
	if tools.OSSBucket != nil {
		return tools.OSSBucket, nil
	}

	confMap := tools.ConfMap

	client, err := tools.getClient()
	if err != nil {
		tools.HandleError(err)
	}

	bucket, err := client.Bucket(confMap["Bucket"])
	if err != nil {
		tools.HandleError(err)
	}

	tools.OSSBucket = bucket

	return bucket, err
}

/**
OSS 图片处理
 */
func (tools Tools) imgProcessSave(objKey, newObjKey, process string) string {
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
		// handle error
	}

	ossDate := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set(HTTPHeaderDate, ossDate)

	tools.signHeader(req, "/static-bstcine/"+objKey+"?x-oss-process")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	return string(body)
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
