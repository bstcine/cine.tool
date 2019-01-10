#!/bin/bash
# 获取工程目录
basePath=$(cd `dirname $0`;pwd);
# 转移执行脚本，执行第三方包管理
${basePath}/gpm.sh;
# 创建build/cine_course_check沙盒
buildPath="${basePath}/build";
mkdir $buildPath;

# 编译oss资源检查文件cine_course_check.go
#echo "开始编译cine_course_check.go";
#cine_check_course_Path="${buildPath}/cine_course_check";
#mkdir $cine_check_course_Path;
#cp ${basePath}/config/cine_course_check.cfg ${cine_check_course_Path};
#go build -o ${cine_check_course_Path}/cine_course_check ${basePath}/src/cine_course_check.go;
#echo "cine_course_check.go 编译完成";

# 编译oss资源下载文件cine_course_download.go
echo "开始编译cine_course_download.go";
cine_course_download_path="${buildPath}/cine_course_download";
mkdir ${cine_course_download_path};
cp ${basePath}/config/cine_course_download.cfg ${cine_course_download_path};
go build -o ${cine_course_download_path}/cine_course_download ${basePath}/src/cine_course_download.go;
echo "cine_course_download.go 编译完成";

# 编译cine_tools.go
echo "开始编译cine_tools.go";
cine_tools_path="${buildPath}/cien_tools";
mkdir ${cine_tools_path};
cp ${basePath}/config/cine_tools.cfg ${cine_tools_path};
go build -o ${cine_tools_path}/cine_tools ${basePath}/src/cine_tools.go;
echo "cine_tools.go 编译完成";

# 编译app_download_word.go
echo "开始编译app_download_word.go";
cine_download_word_path="${buildPath}/cine_download_word";
mkdir ${cine_download_word_path};
go build -o ${cine_download_word_path}/app_download_word ${basePath}/src/app_download_word.go;
echo "app_download_word.go 编译完成";

# 编译cine_media_synthesizer.go
echo "开始编译cine_media_synthesizer.go";
cine_media_synthesizer_path="${buildPath}/cine_media_synthesizer";
mkdir ${cine_media_synthesizer_path};
cp ${basePath}/config/cine_media_synthesizer.cfg ${cine_media_synthesizer_path};
go build -o ${cine_media_synthesizer_path}/cine_media_synthesizer ${basePath}/src/cine_media_synthesizer.go;
echo "cine_media_synthesizer.go 编译完成";

# 编译cine_media_mp4_append.go
echo "开始编译cine_media_mp4_append.go";
cine_media_mp4_append_path="${buildPath}/cine_media_mp4_append";
mkdir ${cine_media_mp4_append_path};
go build -o ${cine_media_mp4_append_path}/cine_media_mp4_append ${basePath}/src/cine_media_mp4_append.go;
echo "cine_media_mp4_append.go 编译完成";


