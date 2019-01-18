package conf

//是否开发模式
const IsDebug = false

// 是否为测试环境，用来判断相应服务器
// 工作环境， 0 为本地环境，1为测试环境，2为线上环境
const HostEnv = 1

const WorkDir = "/Go/Cine/cine.tool/"         //工作空间
const LogFile = "log/cine_tools.log"          //日志文件
const ConfFile = "conf/cine_tools.cfg"        //配置文件
const ConfFileTmp = "assets/cine_tools_tmp.cfg" //临时配置文件

//API 配置
const API_BASE_URL = "http://www.bstcine.com"
const API_BASE_URL_TEST = "https://dev.bstcine.com"
const API_BASE_URL_LOCAL = "http://local.bstcine.com:9000"

const Media_Host_KJ = "oss.bstcine.com"