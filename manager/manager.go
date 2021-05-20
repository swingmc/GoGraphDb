package manager

import (
	"GoGraphDb/conf"
	"GoGraphDb/log"
	"GoGraphDb/memory_cache"
	"GoGraphDb/transaction"
	"GoGraphDb/utils"
	"context"
	"time"
)

func init() {
	log.CtxInfo(context.Background(),"manager init")
	go GcTransactionCount()
}

func GcTransactionCount() {
	i := 0
	for {
		<- transaction.TransactionCounter
		i = i+1
		if i >= conf.GcTransactionNum{
			err := Flush()
			if err != nil {
				panic(err)
			}
			i = 0
		}
	}
}

func Wait() {
	for _, t := range transaction.TransactionGetter {
		<-t.Block
	}
}

func Flush() error{
	log.CtxInfo(context.Background(),"start flush time: %+v", utils.GenTimeStamp())
	if transaction.StopTheWorld != nil{
		if _ ,isClosed := <- transaction.StopTheWorld; isClosed{
			close(transaction.StopTheWorld)
		}
	}
	transaction.StopTheWorld = make(chan int)
	//沉睡100ms防止死锁
	time.Sleep(100*time.Microsecond)
	Wait()
	log.CtxInfo(context.Background(),"all transaction done time: %+v", utils.GenTimeStamp())
	log.UndoBegin()
	//memory_cache.VertexTree.Print()
	memory_cache.VertexTree.Flush()
	//memory_cache.EdgeTree.Flush()
	log.UndoCommit()
	memory_cache.VertexTree.CleanStatus()
	memory_cache.EdgeTree.CleanStatus()
	//清空事务
	transaction.CleanTransaction()
	memory_cache.RefreshBloomFilter()
	close(transaction.StopTheWorld)
	log.CtxInfo(context.Background(),"end flush time: %+v", utils.GenTimeStamp())
	return nil
}