package main

import (
	"./utils"
	"./model"
	"fmt"
	"os"
	"bufio"
	"strings"
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

var oss_download_downloadStatus bool = true

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
	isResourceExist := utils.CreatDirectory(resourcePath)

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

	courseIds,lessonIds := readConfig(configPath)

	// 判断配置lesson数组和course数组是否对应
	if lessonIds != nil && len(courseIds) > 0 && len(courseIds) != len(lessonIds) {

		fmt.Println("配置course 和lesson 数量不对应，如有部分课程全部下载，请配置\"[]\"")

		return
	}

	// 检查是否存在错误日志表
	errorStatus := candownloadErrorLog()

	if errorStatus {
		return
	}

	// 开始登入服务器获取权限
	token := getToken()

	if token == "" {
		return
	}

	// 已正确获取服务器权限，开始访问列表

	for i := 0; i < len(courseIds);i++  {

		courseId := courseIds[i]
		var downloadLessonIds []string

		if lessonIds != nil {
			downloadLessonIds = lessonIds[0]
		}

		downloadCourse(resourcePath,token,courseId,downloadLessonIds)

	}

	if oss_download_downloadStatus {
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

	downloadRows(coursePath, courseId, rows)
}

/// 下载课程列表
/**
 * @param coursePath 课程目录，存放chapter资源
 * @param rows 待下载的chapter资源列表
 */
func downloadRows(coursePath string, courseId string, rows []model.Chapter) {

	for _,chapter := range rows {

		var chapterPath = coursePath + "/" + chapter.Name

		for _,lesson := range chapter.Children {

			var lessonPath = chapterPath + "/" + lesson.Name

			downloadLesson(lessonPath,courseId,lesson)

		}
	}

}

/// 下载课程
/**
 * @param lessonPath 课程目录，存放lesson资源
 * @param courseId lesson所属的课件id
 * @param lesson 待下载的课件
 */
func downloadLesson(lessonPath string, courseId string, lesson model.Lesson) {

	utils.CreatDirectory(lessonPath)

	for i := 0; i < len(lesson.Medias); i++ {

		media := lesson.Medias[i]

		var mediaPaths []string

		if strings.Contains(media.Url,"/f/") {

			mediaPaths = strings.Split(media.Url,"/f/")

		}else if strings.Contains(media.Url,"/kj/") {

			mediaPaths = strings.Split(media.Url,"/kj/")

		}else {

			continue
		}

		mediaPath := mediaPaths[len(mediaPaths) - 1]

		localMedia := lessonPath + "/" + utils.ChangeInt(i)

		// 下载media
		_ = downloadOssResource(localMedia,courseId,mediaPath,false)

		// 下载照片

		for j := 0 ; j < len(media.Images) ; j++ {

			image := media.Images[j]

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

			localImage := localMedia + "_" + utils.ChangeInt(j)

			lessonStatus := downloadOssResource(localImage,courseId,imageParh,isImage)

			if !lessonStatus {
				oss_download_downloadStatus = false
			}

		}

	}

}

/// 下载oss上的资源
/**
 * 从oss下载资源
 * @param courseId 待下载资源对应的课件id
 * @param path 资源相对url路径（"/f/"之后的路径）
 * @param isImage 是否是存储在图片库的资源（用来判断是否统一修改扩展名为".jpg"）
 * @return 资源下载结果
 */
func downloadOssResource(savePath string, courseId string, path string, isImage bool) bool {

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

	fileName := strings.Split(objectKey,".")

	savePath = savePath + "." + fileName[len(fileName)-1]

	downloadStatus := utils.DownloadOssResource(oss_download_endPoint,oss_download_accessKeyId,oss_download_accessKeySecret,oss_download_bucket,savePath,objectKey)

	if !downloadStatus {
		// 将错误信息写入报错日志
		fmt.Println("下载失败：",savePath,objectKey)
		addErrorObject(objectKey,savePath)
	}

	return downloadStatus
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

	for _,errorObject := range errorObjects {

		if !strings.Contains(errorObject,",") {
			continue
		}

		objectComponent := strings.Split(errorObject,",")

		downloadErrorObject(objectComponent[1],objectComponent[0])
	}

	if oss_download_downloadStatus {
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

	fmt.Println(objectKey)

	downloadStatus := utils.DownloadOssResource(oss_download_endPoint,oss_download_accessKeyId,oss_download_accessKeySecret,oss_download_bucket,savePath,objectKey)

	if downloadStatus {

		// 下载完成，移出错误日志
		fmt.Println("下载完成",objectKey)
		removeErrorObject(objectKey)

	}else {

		oss_download_downloadStatus = false

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

			oss_download_account = utils.ClientInputWithMessage("请输入用户名：",'\n')

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
func readConfig(configPath string) ([]string, [][]string) {

	// 读取配置文件
	configMap,err := utils.ReadLines(configPath)

	if err != nil {

		fmt.Println(err)

		return nil, nil
	}

	var courseIds []string
	var lessonIds [][]string

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

	if len(courseIds) <= 0 {

		fmt.Println("没有配置courseId")

		return nil, nil
	}

	if lessonIdString == "" {
		return courseIds,nil
	}

	var lessonArrs []string

	lessonArrs = strings.Split(lessonIdString,",")

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

	fmt.Println(configObject)

	fmt.Println("配置信息读取成功",courseIds,lessonIds)

	return courseIds, lessonIds
}

/// 创建配置文件
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

	addConfig := utils.ClientInputWithMessage("下载配置文件不存在，是否立即配置？yes/no(以 enter 键结束)",'\n')

	if addConfig != "yes" {

		fileWriter := bufio.NewWriter(fileHandle)

		fileWriter.WriteString("courseIds=\nlessonIds=\nBucket=\nEndpoint=\nAccessKeyId=\nAccessKeySecret=\naccount=\npassword=")

		err = fileWriter.Flush()

		fmt.Println("您选择暂不配置下载文，配置模板已自动生成,您可以稍候自行在 oss_download_config.txt文件中输入相关信息，程序已结束。")

		return false
	}
	// 准备courseIds,lessonIds

	var courseIds string
	var lessonIds string

	var i int = 0

	for {

		cid := utils.ClientInputWithMessage("请输入待下载课程 id(courseId),以 enter 键结束,不能包含\",\"等特殊字符：",'\n')

		if strings.Contains(cid,",") {
			fmt.Println("输入错误，不能包含\",\"字符, 请重新输入")
			continue
		}

		if courseIds == "" {
			courseIds = cid
		}else {
			courseIds = courseIds + "," + cid
		}

		lid := utils.ClientInputWithMessage("请为课程指定需要下载的lesson，每个lesson用\",\"隔开，以 enter 键结束，如果需要下载全部lesson，请直接点击 enter 键",'\n')

		lid = "[" + lid + "]"

		if lessonIds == "" {
			lessonIds = lid
		}else {
			lessonIds = lessonIds + "," + lid
		}

		i++

		fmt.Printf("您已经成功配置了%d个课程，是否继续添加待下载课程 y/n ",i)

		addStatus,err := utils.ClientInput('\n')

		if err != nil {
			fmt.Println("标准输入流出错，程序结束")
			return  false
		}

		if addStatus == "Y" || addStatus == "y" || addStatus == "YES" || addStatus == "yes" || addStatus == "Yes" {
			continue
		}

		break
	}

	fmt.Println("开始配置oss参数：")

	bucket := utils.ClientInputWithMessage("请输入Bucket: ",'\n')

	endpoint := utils.ClientInputWithMessage("请输入Endpoint: ",'\n')

	accessKeyId := utils.ClientInputWithMessage("请输入AccessKeyId: ",'\n')

	accessKeySecret := utils.ClientInputWithMessage("请输入AccessKeySecret: ",'\n')

	ids := "courseIds="+courseIds+"\nlessonIds="+lessonIds+"\nBucket="+bucket+"\nEndpoint="+endpoint+"\nAccessKeyId="+accessKeyId+"\nAccessKeySecret="+accessKeySecret

	fileWriter := bufio.NewWriter(fileHandle)

	fileWriter.WriteString(ids)

	err = fileWriter.Flush()

	if err != nil {

		fmt.Println("配置写入失败")

		return  false
	}

	fmt.Println("配置文件创建成功")

	return true
}
