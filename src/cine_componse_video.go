package main

import (
	"./utils"
	"fmt"
	"strings"
	"strconv"
	"os"
	"path/filepath"
	"./conf"
)

type componseVideo struct {
	LessonDir string
	SavePath string
	TmpPath string
	Audios []componseAudio
}

type componseAudio struct {
	Seq int
	OriginPath string
	SavePath string
	Images []componseImage
}

type componseImage struct {
	Seq int
	duration int
	OriginPath string
	TmpPath string
}


func main() {

	var componse_workdir = ""

	if conf.IsDebug {
		componse_workdir = "/Users/lidangkun/Desktop/oss_download"
	}else {

		dir,err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			return
		}

		componse_workdir = strings.Replace(dir,"\\","/",-1)
	}

	fmt.Println(componse_workdir)
	// 生成保存文件夹MP4
	savePath := componse_workdir + "/MP4"

	dealDirectory(componse_workdir,savePath)
}

func dealDirectory(dirPath string,savePath string) {

	// 判断是否是LessonDirectory
	dirComponents := strings.Split(dirPath,"/")
	dirName := dirComponents[len(dirComponents)-1]
	if strings.Contains(dirName,"ls_") {
		fmt.Println("获取lesson目录，开始读取视频及图片",dirPath)
		videoModel := dealLessonDirectory(dirPath,savePath)
		fmt.Println("lesson目录读取成功，开始处理文件",dirPath)
		startComponseVideoModel(videoModel)
		fmt.Println("lesson目录处理结束",dirPath)
	}else {
		// 获取目录中的所有目录，递归访问
		dirNames := utils.GetAllDirectoryNames(dirPath)

		if len(dirNames) == 0 {

			return
		}

		for _,dirName := range dirNames {

			dealDirectory(dirPath+"/"+dirName,savePath+"/"+dirName)
		}
	}
}

func startComponseVideoModel(videoModel componseVideo){

	// 判断是否已存在合成视频
	if utils.Exists(videoModel.SavePath) {
		fmt.Println("该视频已存在，不需要合成",videoModel.SavePath)
		return
	}

	if len(videoModel.Audios) == 0 {
		fmt.Println("该目录下没有需要合成的音频")
		return
	}

	// 创建tmpDirctory
	if !utils.Exists(videoModel.TmpPath) {
		utils.CreatDirectory(videoModel.TmpPath)
	}

	if len(videoModel.Audios) == 1 {

		// 只有一个音频组
		tmpPath := startComponseAudioToVideo(videoModel,videoModel.Audios[0])

		// 拷贝视频到保存位置
		fmt.Println(tmpPath,videoModel.SavePath)
		err := os.Rename(tmpPath,videoModel.SavePath)

		if err != nil {
			fmt.Println(err)
		}else {
			os.RemoveAll(videoModel.TmpPath)
		}
		return
	}

	// 创建临时合并文件组

	var mpegtsArr []string

	for _,audioModel := range videoModel.Audios {

		tmpPath := startComponseAudioToVideo(videoModel,audioModel)

		// 生成可拼接文件类型
		mpegtsPath := strings.Replace(tmpPath,".mp4","",-1) + ".ts"

		mpegtsLine := utils.CreatVideoToMpegtsLines(tmpPath,mpegtsPath)

		_,err := utils.RunCMD(mpegtsLine)

		if err != nil {
			fmt.Println("合并临时视频失败",videoModel.LessonDir)
			return
		}

		mpegtsArr = append(mpegtsArr,mpegtsPath)
	}

	// 将临时mpegts文件合并为目标文件,由于视频文件未知特殊性，需要将音视频分离后再封装起来
	// 如果不执行分离后再封装操作，将可能丢失视频预览图片
	fmt.Println("开始合并临时文件")
	// 生成临时video
	componseVideoPath := videoModel.TmpPath + "/video.mp4"
	componseVideoLine := utils.CreatComponseMpegtsLines(mpegtsArr,componseVideoPath,"-an")
	utils.RunCMD(componseVideoLine)
	// 生成临时audio
	componseAudioPath := videoModel.TmpPath + "/audio.mp4"
	componseAudioLine := utils.CreatComponseMpegtsLines(mpegtsArr,componseAudioPath,"-vn")
	utils.RunCMD(componseAudioLine)
	// 合并临时video和audio
	componseLine := "ffmpeg -i \"" + componseVideoPath + "\" -i \"" + componseAudioPath + "\" \"" + videoModel.SavePath + "\""
	_,err := utils.RunCMD(componseLine)
	if err != nil {
		fmt.Println("合并视频失败",videoModel.SavePath)
	}
	fmt.Println("视频合并成功",videoModel.SavePath)
	os.RemoveAll(videoModel.TmpPath)
}

