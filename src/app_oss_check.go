package main

import (
	"./utils"
	"./model"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strings"
)

/**检查课程资源迁移到oss是否成功
 * 1.读取配置资源信息，获取账户名，密码，待检查的课件id列表，及需要的隐私信息
 * 2.登入系统，获取权限令牌token
 * 3.根据课件id列表，逐一获取lesson列表
 * 4.遍历lesson列表中的资源路径
 * 5.根据资源路径和课件id，访问oss上是否存在这些资源
 * 6.lesson中的资源如果已确认存在，标记lesson检查状态为OK
 * 7.将检查失败的lesson，打印对应的lessonId,courseId
 */

 type checkModel struct {
 	lessonIds []string
 }

 var myarsp map[string]string
 var account string
 var passwords string
 var token string
 var courseIds []string
 var keyId string
 var keySecret string
 var mybucket string
 var endpoint string

func main() {

	// 启动加载配置信息
	launchArsp()

	// 登录账户，获取权限令牌token
	data := make(map[string]interface{})
	data["phone"] = account
	data["password"] = passwords
	_,token = utils.Signin(model.Request{"","cine.web",data})

	// 检查课程列表是否已迁移完毕
	for i := 0;i < len(courseIds) ;i++  {
		checkCourse(courseIds[i])
	}

}

/**
 * 检查课件是否迁移完毕
 * @param courseId 待检查的课件id
 */
func checkCourse(courseId string){

	data := make(map[string]interface{})
	data["cid"] = courseId
	_,rows := utils.ListWithMedias(model.Request{token,"cine.web",data})

	if len(rows) <= 0 {
		fmt.Println("课件id有误：cid=",courseId)
		return
	}

	for i := 0; i < len(rows); i++ {

		lessons := rows[i].Children

		var lesson model.Lesson

		for j := 0; j < len(lessons); j++ {

			lesson = lessons[j]

			if len(lesson.Medias) <= 0 {
				continue
			}

			// 判断lesson中的资源是否存在
			checkStatus := checkLesson(lesson,courseId)

			if checkStatus {
				fmt.Println("剩余",len(rows)-i,"组",len(lessons)-j, "个")
			}else {
				fmt.Println("资源检查失败，",courseId,"——",lesson.Id)
			}

		}
	}

	fmt.Println("课程检查完毕")

}
/**
 * 检查迁移的lesson是否迁移完成
 * @param lesson 课程的结构体对象
 * @param courseId 待检查课程对应的课件id
 * @param return lesson资源检查结果（是否全部资源已迁移完成）
 */
func checkLesson(lesson model.Lesson, courseId string) bool {

	var media model.Media
	var image model.Image
	var path  string

	// 获取lesson中的medias
	for i := 0; i < len(lesson.Medias);i++  {

		media = lesson.Medias[i]

		// 截取media的相对路径
		mediaPaths := strings.Split(media.Url,"/f/")
		path = mediaPaths[len(mediaPaths) - 1]
		// 判断media 是否存在
		isMediaExist := checkOSSexist(courseId,path,false)

		if !isMediaExist {
			return false
		}

		// media 已存在，判断image是否成功

		for j := 0;j < len(media.Images);j++ {

			image = media.Images[j]

			imagePaths := strings.Split(image.Url,"/f/")

			path = imagePaths[len(imagePaths) - 1]

			var isImage bool
			if strings.Contains(image.Url,"www.bstcine.com") {
				isImage = false
			}else {
				isImage = true
			}

			isImageExist := checkOSSexist(courseId,path,isImage)

			if !isImageExist {
				return false
			}
		}
	}

	var data  = make(map[string]interface{})

	data["lesson_ids"] = []string{
		lesson.Id,
	}

	_,status := utils.UpdateLessonCheckStatus(model.Request{token,"cine.web",data})

	if !status {
		fmt.Println("失败：",lesson.Name,lesson.Id)
	}

	return status
}

// 启动加载配置信息，必须在第一不执行，否则配置信息可能为空
func launchArsp(){

	myarsp = utils.GetConfArgs("/Users/lidangkun/Desktop/GoProjects/workspace/cineTool/assets/app_oss_tmp.cfg")

	fmt.Println(myarsp)

	account = myarsp["account"]
	passwords = myarsp["password"]
	courses := myarsp["checkCourses"]
	courseIds = strings.Split(courses,",")
	fmt.Println(courseIds)
	keyId = myarsp["AccessKeyId"]
	keySecret = myarsp["AccessKeySecret"]
	mybucket = myarsp["Bucket"]
	endpoint = myarsp["Endpoint"]

}

/**
 * 从oss读取资源，判断是否已存在
 * @param courseId 待检查资源对应的课件id
 * @param path 资源相对url路径（"/f/"之后的路径）
 * @param isImage 是否是存储在图片库的资源（用来判断是否统一修改扩展名为".jpg"）
 * @return 资源检查的结果(是否已存储在oss中)
 */
func checkOSSexist(courseId string,path string, isImage bool) bool {

	client,err := oss.New(endpoint,keyId,keySecret)

	if err != nil {
		fmt.Println(err)
		fmt.Println("创建客户端对象失败")
		return false
	}

	bucket,err := client.Bucket(mybucket)

	if err != nil {
		fmt.Println(err)
		fmt.Println("创建bucket失败")
		return false
	}

	var objectKey string

	if isImage {
		paths := strings.Split(path,".")
		objectKey = "img/" + courseId + "/" + paths[0] + ".jpg"
	}else {
		objectKey = "kj/" + path
	}
	isExist,err := bucket.IsObjectExist(objectKey)

	if err != nil {
		fmt.Println("资源判断是否存在失败")
		fmt.Println(err)
		return false
	}

	return isExist
}
