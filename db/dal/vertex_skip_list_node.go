package dal

import (
	"GoGraphDb/conf"
	"GoGraphDb/db/db_model"
	errType "GoGraphDb/error"
	"GoGraphDb/log"
	"GoGraphDb/transaction"
	"context"
	"sync/atomic"
	"errors"
	"fmt"
	"unsafe"
)

type VertexSkipListNode struct {
	score             int64
	vertex            *db_model.Vertex
	VersionId         int64
	LastReadVersionId int64
	t                 *transaction.Transaction
	changed           int32
	next              *VertexSkipListNode
	pre               *VertexSkipListNode
	up                *VertexSkipListNode
	down              *VertexSkipListNode
	NextVersion       *VertexSkipListNode
}

func NewVertexNode(versionId int64, id int64, vertexType int32) *VertexSkipListNode{
	v := &VertexSkipListNode{
		score:     id,
		vertex:    db_model.NewVertex(id, vertexType),
		VersionId: versionId,
		t:         transaction.GetTransaction(versionId),
		changed:   conf.Modify_Create,
	}
	return v
}

func (node *VertexSkipListNode) CreateNextVersionNode(versionId int64, vertex *db_model.Vertex) error{
	curNode := node
	var rearNode *VertexSkipListNode
	//查找有没有更新的版本在写入，有的话回滚
	for curNode != nil{
		switch curNode.Writeable(versionId) {
		case conf.DataWriteableStatus_OneTransaction:
			{
				curNode.vertex.Idntifier = vertex.Idntifier
				curNode.vertex.VertexType = vertex.VertexType
				curNode.vertex.InE = vertex.InE
				curNode.vertex.OutE = vertex.OutE
				curNode.vertex.Properties = vertex.Properties
				return nil
			}
		case conf.DataReadableStatus_VersionTooLate:
			return errType.VertexDataVersionTooLate
		case conf.DataWriteableStatus_Writeable:
			{
				rearNode = curNode
				curNode = curNode.NextVersion
				continue
			}
		case conf.DataWriteableStatus_Executing:
			{
				//阻塞，并发控制
				<-node.t.Block
				continue
			}
		}
	}
	if rearNode == nil{
		return errors.New(fmt.Sprintf("create next version error, rear node nil, version: %+v, id: %+v", versionId, node.VersionId))
	}
	/*
	//并发控制 只能有一个事务在链后加版本
	if rearNode.writeMutex == nil {
		return errors.New(fmt.Sprintf("create next version error, writeMutex nil, version: %+v, id: %+v", versionId, node.VersionId))
	}
	<-rearNode.writeMutex
	if rearNode.NextVersion != nil && rearNode.NextVersion.t.Status != conf.TransactionStatus_Canceled {

		return errors.New(fmt.Sprintf("create next version wrong, concurrence problem, version: %+v, id: %+v", versionId, node.VersionId))
	}
	 */
	newNode := &VertexSkipListNode{
		score:             rearNode.score,
		vertex:            vertex,
		VersionId:         versionId,
		t:                 transaction.GetTransaction(versionId),
		changed:           conf.Modify_Changed,
	}
	ok := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(rearNode.NextVersion)), nil, unsafe.Pointer(newNode))
	//原子操作失败重试
	if !ok {
		log.CtxWarn(context.Background(), "create version atomic retry")
		return node.CreateNextVersionNode(versionId, vertex)
	}
	return nil
}

func (node *VertexSkipListNode) Remove(versionId int64) error{
	curNode := node
	var rearNode *VertexSkipListNode
	//查找有没有更新的版本在写入，有的话回滚
	for curNode != nil{
		switch curNode.Writeable(versionId) {
		case conf.DataWriteableStatus_OneTransaction:
			{
				curNode.changed = conf.Modify_Removed
				curNode.vertex = nil
				return nil
			}
		case conf.DataReadableStatus_VersionTooLate:
			return errType.VertexDataVersionTooLate
		case conf.DataWriteableStatus_Writeable:
			{
				rearNode = curNode
				curNode = curNode.NextVersion
				continue
			}
		case conf.DataWriteableStatus_Executing:
			{
				//阻塞，并发控制
				<-node.t.Block
				continue
			}
		}
	}
	if rearNode == nil{
		return errors.New(fmt.Sprintf("remove node error, rear node nil, version: %+v, id: %+v", versionId, node.VersionId))
	}
	newNode := &VertexSkipListNode{
		score:             rearNode.score,
		vertex:            nil,
		VersionId:         versionId,
		t:                 transaction.GetTransaction(versionId),
		changed:           conf.Modify_Removed,
	}
	ok := atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(rearNode.NextVersion)), nil, unsafe.Pointer(newNode))
	//原子操作失败重试
	if !ok {
		log.CtxWarn(context.Background(), "remove node atomic retry")
		return node.Remove(versionId)
	}
	return nil
}

