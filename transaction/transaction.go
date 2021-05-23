package transaction

import (
	"GoGraphDb/conf"
	"GoGraphDb/db/db_model"
	"GoGraphDb/utils"
	"sync"
)

var (
	StopTheWorld = make(chan int)
)

type Transaction struct {
	Version       int64
	Status        int32
	ReadOnly      bool
	Block         chan int
	DataChan      chan interface{}
	vertexSetBind map[string](map[int64]*db_model.Vertex)
	edgeSetBind   map[string](map[int64]*db_model.Edge)
	vertexLock    sync.RWMutex
	edgeLock      sync.RWMutex
}

func NewTransaction() *Transaction{
	//并发控制 数据落盘时避免新的事务开始
	<-StopTheWorld
	t := &Transaction{
		Version:       utils.GenTimeStamp(),
		vertexSetBind: map[string](map[int64]*db_model.Vertex){},
		edgeSetBind:   map[string](map[int64]*db_model.Edge){},
		Block:         make(chan int, 10),
		DataChan:      make(chan interface{}, 100),
	}
	addTransaction(t.Version, t)
	//事务计数
	TransactionCounter <- 0
	return t
}
/*
func NewReadOnlyTransaction() *Transaction {
	return &Transaction{
		Version:    utils.GenTimeStamp(),
		ReadOnly:   true,
		vertexSetBind: map[string](map[int64]*db_model.Vertex){},
		edgeSetBind: map[string](map[int64]*db_model.Edge){},
	}
}
*/

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
	_, ok := t.vertexSetBind[str]
	if !ok {
		return false
	}
	return true
}

func (t *Transaction) IsEdge(str string) bool{
	_, ok := t.edgeSetBind[str]
	if !ok {
		return false
	}
	return true
}

func (t *Transaction) VertexRead(identifier string) (map[int64]*db_model.Vertex, bool){
	t.vertexLock.RLock()
	vertexSet, ok := t.vertexSetBind[identifier]
	t.vertexLock.RUnlock()
	return vertexSet, ok
}

func (t *Transaction) VertexWrite(identifier string, vertexMap map[int64]*db_model.Vertex) {
	t.vertexLock.Lock()
	t.vertexSetBind[identifier] = vertexMap
	t.vertexLock.Unlock()
}

func (t *Transaction) EdgeRead(identifier string) (map[int64]*db_model.Edge, bool){
	t.edgeLock.RLock()
	edgeSet, ok := t.edgeSetBind[identifier]
	t.edgeLock.RUnlock()
	return edgeSet, ok
}

func (t *Transaction) EdgeWrite(identifier string, edgeMap map[int64]*db_model.Edge) {
	t.edgeLock.Lock()
	t.edgeSetBind[identifier] = edgeMap
	t.edgeLock.Unlock()
}