package main

import (
	"./model"
	"./utils"
	"fmt"
	"os"
	"strconv"
	"path"
	"strings"
	"sync"
)

func main() {
	fmt.Println("欢迎使用课程资源下载工具....")

	var debug = false

	var login string
	var password string

	var fileType string
	var courseId string
	var lessonId string
	var courseName string

	var outPutPath string

	if debug {
		login = ""
		password = "123"
		fileType = "0"
		courseId = "d0114987070071216e0wtYHMZj"
		lessonId = ""
		courseName = courseId

		outPutPath = "/Test/课件资源/"
	} else {
		for len(login) <= 0 || len(password) <= 0 || len(fileType) <= 0 || len(courseId) <= 0 {
			login, password, fileType, courseId, lessonId, courseName = getDownloadArags()
		}

		if len(courseName) <= 0 || courseName == "" {
			courseName = courseId
		}

		outPutPath = utils.GetCurPath() + "课件资源" + string(os.PathSeparator)
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
				file.Path = getUrlByType(fileType, media.Url)
				files = append(files, file)

				images := medias[k].Images
				for l := 0; l < len(images); l++ {
					image := images[l]
					//fmt.Println("=>=>=>=>4. time:" + image.Time + ", url:" + image.Url)

					file := model.DownFile{}
					file.Name = courseName + "-" + lessonName + "-" + strconv.Itoa(k+1) + "-" + image.Time + path.Ext(image.Url)
					file.ChapterName = chapterName
					file.LessonName = lessonName
					file.Path = getUrlByType(fileType, image.Url)
					files = append(files, file)
				}
			}
		}
	}

	var maxGo = 8
	for i := 0; i < len(files); i += maxGo {
		tmp := files[i:i+8]
		var wg sync.WaitGroup
		wg.Add(maxGo)
		for i := 0; i < maxGo; i++ {
			file := tmp[i]

			if file.Path == "" {
				wg.Done()
				continue
			}

			go func() {
				downPath := outPutPath+courseName+string(os.PathSeparator)+file.ChapterName+string(os.PathSeparator)+file.LessonName+string(os.PathSeparator)+file.Name
				utils.DownloadFile(file.Path, downPath)
				wg.Done()
			}()
		}
		wg.Wait()
	}

	fmt.Printf("一共有%d个文件，请输入任意键退出", len(files))
	fmt.Scanln()
}

func getDownloadArags() (login, password, fileType, courseId, lessonId, courseName string) {
	fmt.Println("请输入 User,Password,Type,CourseId,LessonId,CourseName 中间用空格隔开!")
	fmt.Println("注释：Type:(0:压缩文件，1：原始文件) 、 LessonId：(-1:下载所有)")
	fmt.Println("例如：user password 0 42 -1 动物农庄")
	fmt.Scanln(&login, &password, &fileType, &courseId, &lessonId, &courseName)
	return login, password, fileType, courseId, lessonId, courseName
}

func getUrlByType(fileType, url string) string {
	if fileType != "1" {
		return url
	}

	url = url[strings.Index(url, "/f/")+3:len(url)]

	lastIndex := strings.LastIndex(url, "/") + 1

	srcPrefix := url[0:lastIndex]
	srcName := url[lastIndex:len(url)]
	srcName = strings.Replace(srcName, ".png.jpg", ".png", -1)
	srcName = strings.Replace(srcName, ".jpg.jpg", ".jpg", -1)

	return "http://www.bstcine.com/f/" + srcPrefix + srcName
}
