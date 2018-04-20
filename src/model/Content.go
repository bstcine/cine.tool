package model

type Course struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type Chapter struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Children []Lesson `json:"children"`
}

type Lesson struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Medias []Media `json:"medias"`
}

type CheckLesson struct {
	ChapterId   string `json:"chapter_id"`
	ChapterName string `json:"chapter_name"`
	Id          string  `json:"lesson_id"`
	Name        string  `json:"lesson_name"`
	Medias      []CheckMedia `json:"medias"`
}

type CheckMedia struct {
	Seq int `json:"seq"`
	CheckStatus int `json:"check_status"`
	Url string `json:"url"`
	Images []Image `json:"images"`
}

type Media struct {
	Id     string  `json:"id"`
	Seq    int     `json:"seq"`
	Type   string  `json:"type"`
	Url    string  `json:"url"`
	Images []Image `json:"images"`
}

type Image struct {
	Time string `json:"time"`
	Url  string `json:"url"`
}
