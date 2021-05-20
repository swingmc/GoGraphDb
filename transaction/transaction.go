package transaction

import (
	"GoGraphDb/conf"
	"GoGraphDb/db/db_model"
	"GoGraphDb/utils"
)

var (
	StopTheWorld = make(chan int)
)

type Transaction struct {
	Version        int64
	Status         int32
	ReadOnly       bool
	Block          chan int
	VertexBind     map[string]int64
	EdgeBind       map[string]int64
	VertexSetBind  map[string](map[int64]*db_model.Vertex)
	EdgeSetBind  map[string](map[int64]*db_model.Edge)
}

func NewTransaction() *Transaction{
	//并发控制 数据落盘时避免新的事务开始
	<-StopTheWorld
	t := &Transaction{
		Version: utils.GenTimeStamp(),
		VertexBind: map[string]int64{},
		EdgeBind: map[string]int64{},
		Block: make(chan int, 10),
	}
	addTransaction(t.Version, t)
	//事务计数
	TransactionCounter <- 0
	return t
}

func NewReadOnlyTransaction() *Transaction {
	return &Transaction{
		Version:    utils.GenTimeStamp(),
		ReadOnly:   true,
		VertexBind: map[string]int64{},
		EdgeBind:   map[string]int64{},
	}
}

func (t *Transaction) End() error{
	t.Status = conf.TransactionStatus_Complete
	close(t.Block)
	return nil
}

func (t *Transaction) RollBack() error{
	t.Status = conf.TransactionStatus_Canceled
	close(t.Block)
	return nil
}

func (t *Transaction) IsVertex(str string) bool{
	_, ok := t.VertexBind[str]
	if !ok {
		return false
	}
	return true
}

func (t *Transaction) IsEdge(str string) bool{
	_, ok := t.EdgeBind[str]
	if !ok {
		return false
	}
	return true
}