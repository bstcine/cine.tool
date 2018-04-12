package main

import (
	"./utils"
	"./conf"
	"./model"
	"fmt"
	"strings"
	"strconv"
	"os"
	"net/http"
)

var checkConfigFile_debug = "/Users/lidangkun/Desktop/oss_checkConfig"
var checkResourceLog_bug = "/Users/lidangkun/Desktop/oss_checkConfig/resourceLog.txt"

var checkWorkDir string
var checkResourceLog string

var checkCount int

var rounterCount = 6

type checkResource struct {

	lessonId    string
	lessonName  string
	objectKey   string
	mediaType   int       // 多媒体资源类型，3代表media,2代表加水印图，1代表原图
	mediaSeq    int       // medis序号
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
		checkResourceLog = checkResourceLog_bug
	}else {
		checkWorkDir = conf.Course_checkWorkDir
		checkResourceLog = conf.Course_check_log
	}

	utils.CreatDirectory(checkWorkDir)

	os.Remove(checkResourceLog)

	// 登录服务器获取权限
	checkToken := utils.GetToken("","")

	courses := getCourses(checkToken)

	for _,courseModel := range courses {

		lessons := getLessons(checkToken,courseModel.Id)

		// 获取待检查资源数量
		var resourceCount = 0

		for _,lessonModel := range lessons {

			resourceCount += len(lessonModel.Medias)

			for _,mediaModel := range lessonModel.Medias {
				resourceCount += (len(mediaModel.Images)*2)
			}
		}

		// 创建一个工作管道，和结果管道
		var jobs  = make(chan checkResource,resourceCount)
		var resultChan  = make(chan checkResult,resourceCount)

		fmt.Println(resourceCount)

		// 开通协程工作
		for i := 1;i <= rounterCount;i++  {
			go startJos(jobs,resultChan)
		}

		// 将工作内容加入管道，开始工作
		go addJobsWithLessons(jobs,courseModel,lessons)

		// 监听工作结果
		updateData := dealResults(resultChan,courseModel,resourceCount)

		close(resultChan)

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

}

// lessons加入工作队列
func addJobsWithLessons(jobs chan checkResource, courseModel model.Course,lessons []model.CheckLesson) {

	for _,lessonModel := range lessons {

		// 调用Media资源加入工作组
		addJobsWithMedias(jobs,courseModel,lessonModel,lessonModel.Medias)

	}

	close(jobs)
}
// media资源加入工作队列
func addJobsWithMedias(jobs chan checkResource,courseModel model.Course,lessonModel model.CheckLesson, medias []model.CheckMedia){

	for _,mediaModel := range medias  {

		// 创建一个media资源
		mediaResource := checkResource{
			lessonId:lessonModel.Id,
			lessonName:lessonModel.Name,
			objectKey:("kj/"+mediaModel.Url),
			mediaType:3,
			mediaSeq:mediaModel.Seq,
		}

		// 将资源加入工作组
		jobs <- mediaResource

		// 调用Image资源加入工作组
		addJobsWithImages(jobs,courseModel,lessonModel,mediaModel,mediaModel.Images)

	}

}
// image资源加入工作队列
func addJobsWithImages(jobs chan checkResource,courseModel model.Course,lessonModel model.CheckLesson, mediaModel model.CheckMedia, images []model.Image){

	for _,imageModel := range images {

		imagePath := imageModel.Url

		// 切断后缀
		imagePathArr := strings.Split(imagePath,".")

		imagePath = imagePathArr[0] + ".jpg@!style_ori"

		// 判断原始图片是否存在
		originPath := "kj/" + imagePath

		usePath := "img/" + courseModel.Id + "/" + imagePath

		imageOriginResource := checkResource{
			lessonId:lessonModel.Id,
			lessonName:lessonModel.Name,
			objectKey:originPath,
			mediaType:1,
			mediaSeq:mediaModel.Seq,
		}

		imageUseResource := checkResource{
			lessonId:lessonModel.Id,
			lessonName:lessonModel.Name,
			objectKey:usePath,
			mediaType:2,
			mediaSeq:mediaModel.Seq,
		}

		jobs <- imageOriginResource
		jobs <- imageUseResource

	}

}

