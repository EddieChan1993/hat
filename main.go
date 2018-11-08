package main

import (
	"flag"
	"fmt"
	ver "github.com/EddieChan1993/hat/version"
	"os"
	"strings"
	"time"
)

const YMD_HIS = "2006-01-02 15:04:05"

var (
	usageStr, command string
	env               = map[string]string{
		ver.COMMAND_B_DEV:  "开发模式",
		ver.COMMAND_B_PROD: "生产模式",
	}
)

func main() {
	app := ver.Folder()
	version := flag.String("v", "none", "programe's version")
	appName := flag.String("n", app, "programe's name")
	flag.Parse()
	flag.Usage = usage
	command = flag.Arg(0)

	switch command {
	case ver.COMMAND_B_DEV:
		buildDev(*version, *appName)
	case ver.COMMAND_B_PROD:
		buildProd(*version, *appName)
	case ver.COMMAND_START:
		nohupApp(*appName)
		ver.WriteStart(command)
		showStatus()
	case ver.COMMAND_STATUS:
		showStatus()
	case ver.COMMAND_RESTART:
		restartApp(*appName)
		ver.WriteStart(command)
		showStatus()
	case ver.COMMAND_STOP:
		stopApp(*appName)
		ver.WriteStop(command)
	case ver.COMMAND_VER_DEV:
		ver.GetVerLog(ver.VER_PROD, command)
	case ver.COMMAND_VER_PROD:
		ver.GetVerLog(ver.VER_DEV, command)
	case ver.COMMAND_VERS:
		ver.GetVerLog(ver.VER_ALL, command)
	case ver.COMMAND_VER:
		ver.GetVerLog(ver.VER_LAST_ONE, command)
	case ver.COMMAND_HELP:
		flag.Usage()
	default:
		flag.Usage()
	}
}

//平滑重启程序
func restartApp(appName string) {
	isExtraAppName(appName)
	isExtraApp(appName)

	c := fmt.Sprintf("ps aux | grep \"%s\" | grep -v grep | awk '{print $2}' | xargs -i kill -1 {}  >> nohup.out 2>&1", appName)
	ver.ExecShell(c)
}

//关闭程序
func stopApp(appName string) {
	ver.IsExtraMain()
	isExtraAppName(appName)
	isExtraApp(appName)
	c := fmt.Sprintf("ps aux | grep \"%s\" | grep -v grep | awk '{print $2}' | xargs -i kill -9 {}  >> nohup.out 2>&1", appName)
	ver.ExecShell(c)
}

func load(buildEnv, appName, v string) {
	fmt.Println(buildEnv)
	brand := ver.Branch()
	commitId := ver.CommitId()
	ver.Spinner(100*time.Millisecond, fmt.Sprintf("正在编译【%s】程序\n分支:%s,提交ID:%s\n版本号:%s,程序名称:%s", env[buildEnv], brand, commitId, v, appName))
}

//编译生成开发环境程序
func buildDev(v, appName string) {
	versionName := v
	buildCond(v, appName)
	v = getBuildVer(v, env[ver.COMMAND_B_DEV], ver.COMMAND_B_DEV)
	v = fmt.Sprintf("v%s", v)
	go load(ver.COMMAND_B_DEV, appName, v)
	versionStr := fmt.Sprintf("-X main._version_=%s", v)
	c := fmt.Sprintf("go build -ldflags \"%s\" -o %s >> nohup.out 2>&1", versionStr, appName)
	ver.ExecShell(c)
	logVersion(versionName, env[ver.COMMAND_B_DEV], ver.COMMAND_B_DEV)
}

//编译生成开发环境程序
func buildProd(v, appName string) {
	versionName := v
	buildCond(v, appName)
	v = getBuildVer(v, env[ver.COMMAND_B_PROD], ver.COMMAND_B_PROD)
	v = fmt.Sprintf("v%s", v)
	go load(ver.COMMAND_B_PROD, appName, v)
	versionStr := fmt.Sprintf("-X main._version_=%s", v)
	c := fmt.Sprintf("go build -ldflags \"%s\" -tags=prod -o %s >> nohup.out 2>&1", versionStr, appName)
	ver.ExecShell(c)
	logVersion(versionName, env[ver.COMMAND_B_PROD], ver.COMMAND_B_PROD)
}

func nohupApp(appName string) {
	isExtraAppName(appName)
	isExtraApp(appName)
	fmt.Println("please CTRL+Z")
	c := fmt.Sprintf("nohup ./%s >> %s 2>&1 &", appName, "nohup.out")
	//fmt.Println(c)
	ver.ExecShell(c)
}

//查看运行状态
func showStatus() {
	c := "tail -f nohup.out"
	ver.ExecCommand(c)
}

//编译条件
func buildCond(version, appName string) {
	isExtraAppName(appName)
	isExtraVersion(version)
	ver.IsExtraMain()
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
	out, _ := ver.ExecShellRes(c)
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
	usageStr += fmt.Sprintf("	%s [version_code] %s [app_name] %s	create %s's program  and eg version_code=1.0\n", "-v", "-n", ver.COMMAND_B_DEV, ver.COMMAND_B_DEV)
	usageStr += fmt.Sprintf("	%s [version_code] %s [app_name] %s	create %s's program\n", "-v", "-n", ver.COMMAND_B_PROD, ver.COMMAND_B_PROD)
	usageStr += fmt.Sprintf("	%s [app_name] %s %25s program and default app_name=basename $PWD,next eq\n", "-n", ver.COMMAND_START, ver.COMMAND_START)
	usageStr += fmt.Sprintf("	%s [app_name] %s %25s program\n", "-n", ver.COMMAND_RESTART, ver.COMMAND_RESTART)
	usageStr += fmt.Sprintf("	%s [app_name] %s %25s program\n", "-n", ver.COMMAND_STOP, ver.COMMAND_STOP)
	usageStr += fmt.Sprintf("	%s [app_name] %s %25s program\n", "-n", ver.COMMAND_STATUS, ver.COMMAND_STATUS)
	usageStr += fmt.Sprintf("	%-27s%25s\n", ver.COMMAND_HELP, "look up help")
	usageStr += fmt.Sprintf("	%-40s%25s\n", ver.COMMAND_VER_DEV, "look up dev's version log")
	usageStr += fmt.Sprintf("	%-40s%25s\n", ver.COMMAND_VER_PROD, "look up prod's version log")
	usageStr += fmt.Sprintf("	%-40s%23s\n", ver.COMMAND_VERS, "look up all version log")
	usageStr += fmt.Sprintf("	%-40s%23s\n", ver.COMMAND_VER, "look up last one version log")
	fmt.Fprintf(os.Stderr, usageStr)
}

//获取编译版本
func getBuildVer(v, mode, cmd string) string {
	dateNow := time.Now().Format(YMD_HIS)
	cmdStr := `git rev-parse --abbrev-ref HEAD`
	branch, err := ver.ExecShellRes(cmdStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	branch = strings.Replace(branch, "\n", "", -1)
	cmdStr = `git log --pretty=format:"%h" -1`
	commitId, err := ver.ExecShellRes(cmdStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	appV := ver.AppVersion{
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
	branch := ver.Branch()
	commitId := ver.CommitId()
	appV := ver.AppVersion{
		Model:    mode,
		Version:  v,
		DateNow:  dateNow,
		Branch:   branch,
		CommitId: commitId}

	version := appV.WriteVersion(cmd)
	fmt.Println("版本序列化 ok")
	return version
}
