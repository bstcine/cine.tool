package utils

import (
	"strconv"
	"fmt"
	"os"
	"image/jpeg"
	"encoding/json"
	"../model"
)

// 需要提供八个对外方法

/// 1.单个图片+指定时间 => 无声的ts文件
/**
 @ param imagepath 图片地址
 @ param duration 生成ts视频时长
 @ param targetPath 生成的额ts视频保存的位置
 @ param mediaModel 合成的视频参数设置信息
 @ return 合成结果
 */
func CreatVideoMpegtsWithImage(imagePath string, duration int,targetPath string, defaultFPS bool, mediaModel model.MediaConfig) bool {

	var rate = mediaModel.Rate

	if !defaultFPS && rate == 1 {
		rate = 25
	}

	frameX,frameY := videoImageFrame(imagePath,int(mediaModel.Width),int(mediaModel.Height))
	frameXStr := strconv.Itoa(frameX)
	frameYStr := strconv.Itoa(frameY)
	videoWidth := strconv.Itoa(int(mediaModel.Width))
	videoHeight := strconv.Itoa(int(mediaModel.Height))

	rataLine := " -r " + strconv.Itoa(rate)
	inputFileLine := " -i " + "\"" + imagePath + "\""
	vcodeLine := " -vcodec libx264 -pix_fmt " + mediaModel.Pix
	durationLine := " -t " + strconv.Itoa(duration)
	vfScaleLine := " -vf scale=" + videoWidth + ":" + videoHeight + ":force_original_aspect_ratio=decrease,"
	vfPadLine := "pad="+ videoWidth + ":" + videoHeight + ":" + frameXStr + ":" + frameYStr
	outputPath := " -y -f mpegts " + "\"" + targetPath + "\""
	lines := "ffmpeg" + rataLine + " -f image2 -loop 1" + inputFileLine + vcodeLine + durationLine + vfScaleLine + vfPadLine + outputPath

	_,err := RunCMD(lines)

	return err == nil
}

/// 2.音频文件mp3生成ts文件
/**
 @ param audioPath 音频文件路径
 @ param savePath 生成的ts文件保存路径
 @ return 执行结果
 */
func CreatAudioMpegtsWithMP3(audioPath string, savePath string) bool {

	inputFileLine := " -i " + "\"" + audioPath + "\""

	acodecLine := " -acodec mp2 -vn -ar 44100 -ac 2 -ab 128K -f mpegts"

	outputPath := " -y " + "\"" + savePath + "\""

	line := "ffmpeg" + inputFileLine + acodecLine + outputPath

	_,err := RunCMD(line)

	return  err == nil
}

/// 3.视频文件提取视频源ts文件
/**
 @param videoPath 待提取的视频文件路径
 @param savePath 提取出的视频源路径
 @param mediaModel 待提取视频参数
 @return 提取结果
 */
func CreatVideoMpegtsWithMP4(videoPath string, savePath string, mediaModel model.MediaConfig) bool {

	frameX,frameY := videoMP4Frame(videoPath,int(mediaModel.Width),int(mediaModel.Height))
	frameXStr := strconv.Itoa(frameX)
	frameYStr := strconv.Itoa(frameY)
	videoWidth := strconv.Itoa(int(mediaModel.Width))
	videoHeight := strconv.Itoa(int(mediaModel.Height))

	inputFileLine := " -i " + "\"" + videoPath + "\""
	vcodecLine := " -vcodec libx264 -pix_fmt yuv420p"
	vfScaleLine := " -an -vf scale=" + videoWidth + ":" + videoHeight + ":force_original_aspect_ratio=decrease,"
	vfPadLine := "pad="+ videoWidth + ":" + videoHeight + ":" + frameXStr + ":" + frameYStr
	outputPath := " -y " + "\"" + savePath + "\""

	line := "ffmpeg" + inputFileLine + vcodecLine + vfScaleLine + vfPadLine + outputPath

	_,err := RunCMD(line)

	return err == nil
}

/// 3.视频文件提取音频源ts文件
/**
 @param videoPath 待提取的视频文件路径
 @param savePath 提取后的音频源路径
 @return 提取结果
 */
func CreatAudioMpegtsWithMP4(videoPath string, savePath string) bool {

	inputFileLine := " -i " + "\"" + videoPath + "\""

	acodecLine := " -acodec mp2 -vn -ar 44100 -ac 2 -ab 128K -f mpegts"

	outputPath := " -y " + "\"" + savePath + "\""

	line := "ffmpeg" + inputFileLine + acodecLine + outputPath

	_,err := RunCMD(line)

	return  err == nil
}

/// 4.将一组ts文件合成为一个ts文件
/**
 @param mepgtsArr 待合成的源文件路径集合
 @param targetPath 合成后的资源文件路径
 @return 合成结果
 */
