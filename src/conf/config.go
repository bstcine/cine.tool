package conf

//是否开发模式
const IsDebug = false

// 是否为测试环境，用来判断相应服务器
// 工作环境， 0 为本地环境，1为测试环境，2为线上环境
const HostEnv = 2

const WorkDir = "/Go/Cine/cine.tool/"           //工作空间
const LogFile = "log/cine_tools.log"            //日志文件
const ConfFile = "/cine_tools.cfg"              //配置文件
const ConfFileTmp = "assets/cine_tools_tmp.cfg" //临时配置文件

//API 配置
const API_BASE_URL = "http://www.bstcine.com"
const API_BASE_URL_TEST = "https://dev.bstcine.com"
const API_BASE_URL_LOCAL = "http://local.bstcine.com:9000"

const Media_Host_KJ = "oss.bstcine.com"

// oss 资源下载配置路径
const Course_downloadWorkDir = "/download_resource"                   // 下载oss工作目录
const Course_download_Config = "/cine_course_download.cfg"            // 下载oss配置文件
const Course_download_errorLog = "/cine_course_download_errorlog.txt" // 下载oss错误日志

// oss资源检查配置文件路径
const Course_checkWorkDir = "/oss_checkConfig"       // 检查oss资源工作目录
const Course_checkConfig = "/cine_course_check.cfg"  // 检查oss资源配置信息
const Course_check_log = "/cine_course_errorlog.txt" // 检查oss资源错误文件

// 非课件资源迁移配置文件相对路径
const FMedia_Move_Config = "/move_fmedia_config.cfg"     // 配置文件
const FMedia_Move_ErrorLog = "/move_fmedia_errorLog.txt" // 错误报告
