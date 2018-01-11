# cine.tool
- ### 文件压缩工具-Build
   - Window(exe)
      ```
      $  cd cine.tool
      $  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/app_compress.exe src/app_compress.go
      ```
   - Mac
      ```
      $  cd cine.tool
      $  go build -o bin/app_compress src/app_compress.go
      ```
- ### 词汇习题下载工具-Build
   - Window(exe)
      ```
      $  cd cine.tool
      $  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/app_download.exe src/app_download.go
      ```
   - Mac
      ```
      $  cd cine.tool
      $  go build -o bin/app_download src/app_download.go
      ```
- ### 课件下载工具-Build
   - Window(exe)
      ```
      $  cd cine.tool
      $  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/app_course_download.exe src/app_course_download.go
      ```
   - Mac
      ```
      $  cd cine.tool
      $  go build -o bin/app_course_download src/app_course_download.go
      ```
