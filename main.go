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
	"strings"
	ver "github.com/EddieChan1993/hat/version"
	"runtime"
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
	app := folder()
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
	case ver.COMMAND_STOP:
		stopApp(*appName)
		ver.WriteStop(command)
	case ver.COMMAND_VER_DEV:
		ver.GetVerAllLog(env[command], command)
	case ver.COMMAND_VER_PROD:
		ver.GetVerAllLog(env[command], command)
	case ver.COMMAND_VERS:
		ver.GetVerAllLog("", command)
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

	c := fmt.Sprintf("ps aux | grep \"%s\" | grep -v grep | awk '{print $2}' | xargs -i kill -1 {}", appName)
	execShell(c)
}

//关闭程序
func stopApp(appName string) {
	isExtraAppName(appName)
	isExtraApp(appName)

	c := fmt.Sprintf("ps aux | grep \"%s\" | grep -v grep | awk '{print $2}' | xargs -i kill -9 {}", appName)
	execShell(c)
}

//编译生成开发环境程序
func buildDev(v, appName string) {
	versionName := v
	fmt.Println(versionName)
	buildCond(v, appName)
	v = getBuildVer(v, env[ver.COMMAND_B_DEV], ver.COMMAND_B_DEV)
	v = fmt.Sprintf("v%s", v)
	go spinner(100*time.Millisecond, fmt.Sprintf("正在编译【%s】程序,版本号:%s,程序名称:%s", env[ver.COMMAND_B_DEV], v, appName))
	versionStr := fmt.Sprintf("-X main._version_=%s", v)
	c := fmt.Sprintf("go build -ldflags \"%s\" -o %s", versionStr, appName)
	execShell(c)
	logVersion(versionName, env[ver.COMMAND_B_DEV], ver.COMMAND_B_DEV)
}

//编译生成开发环境程序
func buildProd(v, appName string) {
	versionName := v
	buildCond(v, appName)
	v = getBuildVer(v, env[ver.COMMAND_B_PROD], ver.COMMAND_B_PROD)
	v = fmt.Sprintf("v%s", v)
	go spinner(100*time.Millisecond, fmt.Sprintf("正在编译【%s】程序,版本号:%s,程序名称:%s", env[ver.COMMAND_B_PROD], v, appName))
	versionStr := fmt.Sprintf("-X main._version_=%s", v)
	c := fmt.Sprintf("go build -ldflags \"%s\" -tags=prod -o %s", versionStr, appName)
	execShell(c)
	logVersion(versionName, env[ver.COMMAND_B_PROD], ver.COMMAND_B_PROD)
}

func nohupApp(appName string) {
	//fmt.Println("please CTRL+D")
	isExtraAppName(appName)
	isExtraApp(appName)
	c := fmt.Sprintf("nohup ./%s > %s 2>&1  &", appName,"nohup.out")
	//fmt.Println(c)
	execShell(c)
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
	isExtraMain()
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
	out, _ := execShellRes(c)
	out = strings.Replace(out, "\n", "", -1)

	out = fmt.Sprintf("%s/%s", out, appName)
	_, err := os.Stat(out)
	if err != nil {
		fmt.Println("执行应用文件不存在")
		os.Exit(1)

	}
}

func isExtraMain() {
	c := "pwd"
	out, _ := execShellRes(c)
	out = strings.Replace(out, "\n", "", -1)

	out = fmt.Sprintf("%s/%s", out, "main.go")
	_, err := os.Stat(out)
	if err != nil {
		fmt.Println("入口文件不存在，无法编译")
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
	fmt.Fprintf(os.Stderr, usageStr)
}

//获取项目名
//sys 考虑当前操作系统
func folder() string {
	c := "basename $PWD"
	out, _ := execShellRes(c)
	out = strings.Replace(out, "\n", "", -1)
	if runtime.GOOS == `windows` {
		out = fmt.Sprintf("%s.exe", out)
	}
	return out
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
func execShell(s string) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("sh", "-c", s)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	checkErr(err, "")
}

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出，有返回值
func execShellRes(s string) (r string, err error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("sh", "-c", s)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err = cmd.Run()

	return out.String(), err
}

//异常处理
func checkErr(err error, out string) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if out == "" {
		fmt.Println("success")
		fmt.Println("=======")
	} else {
		fmt.Println(string(out))
	}
}

//获取编译版本
func getBuildVer(v, mode, cmd string) string {
	dateNow := time.Now().Format(YMD_HIS)
	cmdStr := `git rev-parse --abbrev-ref HEAD`
	branch, err := execShellRes(cmdStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	branch = strings.Replace(branch, "\n", "", -1)
	cmdStr = `git log --pretty=format:"%h" -1`
	commitId, err := execShellRes(cmdStr)
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
	cmdStr := `git rev-parse --abbrev-ref HEAD`
	branch, err := execShellRes(cmdStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	branch = strings.Replace(branch, "\n", "", -1)
	cmdStr = `git log --pretty=format:"%h" -1`
	commitId, err := execShellRes(cmdStr)
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

	version := appV.WriteVersion(cmd)
	fmt.Println("版本序列化 ok")
	return version
}
