package main

import (
	"flag"
	"fmt"
	"os/exec"
	"os"
	"time"
	"io"
	"bufio"
	"bytes"
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
	case COMMAND_B_PROD:
		buildProd(*version, *appName)
	case COMMAND_START:
		nohupApp(*appName)
		showStatus()
	case COMMAND_STATUS:
		showStatus()
	case COMMAND_RESTART:
		restartApp(*appName)
	case COMMAND_STOP:
		stopApp(*appName)
	case COMMAND_HELP:
		flag.Usage()
	default:
		flag.Usage()
	}
}

//平滑重启程序
func restartApp(appName string) {
	isExtraAppName(appName)
	c := fmt.Sprintf("ps aux | grep \"%s\" | grep -v grep | awk '{print $2}' | xargs -i kill -1 {}", appName)
	exec_shell(c)
}

//关闭程序
func stopApp(appName string) {
	isExtraAppName(appName)
	c := fmt.Sprintf("ps aux | grep \"%s\" | grep -v grep | awk '{print $2}' | xargs -i kill {}", appName)
	exec_shell(c)
}

//编译生成开发环境程序
func buildDev(version, appName string) {
	buildCond(version, appName)
	go spinner(100*time.Millisecond, fmt.Sprintf("正在编译【%s】程序,版本号:%s,程序名称:%s", env[COMMAND_B_DEV], version, appName))
	versionStr := fmt.Sprintf("-X main._version_=%s", version)
	c := fmt.Sprintf("go build -ldflags \"%s\" -o %s", versionStr, appName)
	exec_shell(c)
}

//编译生成开发环境程序
func buildProd(version, appName string) {
	buildCond(version, appName)
	go spinner(100*time.Millisecond, fmt.Sprintf("正在编译【%s】程序,版本号:%s,程序名称:%s", env[COMMAND_B_DEV], version, appName))
	versionStr := fmt.Sprintf("-X main._version_=%s", version)
	c := fmt.Sprintf("go build -ldflags \"%s\" -tags=prod -o %s", versionStr, appName)
	exec_shell(c)
}

func nohupApp(appName string) {
	fmt.Println("please CTRL+Z")
	isExtraAppName(appName)
	c := fmt.Sprintf("nohup ./%s &", appName)
	exec_shell(c)
}

//查看运行状态
func showStatus() {
	c := "tail -f nohup.out"
	execCommand(c)
}

//编译条件
func buildCond(version, appName string) {
	isExtraAppName(appName)
	isExtraVersion(version)
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

//需要对shell标准输出的逐行实时进行处理的
func execCommand(s string) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("sh", "-c", s)
	//显示运行的命令
	//fmt.Println(cmd.Args)
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cmd.Start()
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Print(line)
	}
	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	cmd.Wait()
}

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func exec_shell(s string) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("sh", "-c", s)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	checkErr(err, out.String())
}

func checkErr(err error, out string) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if out == "" {
		fmt.Println("\n=======")
		fmt.Println("success")
	} else {
		fmt.Println(string(out))
	}
}
