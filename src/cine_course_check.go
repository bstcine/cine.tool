package main

import (
	"./utils"
	"./conf"
	"./model"
	"fmt"
	"strings"
	"strconv"
	"os"
	//"net/http"
	"path/filepath"
)

var checkConfigFile_debug = "/Users/lidangkun/Desktop/oss_checkConfig"
var checkConfig_debug = "/Users/lidangkun/Desktop/oss_checkConfig/cine_course_check.cfg"
var checkResourceLog_bug = "/Users/lidangkun/Desktop/oss_checkConfig/resourceLog.txt"

var checkWorkDir string
var checkConfig  string
var checkResourceLog string

var checkCount int

var rounterCount = 6

type config struct {

}

type checkResource struct {

	lessonId    string
	lessonName  string
	objectKey   string
	mediaType   int       // 多媒体资源类型，3代表media,2代表加水印图，1代表原图
	mediaSeq    int       // medis序号
	checkStatus int       // 原始检测状态  2表示只没有原图，其余全部检查
}

type checkResult struct {

	lessonId    string
	lessonName  string
	objectKey   string
	mediaType   int       // 多媒体资源类型，3代表media,2代表加水印图，1代表原图
	status      bool      // 检查结果类型，
	mediaSeq    int       // medis序号
}

func main() {

	// 获取配置信息
	if conf.IsDebug {
		checkWorkDir = checkConfigFile_debug
		checkConfig = checkConfig_debug
		checkResourceLog = checkResourceLog_bug
		utils.CreatDirectory(checkWorkDir)
	}else {
		dir,err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			return
		}

		dir = strings.Replace(dir,"\\","/",-1)

		fmt.Println(dir)
		checkWorkDir = dir + conf.Course_checkWorkDir
		checkConfig = dir + conf.Course_checkConfig
		checkResourceLog = dir + conf.Course_check_log
	}

	os.Remove(checkResourceLog)

	// 读取配置服务
    checkAccount,checkPassword,checkCourseIds := getConfigs()

	// 登录服务器获取权限
	checkToken := utils.GetAdminPermission(checkAccount,checkPassword)

	courses := getCourses(checkToken,checkCourseIds)

	for _,courseModel := range courses {

		lessons := getLessons(checkToken,courseModel.Id)

		if len(lessons) == 0 {
			fmt.Println(courseModel.Id,courseModel.Name,"已经检查完毕")
			continue
		}

		// 获取待检查资源数量
		var resources = creatCheckResources(courseModel,lessons)
		var resourceCount = len(resources)

		// 创建一个工作管道，和结果管道
		var jobs  = make(chan checkResource,resourceCount)
		var resultChan  = make(chan checkResult,resourceCount)

		fmt.Println(resourceCount)

		// 开通协程工作
		for i := 1;i <= rounterCount;i++  {
			go startJos(jobs,resultChan)
		}

		// 将工作内容加入管道，开始工作
		go addJobsWithResources(jobs,resources)

		// 监听工作结果
		updateData := dealResults(resultChan,courseModel,resourceCount)

		fmt.Println("获取更新结果：",courseModel.Name,"\n",updateData)

		//// 更新该状态
		//_,status := utils.UpdateLessonCheckStatus(model.Request{checkToken,"cine.web",updateData})
		//
		//if !status {
		//	fmt.Println("更新失败",courseModel.Id,courseModel.Name)
		//}else {
		//	fmt.Println("更新成功",courseModel.Id,courseModel.Name)
		//}

	}

	// 结束

	fmt.Println("工作结束！")
}

