package main

import (
	"utils"
	"fmt"
	"sync"
	"os"
	"strconv"
	"path"
)

var CommonReq = utils.Request{"TGuPYryS", "cine.web", make(map[string]string)}

var courseDownloadWaitGroup sync.WaitGroup

func main() {
	fmt.Println("欢迎使用课程资源下载工具....")

	isdebug := false

	var outPutPath string
	var courseId string

	if isdebug {
		outPutPath = "/Go/Test/课件资源/"
		os.MkdirAll(outPutPath, 0777)
	} else {
		outPutPath = utils.GetCurPath() + string(os.PathSeparator) + "课件资源" + string(os.PathSeparator)
	}

	fmt.Println("请输入课程ID：(注：不填为所有课程)，下载文件到 '习题单词库' 文件夹下")
	fmt.Scanln(&courseId)

	CommonReq.Data = make(map[string]string);
	CommonReq.Data["cid"] = courseId
	result := utils.ListWithMedias(CommonReq)

	rows := result.Result.Rows
	fmt.Println(len(rows))

	var files []utils.DownFile
	for i := 0; i < len(rows); i++ {
		chapterName := rows[i].Name
		fmt.Println("=>1. " + chapterName)
		children := rows[i].Children;
		for j := 0; j < len(children); j++ {
			lessonName := children[j].Name
			fmt.Println("=>=>2. " + lessonName)
			medias := children[j].Medias;
			for k := 0; k < len(medias); k++ {
				media := medias[k]
				fmt.Println("=>=>=>3. type:" + media.Type + ", url:" + media.Url)

				file := utils.DownFile{}
				file.Name = strconv.Itoa(k+1) + path.Ext(media.Url)
				file.ChapterName = chapterName
				file.LessonName = lessonName
				file.Path = media.Url
				files = append(files, file)

				images := medias[k].Images
				for l := 0; l < len(images); l++ {
					image := images[l]
					fmt.Println("=>=>=>=>4. time:" + image.Time + ", url:" + image.Url)

					file := utils.DownFile{}
					file.Name = strconv.Itoa(k+1) + "-" + image.Time + path.Ext(image.Url)
					file.ChapterName = chapterName
					file.LessonName = lessonName
					file.Path = image.Url
					files = append(files, file)
				}
			}
		}
		fmt.Println("")
	}

	fmt.Printf("一共%d个文件\n", len(files))
	for i := 0; i < len(files); i++ {
		file := files[i]

		downCoursePutPath := outPutPath + courseId + string(os.PathSeparator)
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

		courseDownloadWaitGroup.Add(1)
		go HelloDown(downLessonPath+file.Name, file.Path)
	}

	courseDownloadWaitGroup.Wait()
}

func HelloDown(path, url string) {
	utils.DownloadFile(path, url)
	courseDownloadWaitGroup.Done()
	fmt.Println("download ok => name: " + path + " || url: " + url)
}