func ComponseMpegts(mpegtsArr []string,targetPath string) bool {

	if len(mpegtsArr) <= 1 {
		fmt.Println("待合并文件不超过一个，不需要执行合并操作")
		return  false
	}

	var mpegtsConcat = ""
	for _,mpeg := range mpegtsArr {

		if mpegtsConcat == "" {
			mpegtsConcat = "concat:"+mpeg
		}else {
			mpegtsConcat = mpegtsConcat + "|" + mpeg
		}
	}

	inputFileLine := " -i " + "\"" + mpegtsConcat + "\""
	vcodeLine := " -c copy"
	lines := "ffmpeg" + inputFileLine + vcodeLine + " -y " + "\"" + targetPath + "\""

	_,err := RunCMD(lines)

	return err == nil
}

/// 5.将一个视频源和一个音频源合成一个包含视频与音频的ts文件
/**
 @param videoPath 视频源路径
 @param audioPath 音频源路径
 @param savePath 合成后的ts文件路径
 @return 合成结果
 */
func CreatMpegtsWithAudioAndVideo(videoPath string ,audioPath string, savePath string) bool {

	videoLine := " -i " + "\"" + videoPath + "\""

	audioLine := " -i " + "\"" + audioPath + "\""

	outputPath := " -y " + "\"" + savePath + "\""

	line := "ffmpeg" + videoLine + audioLine + " -c copy" + outputPath

	_,err := RunCMD(line)

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

	vcodeLine := " -vcodec libx264 -x264-params \"profile="+mediaModel.Profile+":level=" + mediaModel.Level + "\" -pix_fmt " + mediaModel.Pix

	acodecLine := " -acodec aac"

	maxPacketLine := " -max_muxing_queue_size 9999"

	outputLine := " -y " + "\"" + savePath + "\""

	line := "ffmpeg" + inputLine + vcodeLine + acodecLine + maxPacketLine + outputLine

	fmt.Println(line)
	_,err := RunCMD(line)

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
func videoMP4Frame(mp4Path string, width int, height int) (videoX int,videoY int) {

	originWidth,originHeight := GetVideoSize(mp4Path)

	scale := float64(originWidth) / float64(originHeight)

	var targetWidth float64 = 0
	var targetHeight float64 = 0

	if scale > float64(width)/float64(height) {
		targetWidth = float64(width)
		targetHeight = targetWidth / scale
	}else {
		targetHeight = float64(height)
		targetWidth = targetHeight * scale
	}

	videoX = (width - int(targetWidth)) / 2
	videoY = (height-int(targetHeight)) / 2

	return videoX,videoY
}

/// 获取指定图片的在视频画布中的 x，y 值
/**
 @param imagePath 原图片文件路径
 @param width 画布的宽度
 @param height 画布的高度
 @return videoX 原图片在画布中x坐标
 @return videoY 原图片在画布中y坐标
 */
func videoImageFrame(imagePath string, width int, height int) (videoX int,videoY int) {

	originWidth,originHeight := GetImageSize(imagePath)

	scale := float64(originWidth) / float64(originHeight)

	var targetWidth float64 = 0
	var targetHeight float64 = 0

	if scale > float64(width)/float64(height) {
		targetWidth = float64(width)
		targetHeight = targetWidth / scale
	}else {
		targetHeight = float64(height)
		targetWidth = targetHeight * scale
	}

	videoX = (width - int(targetWidth)) / 2
	videoY = (height-int(targetHeight)) / 2

	return videoX,videoY
}

/// 获取视频尺寸
/**
 @parma videoPath 视频路径
 @return width 视频宽度
 @return height 视频高度
 */
func GetVideoSize(videoPath string) (width int,height int) {

	videoStreams,err := GetJsonFileIndoVideo(videoPath)

	if err != nil {
		return 0,0
	}

	videoWidth := videoStreams["width"].(float64)
	videoHeight := videoStreams["height"].(float64)


	return int(videoWidth),int(videoHeight)
}

/// 获取图片尺寸
/**
 @param imagePath 图片路径
 @return width 图片宽度
 @return height 图片高度
 */
func GetImageSize(imagePath string) (width int,height int) {

	file,err := os.Open(imagePath)
	if err != nil {
		fmt.Println("打开失败")
		return 0,0
	}
	cof,err := jpeg.DecodeConfig(file)

	if err != nil {
		fmt.Println("不是jpeg",err)
		return 0,0
	}

	return cof.Width,cof.Height
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
		return  false
	}

	if secondDuration > (firstDuration + 1) {
		return false
	}

	return  true
}

/// 获取音频时长
/**
 @param audioPath 音频路径
 @return 音频时长
 */
func GetDuration(audioPath string) int {

	res,err := GetJsonFileInfo(audioPath)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	jsonInterface := make(map[string]map[string]interface{})
	err = json.Unmarshal([]byte(res),&jsonInterface)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	durationStr := jsonInterface["format"]["duration"].(string)

	duration,err := strconv.ParseFloat(durationStr,0)

	if err != nil {
		fmt.Println(err)
		return  0
	}

	return int(duration)
}