// 开启工作组开始工作，将工作的结果加入结果队列
func startJos(jobs chan checkResource,results chan checkResult) {

	for resource := range jobs {

		isSave := checkResourceSaveStatus(resource.objectKey)

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

		lessonDict["check_status"] = 2

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

	// 打包检查结果资源

	updateData := make(map[string]interface{})

	updateData["lesson_ids"] = lessonDicts

	return updateData
}

/// 检查Lesson 数组迁移状态，需要返回整组lesson的状态字典
/**
 * @ courseModel 课程对象，便于处理课程id 和课程名称
 * @ lessons 待检查的lesson数组
 * @return 返回检查状态数组，每一个元素都包含一个lessond和一个lesson_status
 *         如果检查存在错误，则还要包含错误的meidas列表
 */
//func checkLessons(courseModel model.Course, lessons []model.CheckLesson) (lessonsResult []map[string]interface{}) {
//
//	var checkResult []map[string]interface{}
//
//	for _,lessonModel := range lessons {
//
//		errorMessage := "#" + courseModel.Id + "-" + courseModel.Name + "-" + lessonModel.Id + "-" + lessonModel.Name + " : "
//
//		errorLesson := ""
//
//		var mediasResult []map[string]interface{}
//
//		for _,mediaModel := range lessonModel.Medias {
//
//			mediaResult := make(map[string]interface{})
//
//			mediaResult["seq"] = mediaModel.Seq
//
//			errorMap := make(map[string]string)
//
//			mediaStatus := 0
//
//			if !checkMediaStatus(mediaModel.Url) {
//				errorMap["seq"] = strconv.Itoa(mediaModel.Seq)
//				errorMap["error_status"] = "3"
//				mediaStatus = 3
//			}
//
//			for _,imageModel := range mediaModel.Images{
//
//				imageStatus := checkImageStatus(courseModel.Id,imageModel.Url)
//
//				if imageStatus != 1 {
//
//					if mediaStatus != 3 {
//						mediaStatus = imageStatus
//					}
//
//					errorMap["seq"] = strconv.Itoa(mediaModel.Seq)
//
//					if errorMap["check_status"] == "" {
//						errorMap["check_status"] = strconv.Itoa(imageStatus)
//					}
//
//				}
//
//			}
//
//			if mediaStatus != 0 {
//				mediaResult["status"] = mediaStatus
//				mediasResult = append(mediasResult,mediaResult)
//			}
//
//			// 将errorMap中的信息拼接
//			if len(errorMap) == 0 {
//				continue
//			}
//			errorString := getErrorString(courseModel.Id,lessonModel.Id,lessonModel.Name,errorMap["seq"],errorMap["check_status"])
//
//			if errorLesson == "" {
//				errorLesson = errorString
//			}else {
//				errorLesson = errorLesson + "\n" + errorString
//			}
//
//		}
//
//		lessonResult := make(map[string]interface{})
//
//		lessonResult["lesson_id"] = lessonModel.Id
//
//		// 判断是否存在错误信息
//		if errorLesson == "" {
//
//			lessonResult["check_status"] = 1
//
//			checkResult = append(checkResult,lessonResult)
//
//			continue
//		}
//
//		lessonResult["check_status"] = 2
//
//		lessonResult["medias"] = mediasResult
//
//		checkResult = append(checkResult,lessonResult)
//
//		utils.AppendStringToFile(checkResourceLog,errorMessage + "\n" + errorLesson + "\n")
//
//	}
//
//	return checkResult
//}

func getErrorString(courseId string,lessonId string,lessonName string,seq string,errorStatus string,objectKey string) string {

	errorString := "seq=" + seq + "," + " error_status=" + errorStatus + "," + " objectKey=" + objectKey

	errorString = errorString + "\n" + "<a href=\"" + "http://www.bstcine.com/learn/lesson/"+ courseId + "?"

	errorString = errorString + "content_id=" + lessonId + "&seq=" + seq + "&error_type=" + errorStatus + "\">" + lessonName + "</a>"

	errorString = errorString + "\n" + "-----------------------------------------"

	return  errorString

}

func checkMediaStatus(mediaPath string) bool {

	mediaPath = "kj/" + mediaPath

	return checkResourceSaveStatus(mediaPath)
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

	originStatus := checkResourceSaveStatus(originPath)

	usePath := "img/" + courseId + "/" + imagePath

	useStatus := checkResourceSaveStatus(usePath)

	if !useStatus {
		 return 3
	}

	if !originStatus {
		return 2
	}

	return  1
}

func checkResourceSaveStatus(objectKey string) bool {

	// 检查图片状态需要检查两次，一次为原图，一次为水印图，两次成功才返回 1
	requestPath := "http://oss.bstcine.com/" + objectKey

	checkCount = 3

	for checkCount >= 0  {

		checkCount -= 1

		resp,err := http.Head(requestPath)

		if err == nil {

			if resp.StatusCode == 200 {

				length := resp.Header["Content-Length"]

				if len(length) > 0 {

					currentLength,err := strconv.ParseInt(length[0],0,64)

					if err != nil {
						return  false
					}

					return (currentLength > 10240)
				}

			}

		}

	}

	return  false
}

//***************************************************************
//*********************                 *************************
//*********************     网络模块     *************************
//*********************                 *************************
//***************************************************************
func getCourses(token string) []model.Course {

	_,courses := utils.ListWithCourses(model.Request{token,"cine.web",nil})

	// 获取课程列表
	clientCourseIds := utils.ClientInputWithMessage("请输入待检查的课程Id，中间用\",\"隔开，点击\"enter\"键结束\n如果要检查全部课件，请直接点击\"endter\"键",'\n')

	if clientCourseIds == "" {
		return courses
	}

	var targetCourse []model.Course

	courseIds := strings.Split(clientCourseIds,",")

	fmt.Println(courseIds)

	for i := 0;i < len(courses);i++  {

		courseModel := courses[i]

		courseIdsCount := len(courseIds)

		if courseIdsCount <= 0 {
			break
		}

		for j := 0;j < len(courseIds) ;j++  {

			courseId := courseIds[j]

			if courseModel.Id == courseId {

				targetCourse = append(targetCourse,courseModel)

				courseIds = append(courseIds[:j],courseIds[j+1:]...)

				break
			}

		}

	}

	return targetCourse

}

func getLessons(token string, courseId string) []model.CheckLesson {

	data := make(map[string]interface{})
	data["cid"] = courseId
	_,lessons := utils.ListWithCheckMedias(model.Request{token,"cine.web",data})

	return lessons
}

