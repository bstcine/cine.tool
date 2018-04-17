# cine.tool
- ## Setup
```shell
$ git clone https://github.com/bstcine/cine.tool.git
$ cd cine.tool
```
      
- ## Build
- ### cine_tools
   - Window(exe)
      ```
      $ go get github.com/nfnt/resize
      $ go github.com/aliyun/aliyun-oss-go-sdk/oss
      $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/cine_tools.exe src/cine_tools.go
      ```
   - Mac
      ```
      $ go get github.com/nfnt/resize
      $ go github.com/aliyun/aliyun-oss-go-sdk/oss
      $ go build -o build/cine_tools src/cine_tools.go
      $ ./bin/cine_tools
      ```
   - Linux
     ```
     $ go get github.com/nfnt/resize
     $ go github.com/aliyun/aliyun-oss-go-sdk/oss
     $ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/cine_tools_linux src/cine_tools.go
     $ ./bin/cine_tools
     ```  

- ### 词汇习题下载工具
   - Window(exe)
      ```
      $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/app_download_word.exe src/app_download_word.go
      ```
   - Mac
      ```
      $ go build -o build/app_download src/app_download_word.go
      ```
      
- ### 课程检查
  - Window(exe)
      ```
      $ mkdir build
      $ cp config/cine_course_check.cfg build/cine_course_check.cfg
      $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/cine_course_check.exe src/cine_course_check.go
      ```
      
   - Mac
      ```
      $ mkdir build
      $ cp config/cine_course_check.cfg build/cine_course_check.cfg
      $ go build -o build/cine_course_check src/cine_course_check.go
      ```
      
   - 注意
      ```
      - 在执行 go build -o 之前，需要更改 cine.tool/src/conf/config.go 中的 IsDebug=false
      - 课件多媒体资源（音视频,水印图,原图）检查工具
      ```
      

     
