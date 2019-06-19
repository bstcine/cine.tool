package main

import (
	"./conf"
	"./model"
	"./utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type componseVideo struct {
	LessonDir      string
	DirName        string
	SavePath       string
	HadOriginVideo bool
	TmpPath        string
	Audios         []componseAudio
}

type componseAudio struct {
	Seq        int
	OriginPath string
	SavePath   string
	Images     []componseImage
}

type componseImage struct {
	Seq        int
	duration   int
	OriginPath string
	TmpPath    string
}

var mediaConfgiModel model.MediaConfig

func main() {

	var componse_workdir = ""
	var mediaSynthesizerConfig string

	if conf.IsDebug {
		componse_workdir = "/Users/lidangkun/Desktop/oss_download"
	} else {

		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			return
		}

		componse_workdir = strings.Replace(dir, "\\", "/", -1)
	}

	mediaSynthesizerConfig = componse_workdir + "/cine_media_synthesizer.cfg"

	// 读取配置参数
	readResult := readMediaSynthesuzerConfig(mediaSynthesizerConfig)

	if !readResult {
		return
	}
	var prefixAdTmpTs string
	if mediaConfgiModel.IsAddPrefix {
		fmt.Println("是否在头部追加推广片段 true")
		// 合成追加的ts,放置到临时目录中，待使用
		adDir := componse_workdir + "/PREFIX/tmp"
		if !utils.Exists(adDir) {
			utils.CreatDirectory(adDir)
		}
		prefixAdTmpTs = adDir + "/tmp.ts"

		imagePath := componse_workdir + "/PREFIX/000.png"
		audioPath := componse_workdir + "/PREFIX/000.mp3"

		res := utils.CreateTsWithImageAudio(imagePath, audioPath, prefixAdTmpTs, mediaConfgiModel)
		if !res {
			fmt.Println("前置推广片段 ts 文件合成失败")
			return
		}
		fmt.Println("前置推广图 ts 合成完毕")
	}

	var suffixAdTmpTs string
	if mediaConfgiModel.IsAddSuffix {
		fmt.Println("是否在结尾追加推广片段 true")
		// 合成追加的ts,放置到临时目录中，待使用
		adDir := componse_workdir + "/SUFFIX/tmp"
		if !utils.Exists(adDir) {
			utils.CreatDirectory(adDir)
		}
		suffixAdTmpTs = adDir + "/tmp.ts"

		imagePath := componse_workdir + "/SUFFIX/000.png"
		audioPath := componse_workdir + "/SUFFIX/000.mp3"

		res := utils.CreateTsWithImageAudio(imagePath, audioPath, suffixAdTmpTs, mediaConfgiModel)
		if !res {
			fmt.Println("后置推广片段 ts 文件合成失败")
			return
		}
		fmt.Println("后置推广图 ts 合成完毕")
	}

	// 生成保存文件夹MP4
	savePath := componse_workdir + "/MP4"

	// 处理工作目录
	dealDirectory(componse_workdir, savePath, prefixAdTmpTs, suffixAdTmpTs)
}

func dealDirectory(dirPath string, savePath string, prefixAdTmpTs string, suffixAdTmpTs string) {

	// 判断是否是LessonDirectory
	dirComponents := strings.Split(dirPath, "/")
	dirName := dirComponents[len(dirComponents)-1]
	if strings.Contains(dirName, "ls_") {
		videoModel := dealLessonDirectory(dirPath, savePath)
		fmt.Println("一个课件扫描完成", videoModel.LessonDir)
		for _, audioModel := range videoModel.Audios {
			fmt.Printf("%+v\n", audioModel)
		}

		startComponseVideoModel(videoModel, prefixAdTmpTs, suffixAdTmpTs)
		fmt.Println("lesson目录处理结束", dirPath)
	} else {
		// 获取目录中的所有目录，递归访问
		dirNames := utils.GetAllDirectoryNames(dirPath)
		fmt.Println("获取目录内容: \n", dirNames)

		if len(dirNames) == 0 {
			fmt.Println("文件内容为空, 开始下一个")
			return
		}

		for _, dirName := range dirNames {

			if dirName == "MP4" || dirName == "SUFFIX" || dirName == "PREFIX" {
				continue
			}
			dealDirectory(dirPath+"/"+dirName, savePath+"/"+dirName, prefixAdTmpTs, suffixAdTmpTs)
		}
	}
}

