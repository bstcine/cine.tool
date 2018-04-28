package main

import (
	"./conf"
	"./utils"
	"./model"
	"fmt"
	"os"
	"bufio"
	"strings"
	"path/filepath"
	"strconv"
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
var oss_download_resources string = conf.Course_downloadWorkDir
var oss_download_resources_debug string = "/Users/lidangkun/Desktop/oss_download"

/// 配置配置文件路径
var oss_download_configFile string = conf.Course_download_Config
var oss_download_configFile_debug = "/Users/lidangkun/Desktop/oss_download/oss_download_config.cfg"

/// 下载错误信息
var oss_download_errorLog = conf.Course_download_errorLog
var oss_download_errorLog_debug = "/Users/lidangkun/Desktop/oss_download/oss_error.txt"
var errorLogPath string

var oss_download_debug bool = conf.IsDebug

var oss_download_endPoint = "oss-cn-shanghai.aliyuncs.com"
var oss_download_accessKeyId string
var oss_download_accessKeySecret string
var oss_download_bucket = "static-bstcine"
var oss_download_account string
var oss_download_password string

var oss_download_av_media = true
var oss_download_image = true
var imageStyle = ""

func main() {

	// 确定资源存放路径和配置文件路径
	var resourcePath string
	var configPath string
	if oss_download_debug {
		resourcePath = oss_download_resources_debug
		configPath = oss_download_configFile_debug
		errorLogPath = oss_download_errorLog_debug
	}else {
		dir,err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			return
		}

		dir = strings.Replace(dir,"\\","/",-1)

		resourcePath = dir + oss_download_resources
		configPath = dir + oss_download_configFile
		errorLogPath = dir + oss_download_errorLog
	}

	// 创建资源存放目录
	isResourceExist := utils.CreatDirectory(resourcePath)

	if !isResourceExist {
		fmt.Println("资源存放文件夹创建失败，程序结束")
		return
	}

	// 创建配置文件
	courseIds,courseAliasNames,lessonIds := readConfig(configPath)

	if len(courseIds) == 0 || len(courseAliasNames) == 0 || len(courseIds) != len(courseAliasNames) {

		fmt.Println("课程 Id 和别名设置出错",courseIds,courseAliasNames)

		return
	}

	// 判断配置lesson数组和course数组是否对应
	if lessonIds != nil && len(courseIds) > 0 && len(courseIds) != len(lessonIds) {

		fmt.Println("配置course 和lesson 数量不对应，如有部分课程全部下载，请配置\"[]\"")

		return
	}

	//// 检查是否存在错误日志表
	errorStatus := candownloadErrorLog()

	if errorStatus {
		return
	}

	// 开始登入服务器获取权限
	token := utils.GetAdminPermission(oss_download_account,oss_download_password)

	if token == "" {
		return
	}

	// 开始下载课件列表
	downloadStatus := downloadCourseList(resourcePath,token,courseIds,courseAliasNames,lessonIds)

	if downloadStatus {
		fmt.Println("本次课件资源全部下载完成")
	}else{
		fmt.Println("本次课件部分资源下载失败，请重新执行本程序，续传失败的文件")
	}

}

//*****************************************************************
//*****************************************************************
//************************** oss下载模块 ***************************
//*****************************************************************
//*****************************************************************

/// 获取下载课件对象列表
func downloadCourseList(resourcePath string,token string,courseIds []string,courseAliasNames []string,lessonIds [][]string) bool {

	_,courseModels := utils.ListWithDownloadCourses(token,courseIds)

	if len(courseModels) != len(courseAliasNames) {

		fmt.Println("网端可下载课件数量与别名不符")
		fmt.Println("可下载课件如下:\n",courseModels)
		fmt.Println("别名如下:\n",courseAliasNames)

		return false
	}

	var realCourseModels []model.Course

	for _,courseId := range courseIds {

		for _,courseModel := range courseModels {

			if courseId == courseModel.Id {
				realCourseModels = append(realCourseModels, courseModel)
				break
			}
		}
	}

	courseModels = realCourseModels

	var downloadStatus bool = true

	// 已正确获取服务器权限，开始访问列表

	for i := 0; i < len(courseModels);i++  {

		courseModel := courseModels[i]



		alias := courseAliasNames[i]

		var downloadLessonIds []string

		if len(lessonIds) > 0 {
			downloadLessonIds = lessonIds[i]
		}

		courseDownloadStatus := downloadCourse(resourcePath,token,courseModel.Id,alias,downloadLessonIds)

		if courseDownloadStatus == false {
			downloadStatus = false
		}
	}

	return downloadStatus
}

