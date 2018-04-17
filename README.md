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
      $ git clone https://github.com/bstcine/cine.tool.git
      $ cd cine.tool
      $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/acine_course_check.exe src/cine_course_check.go
      ```
    - Mac
      ```
      $ git clone https://github.com/bstcine/cine.tool.git
      $ cd cine.tool
      $ go build -o bin/cine_course_check src/cine_course_check.go
      ```