/// 开始合成video模型
/**
@param videomodel 合成视频模型
*/
func startComponseVideoModel(videoModel componseVideo, prefixAdTmpTs string, suffixAdTmpTs string) bool {

	var saveSuffix string
	if mediaConfgiModel.IsTs {
		saveSuffix = ".ts"
	} else {
		saveSuffix = ".mp4"
	}

	componsePath := videoModel.SavePath + "/" + videoModel.DirName + saveSuffix

	// 判断是否已存在合成视频
	if mediaConfgiModel.UseOldFile && utils.Exists(componsePath) {

		fmt.Println("该视频已存在，不需要合成", componsePath)
		return true
	}

	// 判断源文件是否存在
	if len(videoModel.Audios) == 0 {
		fmt.Println("该目录下没有需要合成的音频")
		return true
	}

	// 创建tmpDirctory，保存临时合成文件
	if !utils.Exists(videoModel.TmpPath) {
		utils.CreatDirectory(videoModel.TmpPath)
	}

	fmt.Println("开始处理目录: ", videoModel.LessonDir)

	// 生成临时文件路径
	tmpAVPath := videoModel.TmpPath + "/" + videoModel.DirName + ".ts"
	tmpPath := videoModel.TmpPath + "/" + videoModel.DirName + ".mp4"
	saveAVPath := videoModel.SavePath + "/" + videoModel.DirName + ".ts"
	savePath := videoModel.SavePath + "/" + videoModel.DirName + ".mp4"

	// 如果只有一个音频，则直接采用一个合成命令行
	if len(videoModel.Audios) == 1 && !mediaConfgiModel.IsAddSuffix && !mediaConfgiModel.IsAddPrefix {

		// 只有一个音频组
		avTs := startComponseAudioToVideo(videoModel, videoModel.Audios[0])

		// 移动ts文件到指定位置
		os.Rename(avTs, tmpAVPath)

		if mediaConfgiModel.IsTs {

			fmt.Println("源数据合成完毕!")
			err := os.Rename(tmpAVPath, saveAVPath)

			if err != nil {
				return false
			}

			if !mediaConfgiModel.SaveTmp {
				os.RemoveAll(videoModel.TmpPath)
			}
			return true
		} else {
			isSuc := utils.CreatMp4WithMpegts(tmpAVPath, tmpPath, mediaConfgiModel)
			return isSuc
		}
	}

	// 表示有多个音频需要合成, 创建临时合并文件组
	var mpegtsAVArr []string

	if prefixAdTmpTs != "" {
		mpegtsAVArr = append(mpegtsAVArr, prefixAdTmpTs)
	}

	audioCount := len(videoModel.Audios)
	for index, audioModel := range videoModel.Audios {

		tmsTs := startComponseAudioToVideo(videoModel, audioModel)

		if tmsTs == "" {
			fmt.Println("标准音视频源获取失败")
			return false
		}

		fmt.Println("lesson 音视频源已准备", index+1, "/", audioCount)

		mpegtsAVArr = append(mpegtsAVArr, tmsTs)
	}

	if suffixAdTmpTs != "" {
		mpegtsAVArr = append(mpegtsAVArr, suffixAdTmpTs)
	}

	fmt.Println("开始合并标准多媒体数据源...")
	isSuc := utils.ComponseMpegts(mpegtsAVArr, tmpAVPath, mediaConfgiModel)
	if !isSuc {
		fmt.Println("数据源合并失败！")
		return false
	}

	if mediaConfgiModel.IsTs {
		fmt.Println("源数据合成完毕!")

		err := os.Rename(tmpAVPath, saveAVPath)

		if err != nil {
			return false
		}

		if !mediaConfgiModel.SaveTmp {
			os.RemoveAll(videoModel.TmpPath)
		}
		return true
	}

	// 执行转码
	fmt.Println("源数据合成完毕,开始转码...")
	isSuc = utils.CreatMp4WithMpegts(tmpAVPath, tmpPath, mediaConfgiModel)

	if !isSuc {
		fmt.Println("转码失败")
		return false
	}

	fmt.Println("视频转码成功", videoModel.SavePath)

	err := os.Rename(tmpPath, savePath)

	if err != nil {
		return false
	}

	// 判断是否保存临时文件
	if !mediaConfgiModel.SaveTmp {
		os.RemoveAll(videoModel.TmpPath)
	}

	return true
}

