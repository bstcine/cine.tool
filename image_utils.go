package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"encoding/base64"
	"strings"
	"log"
	"image"
	"image/png"
	"image/jpeg"
	"github.com/nfnt/resize"
	"fmt"
	"runtime"
	"io/ioutil"
	"image/draw"
)

const TestPath string = "/Volumes/Go/test/"
const LogoName string = ".logo.png"
const DoDir string = "doing" + string(os.PathSeparator)

/**
获取当前路径
 */
func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index] + string(os.PathSeparator)
	return ret
}

/**
检查工具（ Magick ）是否存在
 */
func CheckHasMagick() bool {
	command := "magick -version"

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("请安装图片处理工具 https://www.imagemagick.org/script/download.php#windows")
			fmt.Println(err)
			return false
		}
		fmt.Println(string(out))
	} else {
		cmd := exec.Command("bash", "-c", command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("请安装图片处理工具 https://www.imagemagick.org/script/download.php")
			fmt.Println(err)
			return false
		}
		fmt.Println(string(out))
	}

	return true
}

/**
检查水印
 */
func checkLogoFile(logoPath string)  {
	file, err := os.Open(logoPath)
	if err != nil && os.IsNotExist(err) {
		logoPerm, _ := base64.StdEncoding.DecodeString(logoBaseStr) //成图片文件并把文件写入到buffer
		err := ioutil.WriteFile(logoPath, logoPerm, 0666)
		fmt.Println("水印图片创建成功...")
		if err != nil {
			fmt.Println(err)
		}
	}else {
		fmt.Println("水印图片已经存在...")
	}
	defer file.Close()
}

/**
加工图片
 */
func MachiningImage(debug bool) {
	var doPath, toPath, logoPath string
	if debug {
		doPath = TestPath
		toPath = TestPath
		logoPath = doPath + LogoName
	} else {
		doPath = GetCurrPath()
		toPath = doPath + DoDir
		os.MkdirAll(toPath, 0777)
		logoPath = doPath + LogoName
	}

	files, err := ioutil.ReadDir(doPath)
	if err != nil {
		return
	}

	hasMagick := CheckHasMagick()
	checkLogoFile(logoPath)

	fmt.Println("===============   开始图片文件处理   ===============")
	for i := 0; i < len(files); i++ {
		var info = files[i]
		var name = info.Name()

		if strings.HasPrefix(name,".") || strings.HasPrefix(name,"resize-") || strings.HasPrefix(name,"logo-")|| strings.HasPrefix(name,"m-") || strings.HasPrefix(name,"n-"){
			continue
		}

		if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".png") {
			fmt.Print(name + " 处理中...")
			if debug {
				ResizeImgByMagick(doPath, toPath, name, "m-" + name)
				ResizeImg(doPath, doPath, name, "n-"+name)
			} else {
				resizeName := "resize-" + name
				logoName := "logo-" + name
				if hasMagick {
					ResizeImgByMagick(doPath, toPath, name, resizeName)
					LogoImgByMagick(logoPath, toPath, toPath, resizeName+".jpg", logoName+".jpg")
				} else {
					ResizeImg(doPath, toPath, name, resizeName)
					LogoImg(logoPath, toPath, toPath, resizeName,logoName)
				}
			}

			fmt.Println(" 完成.")
		}
	}

	fmt.Println("===============   图片文件处理成功   ===============")
	fmt.Println("请输入 end ,结束本程序...")
}

/**
图片压缩 - nfnf-reszie
 */
func ResizeImg(oidPath, newPath, oidName, newName string) {
	file, err := os.Open(oidPath + oidName)
	if err != nil {
		log.Fatal(err)
	}

	var img image.Image
	if strings.HasSuffix(oidName, ".png") {
		img, err = png.Decode(file)
	} else if strings.HasSuffix(oidName, ".jpg") || strings.HasSuffix(oidName, ".jpeg") {
		img, err = jpeg.Decode(file)
	}
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Resize(0, 0, img, resize.Bilinear)

	out, err := os.Create(newPath + newName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)

	fmt.Print(" 压缩成功 ")
}

/**
图片水印
 */
func LogoImg(logoPath, oidPath, newPath, oidName, newName string) {
	imgb, _ := os.Open(oidPath + oidName)
	img, _ := jpeg.Decode(imgb)
	defer imgb.Close()

	wmb, _ := os.Open(logoPath)
	watermark, _ := png.Decode(wmb)
	defer wmb.Close()

	offset := image.Pt(0, 0)
	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

	imgw, _ := os.Create(newPath + newName)
	jpeg.Encode(imgw, m, &jpeg.Options{jpeg.DefaultQuality})
	defer imgw.Close()

	fmt.Print(" 水印成功 ")
}

/**
图片压缩 - imagick
 */
func ResizeImgByMagick(oidPath, newPath, oidName, newName string) {
	resizeCmd := "magick convert -quality 30 " + oidPath + oidName + " " + newPath + newName + ".jpg"
	status := CineCMD(resizeCmd)
	if status {
		fmt.Print(" 压缩成功 ")
	} else {
		fmt.Print(" 压缩失败 ")
	}
}

/**
图片水印 - imagick
 */
func LogoImgByMagick(logoPath, oidPath, newPath, name, newName string) {
	logoCmd := "magick convert " + oidPath + name + " " + logoPath + " -gravity southwest -geometry +0+0 -composite " + newPath + newName
	status := CineCMD(logoCmd)
	if status {
		fmt.Print(" 水印成功 ")
	} else {
		fmt.Print(" 水印失败 ")
	}
}

/**
运行命令
 */
func CineCMD(command string) bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Print(err)
			return false
		}
		fmt.Print(string(out))
	} else {
		cmd := exec.Command("bash", "-c", command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Print(err)
			return false
		}
		fmt.Print(string(out))
	}
	return true
}
