package main

import (
	"./utils"
	"fmt"
	"os"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"log"
)

type Lesson struct {
	Id       string
	CourseId string
	Duration int
	Medias   []Media
}

type Media struct {
	Type     string  `json:"type"`
	Url      string  `json:"url"`
	Duration int     `json:"duration"`
	Images   []Image `json:"images"`
	Size     int     `json:"size"`
}

type Image struct {
	Time string `json:"time"`
	Url  string `json:"url"`
}

type UpdateArgs struct {
	BaseUrl  string
	Host     string
	User     string
	Password string
	Database string
	CourseId string
}

type Format struct {
	Duration string
	Size string
}

var db *sql.DB

var updateArgs UpdateArgs

func main() {
	///mnt/web/kj.bstcine.com/wwwroot/kj/d011502846526016MpDBrnsR8p/f/2017/11/03/xxxxxx.mp3
	//http://www.bstcine.com/img/
	//baseUrl = "http://www.bstcine.com/f/"
	isdebug := true

	if isdebug {
		updateArgs.BaseUrl = "/Volumes/Go/Test"
		updateArgs.Host = "127.0.0.1:3306"
		updateArgs.User = "root"
		updateArgs.Password = "dev0423wx"
		updateArgs.Database = "cine"
		updateArgs.CourseId = "updatetest"
	} else {
		args := os.Args
		if len(args) < 6 {
			fmt.Println("Please input need arge (/Volumes/Go 127.0.0.1:3306 root 123456 cine courseid)...")
			return
		}

		updateArgs.BaseUrl = args[1]
		updateArgs.Host = args[2]
		updateArgs.User = args[3]
		updateArgs.Password = args[4]
		updateArgs.Database = args[5]
		updateArgs.CourseId = args[6]
	}


	//校验文件路径...
	if !utils.Exists(updateArgs.BaseUrl) {
		fmt.Println("error: Input BasePath is not exists...")
		return
	}

	//校验数据库链接... root:dev0423wx@tcp(127.0.0.1:3306)/cine
	dataSourceName := updateArgs.User + ":" + updateArgs.Password + "@tcp(" + updateArgs.Host + ")/" + updateArgs.Database
	db, _ = sql.Open("mysql", dataSourceName)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("error: Input db args is error...")
		return
	}

	if updateArgs.CourseId == "" {
		fmt.Println("error: Input course_id is error....")
		return
	}

	lessons, _ := DbGetLessons(updateArgs.CourseId)

	SearchLesson(lessons)
}

/**
遍历 Lesson
 */
func SearchLesson(lessons []Lesson) {
	if len(lessons) <=0 {
		fmt.Println("info: Not lesson data need update...")
		return
	}

	for _, lesson := range lessons {
		duration, medias := UpdateMedias(lesson)
		jsonMedias, err := json.Marshal(medias)

		if err == nil {
			DbUpdateLesson(lesson.Id, duration, string(jsonMedias))
		} else {
			fmt.Println(err)
		}
	}
}

/**
更新 Lesson 下 Media 文件时长和大小
 */
func UpdateMedias(lesson Lesson) (newDuration int,newMedias []Media) {
	courseId := lesson.CourseId

	for _, media := range lesson.Medias {
		format := GetMediaFile(courseId, media)

		i, _ := strconv.ParseFloat(format.Duration, 64)
		duration := int(i)
		size,_ := strconv.Atoi(format.Size)

		newDuration += duration
		media.Duration = duration
		media.Size = size/1000
		newMedias = append(newMedias, media)
	}

	return newDuration, newMedias
}

/**
获取 Media 文件的信息
 */
func GetMediaFile(courseId string, media Media) Format {
	path := updateArgs.BaseUrl + string(os.PathSeparator) + courseId + string(os.PathSeparator) + media.Url

	jsonFileInfo := utils.GetJsonFileInfo(path)

	var fileinfo map[string]Format
	json.Unmarshal([]byte(jsonFileInfo), &fileinfo)
	format, _ := fileinfo["format"]

	log.Println(format)

	return format
}

/**
查询数据库
 */
func DbGetLessons(courseId string) (lessons []Lesson, errLessonIds []string) {
	rows, _ := db.Query("select id,lesson_id,medias from t_content where delete_by is null and type = '2'  and check_status = '1' and lesson_id = ? ", courseId)
	defer rows.Close()

	for rows.Next() {
		var lesson Lesson
		var mediasJson string
		err := rows.Scan(&lesson.Id, &lesson.CourseId, &mediasJson)

		if err == nil {
			var medias []Media
			err := json.Unmarshal([]byte(mediasJson), &medias)

			if err == nil {
				lesson.Medias = medias
				lessons = append(lessons, lesson)
			} else {
				errLessonIds = append(errLessonIds, lesson.Id)
			}
		}
	}

	return lessons, errLessonIds
}

/**
更新数据库
 */
func DbUpdateLesson(id string, duration int, medias string) {
	if duration == 0 {
		fmt.Println("error : lesson duration is 0 ....")
		return
	}

	if medias == "" {
		fmt.Println("error : medias is null ....")
		return
	}

	_, err := db.Exec("UPDATE t_content SET duration = ?,medias = ?,check_status = '2' WHERE id = ? and check_status = '1' ", duration, medias, id)
	if err == nil {
		fmt.Println("info: lesson( "+id+") data update success...")
	} else {
		fmt.Println(err)
	}
}