/// 获取一个tag的临时ts数据源，（每个tag包含一个数据源）
/**
@param videoModel 视频数据模型
@param audioModel 待处理的音频模型
@return tmpVideoPath 临时视频源路径
@retutn tmpAudioPath 临时音频源路径
*/
func startComponseAudioToVideo(videoModel componseVideo, audioModel componseAudio) (tmpTS string) {

	// 生成临时路径
	//tmpVideoPath := videoModel.TmpPath + "/" + utils.ChangeIntToThirdStr(audioModel.Seq) + "_video.ts"
	//tmpAudioPath := videoModel.TmpPath + "/" + utils.ChangeIntToThirdStr(audioModel.Seq) + "_audio.ts"
	tmpTS = videoModel.TmpPath + "/" + utils.ChangeIntToThirdStr(audioModel.Seq) + ".ts"

	// 判断临时文件是否已经存在
	if mediaConfgiModel.UseOldFile && utils.Exists(tmpTS) {
		fmt.Println("ts文件已存在，不需要执行合并操作")
		return tmpTS
	}

	// 判断这个音频模式是否是个视频
	if strings.Contains(audioModel.OriginPath, ".mp4") {

		// 将视频MP4转换为ts文件
		isSuc := utils.CreatVideoMpegtsWithMP4(audioModel.OriginPath, tmpTS, mediaConfgiModel)
		if !isSuc {
			return ""
		}
		fmt.Println("源视频处理完成，已生成标准多媒体数据文件，请勿操作计算机，以免造成误删!")
		return tmpTS
	}

	// 如果没有图片，则表示文件夹配置异常
	if len(audioModel.Images) == 0 {

		fmt.Println("序号为", audioModel.Seq, "的音频没有对应图片，已结束合成")
		return ""
	}

	// 开始检测源文件
	fmt.Println(audioModel.OriginPath, "\n", audioModel.Images)

	// 如果只有一张图片，则直接合成一个视频
	if len(audioModel.Images) == 1 {

		// 一张图片+音频合成ts文件
		isSuc := utils.CreateTsWithImageAudio(audioModel.Images[0].OriginPath, audioModel.OriginPath, tmpTS, mediaConfgiModel)
		if !isSuc {
			return ""
		}
		fmt.Println("合成成功，已生成标准多媒体数据文件，请勿操作计算机，以免造成误删!")

		return tmpTS
	}

	// 如果存在多张图片，则需要将多张图片分成合成视频后拼接，最后与音频文件合并，生成最终的临时视频
	// 生成audio临时目录，保存临时文件，合成成功后，移除临时目录
	audioTmpDir := videoModel.TmpPath + "/" + utils.ChangeIntToThirdStr(audioModel.Seq) + "_tmp"
	utils.CreatDirectory(audioTmpDir)

	var imageArrs []map[string]string
	var tmpDir = audioTmpDir
	for _, imageModel := range audioModel.Images {
		imageMap := make(map[string]string)
		imageMap["path"] = imageModel.OriginPath
		imageMap["duration"] = strconv.Itoa(imageModel.duration)
		imageArrs = append(imageArrs, imageMap)
	}
	isSuc := utils.CreateMpegtsWithImagesAudio(imageArrs, tmpDir, audioModel.OriginPath, tmpTS, mediaConfgiModel)
	if !isSuc {
		return ""
	}

	fmt.Println("合成成功，已生成标准多媒体数据文件，请勿操作计算机，以免造成误删!")

	// 合成成功，判断是否移除临时文件
	if !mediaConfgiModel.SaveTmp {
		os.RemoveAll(audioTmpDir)
	}

	return tmpTS
}

