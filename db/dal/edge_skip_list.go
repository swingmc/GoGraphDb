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

type EdgeSkipList struct {
	head   *EdgeSkipListNode
	tail   *EdgeSkipListNode
	size   int
	levels int
}

func NewEdgeSkipList() *EdgeSkipList {
	sl := new(EdgeSkipList)
	sl.head = new(EdgeSkipListNode)
	sl.tail = new(EdgeSkipListNode)
	sl.head.score = math.MinInt64
	sl.tail.score = math.MaxInt64

	sl.head.next = sl.tail
	sl.tail.pre = sl.head

	sl.size = 0
	sl.levels = 1

	return sl
}

func (sl *EdgeSkipList) Size() int {
	return sl.size
}

func (sl *EdgeSkipList) Levels() int {
	return sl.levels
}

func (sl *EdgeSkipList) Get(versionId int64, score int64) *db_model.Edge {
	node := sl.findNode(versionId ,score)
	if node.score == score{
		return node.Read(versionId).Edge
	} else {
		return nil
	}
}

func (sl *EdgeSkipList) Insert(versionId int64, score int64, Edge *db_model.Edge) error{
	f := sl.findNode(versionId, score)
	if f.score == score {
		//f.Edge = Edge
		//f.changed = Changed
		err := f.CreateNextVersionNode(versionId, Edge)
		if err != nil {
			return err
		}
		return nil
	}
	//curNode := new(EdgeSkipListNode)
	curNode := NewEdgeNode(versionId, score, 0, 0)
	curNode.score = score
	curNode.Edge = Edge
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
		tmpNode := &EdgeSkipListNode{score: score}

		curNode.up = tmpNode
		tmpNode.down = curNode
		sl.insertAfter(f, tmpNode)

		curNode = tmpNode
	}

	sl.size++
	return nil
}

func (sl *EdgeSkipList) Remove(versionId int64, score int64) error{
	f := sl.findNode(versionId, score)
	if f.score != score || f.changed == conf.Modify_Removed{
		return nil
	}

	err := f.Remove(versionId)
	if err != nil {
		return err
	}
	for f != nil {
		f.changed = conf.Modify_Removed
		f = f.up
	}
	return nil
}

func (sl *EdgeSkipList) newlevels() {
	nhead := &EdgeSkipListNode{score: math.MinInt64}
	ntail := &EdgeSkipListNode{score: math.MaxInt64}
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

func (sl *EdgeSkipList) insertAfter(pNode *EdgeSkipListNode, curNode *EdgeSkipListNode) {
	curNode.next = pNode.next
	curNode.pre = pNode
	pNode.next.pre = curNode
	pNode.next = curNode
}

func (sl *EdgeSkipList) findNode(versionId int64, score int64) *EdgeSkipListNode {
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

func (sl *EdgeSkipList) Print() {
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

func (sl *EdgeSkipList) Flush() error{
	log.UndoBegin()
	cur := sl.head
	for cur.down != nil {
		cur = cur.down
	}
	for cur != nil {
		latestVersion := cur.FindLatestVersion()
		if latestVersion.changed != Nochange {
			err := cur.Edge.FlushAsUndoBase(latestVersion.VersionId, cur.changed)
			if err != nil {
				log.CtxError(context.Background(),"Edge undolog flush error: %+v", err)
				cur = cur.next
				continue
			}
			if latestVersion.Edge == nil {
				err = cur.Edge.Flush(latestVersion.VersionId, conf.Modify_Removed)
				if err != nil {
					log.CtxError(context.Background(), "Edge flush error: %+v", err)
					cur = cur.next
					continue
				}
			}else {
				err = latestVersion.Edge.Flush(latestVersion.VersionId, latestVersion.changed)
				if err != nil {
					log.CtxError(context.Background(), "Edge flush error: %+v", err)
					cur = cur.next
					continue
				}
			}
		}
		cur = cur.next
	}
	log.UndoCommit()
	return nil
}

//等所有事务执行完毕
func (sl *EdgeSkipList) Wait(){
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

func (sl *EdgeSkipList) CleanStatus() {
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