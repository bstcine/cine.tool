package main

import (
	"io/ioutil"
	"strings"
	"fmt"
	"os"
	"log"
	"image"
	"image/png"
	"image/jpeg"
	"github.com/nfnt/resize"
	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Println("开始扫描图片文件...")

	dirname := GetCurrPath()
	//dirname = TestSrc

	files, err := ioutil.ReadDir(dirname)

	if err != nil {
		return
	}

	for i := 0; i < len(files); i++ {
		var info = files[i]

		if strings.HasPrefix(info.Name(),TestPrefix) || strings.HasPrefix(info.Name(),TestPrefix2) {
			continue
		}

		if strings.HasSuffix(info.Name(),".jpg") || strings.HasSuffix(info.Name(),".png") {
			ResizeImg(dirname,info.Name(),TestPrefix)
			ResizeImgByImagick(dirname,info.Name(),TestPrefix2)
		}
	}

	fmt.Println("图片压缩成功...")
}


const TestSrc string = "/Volumes/Go/test/"
const TestPrefix string = "nfnt-"
const TestPrefix2 string = "imagick-"

/**
获取当前路径
 */
func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret + "/"
}


/**
图片压缩
 */
func ResizeImg(path,name,prefix string) {
	file, err := os.Open(path + name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(file.Name()+" 图片开始压缩...")

	var img image.Image
	if strings.HasSuffix(name,".png") {
		img, err = png.Decode(file)
	}else if strings.HasSuffix(name,".jpg") || strings.HasSuffix(name,".jpeg") {
		img, err = jpeg.Decode(file)
	}
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Resize(0, 0, img, resize.Lanczos3)

	out, err := os.Create(path + prefix + name)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)
}

func ResizeImgByImagick(path, name, prefix string)  {
	resizecmd := "magick convert -quality 30 "+name+" "+prefix+name

	cmd := exec.Command(resizecmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}