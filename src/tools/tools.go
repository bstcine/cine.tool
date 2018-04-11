package tools

import (
	"fmt"
	"os"
)

type Tools struct {
	WorkPath string
	ConfDir  string
	ConfMap  map[string]string
}

/**
错误终止
 */
func (tools Tools) HandleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}
