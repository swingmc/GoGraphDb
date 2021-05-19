package test

import (
	"GoGraphDb/manager"
	"GoGraphDb/memory_cache"
	"GoGraphDb/utils"
	"testing"
	"time"
)

func TestCreateVertex(t *testing.T){
	i := 0
	for i < 10 {
		a := i
		go func() {
			t.Log(a)
			memory_cache.CreateVertex(utils.GenTimeStamp())
		}()
		i++
	}
	time.Sleep(time.Second)
	err := manager.Flush()
	memory_cache.VertexTree.Print()
	if err != nil{
		t.Error(err)
	}
	time.Sleep(time.Second)
}