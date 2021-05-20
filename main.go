package main

import (
	_"GoGraphDb/db/db_model"
	_"GoGraphDb/db/dal"
	_"GoGraphDb/log"
	_"GoGraphDb/db"
	"GoGraphDb/log"
	"GoGraphDb/memory_cache"
	"GoGraphDb/transaction"
	"context"
	"fmt"
	"time"
	_"GoGraphDb/manager"
	_"GoGraphDb/interpreter"
)

func main(){
	//sl := model.NewVertexSkipList()

	/*node := model.Vertex{
		Idntifier:  1128,
		VertexType: 1,
		OutE:       map[model.EdgeIdentifier]bool{22:true},
		InE:        map[model.EdgeIdentifier]bool{11:true},
		Properties: map[string]string{"name\ntest":"test"},
	}
	bytes, err := json.Marshal(node)
	if err != nil {
		log.CtxError(context.Background(), "test error: %+v", node)
	}
	fmt.Println(node)
	fmt.Println(string(bytes))
	fmt.Println(time.Now().UnixNano())
	err = node.Flush(conf.Modify_Create)
	node2 := model.Vertex{}
	err = json.Unmarshal(bytes, &node2)
	fmt.Println(node2)
	fmt.Println(err)
	 */

	/*
	fmt.Println(*memory_cache.VertexTree.Get(1128))
	con := context.Background()

	fmt.Println(con.Value(1))
	log.CtxInfo(con,"test 1: %+v, 2: %+v", 1, 2)
	fmt.Println(fmt.Sprintf("test 1: %+v, 2: %+v", 1, 2))

	 */
	/*
	ctx := context.WithValue(context.Background(), 1, "ddddd")
	log.CtxInfo(ctx, "test")
	fmt.Println(ctx)
	fmt.Println(ctx.Value(1))
	fmt.Println(utils.CtxToString(ctx))
	 */

	/*
	memory_cache.CreateVertex(utils.GenTimeStamp())
	//time.Sleep(time.Nanosecond)
	memory_cache.CreateVertex(utils.GenTimeStamp())
	//time.Sleep(time.Nanosecond)
	memory_cache.CreateVertex(utils.GenTimeStamp())
	//time.Sleep(time.Nanosecond)
	memory_cache.CreateVertex(utils.GenTimeStamp())
	//time.Sleep(time.Nanosecond)
	memory_cache.CreateVertex(utils.GenTimeStamp())
	//time.Sleep(time.Nanosecond)
	memory_cache.CreateVertex(utils.GenTimeStamp())
	//time.Sleep(time.Nanosecond)
	memory_cache.CreateVertex(utils.GenTimeStamp())
	//time.Sleep(time.Nanosecond)
	memory_cache.CreateVertex(utils.GenTimeStamp())
	memory_cache.VertexTree.Print()
	//manager.Flush()
	 */

	timeAvg := int64(0)
	num := int64(10)
	i := int64(0)
	for i<num {
		//j := i
		t := transaction.NewTransaction()
		start := int64(time.Now().Nanosecond())
		go func() {
			_,err := memory_cache.CreateVertex(t.Version)
			if err != nil {
				log.CtxError(context.Background(),"err: %+v",err)
			}
			end := int64(time.Now().Nanosecond())
			length := (end - start)/num
			timeAvg = timeAvg + length
			fmt.Println(timeAvg)
		}()
		i++
	}
	time.Sleep(time.Second*5)
	//memory_cache.VertexTree.Flush()
	fmt.Println(timeAvg)

}