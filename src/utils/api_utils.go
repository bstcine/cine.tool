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

func ListWithMedias(request model.Request) (res model.Response, rows []model.Chapter) {
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
