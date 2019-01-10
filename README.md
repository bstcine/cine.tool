# cine.tool
- ## Setup

```
前置任务: 打开终端
1. 在程序坞(底部程序栏)中点击LaunchPad，选择 "其他"文件夹，即可找到"终端"程序，点击即可运行


第一步: 安装homebrew，在下面的指令行拷贝到终端中，敲击 enter 键执行
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

第二步: 安装成功后，继续安装 ffmpeg，将下面的指令行拷贝到终端中，敲击 enter 键执行
brew install ffmpeg

第三部：安装结束后，将下面的指令行拷贝到终端中，敲击enter键执行
ffmpeg -version

至此，即可显示ffmpeg的相关信息，表示安装成功
```

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
  - #### 构建
  - Window(exe)
    ```
     $ ./gpm.sh
     $ mkdir build
     $ cp config/cine_media_synthesizer.cfg build/cine_media_synthesizer.cfg
     $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/cine_media_synthesizer.exe src/cine_media_synthesizer.go
    ```

  - Mac
    ```
     $ ./gpm.sh
     $ mkdir build
     $ cp config/cine_media_synthesizer.cfg build/cine_media_synthesizer.cfg
     $ go build -o build/cine_media_synthesizer src/cine_media_synthesizer.go
    ```
     
  - 使用：使用合成工具执行某些课程的合成工作时，需要执行以下步骤
    ```
    - 1.将/build/cine_course_download/cine_course_download, /build/cine_course_download/config,
        /build/cine_media_synthesizer/cine_media_synthesizer, 三个文件拷贝到同一个目录下。
    - 2.点击执行 cine_course_download，课件下载成功后，点击执行cine_media_synthesizer 即可。
    注意: 第一次执行合成的时候，需要执行以下语句
    $ ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
    $ brew install ffmpeg
    ```
