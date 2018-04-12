package utils

import (
	"os"
	"net/http"
	"io"
	"time"
	"log"
	"fmt"
	"path"
	"os/exec"
	"path/filepath"
	"strings"
	"bufio"
	"bytes"
)

/// 往文件中追加文字
/**
 * @param filePath 需要追加的文件路径
 * @param content 准备追加的文字信息
 */
func AppendStringToFile(filePath string, content string) bool {

	// 判断文件是否存在
	_,err := os.Stat(filePath)

	if err != nil {
		os.Create(filePath)
	}

	fileHandler,err := os.OpenFile(filePath,os.O_WRONLY|os.O_APPEND,0666)

	if err != nil {

		return  false
	}

	_,err = fileHandler.WriteString(content)

	if err != nil {
		return  false
	}
	return  true
}

/// 创建目录
/**
 * @param path 资源文件路径
 * @return 是否创建完成
 */
func CreatDirectory(path string) bool {

	err := os.MkdirAll(path,0711)

	if err == nil {
		return true
	}

	return false
}

/**
文件是否存在
*/
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

/**
获取当前路径
*/
func GetCurPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index] + string(os.PathSeparator)
	return ret
}

/**
获取输出路径
*/
func GetOutPath(dir string) string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index] + string(os.PathSeparator) + dir + string(os.PathSeparator)
	return ret
}

/**
获取网络文件流
 */
func GetHttpFileBytes(url string) (*bytes.Buffer,error) {

	file := path.Base(url)

	start := time.Now()

	headResp, err := http.Head(url)
	if err != nil {
		return nil,err
	}

	defer headResp.Body.Close()

	resp, err := http.Get(url)
	if err != nil {
		return nil,err
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	len,err := buf.ReadFrom(resp.Body)

	if err != nil {
		return nil,err
	}

	elapsed := time.Since(start)
	log.Printf("get file bytes(%d) %s from %s - completed in %s \n",len, file, url,elapsed)

	return buf,nil
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

//writeLines writes the lines to the given file
func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

/**
获取JSON格式的文件信息
*/
func GetJsonFileInfo(url string) (string, error) {
	var cmd = "ffprobe -v quiet -print_format json -show_format " + url
	result, err := RunCMD(cmd)
	return result, err
}

/**
下载进度
 */
func PrintDownloadPercent(done chan int64, path string, total int64) {

	var stop bool = false

	for {
		select {
		case <-done:
			stop = true
		default:

			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

			size := fi.Size()

			if size == 0 {
				size = 1
			}

			var percent float64 = float64(size) / float64(total) * 100

			fmt.Printf("%.0f", percent)
			fmt.Println("%")
		}

		if stop {
			break
		}

		time.Sleep(time.Second)
	}
}

/**
下载文件
 */
func DownloadFile(url string, outPath string) {

	dirPath := outPath[0:strings.LastIndex(outPath,"/")+1]
	if _, err := os.Stat(dirPath); err != nil {
		os.MkdirAll(dirPath, 0777)
	}

	file := path.Base(url)

	log.Printf("Downloading file %s from %s\n", file, url)

	start := time.Now()

	out, err := os.Create(outPath)

	if err != nil {
		fmt.Println(outPath)
		panic(err)
	}

	defer out.Close()

	headResp, err := http.Head(url)

	if err != nil {
		panic(err)
	}

	defer headResp.Body.Close()

	//size, err := strconv.Atoi(headResp.Header.CommonGet("Content-Length"))

	if err != nil {
		panic(err)
	}

	//done := make(chan int64)

	//go PrintDownloadPercent(done, outPath, int64(size))

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	io.Copy(out, resp.Body)

	//done <- n

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)
}