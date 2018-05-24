package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"fmt"
	"net/http"
	"strconv"
)

/// 下载oss上的资源
/**
 * 从oss下载资源
 * @param endpoint oss创建客户端时需要的 endpoint
 * @param accessKeyId oss创建客户端时需要的 accessKeyId
 * @param accessKeySecret oss创建客户端时需要的 accessKeySecret
 * @param bucketName oss创建bucket时需要的 bucketName
 * @param savePath 下载oss资源时需要保存的路径
 * @param objectKey oss下载资源时需要的objectKey
 * @return bool 资源下载结果
 */

var bucket *oss.Bucket

func DownloadImage(endpoint string, accessKeyId string, accessKeySecret string, bucketName string, savePath string, objectKey string, style string) bool {

	var err error

	if bucket == nil {
		client,err := oss.New(endpoint,accessKeyId,accessKeySecret)

		if err != nil {
			fmt.Println(err)
			fmt.Println("创建客户端对象失败")
			return false
		}

		bucket,err = client.Bucket(bucketName)

		if err != nil {
			fmt.Println(err)
			fmt.Println("创建bucket失败")
			return false
		}
	}

	process := oss.Process(style)

	url,err := bucket.SignURL(objectKey,"GET",300,process)

	if err != nil {
		fmt.Println("url 生成失败")
		return false
	}

	err = bucket.GetObjectToFileWithURL(url,savePath)

	if err != nil {
		fmt.Println("下载失败",err)
		return false
	}
	fmt.Println("签名url",url)
	return true
}

func DownloadOssResource(endpoint string, accessKeyId string, accessKeySecret string, bucketName string, savePath string, objectKey string) bool {

	var err error

	if bucket == nil {
		client,err := oss.New(endpoint,accessKeyId,accessKeySecret)

		if err != nil {
			fmt.Println(err)
			fmt.Println("创建客户端对象失败")
			return false
		}

		bucket,err = client.Bucket(bucketName)

		if err != nil {
			fmt.Println(err)
			fmt.Println("创建bucket失败")
			return false
		}
	}

	for i := 0; i < 3;i++  {

		err = bucket.DownloadFile(objectKey,savePath,100*1024,oss.Routines(3),oss.Checkpoint(true,""))

		if err == nil {

			return true

		}else {

			if i == 2 {

				fmt.Println(err)
			}

		}
	}

	return false
}

func CheckResourceSaveStatus(objectKey string) bool {

	// 检查图片状态需要检查两次，一次为原图，一次为水印图，两次成功才返回 1
	requestPath := "http://oss.bstcine.com/" + objectKey

	checkCount := 3

	for checkCount >= 0  {

		checkCount -= 1

		resp,err := http.Head(requestPath)

		if err == nil {

			if resp.StatusCode == 200 {

				length := resp.Header["Content-Length"]

				if len(length) > 0 {

					currentLength,err := strconv.ParseInt(length[0],0,64)

					if err != nil {
						return  false
					}

					return (currentLength > 10240)
				}

			}else {
				fmt.Println(resp.Status)
			}

		}else {
			fmt.Println("访问错误")
		}

	}

	return  false
}