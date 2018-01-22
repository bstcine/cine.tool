package main

import (
	"utils"
	"fmt"
	"os"
	"strconv"
	"path"
)

var CommonReq = utils.Request{"", "cine.web", nil}

func main() {
	fmt.Println("欢迎使用课程资源下载工具....")

	isdebug := false

	var outPutPath string
	var courseId string
	var courseName string
	var token string

	if isdebug {
		outPutPath = "/Go/Test/课件资源/"
		token = "9f6jP45S"
		courseId = "d011503974382830Tcne3UQckf"
	} else {
		outPutPath = utils.GetCurPath() + string(os.PathSeparator) + "课件资源" + string(os.PathSeparator)

		for len(token) <= 0 || len(courseId) <= 0 {
			token, courseId, courseName = GetArags()
		}
	}

	os.MkdirAll(outPutPath, 0777)

	CommonReq.Token = token
	CommonReq.Data = make(map[string]string)
	CommonReq.Data["cid"] = courseId

	var files []utils.DownFile

	result := utils.ListWithMedias(CommonReq)
	rows := result.Result.Rows
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

				file := utils.DownFile{}
				file.Name = strconv.Itoa(k+1) + path.Ext(media.Url)
				file.ChapterName = chapterName
				file.LessonName = lessonName
				file.Path = media.Url
				files = append(files, file)

				images := medias[k].Images
				for l := 0; l < len(images); l++ {
					image := images[l]
					//fmt.Println("=>=>=>=>4. time:" + image.Time + ", url:" + image.Url)

					file := utils.DownFile{}
					file.Name = courseName + "-" + lessonName + "-" + strconv.Itoa(k+1) + "-" + image.Time + path.Ext(image.Url)
					file.ChapterName = chapterName
					file.LessonName = lessonName
					file.Path = image.Url
					files = append(files, file)
				}
			}
		}
	}

	if len(courseName) <= 0 || courseName == "" {
		courseName = courseId
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

		utils.DownloadFile(file.Path, downLessonPath)
	}

	var code string

	fmt.Printf("一共有%d个文件，请输入任意键退出", len(files))
	fmt.Scanln(&code)
}

func GetArags() (token, courseId, courseName string) {
	fmt.Println("请输入 Token,CourseId 中间用空格隔开! 例如：")
	fmt.Println("TGuPYryS 42 动物农庄")
	fmt.Scanln(&token, &courseId, &courseName)
	return token, courseId, courseName
}
