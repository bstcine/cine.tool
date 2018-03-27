package main

import (
	"./utils"
	"./model"
	"fmt"
	"os"
	"bufio"
	"strings"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

/**
 * @param oss资源下载器
 * @param 下载步骤如下
 * 1.提示用户输入账户和密码，登入服务器（获取数据读取权限token）
 * 2.读取配置文件，获取待下载课件列表（courseIds）或课程列表（lessonIds）
 * 3.如果读取配置文件失败，提示用户输入课件courseId和lessonIds或修改配置文件
 * 4.获取列表成功，通过API获取资源列表
 * 5.使用oss-SDK逐个下载资源列表，计算下载进度及资源数量百分比
 * 6.将下载完成的资源写入对应沙盒中
 */

/// 配置下载资源路径
var oss_download_resources string = "./oss_download"
var oss_download_resources_debug string = "/Users/lidangkun/Desktop/oss_download"

/// 配置配置文件路径
var oss_download_configFile string = "./oss_download_config.txt"
var oss_download_configFile_debug = "/Users/lidangkun/Desktop/oss_download/oss_download_config.txt"

/// 下载错误信息
var oss_download_errorLog = "./oss_error.txt"
var oss_download_errorLog_debug = "/Users/lidangkun/Desktop/oss_download/oss_error.txt"

var oss_download_debug bool = true

var oss_download_endPoint string
var oss_download_accessKeyId string
var oss_download_accessKeySecret string
var oss_download_bucket string
var oss_download_account string
var oss_download_password string

func main() {

	// 确定资源存放路径和配置文件路径
	var resourcePath string
	var configPath string
	if oss_download_debug {
		resourcePath = oss_download_resources_debug
		configPath = oss_download_configFile_debug
	}else {
		resourcePath = oss_download_resources
		configPath = oss_download_configFile
	}

	// 创建资源存放目录
	isResourceExist := makeSandbox(resourcePath)

	if !isResourceExist {
		fmt.Println("资源存放文件夹创建失败，程序结束")
		return
	}

	// 创建配置文件
	isConfigExist := makeConfigFile(configPath)
	if !isConfigExist {
		fmt.Println("配置文件不存在，程序结束")
		return
	}
	// 读取配置文件
	configMap,err := utils.ReadLines(configPath)

	if err != nil {

		fmt.Println(err)

		return
	}

	var courseIds []string
	var lessonIds []string

	var configObject = make(map[string] string)

	for _,configLine := range configMap {

		var lineValues = strings.Split(configLine,"=")

		if len(lineValues) > 1 {

			configObject[lineValues[0]] = lineValues[1]

		}else {

			configObject[lineValues[0]] = ""

		}

	}

	oss_download_bucket = configObject["Bucket"]
	oss_download_endPoint = configObject["Endpoint"]
	oss_download_accessKeyId = configObject["AccessKeyId"]
	oss_download_accessKeySecret = configObject["AccessKeySecret"]
	oss_download_account = configObject["account"]
	oss_download_password = configObject["password"]

	courseIdString:=configObject["courseIds"]
	lessonIdString:=configObject["lessonIds"]

	courseIds = strings.Split(courseIdString,",")
	lessonIds = strings.Split(lessonIdString,",")

	fmt.Println(configObject)

	if len(courseIds) <= 0 {

		fmt.Println("没有配置courseId")

		return
	}

	fmt.Println("配置信息读取成功",courseIds,lessonIds)

	// 开始登入服务器获取权限
	token := getToken()

	if token == "" {
		return
	}

	// 已正确获取服务器权限，开始访问列表

	for _,courseId := range courseIds{

		downloadCourse(resourcePath, token, courseId, lessonIds)

	}
}

/// 下载一个课程
/**
 * @param token 服务器登入权限
 * @param courseId 需要下载的课程Id
 * @param lessonIds 指定的lesson了列表
 */
func downloadCourse(resourcePath string,token string ,courseId string, lessonIds []string){

	data := make(map[string]interface{})
	data["cid"] = courseId
	if len(lessonIds) > 0 {

		if len(lessonIds) > 1 {
			data["filter"] = lessonIds
		}else {
			lessonId := lessonIds[0]

			if lessonId != "" {
				data["filter"] = lessonIds
			}
		}

	}
	_, rows := utils.ListWithMedias(model.Request{token, "cine.web", data})

	var coursePath = resourcePath + "/" + courseId

	for _,chapter := range rows {

		var chapterPath = coursePath + "/" + chapter.Name

		for _,lesson := range chapter.Children {

			var lessonPath = chapterPath + "/" + lesson.Name

			for _,media := range lesson.Medias {

				var mediaPaths []string

				if strings.Contains(media.Url,"/f/") {

					mediaPaths = strings.Split(media.Url,"/f/")

				}else if strings.Contains(media.Url,"/kj/") {

					mediaPaths = strings.Split(media.Url,"/kj/")

				}else {

					continue
				}

				mediaPath := mediaPaths[len(mediaPaths) - 1]

				mediaComponts := strings.Split(mediaPath,"/")

				mediaName := mediaComponts[len(mediaComponts)-1]

				mediaComponts = strings.Split(mediaName,".")

				mediaName = mediaComponts[0]

				localMedia := lessonPath + "/" + mediaName

				makeSandbox(localMedia)

				// 下载media
				_ = downloadOssResource(localMedia,courseId,mediaPath,false)

				// 下载照片

				for _,image := range media.Images {

					var imageUrl = image.Url

					var imagePaths []string
					var isImage bool
					var pre string

					if strings.Contains(imageUrl,"/f/") {

						imagePaths = strings.Split(imageUrl,"/f/")

						if strings.Contains(imageUrl,"www.bstcine.com") {
							isImage = false
						}else {
							isImage = true
						}

					}else if strings.Contains(imageUrl,"/kj/") {

						imagePaths = strings.Split(imageUrl,"/kj/")
						isImage = false

					}else if strings.Contains(imageUrl,"img/") {

						imagePaths = strings.Split(image.Url,"img/")
						isImage = true
						pre = "img/"

					}else {

						continue

					}

					imageParh := pre + imagePaths[len(imagePaths)-1]

					_ = downloadOssResource(localMedia,courseId,imageParh,isImage)

				}

			}

		}
	}
}

/// 登入服务器，获取token
/**
 * @param
 * @return token
 */
func getToken() string {

	for i := 1; i <= 5; i++ {

		var err error

		if oss_download_account == "" {

			// 获取输入账户名
			fmt.Print("请输入用户名：")

			oss_download_account, err = clientInput('\n')

			if err != nil {

				fmt.Println("标准输入流错误，程序结束")

				return ""
			}
		}

		if oss_download_password == "" {

			fmt.Print("请输入密码：")
			oss_download_password, err = clientInput('\n')

			if err != nil {

				fmt.Println("标准输入流错误，程序结束")

				return ""
			}
		}

		data := make(map[string]interface{})
		data["phone"] = oss_download_account

		data["password"] = oss_download_password

		// 登入服务器
		_, token := utils.Signin(model.Request{"", "cine.web", data})

		if len(token) <= 0 || token == "" {

			if i == 5 {
				fmt.Println("您连续输入账户名和密码超出五次，程序结束！\n请认真确定后再重启本程序")
			} else {
				fmt.Println("用户名或密码错误，请重新输入")

				oss_download_account = ""
				oss_download_password = ""
			}

		} else {

			return token
		}

	}

	return ""
}

/// 读取标准用户输入流
/**
 * @param endbyte 字符串结束符 char类型 如 '\n', '\t',' '等
 * @return string 除掉结束字符的输入字符串
 * @return error 标准输入流报错
 */
func clientInput(endbyte byte) (string, error) {

	clientReader := bufio.NewReader(os.Stdin)

	input,err := clientReader.ReadString(endbyte)

	endData := []byte{endbyte,}
	input = strings.Replace(input,string(endData[:]),"",-1)

	return input,err
}

/// 处理配置文件
/**
 * @param path 配置文件路径
 * @return bool 配置文件是否存在
 */
func makeConfigFile(path string) bool {
	_,err := os.Stat(path)

	if err == nil {
		return true
	}

	fileHandle, err := os.Create(path)

	if err != nil {

		return false
	}

	defer  fileHandle.Close()

	fmt.Println("下载配置文件不存在，是否立即配置？yes/no(以 enter 键结束)")

	addConfig,err := clientInput('\n')

	if addConfig != "yes" {

		fileWriter := bufio.NewWriter(fileHandle)

		fileWriter.WriteString("courseIds:\nlessonIds:")

		err = fileWriter.Flush()

		fmt.Println("您选择暂不配置下载文件，您可以稍候自行在 oss_download_config.txt文件中输入相关信息，程序已结束。")

		return false
	}

	// 准备courseIds,lessonIds

	var selectString string

	var i int = 0

	for {

		i++

		fmt.Println("请选择下载方式：1 -> course, 2 -> lesson，以 enter 键结束")
		selectString,err = clientInput('\n')

		if err != nil {
			fmt.Println("标准输入流出错，程序结束")
			return  false
		}

		if selectString == "1" {

			fmt.Print("请输入所有需要执行下载的courseId，以英文下的\",\"分隔，以 enter 键结束：")

			break

		}else if selectString == "2" {

			fmt.Print("请输入所有需要执行下载的lessonId，以英文下的\",\"分隔，以 enter 键结束：")

			break

		}else {

			if i == 5 {
				fmt.Println("您已经连续五次选择下载类型错误，程序结束")
				return false
			}
			fmt.Println("类型选择错误，请重新选择")

		}

	}

	ids,err := clientInput('\n')

	if err != nil {
		fmt.Println("标准输入流出错，程序结束")
		return  false
	}

	if selectString == "1" {

		ids = "courseIds:"+ids+"\nlessonIds:"

	}else {

		ids = "courseIds:"+"\nlessonIds:"+ids

	}

	fileWriter := bufio.NewWriter(fileHandle)

	fileWriter.WriteString(ids)

	err = fileWriter.Flush()

	if err != nil {

		fmt.Println("配置写入失败")

		return  false
	}

	return true
}

/// 创建资源保存文件夹路径
/**
 * @param path 资源文件路径
 * @return 是否创建完成
 */
func makeSandbox(path string) bool {

	_,err := os.Stat(path);
	if err == nil {
		return true
	}

	err = os.MkdirAll(path,0711)

	if err == nil {
		return true
	}

	return false
}


/// 下载oss上的资源
/**
 * 从oss下载资源
 * @param courseId 待下载资源对应的课件id
 * @param path 资源相对url路径（"/f/"之后的路径）
 * @param isImage 是否是存储在图片库的资源（用来判断是否统一修改扩展名为".jpg"）
 * @return 资源下载结果
 */
func downloadOssResource(saveDirect string, courseId string, path string, isImage bool) bool {

	client,err := oss.New(oss_download_endPoint,oss_download_accessKeyId,oss_download_accessKeySecret)

	if err != nil {
		fmt.Println(err)
		fmt.Println("创建客户端对象失败")
		return false
	}

	bucket,err := client.Bucket(oss_download_bucket)

	if err != nil {
		fmt.Println(err)
		fmt.Println("创建bucket失败")
		return false
	}

	if strings.Contains(path,"?") {
		pathComponts := strings.Split(path,"?")
		path = pathComponts[0]
	}

	var objectKey string

	if strings.Contains(path,"img/") || strings.Contains(path,"kj/") {

		objectKey = path

	}else {

		if isImage {
			paths := strings.Split(path,".")
			objectKey = "img/" + courseId + "/" + paths[0] + ".jpg"
		}else {
			objectKey = "kj/" + path
		}

	}

	fileName := strings.Split(objectKey,"/")

	var savePath = saveDirect + "/" + fileName[len(fileName)-1]

	for i := 0; i < 3;i++  {

		err = bucket.DownloadFile(objectKey,savePath,100*1024,oss.Routines(3),oss.Checkpoint(true,""))

		if err == nil {

			fmt.Println("下载成功", objectKey)
			return true

		}else {

			if i == 2 {

				fmt.Println(err)
			}

		}
	}

	// 将获取失败的文件写入日志中
	fmt.Println("获取资源失败",objectKey)

	return false
}

