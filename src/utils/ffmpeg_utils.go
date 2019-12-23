package utils

import (
	"../model"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// 需要提供4个对外方法

/// 单张图片合成为无声视频
func CreateTsWithImage(imagePath string, audioPath string, rate int, width int, height int, targetPath string) bool {
	imagePath = dealPath(imagePath)
	targetPath = dealPath(targetPath)
	duration := GetDuration(audioPath)
	frameX, frameY := videoImageFrame(imagePath, width, height)
	frameXStr := strconv.Itoa(frameX)
	frameYStr := strconv.Itoa(frameY)
	videoWidth := strconv.Itoa(width)
	videoHeight := strconv.Itoa(height)

	rataLine := " -r " + strconv.Itoa(rate)
	loopLine := " -f image2 -loop 1"
	imageLine := " -i \"" + imagePath + "\""
	audioLine := " -i \"" + audioPath + "\""

	vcodecLine := " -vcodec libx264 -x264-params \"profile=high:level=3.0\" -flags +ildct+ilme -pix_fmt yuv420p"
	durationLine := " -t " + strconv.Itoa(duration)
	vfScaleLine := " -vf scale=" + videoWidth + ":" + videoHeight + ":force_original_aspect_ratio=decrease,"
	vfPadLine := "pad=" + videoWidth + ":" + videoHeight + ":" + frameXStr + ":" + frameYStr
	acodecLine := " -acodec aac -ar 48000 -ac 2 -ab 480k -strict -2"
	outPutLine := " -y -f mpegts \"" + targetPath + "\""

	runLine := "ffmpeg" + rataLine + loopLine + imageLine + audioLine +
		vcodecLine + durationLine + vfScaleLine + vfPadLine + acodecLine + outPutLine
	fmt.Println("CreateTsWithImage==>", runLine)
	_, err := RunCMD(runLine)
	if err != nil {
		fmt.Println(runLine)
		fmt.Println(err)
	}
	return err == nil
}

/// MP4视频转换为Ts(不编码，直接拷贝)
func CreateTsWithMp4(mp4 string, targetPath string) bool {
	runLine := "ffmpeg -i \"" + mp4 + "\" -vcodec copy -acodec copy -y \"" + targetPath + "\""
	fmt.Println("CreateTsWithMp4==>", runLine)
	_, err := RunCMD(runLine)
	if err != nil {
		fmt.Println(runLine)
		fmt.Println(err)
	}
	return err == nil
}

/// ts转换为MP4(不编码，字节拷贝)
func ComponseMP4WithTs(mpegtsArr []string, targetPath string) bool {
	tmpTarget := targetPath + "_tmp.ts"
	var mpegtsConcat = ""
	for _, mpeg := range mpegtsArr {
		mpeg = dealPath(mpeg)
		if mpegtsConcat == "" {
			mpegtsConcat = "concat:" + mpeg
		} else {
			mpegtsConcat = mpegtsConcat + "|" + mpeg
		}
	}
	runLine := "ffmpeg -i \"" + mpegtsConcat + "\" -vcodec copy -acodec copy -y \"" + tmpTarget + "\""
	fmt.Println("ComponseMP4WithTs==>", runLine)
	_, err := RunCMD(runLine)
	if err != nil {
		fmt.Println(runLine)
		fmt.Println(err)
	}
	// 将临时文件转换为MP4
	targetRunLine := "ffmpeg -i \"" + tmpTarget +
		"\"  -vcodec libx264 -x264-params \"profile=high:level=3.0\" -flags +ildct+ilme -pix_fmt yuv420p " +
		"-acodec aac -ar 48000 -ac 2 -ab 480k -strict -2 -y \"" +
		targetPath + "\""
	fmt.Println("ComponseMP4WithTs==>", targetRunLine)
	_, err = RunCMD(targetRunLine)

	fmt.Println(targetRunLine)
	if err != nil {
		fmt.Println(runLine)
		fmt.Println(err)
	}
	os.Remove(tmpTarget)
	return err == nil
}

/// 1.单个图片+音频 => ts文件
/**
@ param imagepath 图片地址
@ param duration 生成ts视频时长
@ param targetPath 生成的额ts视频保存的位置
@ param mediaModel 合成的视频参数设置信息
@ return 合成结果
*/
func CreateTsWithImageAudio(imagePath string, audioPath string, targetPath string, mediaModel model.MediaConfig) bool {
	imagePath = dealPath(imagePath)
	audioPath = dealPath(audioPath)
	targetPath = dealPath(targetPath)
	var rate = mediaModel.Rate
	if rate == 1 {
		rate = 25
	}
	var duration = GetDuration(audioPath)
	frameX, frameY := videoImageFrame(imagePath, int(mediaModel.Width), int(mediaModel.Height))
	frameXStr := strconv.Itoa(frameX)
	frameYStr := strconv.Itoa(frameY)
	videoWidth := strconv.Itoa(int(mediaModel.Width))
	videoHeight := strconv.Itoa(int(mediaModel.Height))

	rataLine := " -r " + strconv.Itoa(rate)
	loopLine := " -f image2 -loop 1"
	imageLine := " -i \"" + imagePath + "\""
	audioLine := " -i \"" + audioPath + "\""
	vcodecLine := " -vcodec libx264 -x264-params \"profile=" +
		mediaModel.Profile + ":level=" + mediaModel.Level + "\"" +
		" -flags +ildct+ilme -pix_fmt " + mediaModel.Pix
	durationLine := " -t " + strconv.Itoa(duration)
	//sizeAspect := " -s 1920*1080 -aspect 16:9"
	vfScaleLine := " -vf scale=" + videoWidth + ":" + videoHeight + ":force_original_aspect_ratio=decrease,"
	vfPadLine := "pad=" + videoWidth + ":" + videoHeight + ":" + frameXStr + ":" + frameYStr
	acodecLine := " -acodec aac -ar 48000 -ac 2 -ab 480k -strict -2"
	outPutLine := " -y -f mpegts \"" + targetPath + "\""
	runLine := "ffmpeg" + rataLine + loopLine + imageLine + audioLine +
		vcodecLine + durationLine + vfScaleLine + vfPadLine + acodecLine + outPutLine
	fmt.Println("CreateTsWithImageAudio==>", runLine)
	_, err := RunCMD(runLine)
	if err != nil {
		fmt.Println(runLine)
		fmt.Println(err)
	}
	return err == nil
}

/// 多个图片+音频合成ts文件
func CreateMpegtsWithImagesAudio(images []map[string]string, tmpDir string, audioPath string, targetPath string, mediaModel model.MediaConfig) bool {
	audioPath = dealPath(audioPath)
	targetPath = dealPath(targetPath)
	tmpDir = dealPath(tmpDir)

	var rate = mediaModel.Rate
	if rate == 1 {
		rate = 25
	}
	rataLine := " -r " + strconv.Itoa(rate)
	frameX, frameY := videoImageFrame(images[0]["path"], int(mediaModel.Width), int(mediaModel.Height))
	frameXStr := strconv.Itoa(frameX)
	frameYStr := strconv.Itoa(frameY)
	videoWidth := strconv.Itoa(int(mediaModel.Width))
	videoHeight := strconv.Itoa(int(mediaModel.Height))

	// 创建input.txt文件
	var mpegtsTmpPaths []string
	for index, imageMap := range images {
		imagePath := imageMap["path"]
		imagePath = dealPath(imagePath)
		//fileInfo := "file '" + imagePath + "'\n"
		//durationInfo := "duration " + imageMap["duration"] + "\n"
		pathComponet := strings.Split(imagePath, "/")
		imageComponent := pathComponet[len(pathComponet)-1]
		imageNameComponent := strings.Split(imageComponent, ".")
		cancelType := imageNameComponent[0]
		nameComponent := strings.Split(cancelType, "_")
		name := nameComponent[0]
		startTime := nameComponent[len(nameComponent)-1]
		var tmpIndex string
		if index < 10 {
			tmpIndex = name + "_00" + strconv.Itoa(index) + ".ts"
		} else if index < 100 {
			tmpIndex = name + "_0" + strconv.Itoa(index) + ".ts"
		}
		tmpPath := tmpDir + "/" + tmpIndex
		mpegtsTmpPaths = append(mpegtsTmpPaths, tmpPath)
		//AppendStringToFile(fileName, fileInfo)
		//AppendStringToFile(fileName, durationInfo)
		duration := imageMap["duration"]

		loopLine := " -f image2 -loop 1"
		imageLine := " -i \"" + imagePath + "\""
		audioLine := " -i \"" + audioPath + "\""
		startLine := " -ss " + startTime
		vcodecLine := " -vcodec libx264 -x264-params \"profile=" +
			mediaModel.Profile + ":level=" + mediaModel.Level + "\"" +
			" -flags +ildct+ilme -pix_fmt " + mediaModel.Pix
		durationLine := " -t " + duration
		fmt.Println("startTime: ", startLine)
		fmt.Println("duration: ", duration)
		//sizeAspect := " -s 1920*1080 -aspect 16:9"
		vfScaleLine := " -vf scale=" + videoWidth + ":" + videoHeight + ":force_original_aspect_ratio=decrease,"
		vfPadLine := "pad=" + videoWidth + ":" + videoHeight + ":" + frameXStr + ":" + frameYStr
		acodecLine := " -acodec aac -ar 48000 -ac 2 -ab 480k -strict -2"
		outPutLine := " -y -f mpegts \"" + tmpPath + "\""
		imageRunLine := "ffmpeg" + rataLine + loopLine + imageLine + audioLine + startLine +
			vcodecLine + durationLine + vfScaleLine + vfPadLine + acodecLine + outPutLine
		fmt.Println("CreateMpegtsWithImagesAudio==>", imageRunLine)
		_, err := RunCMD(imageRunLine)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	fmt.Println("图片处理完成了，准备合成整个视频")
	componseSuc := ComponseMpegts(mpegtsTmpPaths, targetPath, mediaModel)

	return componseSuc
}

/// 2. mp4 文件转换为 ts 文件
/**
@param videoPath MP4文件路径
@param savePath ts路径
@param mediaModel 待提取视频参数
@return 提取结果
*/
func CreatVideoMpegtsWithMP4(videoPath string, savePath string, mediaModel model.MediaConfig) bool {
	videoPath = dealPath(videoPath)
	savePath = dealPath(savePath)
	duration := GetDuration(videoPath)
	if duration <= 0 {
		fmt.Println("视频时长异常", duration)
	}

	frameX, frameY := videoMP4Frame(videoPath, int(mediaModel.Width), int(mediaModel.Height))
	frameXStr := strconv.Itoa(frameX)
	frameYStr := strconv.Itoa(frameY)
	videoWidth := strconv.Itoa(int(mediaModel.Width))
	videoHeight := strconv.Itoa(int(mediaModel.Height))

	inputFileLine := " -i " + "\"" + videoPath + "\""
	vcodecLine := " -vcodec libx264 -x264-params \"profile=" +
		mediaModel.Profile + ":level=" + mediaModel.Level + "\"" +
		" -flags +ildct+ilme -pix_fmt " + mediaModel.Pix
	durationLine := " -t " + strconv.Itoa(duration)
	//sizeAspect := " -s 1920*1080 -aspect 16:9"
	vfScaleLine := " -vf scale=" + videoWidth + ":" + videoHeight + ":force_original_aspect_ratio=decrease,"
	vfPadLine := "pad=" + videoWidth + ":" + videoHeight + ":" + frameXStr + ":" + frameYStr
	acodecLine := " -acodec aac -ar 48000 -ac 2 -ab 480k -strict -2"
	outputPath := " -y -f mpegts \"" + savePath + "\""

	line := "ffmpeg " + inputFileLine + vcodecLine + durationLine +
		vfScaleLine + vfPadLine + acodecLine + outputPath

	_, err := RunCMD(line)
	if err == nil {
		fmt.Println(line)
		fmt.Println(err)
	}

	return err == nil
}

/// 3.将一组ts文件合成为一个ts文件
/**
@param mepgtsArr 待合成的源文件路径集合
@param targetPath 合成后的资源文件路径
@return 合成结果
*/
func ComponseMpegts(mpegtsArr []string, targetPath string, mediaModel model.MediaConfig) bool {

	targetPath = dealPath(targetPath)
	if len(mpegtsArr) <= 1 {
		fmt.Println("待合并文件不超过一个，不需要执行合并操作")
		return false
	}
	var rate = mediaModel.Rate
	if rate == 1 {
		rate = 25
	}
	var mpegtsConcat = ""
	for _, mpeg := range mpegtsArr {
		mpeg = dealPath(mpeg)
		if mpegtsConcat == "" {
			mpegtsConcat = "concat:" + mpeg
		} else {
			mpegtsConcat = mpegtsConcat + "|" + mpeg
		}
	}
	inputFileLine := " -i " + "\"" + mpegtsConcat + "\""
	vcodecLine := " -vcodec libx264 -x264-params \"profile=" +
		mediaModel.Profile + ":level=" + mediaModel.Level + "\"" +
		" -flags +ildct+ilme -pix_fmt " + mediaModel.Pix + " -s 1920*1080 -aspect 16:9"
	acodecLine := " -acodec aac -ar 48000 -ac 2 -ab 480k -strict -2"
	lines := "ffmpeg" + inputFileLine + vcodecLine + acodecLine + " -y -f mpegts " + "\"" + targetPath + "\""
	fmt.Println("ComponseMpegts==>", lines)
	_, err := RunCMD(lines)
	if err != nil {

		fmt.Println(lines)
		fmt.Println(err)
	}
	return err == nil
}

/// 6.将一个包含音频源和视频源的ts文件转换为mp4文件
/**
@param avPath 包含音频源和视频源的ts文件路径
@param savePath 转换后的mp4文件路径
@param mediaModel 转换后的mp4资源配置参数信息
@return 转码的结果
*/
func CreatMp4WithMpegts(avPath string, savePath string, mediaModel model.MediaConfig) bool {

	inputLine := " -i " + "\"" + avPath + "\""

	vcodeLine := " -vcodec libx264 -x264-params \"profile=" + mediaModel.Profile + ":level=" + mediaModel.Level + "\" -pix_fmt " + mediaModel.Pix

	acodecLine := " -acodec copy"

	maxPacketLine := " -max_muxing_queue_size 9999"

	outputLine := " -y " + "\"" + savePath + "\""

	line := "ffmpeg" + inputLine + vcodeLine + acodecLine + maxPacketLine + outputLine

	fmt.Println("CreatMp4WithMpegts==>", line)
	_, err := RunCMD(line)

	return err == nil
}

/// 获取指定视频在画布中的 x,y 值
/**
@param mp4Path 原视频文件路径
@param width 画布的宽度
@param height 画布的高度
@return videoX 原视频在画布中的x坐标
@return videoY 原视频在画布中的y坐标
*/
func videoMP4Frame(mp4Path string, width int, height int) (videoX int, videoY int) {

	originWidth, originHeight := GetVideoSize(mp4Path)

	scale := float64(originWidth) / float64(originHeight)

	var targetWidth float64 = 0
	var targetHeight float64 = 0

	if scale > float64(width)/float64(height) {
		targetWidth = float64(width)
		targetHeight = targetWidth / scale
	} else {
		targetHeight = float64(height)
		targetWidth = targetHeight * scale
	}

	videoX = (width - int(targetWidth)) / 2
	videoY = (height - int(targetHeight)) / 2

	return videoX, videoY
}

/// 获取指定图片的在视频画布中的 x，y 值
/**
@param imagePath 原图片文件路径
@param width 画布的宽度
@param height 画布的高度
@return videoX 原图片在画布中x坐标
@return videoY 原图片在画布中y坐标
*/
func videoImageFrame(imagePath string, width int, height int) (videoX int, videoY int) {

	originWidth, originHeight := GetImageSize(imagePath)

	scale := float64(originWidth) / float64(originHeight)

	var targetWidth float64 = 0
	var targetHeight float64 = 0

	if scale > float64(width)/float64(height) {
		targetWidth = float64(width)
		targetHeight = targetWidth / scale
	} else {
		targetHeight = float64(height)
		targetWidth = targetHeight * scale
	}

	videoX = (width - int(targetWidth)) / 2
	videoY = (height - int(targetHeight)) / 2

	return videoX, videoY
}

/// 获取视频尺寸
/**
@parma videoPath 视频路径
@return width 视频宽度
@return height 视频高度
*/
func GetVideoSize(videoPath string) (width int, height int) {

	videoStreams, err := GetJsonFileIndoVideo(videoPath)

	if err != nil {
		return 0, 0
	}

	videoWidth := videoStreams["width"].(float64)
	videoHeight := videoStreams["height"].(float64)

	return int(videoWidth), int(videoHeight)
}

/// 获取图片尺寸
/**
@param imagePath 图片路径
@return width 图片宽度
@return height 图片高度
*/
func GetImageSize(imagePath string) (width int, height int) {

	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("打开失败")
		return 0, 0
	}
	cof, err := jpeg.DecodeConfig(file)

	if err != nil {
		fmt.Println("不是jpeg", err)
		return 0, 0
	}

	return cof.Width, cof.Height
}

/// 判断两个资源文件的时长是否相等（上下偏差不超过1秒）
/**
@param firstPath 第一个多媒体资源的路径
@param secondPath 第二个多媒体资源路径
@return 比较的结果
*/
func JudgeDurationEqual(firstPath string, secondPath string) bool {

	firstDuration := GetDuration(firstPath)
	secondDuration := GetDuration(secondPath)

	if secondDuration < (firstDuration - 1) {
		return false
	}

	if secondDuration > (firstDuration + 1) {
		return false
	}

	return true
}

/// 获取音频时长
/**
@param audioPath 音频路径
@return 音频时长
*/
func GetDuration(audioPath string) int {

	audioPath = dealPath(audioPath)
	res, err := GetJsonFileInfo(audioPath)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	jsonInterface := make(map[string]map[string]interface{})
	err = json.Unmarshal([]byte(res), &jsonInterface)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	durationStr := jsonInterface["format"]["duration"].(string)

	duration, err := strconv.ParseFloat(durationStr, 0)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	return int(duration)
}

var osType = runtime.GOOS

func dealPath(originPath string) string {
	if osType == "windows" {
		newPath := strings.Replace(originPath, "/", "\\", -1)
		return newPath
	}
	return originPath
}
