# cine.tool
- ### 文件压缩工具-Build
   - Window(exe)
      ```
      $  cd cine.tool
      $  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/app_compress src/app_compress.go
      ```
   - Mac
      ```
      $  cd cine.tool
      $  go build -o bin/app_compress src/app_compress.go
      ```
