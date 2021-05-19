package main

import (
	"GoGraphDb/db/db_model"
	"fmt"
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

	v := db_model.Vertex{}
	fmt.Println(v.Idntifier)
}