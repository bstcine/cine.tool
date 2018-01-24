package utils

import (
	"os"
	"net/http"
	"strconv"
	"io"
	"time"
	"log"
	"fmt"
	"path"
	"os/exec"
	"path/filepath"
	"strings"
)

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

	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

	if err != nil {
		panic(err)
	}

	done := make(chan int64)

	go PrintDownloadPercent(done, outPath, int64(size))

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)

	if err != nil {
		panic(err)
	}

	done <- n

	elapsed := time.Since(start)
	log.Printf("Download completed in %s", elapsed)
}
