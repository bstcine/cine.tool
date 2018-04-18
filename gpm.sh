#!/bin/bash
# 获取工程目录
echo "==> 正在准备第三方库";
go get github.com/nfnt/resize;
go get github.com/aliyun/aliyun-oss-go-sdk/oss;
echo "<== 第三方库准备完毕"；
