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

type Media struct {
	Id     string  `json:"id"`
	Type   string  `json:"type"`
	Url    string  `json:"url"`
	Images []Image `json:"images"`
}

type Image struct {
	Time string `json:"time"`
	Url  string `json:"url"`
}
