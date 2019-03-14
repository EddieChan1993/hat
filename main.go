package main

import (
	"flag"
	"fmt"
	"github.com/EddieChan1993/hat/vers"
	"os"
	"strings"
	"time"
)

const YMD_HIS = "2006-01-02 15:04:05"

var (
	usageStr, command string
	env               = map[string]string{
		vers.COMMAND_B_DEV:  "开发模式",
		vers.COMMAND_B_PROD: "生产模式",
	}
)

func main() {
	app := vers.Folder()
	version := flag.String("v", "none", "programe's version")
	appName := flag.String("n", app, "programe's name")
	flag.Parse()
	flag.Usage = usage
	command = flag.Arg(0)

	switch command {
	case vers.COMMAND_B_DEV:
		buildDev(*version, *appName)
	case vers.COMMAND_B_PROD:
		buildProd(*version, *appName)
	case vers.COMMAND_START:
		nohupApp(*appName)
		showStatus()
	case vers.COMMAND_STATUS:
		showStatus()
	case vers.COMMAND_RESTART:
		restartApp(*appName)
		showStatus()
	case vers.COMMAND_STOP:
		stopApp(*appName)
	case vers.COMMAND_VER_DEV:
		vers.GetVerLog(vers.VER_PROD, command)
	case vers.COMMAND_VER_PROD:
		vers.GetVerLog(vers.VER_DEV, command)
	case vers.COMMAND_VERS:
		vers.GetVerLog(vers.VER_ALL, command)
	case vers.COMMAND_VER:
		vers.GetVerLog(vers.VER_LAST_ONE, command)
	case vers.COMMAND_HELP:
		flag.Usage()
	default:
		flag.Usage()
	}
}

//平滑重启程序
func restartApp(appName string) {
	isExtraAppName(appName)
	isExtraApp(appName)

	c := fmt.Sprintf("ps aux | grep -w \"./%s\"  | awk '{print $2}' | xargs -i kill -1 {}  >> nohup.out 2>&1", appName)
	vers.ExecShell(c)
}

//关闭程序
func stopApp(appName string) {
	vers.IsExtraMain()
	isExtraAppName(appName)
	isExtraApp(appName)
	c := fmt.Sprintf("ps aux | grep -w \"./%s\" | awk '{print $2}' | xargs -i kill -9 {}  >> nohup.out 2>&1", appName)
	vers.ExecShell(c)
}

func load(buildEnv, appName, v string) {
	fmt.Println(buildEnv)
	brand := vers.Branch()
	commitId := vers.CommitId()
	vers.Spinner(100*time.Millisecond, fmt.Sprintf("正在编译【%s】程序\n分支:%s,提交ID:%s\n版本号:%s,程序名称:%s", env[buildEnv], brand, commitId, v, appName))
}

//编译生成开发环境程序
func buildDev(v, appName string) {
	versionName := v
	buildCond(v, appName)
	v = getBuildVer(v, env[vers.COMMAND_B_DEV], vers.COMMAND_B_DEV)
	v = fmt.Sprintf("v%s", v)
	go load(vers.COMMAND_B_DEV, appName, v)
	versionStr := fmt.Sprintf("-X main._version_=%s", v)
	c := fmt.Sprintf("go build -ldflags \"%s\" -o %s >> nohup.out 2>&1", versionStr, appName)
	vers.ExecShell(c)
	logVersion(versionName, env[vers.COMMAND_B_DEV], vers.COMMAND_B_DEV)
}

//编译生成开发环境程序
func buildProd(v, appName string) {
	versionName := v
	buildCond(v, appName)
	v = getBuildVer(v, env[vers.COMMAND_B_PROD], vers.COMMAND_B_PROD)
	v = fmt.Sprintf("v%s", v)
	go load(vers.COMMAND_B_PROD, appName, v)
	versionStr := fmt.Sprintf("-X main._version_=%s", v)
	c := fmt.Sprintf("go build -ldflags \"%s\" -tags=prod -o %s >> nohup.out 2>&1", versionStr, appName)
	vers.ExecShell(c)
	logVersion(versionName, env[vers.COMMAND_B_PROD], vers.COMMAND_B_PROD)
}

