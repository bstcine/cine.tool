package conf

//是否开发模式
const IsDebug = true

const WorkDir = "/Go/Cine/cine.tool/"         //工作空间
const LogFile = "log/cine_tools.log"          //日志文件
const ConfFile = "conf/cine_tools.cfg"        //配置文件
const ConfFileTmp = "assets/cine_tools_tmp.cfg" //临时配置文件

//API 配置
const API_BASE_URL = "http://www.bstcine.com"
const API_BASE_URL_TEST = "http://apptest.bstcine.com"

const Media_Host_KJ = "oss.bstcine.com"

// oss 资源下载配置路径
const Course_downloadWorkDir = "./oss_download"                            // 下载oss工作目录
const Course_download_Config = "./oss_download/oss_download_config.txt"    // 下载oss配置文件
const Course_download_errorLog = "./oss_download/oss_error.txt"            // 下载oss错误日志
