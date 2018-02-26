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
	var courseId string
	var outPutPath string

	if debug {
		login = "kim"
		password = "123"
		courseId = "d0114987070071216e0wtYHMZj,d011499750829444sJyNV5gcuc"

		outPutPath = "/Test/"
	} else {
		for len(login) <= 0 || len(password) <= 0 || len(courseId) <= 0 {
			login, password, courseId = getArags()
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

	cids := strings.Split(courseId,",")

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
					files = append(files, regUrl(media.Url))

					images := medias[k].Images
					for l := 0; l < len(images); l++ {
						image := images[l]
						files = append(files, regUrl(image.Url))
					}
				}
			}
		}
	}

	utils.WriteLines(files, outPutPath+"http.list")
}

func getArags() (login, password, courseId string) {
	fmt.Println("请输入 User,Password,CourseId 中间用空格隔开! 例如：")
	fmt.Println("user password 42")
	fmt.Scanln(&login, &password, &courseId)
	return login, password, courseId
}

func regUrl(url string) string {
	url = url[strings.Index(url, "/f/")+3:len(url)]

	lastIndex := strings.LastIndex(url, "/") + 1

	srcPrefix := url[0:lastIndex]
	srcName := url[lastIndex:len(url)]
	srcName = strings.Replace(srcName, ".png.jpg", ".png", -1)
	srcName = strings.Replace(srcName, ".jpg.jpg", ".jpg", -1)

	return "http://www.bstcine.com" + "/f/ " + srcPrefix + srcName
}
