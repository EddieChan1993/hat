#  hat是辅助于**hatGo**的程序工具


##  安装
```
go get github.com/EddieChan1993/hat
```
## 使用
在执行项目的入口文件下执行下面命令
```

    #帮助文档
    hat help
    #编译成开发程序，带上版本号
    hat -v v1.0 dev
    #编译成生产程序，带上版本号
    hat -v v1.0 prod
    #启动程序
    hat start
    #平滑重启
    hat restart
    #停止程序
    hat stop
    #查看程序状态
    hat status

```
```
$ hat
Usage:

        hat [arguments] command

The commands are:

        -v [version_code] -n [app_name|main] dev                create dev's program
        -v [version_code] -n [app_name|main] prod               create prod's program
        -n [app_name|main] start                                start program
        -n [app_name|main] restart                              restart program
        -n [app_name|main] stop                                 stop program
        -n [app_name|main] status                               status program
        help                                                    look up help
        ver_dev                                                 look up dev's version log
        ver_prod                                                look up prod's version log

```
## 功能特色
版本日志记录，判断冲突版本，同时列出已经使用版本


