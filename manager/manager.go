package manager

import (
	"GoGraphDb/conf"
	"GoGraphDb/interpreter"
	"GoGraphDb/log"
	"GoGraphDb/transaction"
	"GoGraphDb/memory_cache"
	"GoGraphDb/utils"
	"context"
)

func init() {
	go GcTransactionCount()
}

func GcTransactionCount() {
	i := 0
	for {
		<- interpreter.TransactionCounter
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

	transaction.StopTheWorld = make(chan int)
	Wait()
	log.CtxInfo(context.Background(),"all transaction done time: %+v", utils.GenTimeStamp())
	log.UndoBegin()
	memory_cache.VertexTree.Print()
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