func nohupApp(appName string) {
	isExtraAppName(appName)
	isExtraApp(appName)
	fmt.Println("HTTP:please CTRL+Z")
	fmt.Println("WS:please CTRL+C")
	c := fmt.Sprintf("nohup ./%s >> %s 2>&1 &", appName, "nohup.out")
	//fmt.Println(c)
	vers.ExecShell(c)
}

//查看运行状态
func showStatus() {
	c := "tail -f nohup.out"
	vers.ExecCommand(c)
}

//编译条件
func buildCond(version, appName string) {
	isExtraAppName(appName)
	isExtraVersion(version)
	vers.IsExtraMain()
}

func isExtraVersion(version string) {
	if version == "" || version == "none" {
		fmt.Println("version is none")
		usage()
		os.Exit(1)
	}
}

func isExtraAppName(appName string) {
	if appName == "" {
		fmt.Println("appName is none")
		usage()
		os.Exit(1)
	}
}

//是否存在执行应用
func isExtraApp(appName string) {
	c := "pwd"
	out, _ := vers.ExecShellRes(c)
	out = strings.Replace(out, "\n", "", -1)

	out = fmt.Sprintf("%s/%s", out, appName)
	_, err := os.Stat(out)
	if err != nil {
		fmt.Println("执行应用文件不存在")
		os.Exit(1)

	}
}

func usage() {
	usageStr = "Usage:\n"
	usageStr += "\n"
	usageStr += "	hat [arguments] command\n"
	usageStr += "\n"
	usageStr += "The commands are:\n"
	usageStr += "\n"
	usageStr += fmt.Sprintf("	%s [version_code] %s [app_name] %s	create %s's program  and eg version_code=1.0\n", "-v", "-n", vers.COMMAND_B_DEV, vers.COMMAND_B_DEV)
	usageStr += fmt.Sprintf("	%s [version_code] %s [app_name] %s	create %s's program\n", "-v", "-n", vers.COMMAND_B_PROD, vers.COMMAND_B_PROD)
	usageStr += fmt.Sprintf("	%s [app_name] %s %25s program and default app_name=basename $PWD,next eq\n", "-n", vers.COMMAND_START, vers.COMMAND_START)
	usageStr += fmt.Sprintf("	%s [app_name] %s %25s program\n", "-n", vers.COMMAND_RESTART, vers.COMMAND_RESTART)
	usageStr += fmt.Sprintf("	%s [app_name] %s %25s program\n", "-n", vers.COMMAND_STOP, vers.COMMAND_STOP)
	usageStr += fmt.Sprintf("	%s [app_name] %s %25s program\n", "-n", vers.COMMAND_STATUS, vers.COMMAND_STATUS)
	usageStr += fmt.Sprintf("	%-27s%25s\n", vers.COMMAND_HELP, "look up help")
	usageStr += fmt.Sprintf("	%-40s%25s\n", vers.COMMAND_VER_DEV, "look up dev's version log")
	usageStr += fmt.Sprintf("	%-40s%25s\n", vers.COMMAND_VER_PROD, "look up prod's version log")
	usageStr += fmt.Sprintf("	%-40s%23s\n", vers.COMMAND_VERS, "look up all version log")
	usageStr += fmt.Sprintf("	%-40s%23s\n", vers.COMMAND_VER, "look up last one version log")
	fmt.Fprintf(os.Stderr, usageStr)
}

//获取编译版本
func getBuildVer(v, mode, cmd string) string {
	dateNow := time.Now().Format(YMD_HIS)
	cmdStr := `git rev-parse --abbrev-ref HEAD`
	branch, err := vers.ExecShellRes(cmdStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	branch = strings.Replace(branch, "\n", "", -1)
	cmdStr = `git log --pretty=format:"%h" -1`
	commitId, err := vers.ExecShellRes(cmdStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	appV := vers.AppVersion{
		Model:    mode,
		Version:  v,
		DateNow:  dateNow,
		Branch:   branch,
		CommitId: commitId}

	version := appV.GetVersion(cmd)
	return version
}

//序列化版本
func logVersion(v, mode, cmd string) string {
	dateNow := time.Now().Format(YMD_HIS)
	branch := vers.Branch()
	commitId := vers.CommitId()
	appV := vers.AppVersion{
		Model:    mode,
		Version:  v,
		DateNow:  dateNow,
		Branch:   branch,
		CommitId: commitId}

	version := appV.WriteVersion(cmd)
	fmt.Println("版本序列化 ok")
	return version
}
