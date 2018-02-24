package main

import (
	"./model"
	"./utils"
	"fmt"
	"os"
	"strconv"
	"path"
)

func main() {
	fmt.Println("欢迎使用课程资源下载工具....")

	var debug = false

	var login string
	var password string

	var courseId string
	var lessonId string
	var courseName string

	var outPutPath string

	if debug {
		login = "xxx"
		password = "123"
		courseId = "d011503974382830Tcne3UQckf"
		lessonId = "d011504248088023Y3ckCRhuyP"
		courseName = courseId

		outPutPath = "/Test/课件资源/"
	} else {
		for len(login) <= 0 || len(password) <= 0 || len(courseId) <= 0 {
			login, password, courseId, lessonId, courseName = getDownloadArags()
		}

		if len(courseName) <= 0 || courseName == "" {
			courseName = courseId
		}

		outPutPath = utils.GetCurPath() + string(os.PathSeparator) + "课件资源" + string(os.PathSeparator)
	}

	os.MkdirAll(outPutPath, 0777)

	sitecode := "cine.web"

	data := make(map[string]interface{})
	data["phone"] = login
	data["password"] = password
	_, token := utils.Signin(model.Request{"", sitecode, data})

	if len(token) <= 0 || token == "" {
		fmt.Println("no token")
		return
	}

	data = make(map[string]interface{})
	data["cid"] = courseId
	if lessonId != "-1" && !(len(lessonId) <= 0 || lessonId == "") {
		data["filter"] = []string{lessonId}
	}
	_, rows := utils.ListWithMedias(model.Request{token, sitecode, data})

	var files []model.DownFile
	for i := 0; i < len(rows); i++ {
		chapterName := rows[i].Name
		//fmt.Println("=>1. " + chapterName)
		children := rows[i].Children;
		for j := 0; j < len(children); j++ {
			lessonName := children[j].Name
			//fmt.Println("=>=>2. " + lessonName)
			medias := children[j].Medias;
			for k := 0; k < len(medias); k++ {
				media := medias[k]
				//fmt.Println("=>=>=>3. type:" + media.Type + ", url:" + media.Url)

				file := model.DownFile{}
				file.Name = strconv.Itoa(k+1) + path.Ext(media.Url)
				file.ChapterName = chapterName
				file.LessonName = lessonName
				file.Path = media.Url
				files = append(files, file)

				images := medias[k].Images
				for l := 0; l < len(images); l++ {
					image := images[l]
					//fmt.Println("=>=>=>=>4. time:" + image.Time + ", url:" + image.Url)

					file := model.DownFile{}
					file.Name = courseName + "-" + lessonName + "-" + strconv.Itoa(k+1) + "-" + image.Time + path.Ext(image.Url)
					file.ChapterName = chapterName
					file.LessonName = lessonName
					file.Path = image.Url
					files = append(files, file)
				}
			}
		}
	}

	for i := 0; i < len(files); i++ {
		file := files[i]

		downCoursePutPath := outPutPath + courseName + string(os.PathSeparator)
		if _, err := os.Stat(downCoursePutPath); err != nil {
			os.MkdirAll(downCoursePutPath, 0777)
		}

		downChapterPath := downCoursePutPath + file.ChapterName + string(os.PathSeparator)
		if _, err := os.Stat(downChapterPath); err != nil {
			os.MkdirAll(downChapterPath, 0777)
		}

		downLessonPath := downChapterPath + file.LessonName + string(os.PathSeparator)
		if _, err := os.Stat(downLessonPath); err != nil {
			os.MkdirAll(downLessonPath, 0777)
		}

		utils.DownloadFile(file.Path, downLessonPath+file.Name)
	}

	var code string

	fmt.Printf("一共有%d个文件，请输入任意键退出", len(files))
	fmt.Scanln(&code)
}

func getDownloadArags() (login, password, courseId, lessonId, courseName string) {
	fmt.Println("请输入 User,Password,CourseId,LessonId,CourseName 中间用空格隔开!")
	fmt.Println("注释：lessonId：-1时，下载所有Lesson")
	fmt.Println("例如：user password 42 -1 动物农庄")
	fmt.Scanln(&login, &password, &courseId, &lessonId, &courseName)
	return login, password, courseId, lessonId, courseName
}
