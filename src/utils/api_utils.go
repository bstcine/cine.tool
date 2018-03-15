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
)

func GetBaseUrl(url string) string {
	if conf.IsDebug {
		url = conf.API_BASE_URL_TEST + url
	} else {
		url = conf.API_BASE_URL + url
	}

	return url
}

func CommonPost(url string, request model.Request, res interface{}) {
	jsonBytes, _ := json.Marshal(request)

	log.Println("======== 网络请求中 > body: " + string(jsonBytes))

	resp, err := http.Post(GetBaseUrl(url), "application/json", bytes.NewBuffer(jsonBytes))
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

func Signin(request model.Request) (res model.Res, token string) {
	CommonPost(conf.APIURL_Auth_Signin, request, &res)

	return res, res.Result["token"].(string)
}

func GetFiles(password, fileType, cid string) (res map[string]interface{}, data []interface{}) {
	url := conf.APIURL_Tool_Files + "?password=" + password + "&type=" + fileType + "&cid=" + cid

	CommonGet(url, &res)

	return res, res["data"].([]interface{})
}

func ListWithMedias(request model.Request) (res model.ResList, rows []model.Chapter) {
	CommonPost(conf.APIURL_Content_Chapter_ListWithMedia, request, &res)

	rowsJson, _ := json.Marshal(res.Result.Rows)
	json.Unmarshal([]byte(rowsJson), &rows)

	return res, rows
}

func UpdateLessonCheckStatus(request model.Request) (res model.ResCheckList, status bool) {

	url := "http://apptest.bstcine.com/api/tool/content/lesson/checkStatus"

	jsonBytes, _ := json.Marshal(request)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))

	if err != nil {
		println("更新请求出错：", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("解析出错：", err)
	}

	json.Unmarshal(body, &res)

	return res, res.Result["status"]
}
