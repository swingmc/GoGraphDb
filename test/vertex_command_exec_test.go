package test

import (
	_"GoGraphDb/db"
	_"GoGraphDb/manager"
	"GoGraphDb/interpreter"
	"os"
	"runtime"
	"testing"
	"time"
)

var (
	file1Path = runtime.GOROOT() + "/src/GoGraphDb/script/test1.txt"
)

func TestExecTestFile1(t *testing.T){
	i := interpreter.Interpreter{}
	file1, err := os.OpenFile(file1Path, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		t.Error(err)
	}
	err = i.ExeDmlFile(file1)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second)
}