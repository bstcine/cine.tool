package utils

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
	"os"
	"strconv"
	"bufio"
	"log"
	"time"
	"math/rand"
	"io"
	"path/filepath"
)

/**
 * @ 生成指定范围内指定数量的不重复的随机数
 * @param start 起始点
 * @param end 终点
 * @param count 随机数的数量
 * @return 计算结果
 */
func GenerateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}
	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start
		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}
/***/

/// 清理字典中的所有值的指定字符（如空格键，','等）
func ClearDictionaryChar(dict map[string]string, char string){

	for key,value := range dict {

		if !strings.Contains(value,char) {
			continue
		}

		value = strings.Replace(value,char,"",-1)
		dict[key] = value
	}
}

/// 拷贝文件夹到指定位置
func CopyDir(originPath string,targetPath string) bool{

	CreatDirectory(targetPath);

	err := filepath.Walk(originPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return  nil;
		}
		relativePath := strings.Replace(path,originPath,"",-1);
		newPath := targetPath + relativePath;

		if info.IsDir() {
			// 创建临时文件目录
			//fmt.Println("目录地址：",path);
			if relativePath != "" {
				CreatDirectory(newPath);
			}
		}else {
			// 创建
			//fmt.Println("新的地址：",dest_new);
			CopyFile(path,newPath);
		}

		return  nil
	})

	if err != nil {
		fmt.Println(err);
		return false;
	}

	return true;
}
/// 拷贝文件到指定位置
func CopyFile(originPath string, targetPath string) bool {

	srcFile,err := os.Open(originPath);
	if err != nil {
		fmt.Println(err);
		return  false
	}
	defer srcFile.Close();

	dstFile,err := os.Create(targetPath);
	if err != nil {
		fmt.Println(err);
		return false
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile,srcFile);
	if err != nil {
		fmt.Println(err);
		return  false
	}

	return true
}

/// 获取目录下的所有目录名称
func GetAllDirectoryNames(dirPath string) []string {

	fileHandler,err := ioutil.ReadDir(dirPath)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	var dirNames []string
	for _,file := range fileHandler {
		if file.IsDir() {
			dirNames = append(dirNames,file.Name())
		}
	}

	return dirNames

}

/// 获取目录下的所有文件名
func GetAllFiloeNames(dirPath string) []string{

	fileHandler,err := ioutil.ReadDir(dirPath)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	var fileNames []string
	var tempNames []string
	for _,file := range fileHandler {
		if file.IsDir() {
			continue
		}
		tempNames = append(tempNames,file.Name())
		name := strings.Replace(file.Name(),"._","",-1)

		if strings.HasPrefix(name,".") {
			continue
		}

		if Contains(fileNames,name) {
			continue
		}
		fileNames = append(fileNames,name)
	}

	fmt.Println("所有文件名称如下: \n",tempNames)
	return fileNames
}

/// 获取目录下的所有mp3文件名
func GetAllMp3FiloeNames(dirPath string) (hadVideo bool, fileNames []string){

	hadVideo = false

	fileHandler,err := ioutil.ReadDir(dirPath)

	if err != nil {
		fmt.Println(err)
		return false,nil
	}

	var tempNames []string

	for _,file := range fileHandler {
		if file.IsDir() {
			continue
		}
		tempNames = append(tempNames,file.Name())
		fileName := strings.Replace(file.Name(),"._","",-1)
		if strings.HasPrefix(fileName,".") {
			continue
		}
		if !strings.Contains(fileName,".mp3") {
			if !strings.Contains(fileName,".mp4") {
				continue
			}
			hadVideo = true
		}

		if Contains(fileNames,fileName) {
			continue
		}
		fileNames = append(fileNames,fileName)
	}

	fmt.Println("所有音视频文件名称如下: \n",tempNames)
	return hadVideo,fileNames
}

/// 判断一个数组中是否包含指定元素
func Contains(sliceEl []string, seq string) bool {

	for _,ele := range sliceEl {
		if ele == seq {
			return true
		}
	}
	return false
}

/// 判断一个字符串是否为float64
func JudgeIsFloat64(judgeStr string) (float64, bool) {
	value,err := strconv.ParseFloat(judgeStr,64)
	if err != nil {
		return 0,false
	}
	return value,true
}

