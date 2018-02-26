package main

import (
	"fmt"
	"strings"
	"os"
	"./utils"
	"./model"
)

func main() {
	var debug = false

	var login string
	var password string
	var fileType string
	var courseId string
	var outPutPath string

	if debug {
		login = "kim"
		password = "123"
		fileType = "0"
		courseId = "d011489545230763dQnpS4PsE0,d0114987070071216e0wtYHMZj"

		outPutPath = "/Test/"
	} else {
		for len(login) <= 0 || len(password) <= 0 || len(fileType) <= 0 || len(courseId) <= 0 {
			login, password, fileType, courseId = getArags()
		}

		outPutPath = utils.GetCurPath() + string(os.PathSeparator)
	}

	data := make(map[string]interface{})
	data["phone"] = login
	data["password"] = password
	req := model.Request{"", "cine.web", data}
	_, token := utils.Signin(req)
	fmt.Println(token)

	var files []string

	cids := strings.Split(courseId, ",")

	for ii := 0; ii < len(cids); ii++ {
		data["cid"] = cids[ii]
		req = model.Request{token, "cine.web", data}
		_, rows := utils.ListWithMedias(req)

		for i := 0; i < len(rows); i++ {
			children := rows[i].Children;
			for j := 0; j < len(children); j++ {
				medias := children[j].Medias;
				for k := 0; k < len(medias); k++ {
					media := medias[k]

					if !strings.Contains(media.Url, "www.bstcine.com") && len(media.Url) > 0 {
						files = append(files, regUrl(fileType, media.Url))
					}

					images := medias[k].Images
					for l := 0; l < len(images); l++ {
						image := images[l]
						if !strings.Contains(image.Url, "www.bstcine.com") && len(image.Url) > 0 {
							files = append(files, regUrl(fileType, image.Url))
						}
					}
				}
			}
		}
	}

	utils.WriteLines(files, outPutPath+"http.list")
}

func getArags() (login, password, fileType, courseId string) {
	fmt.Println("请输入 User,Password,fileType,CourseId 中间用空格隔开! ")
	fmt.Println("注释：Type:(0:压缩文件，1：原始文件) 、 LessonId：(-1:下载所有)")
	fmt.Println("例如：user password 0 42")
	fmt.Scanln(&login, &password, &fileType, &courseId)
	return login, password, fileType, courseId
}

func regUrl(fileType, url string) string {
	if fileType == "0" {
		urls := strings.Split(url, "/")
		url = urls[0] + "//" + urls[2] + "/" + urls[3] + "/ " + urls[4] + "/" + urls[5] + "/" + urls[6] + "/" + urls[7] + "/" + urls[8] + "/" + urls[9]
	} else {
		url = url[strings.Index(url, "/f/")+3:len(url)]

		lastIndex := strings.LastIndex(url, "/") + 1

		srcPrefix := url[0:lastIndex]
		srcName := url[lastIndex:len(url)]
		srcName = strings.Replace(srcName, ".png.jpg", ".png", -1)
		srcName = strings.Replace(srcName, ".jpg.jpg", ".jpg", -1)

		url =  "http://www.bstcine.com" + "/ f/" + srcPrefix + srcName
	}

	fmt.Println(url)
	return url
}
