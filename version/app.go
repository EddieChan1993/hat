package version

import (
	"os"
	"log"
	"encoding/json"
	"fmt"
)

type AppVersion struct {
	Model    string `json:"model"`
	Version  string `json:"version"`
	DateNow  string `json:"date_now"`
	Branch   string `json:"branch"`
	CommitId string `json:"commit_id"`
}

//获取所有dev模式下的版本记录
func GetVerAllLog(mode string) {
	fileName, file := getLogFilePullPath("version", "app")
	defer file.Close()
	av := jsonRead(fileName)
	fmt.Println(fmt.Sprintf("【%2s】",mode))
	fmt.Println(fmt.Sprintf("%2s%s%9s%9s", "", "版本号", "提交ID", "分支"))
	for _, v := range av {
		if v.Model == mode {
			fmt.Println(fmt.Sprintf("%2s%s%14s%12s", "", v.Version, v.CommitId, v.Branch))
		}
	}
}

//版本日志记录
func (this *AppVersion) WriteVersion() {
	fileName, file := getLogFilePullPath("version", "app")
	defer file.Close()
	u := jsonRead(fileName)
	this.isExtraVersion(u)

	av := AppVersion{this.Model, this.Version, this.DateNow, this.Branch, this.CommitId}
	u = append(u, av)
	data, err := json.MarshalIndent(u, "", "	 ")
	if err != nil {
		log.Fatalln(err)
	}
	jsonWrite(file, data)
}

//是否已经使用当前版本
func (this *AppVersion) isExtraVersion(av []AppVersion) {
	for _, v := range av {
		if v.Model == this.Model && (v.Version == this.Version || v.CommitId == this.CommitId) {
			fmt.Println(" 版本号或提交ID冲突,已发布程序记录:\n【", this.Version, "】【", this.CommitId, "】【", this.Model, "】")
			this.getAllVersion(av)
			os.Exit(1)
		}
	}
}

//获取所有版本
func (this *AppVersion) getAllVersion(av []AppVersion) {
	fmt.Println("")
	fmt.Println(fmt.Sprintf("%2s%s%9s%9s", "", "版本号", "提交ID", "分支"))
	for _, v := range av {
		if v.Model == this.Model {
			fmt.Println(fmt.Sprintf("%2s%s%14s%12s", "", v.Version, v.CommitId, v.Branch))
		}
	}
}


