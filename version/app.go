package version

import (
	"os"
	"log"
	"encoding/json"
	"fmt"
)

type AppVersion struct {
	Model string `json:"model"`
	Version string `json:"version"`
	DateNow string `json:"date_now"`
}

//版本日志记录
func (this *AppVersion) WriteVersion() {
	fileName, file := getLogFilePullPath("version", "app")
	defer file.Close()
	u:=jsonRead(fileName)
	this.isExtraVersion(u)

	av:=AppVersion{this.Model,this.Version,this.DateNow}
	u=append(u,av)
	data,err:=json.MarshalIndent(u, "", "	 ")
	if err !=nil{
		log.Fatalln(err)
	}
	jsonWrite(file,data)
}

//是否已经使用当前版本
func (this *AppVersion)isExtraVersion(av []AppVersion)  {
	for _,v:=range av{
		if v.Version == this.Version &&v.Model==this.Model{
			fmt.Println("【",this.Model,"】版本号冲突")
			getAllVersion(av,this.Model)
			os.Exit(1)
		}
	}
}

//获取所有版本
func getAllVersion(av []AppVersion,mode string) {
	fmt.Println("【",mode,"】所有已用版本：")
	for _,v:=range av{
		if v.Model==mode {
			fmt.Println(v.Version)
		}
	}
}