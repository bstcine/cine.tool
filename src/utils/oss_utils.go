package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"fmt"
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
func DownloadOssResource(endpoint string, accessKeyId string, accessKeySecret string, bucketName string, savePath string, objectKey string) bool {

	client,err := oss.New(endpoint,accessKeyId,accessKeySecret)

	if err != nil {
		fmt.Println(err)
		fmt.Println("创建客户端对象失败")
		return false
	}

	bucket,err := client.Bucket(bucketName)

	if err != nil {
		fmt.Println(err)
		fmt.Println("创建bucket失败")
		return false
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