func dealLessonDirectory(dirPath string, saveDirPath string) (videoModel componseVideo) {

	fmt.Println("开始扫描lesson 目录", dirPath, "\n", saveDirPath)

	hadVideo, allMp3Names := utils.GetAllMp3FiloeNames(dirPath)

	if len(allMp3Names) == 0 {
		return videoModel
	}

	//
	allFileNames := utils.GetAllFiloeNames(dirPath)

	fmt.Println("所有文件名称", allFileNames)
	utils.CreatDirectory(saveDirPath)

	lessonComponents := strings.Split(dirPath, "/")
	lessonName := lessonComponents[len(lessonComponents)-1]

	videoModel = componseVideo{
		LessonDir:      dirPath,
		HadOriginVideo: hadVideo,
		SavePath:       saveDirPath,
		TmpPath:        dirPath + "/" + "tmp_cine",
		DirName:        lessonName,
	}

	var audioModels []componseAudio

	// 生成临时文件夹，用来放置过渡文件

	for _, audioName := range allMp3Names {

		audioSeqStr := strings.Replace(audioName, ".mp3", "", -1)
		audioSeqStr = strings.Replace(audioSeqStr, ".mp4", "", -1)
		audioSeq, _ := strconv.Atoi(audioSeqStr)

		audioModel := componseAudio{
			Seq:        audioSeq,
			OriginPath: dirPath + "/" + audioName,
		}

		if !strings.Contains(audioName, ".mp4") {

			var audioImages []componseImage

			for _, fileName := range allFileNames {

				preName := audioSeqStr + "_"
				if !strings.Contains(fileName, preName) {
					continue
				}

				fileSeqStr := strings.Replace(fileName, preName, "", -1)
				fileSeqStr = strings.Replace(fileSeqStr, ".jpg", "", -1)
				fileSeq, _ := strconv.Atoi(fileSeqStr)

				imageModel := componseImage{
					Seq:        fileSeq,
					OriginPath: dirPath + "/" + fileName,
				}

				if len(audioImages) == 0 {

					audioImages = append(audioImages, imageModel)
				} else {

					var hadInsert = false

					for i := 0; i < len(audioImages); i++ {
						if audioImages[i].Seq <= imageModel.Seq {
							continue
						}
						// 插入当前位置，并break
						var preImages []componseImage
						preImages = append(preImages, audioImages[:i]...)

						lastImages := audioImages[i:]

						preImages = append(preImages, imageModel)
						audioImages = append(preImages, lastImages...)

						hadInsert = true
						break
					}
					if !hadInsert {
						audioImages = append(audioImages, imageModel)
					}
				}

			}

			audioModel.Images = audioImages
		}

		if len(audioModels) == 0 {
			audioModels = append(audioModels, audioModel)
		} else {

			var hadInsert = false

			for index, oldAudio := range audioModels {
				if oldAudio.Seq <= audioModel.Seq {
					continue
				}
				var preAudios []componseAudio
				preAudios = append(preAudios, audioModels[:index]...)
				lastAudios := audioModels[index:]
				preAudios = append(preAudios, audioModel)
				audioModels = append(preAudios, lastAudios...)

				hadInsert = true
				break
			}

			if !hadInsert {
				audioModels = append(audioModels, audioModel)
			}
		}
	}

	// 遍历音频数组，为每一个音频对应的图片设置持续时长，如果只有一张图片，则不需要设置时长
	for _, audioModel := range audioModels {

		if len(audioModel.Images) <= 1 {
			continue
		}

		count := len(audioModel.Images)

		for index, _ := range audioModel.Images {
			if index == count-1 {
				// 获取音频时长
				audioDuration := utils.GetDuration(audioModel.OriginPath)

				if audioDuration == 0 {
					fmt.Println("音频异常，没有合适的时长，", audioModel.OriginPath)
					return videoModel
				}

				audioModel.Images[index].duration = audioDuration - audioModel.Images[index].Seq

			} else {

				audioModel.Images[index].duration = audioModel.Images[index+1].Seq - audioModel.Images[index].Seq
			}

		}

	}

	videoModel.Audios = audioModels

	return videoModel
}

/// 读取多媒体配置信息
func readMediaSynthesuzerConfig(path string) bool {

	configModel := model.MediaConfig{
		false,
		false,
		false,
		1,
		"1920*10810",
		1920,
		1080,
		1920.0 / 1080.0,
		"baseline",
		"3.0",
		"yuv420p",
		false,
		false,
	}

	fmt.Println("开始读取配置文件: ", path)

	configArg := utils.GetConfArgs(path)

	if configArg == nil {
		fmt.Println("配置文件异常，读取失败，启用默认配置：\n", configModel)
	} else {

		// 清理字典元素的空格部分
		utils.ClearDictionaryChar(configArg, " ")
		if configArg["saveTmp"] == "true" {
			configModel.SaveTmp = true
		}
		if configArg["useOldFile"] == "true" {
			configModel.UseOldFile = true
		}
		if configArg["isTs"] == "true" {
			configModel.IsTs = true
		}
		if configArg["isAddSuffix"] == "true" {
			configModel.IsAddSuffix = true
		}
		if configArg["isAddPrefix"] == "true" {
			configModel.IsAddPrefix = true
		}
		rate, isInt := utils.JudgeIsInt(configArg["rate"])
		if isInt && rate > 0 {
			configModel.Rate = rate
		}
		size := configArg["size"]

		if size != "" && strings.Contains(size, "*") {
			sizeValues := strings.Split(size, "*")
			if len(sizeValues) == 2 {
				width, _ := utils.JudgeIsInt(sizeValues[0])
				height, _ := utils.JudgeIsInt(sizeValues[1])
				if width > 0 && height > 0 {
					configModel.Width = float64(width)
					configModel.Height = float64(height)
					configModel.Scale = configModel.Width / configModel.Height
				}
			}
		}

		if configArg["profile"] != "" {
			configModel.Profile = configArg["profile"]
		}
		if configArg["level"] != "" {
			configModel.Level = configArg["level"]
		}
		if configArg["pix"] != "" {
			configModel.Pix = configArg["pix"]
		}

		fmt.Println("配置文件读取完毕，配置信息如下：\n")
	}

	fmt.Printf("%+v\n", configModel)
	//inputStr := utils.ClientInputWithMessage("请核对配置信息，按Enter键结束输入，y(确认配置)/n(重新配置) \n",'\n')
	//
	//if inputStr == "n" || inputStr == "no" {
	//	return false
	//}

	mediaConfgiModel = configModel

	return true
}
