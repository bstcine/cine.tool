
### 前置任务: 预先安装ffmpeg （视频合成需要）
```shell
第一步: 底部程序栏, 点击"Terminal"打开终端程序


第二步: 安装homebrew，拷贝下面指令，敲"enter"执行
$ ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

第三步: 安装ffmpeg
$ brew install ffmpeg

第四部：查看ffmpeg是否安装成功
ffmpeg -version

至此，即可显示ffmpeg的相关信息，表示安装成功
```

<br>
<br>
<br>

## cine.tool
### Setup

```shell
$ git clone https://github.com/bstcine/cine.tool.git
$ cd cine.tool
```

<br>
      
### Build

#### 0. Build所有工具
```
$ ./gbuild_mac.sh
```

<br>

#### 1. cine_tools
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

<br>

#### 2. 课程资源下载（Lesson的课件学习资源）
  - Window(exe)
    ```
     $ ./gpm.sh
     $ mkdir build
     $ mkdir build/cine_course_download
     $ cp config/cine_course_download.cfg build/cine_course_download/cine_course_download.cfg
     $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/cine_course_download/cine_course_download.exe src/cine_course_download.go
    ```
  
  - Mac
    ```
     $ ./gpm.sh
     $ mkdir build
     $ mkdir build/cine_course_download
     $ cp config/cine_course_download.cfg build/cine_course_download/cine_course_download.cfg
     $ go build -o build/cine_course_download/cine_course_download src/cine_course_download.go
    ```

<br>

- #### 3. 音视频合成
  - Window(exe)
    ```
     $ ./gpm.sh
     $ mkdir build
     $ mkdir build/cine_media_synthesizer
     $ cp config/cine_media_synthesizer.cfg build/cine_media_synthesizer/cine_media_synthesizer.cfg
     $ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/cine_media_synthesizer/cine_media_synthesizer.exe src/cine_media_synthesizer.go
    ```

  - Mac
    ```
     $ ./gpm.sh
     $ mkdir build
     $ mkdir build/cine_media_synthesizer
     $ cp config/cine_media_synthesizer.cfg build/cine_media_synthesizer/cine_media_synthesizer.cfg
     $ go build -o build/cine_media_synthesizer/cine_media_synthesizer src/cine_media_synthesizer.go
    ```
  
  <br>
  
#### 常规使用：资源下载 & 音视频合成
    ```
    - 1. 生成 “资源下载” & “音视频合成” 执行文件。（参考上面 2 & 3 说明）
    - 2. 点击执行 cine_course_download，待课件下载成功后. 然后点击执行 cine_media_synthesizer， 即开始合成。
    - 3. 待执行完毕后，即可在 /MP4/ 文件夹下看到所有合成完毕的课件。
    ```
