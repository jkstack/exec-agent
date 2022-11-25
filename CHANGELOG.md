# CHANGELOG

## 1.0.0

1. 实现文件上传下载功能
2. 实现脚本或命令执行功能

## 1.0.1

1. 实现exec中的deferRM功能
2. 修改部分log内容
3. 修正windows下进程重复启动问题
4. go版本升级到1.19.3
5. 升级第三方库版本

## 1.0.2

1. 修正deb和rpm安装包中的epoch问题
2. 增加manifest.yaml配置文件描述信息

## 1.0.3

1. 修改配置文件中的配置项，增加basic前缀
2. 修正manifest.yaml中的id字段缺少类型的问题
3. 修正manifest.yaml中的log.target字段默认值类型问题

## 1.0.4

修改命令行交互方式

## 1.0.5

修正deb和rpm包安装时注册系统服务的问题

## 1.0.6

1. 修正windows下的服务启动问题
2. 去除msi安装包