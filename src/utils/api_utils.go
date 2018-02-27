package utils

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"../model"
)

const API_BASE_URL = "http://www.bstcine.com"

func Signin(request model.Request)(res model.Res,token string){
	url := API_BASE_URL + "/api/auth/signin"

	jsonBytes, _ := json.Marshal(request)

	fmt.Println("======== 网络请求中 > body: " + string(jsonBytes))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("======== 网络请求成功 > body: " + string(body))

	json.Unmarshal(body, &res)

	return res, res.Result["token"].(string)
}

func GetFiles(password,fileType,cid string) (res map[string] interface{},data []interface{}) {
	url := API_BASE_URL + "/api/tool/files?password="+password+"&type="+fileType+"&cid="+cid

	fmt.Println("======== 网络请求中 > ")

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("======== 网络请求成功 > body: " + string(body))

	json.Unmarshal(body, &res)

	return res, res["data"].([]interface{})
}

func ListWithMedias(request model.Request) (res model.ResList, rows []model.Chapter) {
	url := API_BASE_URL + "/api/content/chapter/listWithMedia"

	jsonBytes, _ := json.Marshal(request)

	fmt.Println("======== 网络请求中 > body: " + string(jsonBytes))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("======== 网络请求成功 > body: " + string(body))

	json.Unmarshal(body, &res)

	rowsJson, _ := json.Marshal(res.Result.Rows)
	json.Unmarshal([]byte(rowsJson), &rows)

	return res, rows
}