func (node *VertexSkipListNode) Read(versionId int64) *VertexSkipListNode{
	//找到版本号前最后一个已提交事务的版本,如果有更新的已写入完成的版本
	curNode := node
	var nodeReaded *VertexSkipListNode
	for curNode != nil {
		switch curNode.Readable(versionId){
		case conf.DataReadableStatus_OneTransaction:
			return curNode
		case conf.DataReadableStatus_VersionTooLate:
			{
				if nodeReaded == nil {
					log.CtxWarn(context.Background(), "read nil node")
					return nil
				}
				nodeReaded.LastReadVersionId = versionId
				return nodeReaded
			}
		case conf.DataReadableStatus_Readable:
			{
				nodeReaded = curNode
				curNode = curNode.NextVersion
				continue
			}
		case conf.DataReadableStatus_Canceled:
			{
				curNode = curNode.NextVersion
				continue
			}
		case conf.DataReadableStatus_Executing:
			{
			//并发控制 在事务未执行完成前阻塞
			<-curNode.t.Block
			//重来，再次判断
			continue
			}
		}
	}
	if nodeReaded == nil {
		log.CtxWarn(context.Background(), "read nil node")
		return nil
	}
	nodeReaded.LastReadVersionId = versionId
	return nodeReaded
}

func (node *VertexSkipListNode) Readable(versionId int64) int32{
	if node.VersionId == versionId {
		return conf.DataReadableStatus_OneTransaction
	}
	if node.VersionId <= versionId && (node.t == nil || node.t.Status == conf.TransactionStatus_Complete) {
		return conf.DataReadableStatus_Readable
	}
	if node.t != nil && node.t.Status == conf.TransactionStatus_Canceled {
		return conf.DataReadableStatus_Canceled
	}
	if node.VersionId > versionId {
		return conf.DataReadableStatus_VersionTooLate
	}
	if node.t != nil && node.t.Status == conf.TransactionStatus_Executing {
		return conf.DataReadableStatus_Executing
	}
	err := errors.New(fmt.Sprintf("wrong data readble status! node data: %+v", *node))
	panic(err)
}

func (node *VertexSkipListNode) Writeable(versionId int64) int32{
	if node.VersionId == versionId {
		return conf.DataWriteableStatus_OneTransaction
	}
	if node.VersionId <= versionId && node.LastReadVersionId <= versionId && (node.t == nil || node.t.Status != conf.TransactionStatus_Executing) {
		return conf.DataWriteableStatus_Writeable
	}
	if node.t != nil && node.t.Status == conf.TransactionStatus_Canceled {
		return conf.DataWriteableStatus_Writeable
	}
	if node.VersionId > versionId || node.LastReadVersionId > versionId{
		return conf.DataWriteableStatus_VersionTooLate
	}
	if node.t != nil && node.t.Status == conf.TransactionStatus_Executing {
		return conf.DataWriteableStatus_Executing
	}
	err := errors.New(fmt.Sprintf("wrong data writeble status! node data: %+v, transaction: %+v", *node, *node.t))
	panic(err)
}

func (node *VertexSkipListNode) FindLatestVersion() *VertexSkipListNode{
	curNode := node
	var nodeReaded *VertexSkipListNode
	for curNode != nil {
		if curNode.t == nil || curNode.t.Status != conf.TransactionStatus_Canceled{
			nodeReaded = curNode
		}
		curNode = curNode.NextVersion
	}
	if nodeReaded == nil{
		log.CtxWarn(context.Background(), "find latest version nil node")
		return nil
	}
	if nodeReaded.t != nil {
		<-nodeReaded.t.Block
	}
	return nodeReaded
}