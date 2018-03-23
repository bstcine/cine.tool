package tools

import (
	"fmt"
	"os"
	"log"
	"strings"
	"../conf"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Tools struct {
	WorkPath string
	ConfDir  string
	ConfMap  map[string]string
	OSSClient *oss.Client
	OSSBucket *oss.Bucket
}

var logger *log.Logger

/**
错误终止
 */
func (tools Tools) HandleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

/**
获取日志器
 */
func (tools Tools) GetLogger() *log.Logger {
	if logger == nil {
		outPath := tools.WorkPath + conf.LogFile

		dirPath := outPath[0:strings.LastIndex(outPath, "/")+1]
		if _, err := os.Stat(dirPath); err != nil {
			os.MkdirAll(dirPath, 0777)
		}

		logFile, err := os.Create(outPath)

		if err != nil {
			log.Fatalln("open file error !")
		}

		logger = log.New(logFile, "[Info]", log.Llongfile)
	}

	return logger
}
