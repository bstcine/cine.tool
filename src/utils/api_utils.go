package utils

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"../model"
	"../conf"
	"log"
	"strings"
)

func GetBaseUrl(url string) string {

	if conf.IsTestHost {
		url = conf.API_BASE_URL_TEST + url
	} else {
		url = conf.API_BASE_URL + url
	}

	return url
}

func CommonPost(url string, request model.Request, res interface{}) {
	jsonBytes, _ := json.Marshal(request)

	log.Println("======== 网络请求中 > body: " + string(jsonBytes))

	if !strings.Contains(url,"bstcine.com") {
		url = GetBaseUrl(url)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println("======== 网络请求成功 > body: " + string(body))

	json.Unmarshal(body, &res)
}

func CommonGet(url string, res interface{}) {
	log.Println("======== 网络请求中 > ")

	resp, err := http.Get(GetBaseUrl(url))
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println("======== 网络请求成功 > body: " + string(body))

	json.Unmarshal(body, &res)
}

func SigninPermission(request model.RequestPermission) (res model.ResultPersion, token string) {

	jsonBytes, _ := json.Marshal(request)

	log.Println("======== 网络请求中 > body: " + string(jsonBytes))

	resp, err := http.Post(GetBaseUrl(conf.APIURL_Content_Permission_Sigin), "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println("======== 网络请求成功 > body: " + string(body))

	json.Unmarshal(body, &res)

	return res, res.Token
}

func Signin(request model.Request) (res model.Res, token string) {
	CommonPost(conf.APIURL_Auth_Signin, request, &res)

	return res, res.Result["token"].(string)
}

func GetFiles(password, fileType, cid string) (res map[string]interface{}, data []interface{}) {
	url := conf.APIURL_Tool_Files + "?password=" + password + "&type=" + fileType + "&cid=" + cid

	CommonGet(url, &res)

	return res, res["data"].([]interface{})
}

func ListWithCourses(request model.Request) (res model.ResList, courses []model.Course) {

	CommonPost(conf.APIURL_Content_Course_List,request,&res)

	coursesJson,_ := json.Marshal(res.Result.Rows)
	json.Unmarshal([]byte(coursesJson),&courses)

	return res, courses
}

func ListWithMedias(request model.Request) (res model.ResList, rows []model.Chapter) {

	CommonPost(conf.APIURL_Content_Chapter_ListWithMedia, request, &res)

	rowsJson, _ := json.Marshal(res.Result.Rows)
	json.Unmarshal([]byte(rowsJson), &rows)

	return res, rows
}

func ListWithCheckCourses(request model.Request) (res model.ResList, courses []model.Course) {

	CommonPost(conf.APIURL_Content_Lesson_CheckCourseList,request,&res)

	coursesJson,_ := json.Marshal(res.Result.Rows)
	json.Unmarshal([]byte(coursesJson),&courses)

	return res, courses

}

//获取待检查的lesson列表
func ListWithCheckMedias(request model.Request) (res model.ResList, lessons []model.CheckLesson) {

	CommonPost(conf.APIURL_Content_Lesson_CheckListWithLessons, request, &res)

	rowsJson,_ := json.Marshal(res.Result.Rows)

	json.Unmarshal([]byte(rowsJson),&lessons)

	return res, lessons
}

//更新Lesson检查状态
func UpdateLessonCheckStatus(request model.Request) (res model.ResCheckList, status bool) {

	CommonPost(conf.APIURL_Content_Lesson_UpdateCheckStatus,request,&res)

	fmt.Println(res)

	return res, res.Result["status"]
}

//*****************************************************************
//*****************************************************************
//************************ 获取服务器权限 ***************************
//*****************************************************************
//*****************************************************************
/// 登入服务器，获取token
/**
 * @param
 * @return token
 */
func GetToken(account string,password string) (token string) {

	if account == "" || password == "" {
		account,password = clientInputUser()
	}

	var data = make(map[string]interface{})

	data["phone"] = account
	data["password"] = password

	// 登入服务器
	_, token = Signin(model.Request{"","cine.web",data})

	return token
}

func GetAdminPermission(account string,password string) (token string) {

	if account == "" || password == "" {
		account,password = clientInputUser()
	}

	// 登入服务器
	_, token = SigninPermission(model.RequestPermission{account,password})

	return token
}

func clientInputUser() (account string,password string) {

	for i := 1; i <= 5; i++ {

		if account == "" {

			// 获取输入账户名
			account = ClientInputWithMessage("请输入用户名：",'\n')

			if account == "" {

				continue
			}
		}

		if password == "" {

			password = ClientInputWithMessage("请输入密码：",'\n')

			if password == "" {

				continue
			}
		}

		return account,password

	}

	return "",""

}