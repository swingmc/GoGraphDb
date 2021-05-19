package log

import (
	"fmt"
	"context"
	"GoGraphDb/conf"
	"GoGraphDb/utils"
	"os"
)

var (
	infoLogFile  *os.File
	warnLogFile  *os.File
	errorLogFile *os.File
	UndoLogFile  *os.File
)

func init() {
	var err error
	infoLogFile, err = os.OpenFile(conf.ProjectRootPath + conf.InfoLogPath, os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	warnLogFile, err = os.OpenFile(conf.ProjectRootPath + conf.WarnLogPath, os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	errorLogFile, err = os.OpenFile(conf.ProjectRootPath + conf.ErrorLogPath, os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	UndoLogFile, err = os.OpenFile(conf.ProjectRootPath + conf.UndoLogPath, os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Hello World")
}

func CtxInfo(ctx context.Context, format string, obj ...interface{}) {
	infoLogFile.WriteString(utils.LogTime() + conf.Splitor + fmt.Sprintf(format, obj...) + conf.Splitor + utils.CtxToString(ctx) + "\n")
}

func CtxWarn(ctx context.Context, format string, obj ...interface{}){
	warnLogFile.WriteString(utils.LogTime() + conf.Splitor + fmt.Sprintf(format, obj...) + "\n")
}

func CtxError(ctx context.Context, format string, obj ...interface{}){
	errorLogFile.WriteString(utils.LogTime() + conf.Splitor + fmt.Sprintf(format, obj...) + "\n")
}


func UndoBegin() {
	UndoLogFile.WriteString(conf.InterpreterCommand_StartTransaction + conf.Splitor + utils.LogTime() + "\n")
}

func UndoCommit() {
	UndoLogFile.WriteString(conf.InterpreterCommand_EndTransaction + conf.Splitor + utils.LogTime() + "\n")
}