/// 判断一个字符串是否为int
func JudgeIsInt(judgeStr string) (int, bool) {
	value,err := strconv.ParseInt(judgeStr,0,32)
	if err != nil {
		return 0,false
	}
	return int(value),true
}

/// 判断是否为int64类型
func JudgeIsInt64(judgeStr string) (int64, bool) {
	value,err := strconv.ParseInt(judgeStr,0,64)
	if err != nil {
		return 0,false
	}
	return value,true
}

func ChangeIntToThirdStr(value int) string {
	s := strconv.Itoa(value)

	if value < 10 {

		s = "00"+s

	}else if value < 100 {

		s = "0"+s
	}

	return  s
}

/// 将100以内的int数据转换为string（显示两位）
func ChangeInt(value int) string {

	s := strconv.Itoa(value)

	if value < 10 {

		s = "0"+s

	}

	return  s
}

/// 读取标准用户输入流（包含提示用户的内容）
/**
 *
 */
func ClientInputWithMessage(message string,endbyte byte) string {

	// 提示用户需要输入信息
	fmt.Println(message)

	value,err := ClientInput(endbyte)

	if err != nil {

		// 告知用户标准输入流异常
		fmt.Println("标准输入流异常，输入信息获取失败")

		return  ""
	}

	return  value
}

/// 读取标准用户输入流(提示用户输入信息)
/**
 * @param endbyte 字符串结束符 char类型 如 '\n', '\t',' '等
 * @return string 除掉结束字符的输入字符串
 * @return error 标准输入流报错
 */
func ClientInput(endbyte byte) (string, error) {

	clientReader := bufio.NewReader(os.Stdin)

	input,err := clientReader.ReadString(endbyte)

	endData := []byte{endbyte,}
	input = strings.Replace(input,string(endData[:]),"",-1)

	return input,err
}

/**
获取日志打印器
 */
func GetLogger(outPath string) *log.Logger {
	dirPath := outPath[0:strings.LastIndex(outPath, "/")+1]
	if _, err := os.Stat(dirPath); err != nil {
		os.MkdirAll(dirPath, 0777)
	}

	logFile, err := os.Create(outPath)

	if err != nil {
		log.Fatalln("open file error !")
	}

	logger := log.New(logFile, "[Info]", log.Lshortfile)
	return logger
}

/**
获取配置文件参数
 */
func GetConfArgs(path string) (argsMap map[string]string) {
	if _, err := os.Stat(path); err != nil {
		argsMap = nil
	} else {
		args, _ := ReadLines(path)
		argsMap = make(map[string]string)

		for i := 0; i < len(args); i++ {
			arg := args[i]
			if !strings.Contains(arg, "#") {
				argSplit := strings.Split(arg, "=")
				if len(argSplit) > 1 {
					argsMap[argSplit[0]] = argSplit[1]
				} else {
					argsMap[argSplit[0]] = ""
				}
			}
		}
	}
	return argsMap
}

/**
获取图片和音频
*/
func GetImageAudio(path string) (images, audios []string) {
	//读取当前目录文件列表
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, nil
	}

	for i := 0; i < len(files); i++ {
		var info = files[i]
		var name = info.Name()

		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "resize-") || strings.HasPrefix(name, "logo-") ||
			strings.HasPrefix(name, "m-") || strings.HasPrefix(name, "n-") || strings.HasPrefix(name, "f-") {
			continue
		}

		//图片处理
		if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".png") {
			images = append(images, name)
		}

		//音频处理
		if strings.HasSuffix(name, ".mp3") {
			audios = append(audios, name)
		}
	}

	return images, audios
}

/**
运行命令
*/
func CineCMD(command string) bool {
	var result bool = true

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			result = false
			fmt.Print(err)
		}
		fmt.Print(string(out))
	} else {
		cmd := exec.Command("bash", "-c", command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			result = false
			fmt.Print(err)
		}
		fmt.Print(string(out))
	}
	return result
}

/**
运行命令
*/
func RunCMD(command string) (string, error) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", command)
		out, err := cmd.CombinedOutput()
		return string(out), err
	} else {
		cmd := exec.Command("bash", "-c", command)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
}