// 生成工作资源
func creatCheckResources(courseModel model.Course,lessons []model.CheckLesson) []checkResource {

	var resources []checkResource

	for _,lessonModel := range lessons {
		resources = addLessons(resources,courseModel,lessonModel)
	}

	return  resources
}
// 将lesson资源加入队列中
func addLessons(checkResources []checkResource,courseModel model.Course,lessonModel model.CheckLesson) (resources []checkResource) {

	for _,mediaModel := range lessonModel.Medias {

		if mediaModel.Url == "" {
			continue
		}

		// 当mediamodel的检查状态为2时，表示只没有原图，此时不再检查音频
		if mediaModel.CheckStatus != 2 {

			// 将mediaModel生成资源
			mediaResource := checkResource{
				lessonId:lessonModel.Id,
				lessonName:lessonModel.Name,
				objectKey:("kj/"+mediaModel.Url),
				mediaType:3,
				mediaSeq:mediaModel.Seq,
			}

			checkResources = append(checkResources,mediaResource)
		}

		if len(mediaModel.Images) == 0 {
			continue
		}

		// 将图片资源加入工作队列
		for _,imageModel := range mediaModel.Images  {

			if imageModel.Url == "" {
				continue
			}

			imagePath := imageModel.Url
			// 切断后缀
			imagePathArr := strings.Split(imagePath,".")
			imagePath = imagePathArr[0] + ".jpg@!style_ori"

			if mediaModel.CheckStatus != 2 {
				// 将线上使用的图加入工作队列
				usePath := "img/" + courseModel.Id + "/" + imagePath

				imageUseResource := checkResource{
					lessonId:lessonModel.Id,
					lessonName:lessonModel.Name,
					objectKey:usePath,
					mediaType:2,
					mediaSeq:mediaModel.Seq,
				}
				checkResources = append(checkResources,imageUseResource)
			}

			// 将原图加入工作队列
			originPath := "kj/" + imagePath
			imageOriginResource := checkResource{
				lessonId:lessonModel.Id,
				lessonName:lessonModel.Name,
				objectKey:originPath,
				mediaType:1,
				mediaSeq:mediaModel.Seq,
			}
			checkResources = append(checkResources,imageOriginResource)

		}

	}

	return checkResources
}
// 将检车资源加入工作队列
func addJobsWithResources(jobs chan checkResource, resources []checkResource){

	for _,resource := range resources {

		jobs <- resource

	}

	close(jobs)

}

// 开启工作组开始工作，将工作的结果加入结果队列
func startJos(jobs chan checkResource,results chan checkResult) {

	for resource := range jobs {

		isSave := utils.CheckResourceSaveStatus(resource.objectKey)

		result := checkResult{
			lessonId:resource.lessonId,
			lessonName:resource.lessonName,
			status:isSave,
			mediaType:resource.mediaType,
			mediaSeq:resource.mediaSeq,
			objectKey:resource.objectKey,
		}

		results <- result
	}

}

// 监听工作工作结果星道
func dealResults(results chan checkResult,courseModel model.Course,resourceCount int) (map[string]interface{}){

	// 以LessonId为key,存放lesson字典
	var lessonDicts = make(map[string]map[string]interface{})

	for i:= 0; i < resourceCount;i++ {

		result := <-results

		lessonDict := lessonDicts[result.lessonId]

		if lessonDict == nil {
			lessonDict = make(map[string]interface{})
			lessonDict["lesson_id"] = result.lessonId
			lessonDict["check_status"] = 1
			lessonDicts[result.lessonId] = lessonDict
		}

		if result.status {
			fmt.Println(result,i,"/",resourceCount)
			continue
		}

		lessonDict["check_status"] = 0

		var mediaDicts []map[string]interface{}

		if lessonDict["medias"] != nil {
			mediaDicts = lessonDict["medias"].([]map[string]interface{})
		}

		mediaDict := make(map[string]interface{})

		mediaDict["seq"] = result.mediaSeq

		switch result.mediaType {

		case 1:

			mediaDict["check_status"] = 2

			break

		case 2:

			mediaDict["check_status"] = 3

			break

		case 3:

			mediaDict["check_status"] = 3

			break

		default:

			fmt.Println("类型标识码错误")

			break

		}

		mediaDicts = append(mediaDicts,mediaDict)

		lessonDict["medias"] = mediaDicts

		fmt.Println(result,i,"/",resourceCount)

		errorMessage := "#" + courseModel.Id + "-" + courseModel.Name + "-" + result.lessonId + "-" + result.lessonName + " : "

		fmt.Println(mediaDict)

		errorLesson := getErrorString(courseModel.Id,result.lessonId,result.lessonName,strconv.Itoa(result.mediaSeq),strconv.Itoa(mediaDict["check_status"].(int)),result.objectKey)

		utils.AppendStringToFile(checkResourceLog,errorMessage + "\n" + errorLesson + "\n")

	}

	fmt.Println("course结束")

	var lessonArray []map[string]interface{}

	// 打包检查结果资源
	for _,value := range lessonDicts{
		lessonArray = append(lessonArray,value)
	}

	updateData := make(map[string]interface{})

	updateData["lesson_ids"] = lessonArray

	close(results)

	return updateData
}

