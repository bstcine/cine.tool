package model

type Request struct {
	Token    string            `json:"token"`
	Sitecode string            `json:"sitecode"`
	Data     map[string]interface{} `json:"data"`
}
