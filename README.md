# cine.tool
- ### cine_tools - Build
   - Window(exe)
      ```
      $ git clone https://github.com/bstcine/cine.tool.git
      $ cd cine.tool
      $ go get github.com/nfnt/resize
      $ go github.com/aliyun/aliyun-oss-go-sdk/oss
      $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/cine_tools.exe src/cine_tools.go
      ```
   - Mac
      ```
      $ git clone https://github.com/bstcine/cine.tool.git
      $ cd cine.tool
      $ go get github.com/nfnt/resize
      $ go github.com/aliyun/aliyun-oss-go-sdk/oss
      $ go build -o bin/cine_tools src/cine_tools.go
      $ ./bin/cine_tools
      ```
   - Linux
     ```
     $ git clone https://github.com/bstcine/cine.tool.git
     $ cd cine.tool
     $ go get github.com/nfnt/resize
     $ go github.com/aliyun/aliyun-oss-go-sdk/oss
     $ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/cine_tools_linux src/cine_tools.go
     $ ./bin/cine_tools
     ```  

- ### 词汇习题下载工具-Build
   - Window(exe)
      ```
      $ git clone https://github.com/bstcine/cine.tool.git
      $ cd cine.tool
      $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/app_download_word.exe src/app_download_word.go
      ```
   - Mac
      ```
      $ git clone https://github.com/bstcine/cine.tool.git
      $ cd cine.tool
      $ go build -o bin/app_download src/app_download_word.go
      ```
- ### 课件多媒体资源（音视频,水印图,原图）检查工具-Build
    - Windows
      ```
      $ cd /Users/userrName/Desktop/CineTool
      $ git clone https://github.com/bstcine/cine.tool.git
      $ mkdir cine_course_check
      $ cp cine.tool/assets/cine_course_check.cfg cine_course_check/cine_course_check.cfg
      $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o cine_course_check/cine_course_check.exe cine.tool/src/cine_course_check.go
      ```
    - Mac
      ```
      $ cd /Users/userrName/Desktop/CineTool
      $ git clone https://github.com/bstcine/cine.tool.git
      $ mkdir cine_course_check
      $ cp cine.tool/assets/cine_course_check.cfg cine_course_check/cine_course_check.cfg
      $ go build -o cine_course_check/cine_course_check cine.tool/src/cine_course_check.go
      ```
    - 注意：
      ```
      在执行 go build -o 之前，需要更改 cine.tool/src/conf/config.go 中的 IsDebug=false
      ```
