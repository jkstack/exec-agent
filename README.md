# exec-agent

[![version](https://img.shields.io/github/v/release/jkstack/exec-agent)](https://github.com/jkstack/exec-agent/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/jkstack/exec-agent)](https://goreportcard.com/report/github.com/jkstack/exec-agent)
[![go-mod](https://img.shields.io/github/go-mod/go-version/jkstack/exec-agent)](https://github.com/jkstack/exec-agent)
[![license](https://img.shields.io/github/license/jkstack/exec-agent)](https://www.gnu.org/licenses/agpl-3.0.txt)
![downloads](https://img.shields.io/github/downloads/jkstack/exec-agent/total)

远程脚本执行agent

## 如何编译

1. 下载源代码

       https://github.com/jkstack/exec-agent.git

2. 使用以下命令编译，编译成功后会在当前目录下生成`exec`或`exec.exe`文件

       go build

## linux系统部署

1. 根据当前操作系统下载`deb`或`rpm`安装包，[下载地址](https://github.com/jkstack/exec-agent/releases/latest)
2. 使用`rpm`或`dpkg`命令安装该软件包，程序将被安装到`/opt/exec-agent`目录下
3. 按需修改配置文件，配置文件将被安装在`/opt/exec-agent/conf/agent.conf`目录下，建议修改以下配置内容
   - basic.id: 客户端ID，在该集群下不可重复
   - basic.server: 服务器端地址
4. 使用以下命令启动客户端程序

       /opt/exec-agent/bin/exec-agent start
5. 检查当前客户端是否连接成功

       curl http://<服务端IP>:<端口号(默认13081)>/api/agents/<客户端ID>

## windows系统部署

1. 根据当前操作系统下载`exe`或`msi`安装包，[下载地址](https://github.com/jkstack/exec-agent/releases/latest)
2. 安装该安装包，程序默认会被安装到`C:\Program Files (x86)\exec-agent`目录下
3. 按需修改配置文件，配置文件将默认被安装在`C:\Program Files (x86)\exec-agent\conf\agent.conf`目录下，建议修改以下配置内容
   - basic.id: 客户端ID，在该集群下不可重复
   - basic.server: 服务器端地址
4. 使用以下命令打开系统服务管理面板，找到`exec-agent`服务并启动

       services.msc
5. 检查当前客户端是否连接成功

       curl http://<服务端IP>:<端口号(默认13081)>/api/agents/<客户端ID>

## 开源社区

文档知识库：https://jkstack.github.io/

<img src="docs/wechat_QR.jpg" height=200px width=200px />