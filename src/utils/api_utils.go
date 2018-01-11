package utils

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Request struct {
	Token string `json:"token"`
	Sitecode string `json:"sitecode"`
	Data map[string]string `json:"data"`
}

type Response struct {
	Code string `json:"code"`
	Code_desc string `json:"code_desc"`
	Except_case string `json:"except_case"`
	Except_case_desc string `json:"except_case_desc"`
	Result ResultList `json:"result"`
}

type ResultList struct {
	Cur_page string `json:"cur_page"`
	Max_page string `json:"max_page"`
	Total_rows string `json:"total_rows"`
	Rows []string `json:"rows"`
}

type ResListWithMedia struct {
	Code string `json:"code"`
	Code_desc string `json:"code_desc"`
	Except_case string `json:"except_case"`
	Except_case_desc string `json:"except_case_desc"`
	Result ResultListWithMedia `json:"result"`
}

type ResultListWithMedia struct {
	Cur_page string `json:"cur_page"`
	Max_page string `json:"max_page"`
	Total_rows string `json:"total_rows"`
	Rows []Chapter `json:"rows"`
}


type Chapter struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Children []Lesson `json:"children"`
}

type Lesson struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Medias []Media `json:"medias"`
}

type Media struct {
	Id string `json:"id"`
	Type string `json:"type"`
	Url string `json:"url"`
	Images []Image `json:"images"`
}

type Image struct {
	Time string `json:"time"`
	Url string `json:"url"`
}

type DownFile struct {
	Path   string
	ChapterName   string
	LessonName   string
	Name   string
}

const API_BASE_URL string = "http://www.bstcine.com"

func ListWithMedias(request Request) ResListWithMedia {
	url := API_BASE_URL + "/api/content/chapter/listWithMedia"

	jsonBytes, _ := json.Marshal(request)

	fmt.Println(string(jsonBytes))

	resp, err := http.Post(url,"application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var res ResListWithMedia
	json.Unmarshal(body, &res)
	return res
}