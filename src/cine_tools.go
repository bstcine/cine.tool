package main

import (
	"fmt"
	"./utils"
	"./conf"
	"./tools"
)

func main() {
	var workDir string            //工作目录
	var confFile string           //配置文件
	var confMap map[string]string //配置参数

	if conf.IsDebug {
		workDir = conf.WorkDir
		confFile = conf.ConfFileTmp
	} else {
		workDir = utils.GetCurPath()
		confFile = conf.ConfFile
	}

	confMap = utils.GetConfArgs(workDir + confFile)
	if confMap == nil || len(confMap) <= 0 {
		fmt.Println("配置文件不存在")
		return
	}

	tool := tools.Tools{WorkPath: workDir, ConfDir: confFile, ConfMap: confMap}

	switch confMap["srcType"] {
	case "oss_migrate": //资源迁移
		tool.MigrateObject()
	case "oss_set_acl": //资源授权
		tool.SetObjectACL()
	case "oss_migrate_check": //资源迁移校验
		tool.MigrateCheck()
	case "oss_img_format":
		tool.ImgFormatJPG()
	case "oss_img_watermark":
		tool.ImgWaterMark()
	default:
		fmt.Println("无效参数")
	}

}
