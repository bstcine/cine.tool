# cine.tool
- ## Setup
```shell
$ git clone https://github.com/bstcine/cine.tool.git
$ cd cine.tool
```
      
- ## Build

- ### Build_mac_all
```
$ ./gbuild_mac.sh
```

- ### cine_tools
  - Window(exe)
     ```
     $ go get github.com/nfnt/resize
     $ go get github.com/aliyun/aliyun-oss-go-sdk/oss
     $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/cine_tools.exe src/cine_tools.go
     ```
  - Mac
     ```
     $ go get github.com/nfnt/resize
     $ go get github.com/aliyun/aliyun-oss-go-sdk/oss
     $ go build -o build/cine_tools src/cine_tools.go
     $ ./bin/cine_tools
     ```
  - Linux
     ```
     $ go get github.com/nfnt/resize
     $ go get github.com/aliyun/aliyun-oss-go-sdk/oss
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
     $ ./gpm.sh
     $ mkdir build
     $ cp config/cine_course_check.cfg build/cine_course_check.cfg
     $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/cine_course_check.exe src/cine_course_check.go
     ```
      
  - Mac
     ```
     $ ./gpm.sh
     $ mkdir build
     $ cp config/cine_course_check.cfg build/cine_course_check.cfg
     $ go build -o build/cine_course_check src/cine_course_check.go
     ```
      
  - 注意
     ```
     - 在执行 go build -o 之前，需要更改 cine.tool/src/conf/config.go 中的 IsDebug=false
     - 课件多媒体资源（音视频,水印图,原图）检查工具
     ```

- ### 音视频合成
  ```
  如果已经执行了Build_mac_all 即已经执行了 ./gbuild_mac.sh脚本，则可以忽略以下步骤，请看注意事项
  ```
  - Window(exe)
    ```
     $ ./gpm.sh
     $ mkdir build
     $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/cine_media_synthesizer.exe src/cine_media_synthesizer.go
    ```

  - Mac
    ```
     $ ./gpm.sh
     $ mkdir build
     $ go build -o build/cine_media_synthesizer src/cine_media_synthesizer.go
    ```
     
  - 注意
    ```
    - 0.如果使用了 ./gbuild_mac.sh脚本,cine_media_synthesizer在build/cine_media_synthesizer目录下，
        如果使用 音视频合成->Mac 中的方法，则直接在build中
    - 1.在工具在使用时，需要与待检查目录平级放置，如 courseList/the cat in the hat/chapter1/ls_lesson1/... ,
        courseList/cine_media_synthesier
    - 2.待检查目录可是是course目录，chapter目录，lesson目录，合成工具平级放置，
        例如courselist/cine_media_synthesizer, courselist/the cat in the hat/cine_media_synthesizer,
        courselist/the cat in the hat/chapter1/cine_media_synthesizer
    - 3.放置音视频及图片的目名称录必须以 "ls_" 开头，如ls_lesson1,ls_introduction等
    - 4.音视频资源必须为mp3 和mp4 资源，扩展名为 ",mp3", ".mp4", 必须要以三位整数顺序命名，
        如000.mp3, 001.mp3, 002.mp4, 003.mp3, 004.mp4 等
    - 5.图片资源必须为jpg资源，扩展名为 ".jpg"，以对应音频资源同名，加上三位整数表示出现的时间点，
        如 000_000.jpg, 000_150.jpg, 000_155.jpg, 000_341.jpg,001_000.jpg,003_000.jpg等
    - 6.点击工具启动后，工具会自动检测所有平级目录及其下级目录，获取所有 "ls_" 开头的目录作为合成源
        检测成功后，工具会自动创建一个平级的MP4/目录,逐级放置合成后的视频
    - 7.如果参与合成的多媒体源文件包含一个mp4文件，则输出 mp4 文件的帧率为28，否则帧率为1
    ```