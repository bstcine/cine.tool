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
func CreatAudioMpegtsWithMP4(videoPath string, savePath string) bool {

	inputFileLine := " -i " + "\"" + videoPath + "\""

	acodecLine := " -acodec mp2 -vn -ar 44100 -ac 2 -ab 128K -f mpegts"

	outputPath := " -y " + "\"" + savePath + "\""

	line := "ffmpeg" + inputFileLine + acodecLine + outputPath

	_,err := RunCMD(line)

	return  err == nil
}

/// 4.将一组ts文件合成为一个ts文件
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
func CreatMpegtsWithAudioAndVideo(videoPath string ,audioPath string, savePath string) bool {

	videoLine := " -i " + "\"" + videoPath + "\""

	audioLine := " -i " + "\"" + audioPath + "\""

	outputPath := " -y " + "\"" + savePath + "\""

	line := "ffmpeg" + videoLine + audioLine + " -c copy" + outputPath

	_,err := RunCMD(line)

	return err == nil
}

/// 6.将一个包含音频源和视频源的ts文件转换为mp4文件
func CreatMp4WithMpegts(avPath string, savePath string, mediaModel model.MediaConfig) bool {

	inputLine := " -i " + "\"" + avPath + "\""

	vcodeLine := " -vcodec libx264 -x264-params \"profile="+mediaModel.Profile+":level=" + mediaModel.Level + "\" -pix_fmt " + mediaModel.Pix

	acodecLine := " -acodec copy"

	maxPacketLine := " -max_muxing_queue_size 9999"

	outputLine := " -y " + "\"" + savePath + "\""

	line := "ffmpeg" + inputLine + vcodeLine + acodecLine + maxPacketLine + outputLine

	fmt.Println(line)
	_,err := RunCMD(line)

	return err == nil
}

///// 4.单张图片加单个音频，生成一段视频，帧率默认为1,时长默认与音频相同,需要执行的命令行
//func CreatOneImageAudioLines(imagePath string, audioPath string,targetPath string, defaultFPS bool, mediaModel model.MediaConfig) string {
//
//	var rate = mediaModel.Rate
//
//	if !defaultFPS && rate == 1 {
//		rate = 25
//	}
//
//	audioDuration := GetDuration(audioPath)
//
//	frameX,frameY := videoImageFrame(imagePath,int(mediaModel.Width),int(mediaModel.Height))
//	frameXStr := strconv.Itoa(frameX)
//	frameYStr := strconv.Itoa(frameY)
//	videoWidth := strconv.Itoa(int(mediaModel.Width))
//	videoHeight := strconv.Itoa(int(mediaModel.Height))
//
//	rataLine := " -r " + strconv.Itoa(rate)
//	inputImageLine := " -i " + "\"" + imagePath + "\""
//	inputAudioLine := " -i " + "\"" + audioPath + "\"" + " -t " + strconv.Itoa(audioDuration)
//	vcodeLine := " -vcodec libx264 -x264-params \"profile="+mediaModel.Profile+":level=" + mediaModel.Level + "\" -pix_fmt " + mediaModel.Pix
//	vfScaleLine := " -vf scale=" + videoWidth + ":" + videoHeight + ":force_original_aspect_ratio=decrease,"
//	vfPadLine := "pad="+ videoWidth + ":" + videoHeight + ":" + frameXStr + ":" + frameYStr
//	outputPath := " -y " + "\"" + targetPath + "\""
//	lines := "ffmpeg" + rataLine + " -f image2 -loop 1" + inputImageLine + inputAudioLine + vcodeLine + vfScaleLine + vfPadLine + outputPath
//
//	return lines
//}
//
///// 将一个mp4视频转换为需要大小和清晰度的MP4
//func ChangeVideoSize(videoPath string, savePath string, mediaModel model.MediaConfig) string {
//
//	var rate = mediaModel.Rate
//	if rate == 1 {
//		rate = 25
//	}
//
//	inputLine := " -i " + "\"" + videoPath + "\""
//	rateLine := " -r " + strconv.Itoa(rate)
//	paramsLine := " -x264-params \"profile=" + mediaModel.Profile + ":level=" + mediaModel.Level + "\""
//	pixLine := " -pix_fmt " + mediaModel.Pix + " -s " + mediaModel.Size
//	vcodeLine := " -vcodec libx264 " + paramsLine + pixLine
//	outputLine := " -y \"" + savePath + "\""
//
//	line := "ffmpeg" + inputLine + rateLine + vcodeLine + outputLine
//
//	fmt.Println(line)
//
//	return line
//}
//
///// 将一组可加入concat的ts文件执行合并为可用视频的命令行,
///**
// * @param mpegsAtt 扩展名为.ts路径集合
// * @param targetPath 保存路径
// * @param verbStr 忽略参数，目前仅可以写作 "-an","-vn"两种
// */
//func CreatComponseMpegtsLines(mpegtsArr []string,targetPath string, verbStr string) string {
//
//	if len(mpegtsArr) <= 1 {
//		fmt.Println("待合并文件不超过一个，不需要执行合并操作")
//		return  ""
//	}
//
//	var mpegtsConcat = ""
//	for _,mpeg := range mpegtsArr {
//
//		if mpegtsConcat == "" {
//			mpegtsConcat = "concat:"+mpeg
//		}else {
//			mpegtsConcat = mpegtsConcat + "|" + mpeg
//		}
//	}
//
//	inputFileLine := " -i " + "\"" + mpegtsConcat + "\""
//	vcodeLine := " -c copy -absf aac_adtstoasc" + " " + verbStr
//	lines := "ffmpeg" + inputFileLine + vcodeLine + " -y " + "\"" + targetPath + "\""
//
//	fmt.Println(targetPath)
//	fmt.Println(lines)
//	return lines
//}
//
///// 将一个视频转换为可加入concat的合并文件.ts
//func CreatVideoToMpegtsLines(videoPath string,targetPath string) string {
//
//	// 判断targetPath是否为.ts
//	if !strings.Contains(targetPath,".ts") {
//		fmt.Println("转换后的文件必须以.ts为扩展名")
//		return ""
//	}
//	//
//	inputVideoLine := " -i " + "\"" + videoPath + "\""
//
//	vcodeLine := " -c copy -bsf:v h264_mp4toannexb -f mpegts"
//	outputLine := " -y " + "\"" + targetPath + "\""
//
//	lines := "ffmpeg" + inputVideoLine + vcodeLine + outputLine
//
//	return  lines
//}
//
///// 将一段音频合成到视频中的命令行
//func CreatAudioToVideoLines(audioPath string, videoPath string,targetPath string,mediaModel model.MediaConfig) string {
//
//	inputVideoLine := " -i " + "\"" + videoPath + "\""
//	inputAudioLine := " -i " + "\"" + audioPath + "\""
//	paramsLine := " -x264-params \"profile=" + mediaModel.Profile + ":level=" + mediaModel.Level + "\" -pix_fmt " + mediaModel.Pix
//	vcodeLine := " -vcodec libx264"
//	outputLine := " -y " + "\"" + targetPath + "\""
//	lines := "ffmpeg" + inputAudioLine + inputVideoLine + vcodeLine + paramsLine + outputLine
//
//	return lines
//}

/// 获取指定视频在画布中的 x,y 值
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

/// 获取音频尺寸
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