func startComponseAudioToVideo(videoModel componseVideo,audioModel componseAudio) (tmpVideoPath string) {

	// 如果没有图片，则表示文件夹配置异常
	if len(audioModel.Images) == 0 {

		fmt.Println("序号为",audioModel.Seq,"的音频没有对应图片，已结束合成")
		return tmpVideoPath
	}

	// 生成临时路径
	tmpVideoPath = videoModel.TmpPath + "/" + utils.ChangeIntToThirdStr(audioModel.Seq) + ".mp4"

	// 如果只有一张图片，则直接合成一个视频
	if len(audioModel.Images) == 1 {

		fmt.Println("开始处理音频与图片")
		// 生成合成命令行
		componseLine := utils.CreatOneImageAudioLines(audioModel.Images[0].OriginPath,audioModel.OriginPath,tmpVideoPath)
		// 执行命令行
		_,status := utils.RunCMD(componseLine)

		if status != nil {
			tmpVideoPath = ""
		}

		fmt.Println("处理完成，已生成临时文件,请勿操作计算机，以免造成误删除")
		return tmpVideoPath
	}

	// 如果存在多张图片，则需要将多张图片分成合成视频后拼接，最后与音频文件合并，生成最终的临时视频
	// 生成audio临时目录，保存临时文件，合成成功后，移除临时目录
	audioTmpDir := videoModel.TmpPath + "/" + utils.ChangeIntToThirdStr(audioModel.Seq) + "_tmp"
	utils.CreatDirectory(audioTmpDir)

	// 创建中间文件数组
	var mpegtsArr []string

	var imagesChangeSuc = true

	for _,imageModel := range audioModel.Images  {
		// 生成指定时长的视频
		videoPath := audioTmpDir + "/" + utils.ChangeIntToThirdStr(imageModel.Seq) + ".mp4"
		mpegtspATH := audioTmpDir + "/" + utils.ChangeIntToThirdStr(imageModel.Seq) + ".ts"

		fmt.Println("开始处理图片",imageModel.OriginPath)
		creatVideoLine := utils.CreatLines(imageModel.OriginPath,imageModel.duration,videoPath)
		_,status := utils.RunCMD(creatVideoLine)
		if status != nil {
			fmt.Println("图片转换视频失败，本视频合成结束",creatVideoLine)
			imagesChangeSuc = false
			break
		}

		fmt.Println("图片开始转换视频文件",imageModel.OriginPath)
		mpegsLine := utils.CreatVideoToMpegtsLines(videoPath,mpegtspATH)
		_,status = utils.RunCMD(mpegsLine)
		if status != nil {
			fmt.Println("视频转换中间文件失败，本视频合成失败",mpegsLine)
			imagesChangeSuc = false
			break
		}

		fmt.Println("转换视频文件成功",imageModel.OriginPath)
		mpegtsArr = append(mpegtsArr,mpegtspATH)
	}

	if !imagesChangeSuc {
		return ""
	}

	// 合成一个componseVideo
	fmt.Println("开始拼接视频源 tag:",audioModel.Seq)
	componsePath := audioTmpDir + "/" + "componse.mp4"
	componseVideoLine := utils.CreatComponseMpegtsLines(mpegtsArr,componsePath,"")
	_,status := utils.RunCMD(componseVideoLine)
	if status != nil {
		fmt.Println("临时视频文件合并失败，本视频结束合成")
		return ""
	}
	fmt.Println("视频源拼接完成，开始导入音频源")
	// 将componseVideo添加audio文件
	tmpVideoLine := utils.CreatAudioToVideoLines(audioModel.OriginPath,componsePath,tmpVideoPath)
	_,status = utils.RunCMD(tmpVideoLine)
	if status != nil {
		fmt.Println("临时视频添加音频失败。本视频结束合成")
		return ""
	}
	fmt.Println("已合成临时视频文件",tmpVideoPath)
	// 合成成功，移除临时文件
	os.RemoveAll(audioTmpDir)

	return tmpVideoPath
}

