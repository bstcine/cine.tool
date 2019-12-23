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

var imageName string = "000.png"
var audioName string = "000.mp3"
var saveName string = "Target"

func main() {

	fmt.Print("git log => c80f")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		return
	}
	fmt.Println(dir)
	componse_workdir := strings.Replace(dir, "\\", "/", -1)

	componseConfig := componse_workdir + "/cine_media_mp4_append.cfg"
	fmt.Println(componseConfig)
	readConfig(componseConfig)

	temporaryPath := componse_workdir + "/temporary"
	if !utils.Exists(temporaryPath) {
		utils.CreatDirectory(temporaryPath)
	}

	prefixPath := componse_workdir + "/PREFIX"
	suffixPath := componse_workdir + "/SUFFIX"

	isHasPrefix := utils.Exists(prefixPath)
	isHasSuffix := utils.Exists(suffixPath)
	if !isHasPrefix && !isHasSuffix {
		fmt.Println("PREFIX、SUFFIX必须至少提供一个")
		return
	}

	var prefixTSPath string
	var suffixTSPath string
	if isHasPrefix {
		prefixTSPath = prefixPath + "/temporary.ts"
		pImagePath := prefixPath + "/" + imageName
		pAudioPath := prefixPath + "/" + audioName
		utils.CreateTsWithImage(pImagePath, pAudioPath, 25, 1920, 1080, prefixTSPath)
	}
	if isHasSuffix {
		suffixTSPath = suffixPath + "/temporary.ts"
		sImagePath := suffixPath + "/" + imageName
		sAudioPath := suffixPath + "/" + audioName
		utils.CreateTsWithImage(sImagePath, sAudioPath, 25, 1920, 1080, suffixTSPath)
	}

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
		if isHasPrefix {
			mpegtsArr = append(mpegtsArr, prefixTSPath)
		}
		mpegtsArr = append(mpegtsArr, tmpTs)
		if isHasSuffix {
			mpegtsArr = append(mpegtsArr, suffixTSPath)
		}
		isSuc = utils.ComponseMP4WithTs(mpegtsArr, savePath)
		if !isSuc {
			fmt.Println("合成失败")
		}
		os.Remove(tmpTs)
	}
	if isHasPrefix {
		os.Remove(prefixTSPath)
	}
	if isHasSuffix {
		os.Remove(suffixTSPath)
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
		if dirName == saveName || dirName == "temporary" || dirName == "SUFFIX" || dirName == "PREFIX" {
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
