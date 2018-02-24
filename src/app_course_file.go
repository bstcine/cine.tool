package main

import (
	"fmt"
	"strings"
	"./utils"
	"./model"
	"os"
)

func main() {
	var debug = true
	var inputArgs model.InputArgs

	if debug {
		inputArgs.LocalPath = "/Test/"
		inputArgs.OutputPath = "/Test/"
	} else {
		inputArgs.LocalPath = utils.GetCurPath()
		inputArgs.OutputPath = utils.GetCurPath()
	}

	data := make(map[string]string)
	data["phone"] = "kim"
	data["password"] = "123"
	req := model.Request{"", "cine.web", data}
	_, token := utils.Signin(req)
	fmt.Println(token)

	data["cid"] = "42"
	req = model.Request{token, "cine.web", data}
	_, rows := utils.ListWithMedias(req)

	var files []string

	for i := 0; i < len(rows); i++ {
		//fmt.Println("=>1. " + chapterName)
		children := rows[i].Children;
		for j := 0; j < len(children); j++ {
			//fmt.Println("=>=>2. " + lessonName)
			medias := children[j].Medias;
			for k := 0; k < len(medias); k++ {
				media := medias[k]
				files = append(files, regUrl(media.Url))
				//fmt.Println("=>=>=>3. type:" + media.Type + ", url:" + media.Url)

				images := medias[k].Images
				for l := 0; l < len(images); l++ {
					image := images[l]
					files = append(files, regUrl(image.Url))
					//fmt.Println("=>=>=>=>4. time:" + image.Time + ", url:" + image.Url)
				}
			}
		}
	}

	fmt.Println(files)

	utils.WriteLines(files, inputArgs.OutputPath+string(os.PathSeparator)+"http.list")
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
