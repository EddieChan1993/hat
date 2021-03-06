package vers

import (
	"fmt"
	"os"
	"log"
	"encoding/json"
	"io/ioutil"
)

const mb int64 = 1 << (10 * 2)

var (
	logSavePath = "runtime"
	logFileExt  = "json"
)

func getLogFilePath(logFileName string) string {
	return fmt.Sprintf("%s/%s", logSavePath, logFileName)
}

//获取全路径
func getLogFilePullPath(logPathName, logFileName, cmd string) (string, *os.File) {
	prefixPath := getLogFilePath(logPathName)
	suffixPath := fmt.Sprintf("%s.%s", logFileName, logFileExt)

	filePath := fmt.Sprintf("%s/%s", prefixPath, suffixPath)
	file := openLogFile(logPathName, filePath, cmd)
	return filePath, file
}

//判断文件路径是否正确
func openLogFile(logPathName, filePath, cmd string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		if cmd == COMMAND_VER_DEV || cmd == COMMAND_VER_PROD || cmd == COMMAND_VERS {
			fmt.Println("不存在版本文件，请查看当前路径")
			os.Exit(1)
		} else {
			mkDir(getLogFilePath(logPathName))
		}
	case os.IsPermission(err):
		log.Fatalf("Permission:%v", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatalf("Fail to CreateFile:%v", err)
	}
	return file
}

//创建目录
func mkDir(filePath string) {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+filePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

//读取json文件
func jsonRead(filename string) []AppVersion {
	var appV []AppVersion
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}

	//读取的数据为json格式，需要进行解码
	json.Unmarshal(data, &appV)
	return appV
}

//写入json文件
func jsonWrite(fp *os.File, data []byte) {
	_, err := fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func jsonWriteReal(fp string, data []byte) {
	err:=ioutil.WriteFile(fp,data,os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
}