func getErrorString(courseId string,lessonId string,lessonName string,seq string,errorStatus string,objectKey string) string {

	errorString := "seq=" + seq + "," + " error_status=" + errorStatus + "," + " objectKey=" + objectKey

	errorString = errorString + "\n" + "<a href=\"" + "http://www.bstcine.com/learn/lesson/"+ courseId + "?"

	errorString = errorString + "content_id=" + lessonId + "&seq=" + seq + "&error_type=" + errorStatus + "\">" + lessonName + "</a>"

	errorString = errorString + "\n" + "-----------------------------------------" + "<br>"

	return  errorString

}

func checkMediaStatus(mediaPath string) bool {

	mediaPath = "kj/" + mediaPath

	return utils.CheckResourceSaveStatus(mediaPath)
}

func checkImageStatus(courseId string,imagePath string) int {

	if strings.Contains(imagePath,"png") {
		fmt.Println(imagePath)
	}

	// 切断后缀
	imagePathArr := strings.Split(imagePath,".")

	imagePath = imagePathArr[0] + ".jpg@!style_ori"

	// 判断原始图片是否存在
	originPath := "kj/" + imagePath

	originStatus := utils.CheckResourceSaveStatus(originPath)

	usePath := "img/" + courseId + "/" + imagePath

	useStatus := utils.CheckResourceSaveStatus(usePath)

	if !useStatus {
		 return 3
	}

	if !originStatus {
		return 2
	}

	return  1
}

//func checkResourceSaveStatus(objectKey string) bool {
//
//	// 检查图片状态需要检查两次，一次为原图，一次为水印图，两次成功才返回 1
//	requestPath := "http://oss.bstcine.com/" + objectKey
//
//	checkCount = 3
//
//	for checkCount >= 0  {
//
//		checkCount -= 1
//
//		resp,err := http.Head(requestPath)
//
//		if err == nil {
//
//			if resp.StatusCode == 200 {
//
//				length := resp.Header["Content-Length"]
//
//				if len(length) > 0 {
//
//					currentLength,err := strconv.ParseInt(length[0],0,64)
//
//					if err != nil {
//						return  false
//					}
//
//					return (currentLength > 10240)
//				}
//
//			}
//
//		}
//
//	}
//
//	return  false
//}

//***************************************************************
//*********************                 *************************
//*********************     网络模块     *************************
//*********************                 *************************
//***************************************************************
func getCourses(token string,courseIds []string) []model.Course {

	var data = make(map[string]interface{})

	data["course_ids"] = courseIds

	_,courses := utils.ListWithCheckCourses(model.Request{token,"cine.web",data})

	return courses
}

func getLessons(token string, courseId string) []model.CheckLesson {

	data := make(map[string]interface{})
	data["cid"] = courseId
	_,lessons := utils.ListWithCheckMedias(model.Request{token,"cine.web",data})

	return lessons
}

func getConfigs() (account string,password string,courseIds []string) {

	armConfigs := utils.GetConfArgs(checkConfig)

	fmt.Println(armConfigs)

	account = armConfigs["Account"]
	password = armConfigs["Password"]
	cids := armConfigs["Course_ids"]

	if cids == "" {
		return account,password,nil
	}

	// 清理配置中的空格键 " "
	account = strings.Replace(account," ","",-1)
	password = strings.Replace(password," ","",-1)
	cids = strings.Replace(cids," ","",-1)

	if !strings.Contains(cids,",") {
		courseIds = append(courseIds, cids)
	}else {
		courseIds = strings.Split(cids,",")
	}

	return account,password,courseIds
}