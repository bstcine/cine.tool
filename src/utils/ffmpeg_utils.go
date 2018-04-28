package utils

import (
	"strconv"
	"../conf"
	"fmt"
	"os"
	"image/jpeg"
	"encoding/json"
	"strings"
)

/// 将一组可加入concat的ts文件执行合并为可用视频的命令行,
/**
 * @param mpegsAtt 扩展名为.ts路径集合
 * @param targetPath 保存路径
 * @param verbStr 忽略参数，目前仅可以写作 "-an","-vn"两种
 */
func CreatComponseMpegtsLines(mpegtsArr []string,targetPath string, verbStr string) string {

	if len(mpegtsArr) <= 1 {
		fmt.Println("待合并文件不超过一个，不需要执行合并操作")
		return  ""
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
	vcodeLine := " -c copy -absf aac_adtstoasc" + " " + verbStr
	lines := "ffmpeg" + inputFileLine + vcodeLine + " -y " + "\"" + targetPath + "\""

	fmt.Println(targetPath)
	fmt.Println(lines)
	return lines
}

/// 将一个视频转换为可加入concat的合并文件.ts
func CreatVideoToMpegtsLines(videoPath string,targetPath string) string {

	// 判断targetPath是否为.ts
	if !strings.Contains(targetPath,".ts") {
		fmt.Println("转换后的文件必须以.ts为扩展名")
		return ""
	}
	//
	inputVideoLine := " -i " + "\"" + videoPath + "\""

	vcodeLine := " -c copy -bsf:v h264_mp4toannexb -f mpegts"
	outputLine := " -y " + "\"" + targetPath + "\""

	lines := "ffmpeg" + inputVideoLine + vcodeLine + outputLine

	return  lines
}

/// 将一段音频合成到视频中的命令行
func CreatAudioToVideoLines(audioPath string, videoPath string,targetPath string) string {

	inputVideoLine := " -i " + "\"" + videoPath + "\""
	inputAudioLine := " -i " + "\"" + audioPath + "\""
	vcodeLine := " -vcodec libx264 -x264-params \"profile=baseline:level=3.0\" -pix_fmt yuv420p"
	outputLine := " -y " + "\"" + targetPath + "\""
	lines := "ffmpeg" + inputAudioLine + inputVideoLine + vcodeLine + outputLine

	return lines
}

/// 单张图片加单个音频，生成一段视频，帧率默认为1,时长默认与音频相同,需要执行的命令行
func CreatOneImageAudioLines(imagePath string, audioPath string,targetPath string) string {

	audioDuration := GetDuration(audioPath)

	frameX,frameY := videoImageFrame(imagePath)
	frameXStr := strconv.Itoa(frameX)
	frameYStr := strconv.Itoa(frameY)
	videoWidth := strconv.Itoa(int(conf.FFMPEG_videoWidth))
	videoHeight := strconv.Itoa(int(conf.FFMPEG_videoHeight))

	inputImageLine := " -i " + "\"" + imagePath + "\""
	inputAudioLine := " -i " + "\"" + audioPath + "\"" + " -t " + ChangeIntToThirdStr(audioDuration)
	vcodeLine := " -vcodec libx264 -x264-params \"profile="+conf.FFMPEG_videoProfile+":level=3.0\" -pix_fmt yuv420p"
	vfScaleLine := " -vf scale=" + videoWidth + ":" + videoHeight + ":force_original_aspect_ratio=decrease,"
	vfPadLine := "pad="+ videoWidth + ":" + videoHeight + ":" + frameXStr + ":" + frameYStr
	outputPath := " -y " + "\"" + targetPath + "\""
	lines := "ffmpeg -r 1 -f image2 -loop 1" + inputImageLine + inputAudioLine + vcodeLine + vfScaleLine + vfPadLine + outputPath

	return lines
}

/// 单张图片指定时长，生成一段无声的视频,帧率默认为1，需要执行的命令行
func CreatLines(imagePath string, duration int,targetPath string) string {

	frameX,frameY := videoImageFrame(imagePath)
	frameXStr := strconv.Itoa(frameX)
	frameYStr := strconv.Itoa(frameY)
	videoWidth := strconv.Itoa(int(conf.FFMPEG_videoWidth))
	videoHeight := strconv.Itoa(int(conf.FFMPEG_videoHeight))

	inputFileLine := " -i " + "\"" + imagePath + "\""
	vcodeLine := " -vcodec libx264 -x264-params \"profile="+conf.FFMPEG_videoProfile+":level=3.0\" -pix_fmt yuv420p"
	durationLine := " -t " + strconv.Itoa(duration)
	vfScaleLine := " -vf scale=" + videoWidth + ":" + videoHeight + ":force_original_aspect_ratio=decrease,"
	vfPadLine := "pad="+ videoWidth + ":" + videoHeight + ":" + frameXStr + ":" + frameYStr
	outputPath := " -y " + "\"" + targetPath + "\""
	lines := "ffmpeg -r 1 -f image2 -loop 1" + inputFileLine + vcodeLine + durationLine + vfScaleLine + vfPadLine + outputPath

	return lines
}
/// 获取指定图片的在视频画布中的 x，y 值
func videoImageFrame(imagePath string) (videoX int,videoY int) {

	originWidth,originHeight := GetImageSize(imagePath)

	scale := float64(originWidth) / float64(originHeight)

	var targetWidth float64 = 0
	var targetHeight float64 = 0

	if scale > conf.FFMPEG_videoScale {
		targetWidth = float64(conf.FFMPEG_videoWidth)
		targetHeight = targetWidth / scale
	}else {
		targetHeight = float64(conf.FFMPEG_videoHeight)
		targetWidth = targetHeight * scale
	}

	videoX = int((conf.FFMPEG_videoWidth - targetWidth) / 2)
	videoY = int((conf.FFMPEG_videoHeight-targetHeight) / 2)

	return videoX,videoY
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