/// 下载一个课程
/**
 * @param token 服务器登入权限
 * @param courseId 需要下载的课程Id
 * @param lessonIds 指定的lesson了列表
 */
func downloadCourse(resourcePath string,token string ,courseId string,alias string, lessonIds []string) bool {

	_, lessons := utils.ListWithDownloadMedias(token,courseId,lessonIds)

	var coursePath = resourcePath + "/" + alias

	var downloadStatus bool = true

	for _,lessonModel := range lessons {

		chapterAlias := strings.Replace(lessonModel.ChapterName,"/","#",-1)

		if strings.Contains(lessonModel.ChapterName,"/") {
			lessonModel.ChapterName = strings.Replace(lessonModel.ChapterName,"/","_",-1)
		}

		if strings.Contains(lessonModel.Name,"/") {
			lessonModel.Name = strings.Replace(lessonModel.Name,"/","_",-1)
		}

		// 生成下载路径
		downloadPath := coursePath + "/" + lessonModel.ChapterName + "/" + "ls_" + lessonModel.Name

		utils.CreatDirectory(downloadPath)

		lessonDownloadStatus := downloadLesson(downloadPath,courseId,alias,chapterAlias,lessonModel)

		if lessonDownloadStatus == false {
			downloadStatus = false
		}
	}

	return downloadStatus
}

func downloadLesson(downloadPath string,courseId string,alias string,chapterAlias string,lessonModel model.CheckLesson) bool {

	for index,mediaModel := range lessonModel.Medias {

		downloadAVMedia(courseId,alias,chapterAlias,downloadPath,mediaModel,index)

		for _,imageModel := range mediaModel.Images {

			if len(imageModel.Time) == 1 {
				imageModel.Time = "00"+imageModel.Time
			}else if len(imageModel.Time) == 2 {
				imageModel.Time = "0"+imageModel.Time
			}

			downloadImage(courseId,alias,chapterAlias,downloadPath,imageModel,index)
		}

	}

	return true
}


func downloadAVMedia(courseId string,alias string,chapterAlias string,lessonPath string,media model.CheckMedia, seq int) bool {

	if !oss_download_av_media {
		return  true
	}

	if media.Url == "" {
		return  true
	}

	urlComponents := strings.Split(media.Url,".")

	extension := "." + urlComponents[1]

	// + alias + "_" + chapterAlias + "_"
	savePath := lessonPath + "/" + downloadFileName(seq) + extension

	objectKey := "kj/" + media.Url

	fmt.Println("开始下载资源:",objectKey)

	downloadStatus := downloadOssResource(savePath,objectKey)

	if !downloadStatus {
		fmt.Println("下载失败:",objectKey)
		// 生成错误信息
		errorMes := savePath + "," + objectKey + "\n"

		utils.AppendStringToFile(errorLogPath,errorMes)
	}else {
		fmt.Println("下载完毕:",objectKey)
	}

	return downloadStatus
}

func downloadImage(courseId string,alias string,chapterAlias string,lessonPath string,image model.Image, mediaSeq int) bool {

	if !oss_download_image {
		return true
	}

	if image.Url == "" {
		return  true
	}

	// 将image的url扩展名更换为.jpg
	imageUrlComponents := strings.Split(image.Url,".")
	image.Url = imageUrlComponents[0] + ".jpg"

	// + alias + "_" + chapterAlias + "_"
	savePath := lessonPath + "/" + downloadFileName(mediaSeq) + "_" + image.Time + ".jpg"

	if utils.Exists(savePath) {

		fmt.Println("待下载文件已存在:",savePath)

		return  true
	}

	objectKey := "kj/" + image.Url + imageStyle
	url := "http://oss.bstcine.com/" + objectKey

	fmt.Println("开始下载资源:",objectKey)

	if !utils.CheckResourceSaveStatus(objectKey) {
		fmt.Println("资源不存在:",objectKey)
		return false
	}

	downloadStatus := utils.DownloadFile(url,savePath)

	if !downloadStatus {

		fmt.Println("下载失败:",objectKey)
		errorMes := savePath + "," + objectKey + "\n"
		utils.AppendStringToFile(errorLogPath,errorMes)
	}else {
		fmt.Println("下载完毕:",objectKey)
	}

	return downloadStatus
}

