package vers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const COMMAND_B_DEV = "dev"
const COMMAND_B_PROD = "prod"
const COMMAND_START = "start"
const COMMAND_STATUS = "status"
const COMMAND_RESTART = "restart"
const COMMAND_STOP = "stop"
const COMMAND_HELP = "help"
const COMMAND_VER_DEV = "ver_dev"
const COMMAND_VER_PROD = "ver_prod"
const COMMAND_VERS = "vers"
const COMMAND_VER = "ver"

const (
	VER_PROD     = "开发模式"     //生产
	VER_DEV      = "生产模式"      //开发
	VER_LAST_ONE = "最新版本" //最后一个版本
	VER_ALL      = "全部版本"      //所有版本
)

//分支
func Branch() string {
	cmdStr := `git rev-parse --abbrev-ref HEAD`
	branch, err := ExecShellRes(cmdStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return strings.Replace(branch, "\n", "", -1)
}

//提交序列号
func CommitId() string {
	cmdStr := `git log --pretty=format:"%h" -1`
	commitId, err := ExecShellRes(cmdStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return commitId
}

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出，有返回值
func ExecShellRes(s string) (r string, err error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("sh", "-c", s)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err = cmd.Run()

	return out.String(), err
}

//需要对shell标准输出的逐行实时进行处理的
func ExecCommand(s string) {
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
func ExecShell(s string) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("sh", "-c", s)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	checkErr(err, "")
}

//查看运行状态
func showStatus() {
	c := "tail -f nohup.out"
	ExecCommand(c)
}

//异常处理
func checkErr(err error, out string) {
	if err != nil {
		showStatus()
		os.Exit(1)
	}
	if out == "" {
		fmt.Println("success")
		fmt.Println("=======")
	} else {
		fmt.Println(string(out))
	}
}

//获取项目名
//sys 考虑当前操作系统
func Folder() string {
	c := "basename $PWD"
	out, _ := ExecShellRes(c)
	out = strings.Replace(out, "\n", "", -1)
	if runtime.GOOS == `windows` {
		out = fmt.Sprintf("%s.exe", out)
	}
	return out
}

//编译加载进度
func Spinner(delay time.Duration, title string) {
	fmt.Printf("%s\n", title)
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
