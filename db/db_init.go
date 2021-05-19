package db

import (
	"GoGraphDb/conf"
	"fmt"
	"os"
)

var(
	DbWriter *os.File
	DbReader *os.File
)


func init() {
	//初始化db
	var err error
	DbWriter, err = os.OpenFile(conf.ProjectRootPath + conf.DataDircPath, os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	DbReader, err = os.OpenFile(conf.ProjectRootPath + conf.DataDircPath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	//关闭文件
	//defer DbWriter.Close()
	//defer DbReader.Close()
}