func downloadOssResource(savePath string, objectKey string) bool {

	// 判断是否已经下载过
	isHadDownload := utils.Exists(savePath)

	if isHadDownload {
		fmt.Println("资源已下载过",savePath)
		return true
	}

	// 判断资源是否存在
	isResourceHad := utils.CheckResourceSaveStatus(objectKey)

	if !isResourceHad {

		fmt.Println("资源不存在:",objectKey)

		return false
	}

	downloadStatus := utils.DownloadOssResource(oss_download_endPoint,oss_download_accessKeyId,oss_download_accessKeySecret,oss_download_bucket,savePath,objectKey)

	return downloadStatus

}

func downloadFileName(number int) string {

	s := strconv.Itoa(number)

	if number >= 100 {
		return s
	}

	if number >= 10 {
		return "0"+s
	}

	return "00"+s
}

//*****************************************************************
//*****************************************************************
//*********************** 管理下载错误对象 **************************
//*****************************************************************
//*****************************************************************
/// 下载错误日志中的对象
func candownloadErrorLog() bool {

	var errorPath string

	if oss_download_debug {

		errorPath = oss_download_errorLog_debug

	}else  {

		errorPath = oss_download_errorLog

	}

	_,err := os.Stat(errorPath)

	if err != nil {
		return  false
	}

	// 读取文件中的数据
	errorObjects,err := utils.ReadLines(errorPath)

	if err != nil {
		// 读取失败的处理
		fmt.Println("读取错误日志失败")
		return false
	}

	if len(errorObjects) == 0 {

		// 删除错误日志文件
		os.Remove(errorPath)

		return  false
	}

	var downloadStatus = true

	for _,errorObject := range errorObjects {

		if !strings.Contains(errorObject,",") {
			continue
		}

		objectComponent := strings.Split(errorObject,",")

		downloadErrorStatus := downloadErrorObject(objectComponent[1],objectComponent[0])

		if downloadErrorStatus == false {
			downloadStatus = false
		}
	}

	if downloadStatus {
		// 删除错误日志文件
		os.Remove(errorPath)
	}

	return  true

}
// 下载错误对象
/**
 * @param savePath 保存路径
 * @param objectKey 下载对象
 */
func downloadErrorObject(savePath string,objectKey string) bool {

	downloadStatus := utils.DownloadOssResource(oss_download_endPoint,oss_download_accessKeyId,oss_download_accessKeySecret,oss_download_bucket,savePath,objectKey)

	if downloadStatus {

		// 下载完成，移出错误日志
		fmt.Println("下载完成",objectKey)
		removeErrorObject(objectKey)

	}

	return downloadStatus
}

/// 将下载失败的对象保存到错误日志中
/**
 @param objectKey 下载失败的对象
 @param savePath 保存地址
 */
func addErrorObject(objectKey string, savePath string) {

	var errorPath string

	if oss_download_debug {

		errorPath = oss_download_errorLog_debug

	}else  {

		errorPath = oss_download_errorLog

	}

	// 判断文件是否存在
	_,err := os.Stat(errorPath)

	if err != nil {
		// 创建一个错误日志文件
		_,err = os.Create(errorPath)
		if err == nil {
			fmt.Println("错误日志创建成功")
		}
	}

	// 读取文件中的数据
	errorObjects,err := utils.ReadLines(errorPath)

	if err != nil {
		// 读取失败的处理
		fmt.Println("读取错误日志失败")
		return
	}

	var errorStrings string

	if len(errorObjects) == 0 {

		errorStrings = objectKey+","+savePath

	}else {

		for _,errorObject := range errorObjects {

			if !strings.Contains(errorObject,",") {
				continue
			}

			if strings.Contains(errorObject,objectKey) {
				return
			}

			// 正常数据
			if errorStrings == "" {
				errorStrings = errorObject
			}else {
				errorStrings = errorStrings + "\n" + errorObject
			}
		}

		errorStrings = errorStrings+"\n"+objectKey+","+savePath
	}

	fileHandler,err := os.Create(errorPath)
	// 存储错误信息
	writer := bufio.NewWriter(fileHandler)

	writer.WriteString(errorStrings)

	err = writer.Flush()

}

