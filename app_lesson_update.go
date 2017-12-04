package main

import (
	"./utils"
	"fmt"
	"os"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Lesson struct {
	Id         string
	CourseId   string
	CourseName string
	ParentName string
	Name       string
	Medias     string
	Duration   int
	Msg        string
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
	Size     string
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
		updateArgs.Password = "cine"
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

	lessons := DbGetAllLessons()

	SearchLesson(lessons)
}

/**
遍历 Lesson
 */
func SearchLesson(lessons []Lesson) {
	if len(lessons) <= 0 {
		fmt.Println("info: Not lesson data need update...")
		return
	}

	var errLessons []Lesson

	for _, lesson := range lessons {
		newLesson := CheckMedias(lesson)

		if newLesson.Msg == "" {
			//DbUpdateLesson(lesson.Id, lesson.Duration)
		} else {
			errLessons = append(errLessons, newLesson)
		}
	}

	fmt.Println(errLessons)
}

/**
更新 Lesson 下 Media 文件时长和大小
 */
func CheckMedias(lesson Lesson) (newLesson Lesson) {
	fmt.Println("check lesson:" + lesson.Id)

	var medias []Media
	err := json.Unmarshal([]byte(lesson.Medias), &medias)

	lesson.Medias = ""
	newLesson = lesson

	if err != nil {
		newLesson.Msg = "medias json error"
		return newLesson
	}

	var msg string
	var allDuration int

	for index, media := range medias {
		format, err := ReadMediaFile(lesson.CourseId + string(os.PathSeparator) + media.Url)

		if err == nil {
			i, _ := strconv.ParseFloat(format.Duration, 64)
			duration := int(i)
			size, _ := strconv.Atoi(format.Size)

			allDuration += duration
			media.Size = size / 1000

			if media.Duration != duration {
				msg += "第 " + strconv.Itoa(index+1) + " 个课件的音视频时长应为：" + strconv.Itoa(duration) + ";<br>"
			}
		} else {
			msg += "第 " + strconv.Itoa(index+1) + " 个课件的音视频不存在;<br>"
		}
	}

	newLesson.Duration = allDuration
	newLesson.Msg = msg
	return newLesson
}

/**
获取 Media 文件的信息
 */
func ReadMediaFile(path string) (Format, error) {
	result, err := utils.GetJsonFileInfo(updateArgs.BaseUrl + string(os.PathSeparator) + path)

	var fileinfo map[string]Format
	json.Unmarshal([]byte(result), &fileinfo)
	format, _ := fileinfo["format"]

	return format, err
}

/**
查询数据库
 */
func DbGetAllLessons() (lessons []Lesson) {
	var sql = "select a.id,a.lesson_id,b.name lesson_name,c.name parent_name,a.name,a.medias " +
		"from t_content a left join t_lesson b on a.lesson_id = b.id left join t_content c on a.parent_id = c.id " +
		"where a.delete_by is null and a.parent_id <> '1' and a.type = '2' and a.check_status = '1' and b.delete_by is null and b.id is not null " +
		"order by a.lesson_id,a.parent_id,a.seq asc"
	rows, _ := db.Query(sql)
	defer rows.Close()

	for rows.Next() {
		var lesson Lesson
		err := rows.Scan(&lesson.Id, &lesson.CourseId, &lesson.CourseName, &lesson.ParentName, &lesson.Name, &lesson.Medias)

		if err == nil {
			lessons = append(lessons, lesson)
		}
	}

	return lessons
}

/**
更新数据库
 */
func DbUpdateLesson(id string, duration int) {
	if duration == 0 {
		fmt.Println("error : lesson duration is 0 ....")
		return
	}

	_, err := db.Exec("UPDATE t_content SET duration = ?,check_status = '2' WHERE id = ? and check_status = '1' ", duration, id)
	if err == nil {
		fmt.Println("info: lesson( " + id + ") data update success...")
	} else {
		fmt.Println(err)
	}
}