func dealLessonDirectory(dirPath string,saveDirPath string) (videoModel componseVideo) {

	allMp3Names := utils.GetAllMp3FiloeNames(dirPath)

	if len(allMp3Names) == 0 {
		return videoModel
	}

	//
	allFileNames := utils.GetAllFiloeNames(dirPath)
	utils.CreatDirectory(saveDirPath)

	lessonComponents := strings.Split(dirPath,"/")
	lessonName := lessonComponents[len(lessonComponents)-1]

	videoModel = componseVideo{
		LessonDir:dirPath,
		SavePath:saveDirPath+"/"+lessonName+".mp4",
		TmpPath:dirPath+"/"+"tmp_cine",
	}

	var audioModels []componseAudio

	// 生成临时文件夹，用来防止过渡文件

	for _,audioName := range allMp3Names {

		audioSeqStr := strings.Replace(audioName,".mp3","",-1)
		audioSeq,_ := strconv.Atoi(audioSeqStr)

		audioModel := componseAudio{
			Seq:audioSeq,
			OriginPath:dirPath+"/"+audioName,
		}

		var audioImages []componseImage

		for _,fileName := range allFileNames {

			preName := audioSeqStr + "_"
			if !strings.Contains(fileName,preName) {
				continue
			}

			fileSeqStr := strings.Replace(fileName,preName,"",-1)
			fileSeqStr = strings.Replace(fileSeqStr,".jpg","",-1)
			fileSeq,_ := strconv.Atoi(fileSeqStr)

			imageModel := componseImage{
				Seq:fileSeq,
				OriginPath:dirPath+"/"+fileName,
			}

			if len(audioImages) == 0 {

				audioImages = append(audioImages,imageModel)
			}else {

				var hadInsert = false

				for i := 0;i < len(audioImages) ;i++  {
					if audioImages[i].Seq <= imageModel.Seq {
						continue
					}
					// 插入当前位置，并break
					var preImages []componseImage
					preImages = append(preImages,audioImages[:i]...)

					lastImages := audioImages[i:]

					preImages = append(preImages,imageModel)
					audioImages = append(preImages,lastImages...)

					hadInsert = true
					break
				}
				if !hadInsert {
					audioImages = append(audioImages,imageModel)
				}
			}

		}

		audioModel.Images = audioImages

		if len(audioModels) == 0 {
			audioModels = append(audioModels,audioModel)
		}else {

			var hadInsert = false

			for index,oldAudio := range audioModels {
				if oldAudio.Seq <= audioModel.Seq {
					continue
				}
				var preAudios []componseAudio
				preAudios = append(preAudios,audioModels[:index]...)
				lastAudios := audioModels[index:]
				preAudios = append(preAudios, audioModel)
				audioModels = append(preAudios,lastAudios...)

				hadInsert = true
				break
			}

			if !hadInsert {
				audioModels = append(audioModels,audioModel)
			}
		}
	}

	// 遍历音频数组，为每一个音频对应的图片设置持续时长，如果只有一张图片，则不需要设置时长
	for _,audioModel := range audioModels {

		if len(audioModel.Images) <= 1 {
			continue
		}

		count := len(audioModel.Images)

		for index,_ := range audioModel.Images {
			if index == count - 1 {
				// 获取音频时长
				audioDuration := utils.GetDuration(audioModel.OriginPath)

				if audioDuration == 0 {
					fmt.Println("音频异常，没有合适的时长，",audioModel.OriginPath)
					return videoModel
				}

				audioModel.Images[index].duration = audioDuration - audioModel.Images[index].Seq

			}else {

				audioModel.Images[index].duration = audioModel.Images[index + 1].Seq - audioModel.Images[index].Seq
			}

		}

	}

	videoModel.Audios = audioModels

	return videoModel
}

