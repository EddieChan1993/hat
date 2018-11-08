package version

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type AppVersion struct {
	Model    string `json:"model"`
	Version  string `json:"version"`
	DateNow  string `json:"date_now"`
	Branch   string `json:"branch"`
	CommitId string `json:"commit_id"`
	IsUsed   bool   `json:"is_used"`   //是否正在使用
	IsStatus bool   `json:"is_status"` //当前所处版本
}

//获取所有dev模式下的版本记录
func GetVerLog(mode, cmd string) {
	IsExtraMain()
	fileName, file := getLogFilePullPath("version", "app", cmd)
	defer file.Close()
	av := jsonRead(fileName)
	fmt.Println(mode)
	fmt.Printf("%2s%s%9s%9s%9s%9s%7s%20s\n", "", "版本号", "提交ID", "分支", "当前版本", "正在使用", "时间", "模式")
	if mode == VER_ALL {
		for _, v := range av {
			fmt.Printf("%2s%-11s%-13s%-9s%-13t%-13t%-22s%s\n", "", v.Version, v.CommitId, v.Branch, v.IsStatus, v.IsUsed, v.DateNow, v.Model)
		}
	} else if mode == VER_LAST_ONE {
		if len(av)==0{
			fmt.Println("暂无版本记录")
		}else{
			v := av[len(av)-1]
			fmt.Printf("%2s%-11s%-13s%-9s%-13t%-13t%-22s%s\n", "", v.Version, v.CommitId, v.Branch, v.IsStatus, v.IsUsed, v.DateNow, v.Model)
		}
	} else {
		for _, v := range av {
			if v.Model == mode {
				fmt.Printf("%2s%-11s%-13s%-9s%-13t%-13t%-22s%s\n", "", v.Version, v.CommitId, v.Branch, v.IsStatus, v.IsUsed, v.DateNow, v.Model)
			}
		}
	}
}

//记录运行程序
func WriteStart(cmd string) {
	fileName, file := getLogFilePullPath("version", "app", cmd)
	defer file.Close()
	av := jsonRead(fileName)
	switchStart(av)
	data, err := json.MarshalIndent(av, "", "	 ")
	if err != nil {
		log.Fatalln(err)
	}
	jsonWriteReal(fileName, data)
	//jsonWrite(file, data)
}

//停止应用
func WriteStop(cmd string) {
	fileName, file := getLogFilePullPath("version", "app", cmd)
	defer file.Close()
	av := jsonRead(fileName)
	switchStop(av)
	data, err := json.MarshalIndent(av, "", "	 ")
	if err != nil {
		log.Fatalln(err)
	}
	jsonWriteReal(fileName, data)
	//jsonWrite(file, data)
}

//修改运行版本状态
func switchStart(appV []AppVersion) {
	count := len(appV)
	for i := 0; i < count; i++ {
		if appV[i].IsStatus == true {
			appV[i].IsUsed = true
		} else {
			appV[i].IsUsed = false
		}
	}
}

//修改运行版本状态
func switchStop(appV []AppVersion) {
	count := len(appV)
	for i := 0; i < count; i++ {
		if appV[i].IsUsed == true {
			appV[i].IsUsed = false
		}
	}
}

//修改当前编译版本
func switchStatus(appV []AppVersion, av AppVersion) ([]AppVersion, string) {
	version := av.Version
	count := len(appV)
	flag := false //是否是之前提交ID
	for i := 0; i < count; i++ {
		if appV[i].CommitId == av.CommitId && appV[i].Model == av.Model {
			//模式和提交版本存在一致
			//标记为当前所处版本
			appV[i].IsStatus = true
			appV[i].DateNow = av.DateNow
			flag = true
			version = appV[i].Version
		} else {
			appV[i].IsStatus = false
		}
	}

	if !flag {
		//新的提交记录
		appV = append(appV, av)
	}
	return appV, version
}

//获取当前编译版本
func (this *AppVersion) GetVersion(cmd string) string {
	fileName, file := getLogFilePullPath("version", "app", cmd)
	defer file.Close()
	u := jsonRead(fileName)
	this.isExtraVersion(u)
	av := AppVersion{
		Model:    this.Model,
		Version:  this.Version,
		DateNow:  this.DateNow,
		Branch:   this.Branch,
		CommitId: this.CommitId,
		IsStatus: true,
		IsUsed:   false}

	_, version := switchStatus(u, av)
	return version
}

//版本日志记录
func (this *AppVersion) WriteVersion(cmd string) string {
	fileName, file := getLogFilePullPath("version", "app", cmd)
	defer file.Close()
	u := jsonRead(fileName)
	this.isExtraVersion(u)
	av := AppVersion{
		Model:    this.Model,
		Version:  this.Version,
		DateNow:  this.DateNow,
		Branch:   this.Branch,
		CommitId: this.CommitId,
		IsStatus: true,
		IsUsed:   false}

	u, version := switchStatus(u, av)
	data, err := json.MarshalIndent(u, "", "	 ")
	if err != nil {
		log.Fatalln(err)
	}
	jsonWriteReal(fileName, data)
	//jsonWrite(file, data)
	return version
}

//是否已经使用当前版本
func (this *AppVersion) isExtraVersion(av []AppVersion) {
	for _, v := range av {
		if v.Model == this.Model && v.CommitId == this.CommitId {
			break
		}
		if v.Model == this.Model && v.Version == this.Version {
			fmt.Println(" 版本号冲突,已发布程序记录:\n【", this.Version, "】【", this.CommitId, "】【", this.Model, "】")
			this.getAllVersion(av)
			os.Exit(1)
		}
	}
}

//获取所有版本
func (this *AppVersion) getAllVersion(av []AppVersion) {
	fmt.Println("")
	fmt.Printf("%2s%s%9s%9s%9s%9s%7s\n", "", "版本号", "提交ID", "分支", "当前版本", "正在使用", "时间")
	for _, v := range av {
		if v.Model == this.Model {
			fmt.Printf("%2s%-11s%-13s%-9s%-13t%-13t%s\n", "", v.Version, v.CommitId, v.Branch, v.IsStatus, v.IsUsed, v.DateNow)
		}
	}
}

func IsExtraMain() {
	c := "pwd"
	out, _ := ExecShellRes(c)
	out = strings.Replace(out, "\n", "", -1)

	out = fmt.Sprintf("%s/%s", out, "main.go")
	_, err := os.Stat(out)
	if err != nil {
		fmt.Println("入口文件不存在")
		os.Exit(1)

	}
}