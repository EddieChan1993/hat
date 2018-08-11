package main

import (
	"flag"
	"fmt"
	"os/exec"
	"os"
	"time"
)

const ERROR_MSG = "something error"

const COMMAND_B_DEV = "dev"
const COMMAND_B_PROD = "prod"
const COMMAND_START = "start"
const COMMAND_STATUS = "status"
const COMMAND_RESTART = "restart"
const COMMAND_STOP = "stop"
const COMMAND_HELP = "help"

var env = map[string]string{
	COMMAND_B_DEV:  "开发环境",
	COMMAND_B_PROD: "生产环境",
}

var (
	usageStr, command string
	out               []byte
	err               error
)

func main() {
	version := flag.String("v", "none", "programe's version")
	appName := flag.String("n", "main", "programe's name")

	flag.Parse()
	flag.Usage = usage
	command = flag.Arg(0)

	switch command {
	case COMMAND_B_DEV:
		buildDev(*version, *appName)
		break
	case COMMAND_B_PROD:
		buildProd(*version, *appName)
		break
	case COMMAND_START:
		nohupApp(*appName)
		showStatus()
		break
	case COMMAND_STATUS:
		showStatus()
		break
	case COMMAND_RESTART:
		restartApp(*appName)
		break
	case COMMAND_STOP:
		stopApp(*appName)
		break
	case COMMAND_HELP:
		flag.Usage()
		break
	default:
		flag.Usage()
	}
}

//平滑重启程序
func restartApp(appName string) {
	isExtraAppName(appName)
	c := fmt.Sprintf("ps aux | grep \"%s\" | grep -v grep | awk '{print $2}' | xargs -i kill -1 {}", appName)
	cmd := exec.Command("sh", "-c", c)
	out, err = cmd.Output()

	checkErr(err,[]byte("\nsuccess"))
}

//关闭程序
func stopApp(appName string) {
	isExtraAppName(appName)
	c := fmt.Sprintf("ps aux | grep \"%s\" | grep -v grep | awk '{print $2}' | xargs -i kill {}", appName)
	cmd := exec.Command("sh", "-c", c)
	out, err = cmd.Output()

	checkErr(err,[]byte("\nsuccess"))
}

//编译生成开发环境程序
func buildDev(version, appName string) {
	buildCond(version, appName)
	go spinner(100*time.Millisecond, fmt.Sprintf("正在编译【%s】程序,版本号:%s,程序名称:%s", env[COMMAND_B_DEV], version, appName))
	versionStr := fmt.Sprintf("-X main._version_=%s", version)
	c := fmt.Sprintf("go build -ldflags \"%s\" -o %s", versionStr, appName)
	cmd := exec.Command("sh", "-c", c)
	out, err = cmd.Output()

	checkErr(err,[]byte("\nsuccess"))
}

//编译生成开发环境程序
func buildProd(version, appName string) {
	buildCond(version, appName)
	go spinner(100*time.Millisecond, fmt.Sprintf("正在编译【%s】程序,版本号:%s,程序名称:%s", env[COMMAND_B_DEV], version, appName))
	versionStr := fmt.Sprintf("-X main._version_=%s", version)
	c := fmt.Sprintf("go build -ldflags \"%s\" -tags=prod -o %s", versionStr, appName)
	cmd := exec.Command("sh", "-c", c)
	out, err = cmd.Output()

	checkErr(err,[]byte("\nsuccess"))
}

func nohupApp(appName string) {
	isExtraAppName(appName)
	c := fmt.Sprintf("nohup ./%s &", appName)
	cmd := exec.Command("sh", "-c", c)
	out, err = cmd.Output()

	checkErr(err,[]byte("success\nplease CTRL+Z"))
}

//查看运行状态
func showStatus() {
	c := "tail -f nohup.out"
	cmd := exec.Command("sh", "-c", c)
	out, err = cmd.Output()

	checkErr(err,out)
}

//编译条件
func buildCond(version, appName string) {
	isExtraAppName(appName)
	isExtraVersion(version)
}

func isExtraVersion(version string) {
	if version == ""||version=="none" {
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

func usage() {
	usageStr = "Usage:\n"
	usageStr += "	hat [arguments] command\n"
	usageStr += "The commands are:\n"
	usageStr += fmt.Sprintf("	%s [version_code] %s [app_name|main] %s		create %s's programe\n", "-v", "-n", COMMAND_B_DEV, COMMAND_B_DEV)
	usageStr += fmt.Sprintf("	%s [version_code] %s [app_name|main] %s		create %s's programe\n", "-v", "-n", COMMAND_B_PROD, COMMAND_B_PROD)
	usageStr += fmt.Sprintf("	%s [app_name|main] %s				%s programe\n", "-n", COMMAND_START, COMMAND_START)
	usageStr += fmt.Sprintf("	%s [app_name|main] %s				%s programe\n", "-n", COMMAND_RESTART, COMMAND_RESTART)
	usageStr += fmt.Sprintf("	%s [app_name|main] %s					%s programe\n", "-n", COMMAND_STOP, COMMAND_STOP)
	usageStr += fmt.Sprintf("	%s [app_name|main] %s				%s programe\n", "-n", COMMAND_STATUS, COMMAND_STATUS)
	usageStr += fmt.Sprintf("	%s							look up help\n", COMMAND_HELP)
	fmt.Fprintf(os.Stderr, usageStr)
}

//编译加载进度
func spinner(delay time.Duration, title string) {
	fmt.Printf("%s\n", title)
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func checkErr(err error,out []byte) {
	if err != nil {
		fmt.Println(ERROR_MSG)
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