func removeErrorObject(objectKey string){

	var errorPath string

	if oss_download_debug {

		errorPath = oss_download_errorLog_debug

	}else  {

		errorPath = oss_download_errorLog

	}

	// 判断文件是否存在
	_,err := os.Stat(errorPath)

	if err != nil {

		return
	}

	// 读取文件中的数据
	errorObjects,err := utils.ReadLines(errorPath)

	if err != nil {
		// 读取失败的处理
		fmt.Println("读取错误日志失败")
		return
	}

	if len(errorObjects) == 0 {
		return
	}

	var errorStrings string

	for _,errorObject := range errorObjects {

		if !strings.Contains(errorObject,",") {
			continue
		}

		if strings.Contains(errorObject,objectKey) {

			objectComponent := strings.Split(errorObject,",")

			if objectComponent[0] == objectKey {
				continue
			}

		}

		if errorStrings == "" {
			errorStrings = errorObject
		}else {
			errorStrings = errorStrings + "\n" + errorObject
		}
	}

	fileHandler,err := os.Create(errorPath)
	// 存储错误信息
	writer := bufio.NewWriter(fileHandler)

	writer.WriteString(errorStrings)

	err = writer.Flush()

}

//*****************************************************************
//*****************************************************************
//************************ 获取服务器权限 ***************************
//*****************************************************************
//*****************************************************************
/// 登入服务器，获取token
/**
 * @param
 * @return token
 */
func getToken() string {

	for i := 1; i <= 5; i++ {

		var isConfigAccount bool = true

		if oss_download_account == "" {

			// 获取输入账户名
			isConfigAccount = false

			oss_download_account = utils.ClientInputWithMessage("请输入管理员账户：",'\n')

			if oss_download_account == "" {

				return ""
			}
		}

		if oss_download_password == "" {

			isConfigAccount = false

			oss_download_password = utils.ClientInputWithMessage("请输入密码：",'\n')

			if oss_download_password == "" {

				return ""
			}
		}

		token := utils.GetAdminPermission(oss_download_account,oss_download_password)

		if len(token) <= 0 || token == "" {

			if i == 5 {
				fmt.Println("您连续输入账户名和密码超出五次，程序结束！\n请认真确定后再重启本程序")
			} else {
				fmt.Println("用户名或密码错误，请重新输入")

				oss_download_account = ""
				oss_download_password = ""
			}

		} else {

			fmt.Println(isConfigAccount)
			return token
		}

	}

	return ""
}

//*****************************************************************
//*****************************************************************
//************************* 管理配置信息 ****************************
//*****************************************************************
//*****************************************************************

/// 读取配置文件
func readConfig(configPath string) (courseIds []string, courseAlias []string, lessonIds [][]string) {

	// 读取配置文件
	configMap,err := utils.ReadLines(configPath)

	if err != nil {

		fmt.Println(err)

		return nil, nil,nil
	}

	var configObject = make(map[string] string)

	for _,configLine := range configMap {

		var lineValues = strings.Split(configLine,"=")

		if len(lineValues) > 1 {

			configObject[lineValues[0]] = lineValues[1]

		}else {

			configObject[lineValues[0]] = ""

		}

	}

	oss_download_accessKeyId = configObject["AccessKeyId"]
	oss_download_accessKeySecret = configObject["AccessKeySecret"]
	oss_download_account = configObject["account"]
	oss_download_password = configObject["password"]

	imageStyle = configObject["style"]
	download_image := configObject["download_image"]
	download_av_media := configObject["download_av_media"]

	if download_image == "false" {
		oss_download_image = false
	}

	if imageStyle == "" {
		imageStyle = "@!style_ori"
	}

	if download_av_media == "false" {
		oss_download_av_media = false
	}

	courseIdString:=configObject["courseIds"]
	courseAliasNameString:=configObject["aliasNames"]
	lessonIdString:=configObject["lessonIds"]

	courseIdString = strings.Replace(courseIdString," ","",-1)
	lessonIdString = strings.Replace(lessonIdString," ","",-1)
	courseAliasNameString = strings.Replace(courseAliasNameString," ","",-1)

	courseIds = strings.Split(courseIdString,",")
	courseAlias = strings.Split(courseAliasNameString,",")

	if len(courseIds) <= 0 {

		fmt.Println("没有配置courseId")

		return nil, nil,nil
	}

	if lessonIdString == "" {
		return courseIds,courseAlias,nil
	}

	var lessonArrs []string

	lessonArrs = strings.Split(lessonIdString,"],")

	for _,lessonArr := range lessonArrs {

		lessonArr = strings.Replace(lessonArr,"[","",-1)
		lessonArr = strings.Replace(lessonArr,"]","",-1)

		var lessonArray []string

		if lessonArr == "" {
			lessonArray = make([]string,0)
		}else {
			lessonArray = strings.Split(lessonArr,",")
		}

		lessonIds = append(lessonIds,lessonArray)
	}

	fmt.Println(lessonIds)
	return courseIds,courseAlias,lessonIds
}
