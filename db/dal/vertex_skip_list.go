package dal

import (
	"GoGraphDb/conf"
	"GoGraphDb/db/db_model"
	"GoGraphDb/log"
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const UP_LEVELS_ABILITY = 500
const UP_LEVELS_TOTAL = 1000

const (
	Nochange = 0
	Create	 = 1
	Changed  = 2
	Removed  = 3
)

type VertexSkipList struct {
	head   *VertexSkipListNode
	tail   *VertexSkipListNode
	size   int
	levels int
}

func NewVertexSkipList() *VertexSkipList {
	sl := new(VertexSkipList)
	sl.head = new(VertexSkipListNode)
	sl.tail = new(VertexSkipListNode)
	sl.head.score = math.MinInt64
	sl.tail.score = math.MaxInt64

	sl.head.next = sl.tail
	sl.tail.pre = sl.head

	sl.size = 0
	sl.levels = 1

	return sl
}

func (sl *VertexSkipList) Size() int {
	return sl.size
}

func (sl *VertexSkipList) Levels() int {
	return sl.levels
}

func (sl *VertexSkipList) Get(versionId int64, score int64) *db_model.Vertex {
	node := sl.findNode(versionId ,score)
	if node.score == score{
		return node.Read(versionId).vertex
	} else {
		return nil
	}
}

func (sl *VertexSkipList) Insert(versionId int64, score int64, vertex *db_model.Vertex) error{
	f := sl.findNode(versionId, score)
	if f.score == score {
		//f.vertex = vertex
		//f.changed = Changed
		err := f.CreateNextVersionNode(versionId, vertex)
		if err != nil {
			return err
		}
		f.changed = conf.Modify_Changed
		return nil
	}
	//curNode := new(VertexSkipListNode)
	curNode := NewVertexNode(versionId, score, 0)
	curNode.score = score
	curNode.vertex = vertex
	curNode.changed = Create

	sl.insertAfter(f, curNode)

	rander := rand.New(rand.NewSource(time.Now().UnixNano()))

	curlevels := 1
	for rander.Intn(UP_LEVELS_TOTAL) < UP_LEVELS_ABILITY {
		curlevels++
		if curlevels > sl.levels {
			sl.newlevels()
		}

		for f.up == nil {
			f = f.pre
		}
		f = f.up
		tmpNode := &VertexSkipListNode{score: score}

		curNode.up = tmpNode
		tmpNode.down = curNode
		sl.insertAfter(f, tmpNode)

		curNode = tmpNode
	}

	sl.size++
	return nil
}

func (sl *VertexSkipList) Remove(versionId int64, score int64) error{
	f := sl.findNode(versionId, score)
	if f.score != score || f.changed == conf.Modify_Removed{
		return nil
	}

	err := f.Remove(versionId)
	if err != nil{
		return err
	}
	for f != nil {
		f.changed = conf.Modify_Removed
		f = f.up
	}
	return nil
}

func (sl *VertexSkipList) newlevels() {
	nhead := &VertexSkipListNode{score: math.MinInt64}
	ntail := &VertexSkipListNode{score: math.MaxInt64}
	nhead.next = ntail
	ntail.pre = nhead

	sl.head.up = nhead
	nhead.down = sl.head
	sl.tail.up = ntail
	ntail.down = sl.tail

	sl.head = nhead
	sl.tail = ntail
	sl.levels++
}

func (sl *VertexSkipList) insertAfter(pNode *VertexSkipListNode, curNode *VertexSkipListNode) {
	curNode.next = pNode.next
	curNode.pre = pNode
	pNode.next.pre = curNode
	pNode.next = curNode
}

func (sl *VertexSkipList) findNode(versionId int64, score int64) *VertexSkipListNode {
	p := sl.head

	for p != nil {
		if p.score == score {
			if p.down == nil {
				return p
			}
			p = p.down
		} else if p.score < score {
			if p.next.score > score {
				if p.down == nil {
					return p
				}
				p = p.down
			} else {
				p = p.next
			}
		}
	}
	return p.Read(versionId)
}

func (sl *VertexSkipList) Print() {
	mapScore := make(map[int64]int)

	p := sl.head
	for p.down != nil {
		p = p.down
	}
	index := 0
	for p != nil {
		mapScore[p.score] = index
		p = p.next
		index++
	}
	p = sl.head
	for i := 0; i < sl.levels; i++ {
		q := p
		preIndex := 0
		for q != nil {
			s := q.score
			if s == math.MinInt64 {
				fmt.Printf("%s", "BEGIN")
				q = q.next
				continue
			}
			index := mapScore[s]
			c := (index - preIndex - 1) * 12
			for m := 0; m < c; m++ {
				fmt.Print("-")
			}
			if s == math.MaxInt64 {
				fmt.Printf("-->%s\n", "END")
			} else {
				fmt.Printf("-->%9d", s)
				preIndex = index
			}
			q = q.next
			for q != nil && q.changed == Removed {
				q = q.next
			}
		}
		p = p.down
	}
}

func (sl *VertexSkipList) Flush() error{
	log.CtxInfo(context.Background(),"vertexTree undo log start")
	cur := sl.head
	for cur.down != nil {
		cur = cur.down
	}
	for cur != nil {
		latestVersion := cur.FindLatestVersion()
		if latestVersion.changed != Nochange {
			var err error
			if cur.changed != conf.Modify_Nochange {
				err = cur.vertex.FlushAsUndoBase(latestVersion.VersionId, conf.Modify_Create)
			}else{
				err = cur.vertex.FlushAsUndoBase(latestVersion.VersionId, latestVersion.changed)
			}
			if err != nil {
				log.CtxError(context.Background(),"vertex undolog flush error: %+v", err)
				cur = cur.next
				continue
			}
			if latestVersion.vertex == nil {
				err = cur.vertex.Flush(latestVersion.VersionId, conf.Modify_Removed)
				if err != nil {
					log.CtxError(context.Background(), "vertex flush error: %+v", err)
					cur = cur.next
					continue
				}
			}else {
				err = latestVersion.vertex.Flush(latestVersion.VersionId, latestVersion.changed)
				if err != nil {
					log.CtxError(context.Background(), "vertex flush error: %+v", err)
					cur = cur.next
					continue
				}
			}
		}
		cur = cur.next
	}
	log.CtxInfo(context.Background(),"vertexTree undo log end")
	return nil
}

//等所有事务执行完毕
func (sl *VertexSkipList) Wait(){
	p := sl.head
	for p.down != nil {
		p = p.down
	}
	for p != nil {
		t := p.FindLatestVersion().t
		if t != nil{
			<-t.Block
		}
		p = p.next
	}
}

func (sl *VertexSkipList) CleanStatus() {
	p := sl.head
	for p.down != nil {
		p = p.down
	}
	//删去已标记删除的节点
	for p != nil {
		q := p.FindLatestVersion()
		if q.changed == conf.Modify_Removed {
			cur := p
			for cur != nil {
				p.pre.next = p.next
				p.next.pre = p.pre
				cur = cur.up
			}
		}
		p = p.next
	}

	//选取最新版本 作为基础版本
	p = sl.head
	for p.down != nil {
		p = p.down
	}
	for p != nil {
		q := p.FindLatestVersion()
		if q.changed == conf.Modify_Create || q.changed == conf.Modify_Changed {
			q.pre = p.pre
			q.next = p.next
			q.up = p.up
			q.down = p.down

			p.pre.next = q
			p.next.pre = q
			if p.up != nil {
				p.up.down = q
			}
			if p.down != nil {
				p.down.up = q
			}
		}
		p = p.next
	}

	//重置修改状态
	p = sl.head
	for p.down != nil {
		p = p.down
	}
	for p != nil {
		q := p
		for q != nil {
			q.changed = conf.Modify_Nochange
			q.t = nil
			q.VersionId = 0
			q.LastReadVersionId = 0
			q.NextVersion = nil

			q = q.next
		}
		p = p.down
	}
}