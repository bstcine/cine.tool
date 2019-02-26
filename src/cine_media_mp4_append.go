package main

import (
	"./utils"
	"bytes"
	//"path/filepath"
	//"os"
	//"strings"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var imageName string = "SUFFIX/000.png"
var audioName string = "SUFFIX/000.mp3"
var saveName string = "Target"

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		return
	}
	fmt.Println(dir)
	componse_workdir := strings.Replace(dir, "\\", "/", -1)

	componseConfig := componse_workdir + "/cine_media_mp4_append.cfg"
	fmt.Println(componseConfig)
	readConfig(componseConfig)

	// 创建临时地址
	temporaryPath := componse_workdir + "/temporary"
	if !utils.Exists(temporaryPath) {
		utils.CreatDirectory(temporaryPath)
	}
	appendTsPath := temporaryPath + "/temporary.ts"
	// 将图片合成能为无声的音频
	imagePath := componse_workdir + "/" + imageName
	audioPath := componse_workdir + "/" + audioName
	utils.CreateTsWithImage(imagePath, audioPath, 25, 1920, 1080, appendTsPath)

	// 创建目标地址
	targetDir := componse_workdir + "/" + saveName
	if !utils.Exists(targetDir) {
		utils.CreatDirectory(targetDir)
	}
	// 扫描工作目录
	sourcePaths := scanDirectory(componse_workdir, "", componse_workdir)
	if sourcePaths == nil || len(sourcePaths) == 0 {
		fmt.Println("没有原始视频文件，不能执行添加功能")
		return
	}
	fmt.Println(sourcePaths)

	for _, mp4Name := range sourcePaths {

		tmpTs := componse_workdir + "/temporary/" + mp4Name
		tmpTs = strings.Replace(tmpTs, ".mp4", ".ts", -1)

		isSuc := utils.CreateTsWithMp4(componse_workdir+"/"+mp4Name, tmpTs)
		if !isSuc {
			fmt.Println("转换失败，", mp4Name)
			continue
		}
		savePath := componse_workdir + "/" + saveName + "/" + mp4Name
		var mpegtsArr []string
		mpegtsArr = append(mpegtsArr, tmpTs)
		mpegtsArr = append(mpegtsArr, appendTsPath)
		isSuc = utils.ComponseMP4WithTs(mpegtsArr, savePath)
		if !isSuc {
			fmt.Println("合成失败")
		}
		os.Remove(tmpTs)
	}
	os.RemoveAll(temporaryPath)

}

func scanDirectory(dirPath string, name string, componse_workdir string) []string {

	var sourcePaths []string

	fmt.Println("")
	// 获取目录中的所有目录，递归访问
	dirNames, mp4Names := utils.GetAllMp4FileeNames(dirPath)

	// 将所有MP4文件拷贝到数组中
	if mp4Names != nil && len(mp4Names) > 0 {
		for _, mp4Name := range mp4Names {
			if name == "" {
				sourcePaths = append(sourcePaths, mp4Name)
			} else {
				sourcePaths = append(sourcePaths, name+"/"+mp4Name)
			}
		}
	}

	if len(dirNames) == 0 {
		return sourcePaths
	}
	cachePath := strings.Replace(dirPath, componse_workdir, "", -1)
	cachePath = strings.Replace(cachePath, " ", "", -1)

	for _, dirName := range dirNames {
		if dirName == saveName || dirName == "temporary" || dirName == "SUFFIX" {
			continue
		}
		var newPathBuffer bytes.Buffer
		var newTmpPathBuffer bytes.Buffer

		if cachePath == "" {
			newPathBuffer.WriteString(componse_workdir)
			newPathBuffer.WriteString("/")
			newPathBuffer.WriteString(saveName)
			fmt.Println("cache1: ", newPathBuffer.String())
			newPathBuffer.WriteString("/")
			newPathBuffer.WriteString(dirName)

			newTmpPathBuffer.WriteString(componse_workdir)
			newTmpPathBuffer.WriteString("/temporary/")
			newTmpPathBuffer.WriteString(dirName)
		} else {

			newPathBuffer.WriteString(componse_workdir)
			newPathBuffer.WriteString("/")
			newPathBuffer.WriteString(saveName)
			fmt.Println("cache2: ", newPathBuffer.String())
			fmt.Println("dizi", componse_workdir)
			newPathBuffer.WriteString("/")
			newPathBuffer.WriteString(cachePath)
			newPathBuffer.WriteString("/")
			newPathBuffer.WriteString(dirName)

			newTmpPathBuffer.WriteString(componse_workdir)
			newTmpPathBuffer.WriteString("/temporary/")
			newTmpPathBuffer.WriteString(cachePath)
			newTmpPathBuffer.WriteString("/")
			newTmpPathBuffer.WriteString(dirName)
		}
		var newPath string = newPathBuffer.String()
		var newTmpPath string = newTmpPathBuffer.String()
		if !utils.Exists(newPath) {
			utils.CreatDirectory(newPath)
		}
		if !utils.Exists(newTmpPath) {
			utils.CreatDirectory(newTmpPath)
		}
		fmt.Println("cachePath: ", cachePath, componse_workdir)
		fmt.Println(newPath)
		fmt.Println(newTmpPath)
		nextSourcePaths := scanDirectory(dirPath+"/"+dirName, dirName, componse_workdir)
		sourcePaths = append(sourcePaths, nextSourcePaths...)
	}
	return sourcePaths
}

func readConfig(path string) {

	configArg := utils.GetConfArgs(path)

	if configArg == nil {
		return
	}
	imageName = configArg["imageName"]
}
