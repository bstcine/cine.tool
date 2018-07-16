package main
import (
	"./utils"
	"./conf"
	"./model"
	"fmt"
	"path/filepath"
	"os"
	"strings"
)

/**
 * @ 非课件资源的迁移（首次迁移头像和课程封面）
 * @ 通过api分别读取需要迁移的用户头像和课程封面
 * @ 将读取的用户头像和课程封面，上传到oss上
 */

type FMediaModel struct {
	id           string    // 非课件资源id号码
	originalPath string
	targetPath   string
	moveStatus   bool
}

/// 临时错误文件保存路径
var fmedia_errorLog string
var fmedia_config string
func main() {

	// 创建配置文件路径信息
	isSuc := makeConfigs()
	if !isSuc {
		return
	}

	// 读取配置文件，生成配置模型
	configMode,isSuc := getConfigModel()
	if !isSuc {
		return
	}

	// 通过api获取资源路径
	getOriginHeaderImage(1, configMode);
	getOriginHeaderImage(2, configMode);

}

func makeConfigs() bool {
	if conf.IsDebug {
		fmedia_errorLog = "/Users/lidangkun/Desktop/oss_fmedia_move/move_fmedia_errorLog.txt"
		fmedia_config = "/Users/lidangkun/Desktop/oss_fmedia_move/move_fmedia_config.cfg"
	}else {
		dir,err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return false
		}
		dir = strings.Replace(dir,"\\","/",-1)
		fmedia_errorLog = dir + conf.FMedia_Move_ErrorLog
		fmedia_config = dir + conf.FMedia_Move_Config
	}
	fmt.Println("配置文件路径：",fmedia_config)
	fmt.Println("错误日志路径：",fmedia_errorLog)
	return true
}

func getConfigModel() (model.OSSConfig,bool) {
	// 读取配置文件
	configModel := model.OSSConfig{}
	configMap := utils.GetConfArgs(fmedia_config);
	if configMap == nil {
		fmt.Println("没有配置文件")
		return configModel,false
	}
	configModel.KeyId = configMap["AccessKeyId"]
	configModel.KeySecret = configMap["AccessKeySecret"]
	return  configModel,true
}

func test(){
	//resp,_ := http.Get("http://www.bstcine.com/f/2016/10/08/163623548SzYEGzy.jpg");
	//resp,_ := http.Get("http://www.bstcine.com/f/2016/10/08/163623548SzYEGzy.jpg");
	//if resp.StatusCode != 200 {
	//	return
	//}
	//contentType := resp.Header["Content-Type"]
	//if len(contentType) == 1 && contentType[0] == "text/html; charset=utf-8"{
	//	fmt.Println("就是他不好使")
	//}
	//fmt.Println(contentType,len(contentType))
	//fmt.Println(resp.StatusCode,"123",resp.Status)
	fmt.Println("测试方法")
}

func getOriginHeaderImage(resType uint, configModel model.OSSConfig){

	// 通过循环获取
	for i := 1; i > 0; i++ {
		var res model.ResList
		param := make(map[string]interface{})
		param["type"] = resType;
		param["current_page"] = i;
		utils.CommonPost(conf.APIURL_Tool_FMedia,model.Request{"","",param},&res);
		// 将读取到的头像迁移到oss上
		rows := res.Result.Rows.([]interface{})
		moveFMediaResource(rows, configModel);
		// 判断是否读取完毕
		if len(rows) < 100 {
			break
		}
	}

}

// 迁移资源到oss上
func moveFMediaResource(rows []interface{}, configModel model.OSSConfig) bool {

	for j := 0; j < len(rows);j++  {

		resObj := rows[j].(map[string]interface{})
		id := resObj["id"].(string)
		var imagePath = resObj["image"]
		var image,originPath,targetPath string
		if imagePath != nil {
			image = imagePath.(string)
			if image != "" {
				originPath = "http://www.bstcine.com/f/" + image
				targetPath = image
			}
		}
		fmedia := FMediaModel{id,originPath,targetPath,false}
		moveResource(fmedia, configModel)
	}

	return true
}

// 迁移资源
func moveResource(fmedia FMediaModel, configModel model.OSSConfig) bool {

	if fmedia.originalPath == "" {
		return true
	}

	endPoint := "oss-cn-beijing.aliyuncs.com"
	keyId := configModel.KeyId
	secret := configModel.KeySecret
	bucketName := "static-cine-pub";
	saveObject := "f/"

	// 将声音资源写入到oss上
	err1 := utils.PutFile(fmedia.originalPath,endPoint,keyId,secret,bucketName,saveObject+fmedia.targetPath);
	if err1 != nil {
		utils.AppendStringToFile(fmedia_errorLog,fmedia.targetPath + "  " + err1.Error() + "\n")
	}
	fmt.Println(err1)
	return  err1 == nil
}