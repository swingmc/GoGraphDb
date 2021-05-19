package memory_cache

import (
	"GoGraphDb/conf"
	"GoGraphDb/db"
	"GoGraphDb/db/db_model"
	"GoGraphDb/db/dal"
	"GoGraphDb/log"
	"GoGraphDb/utils"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

var (
	VertexTree *dal.VertexSkipList
	EdgeTree *dal.EdgeSkipList
)

func init() {
	reader := bufio.NewReader(db.DbReader)
	//读取lsm_tree数据，构建skip_list
	VertexTree = dal.NewVertexSkipList()
	EdgeTree = dal.NewEdgeSkipList()
	for {
		row, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.CtxError(context.Background(), "lsm_tree init wrong, read data file error: %+v", err)
			}
			break
		}
		units := strings.Split(row, conf.Splitor)
		switch units[conf.InstructionPointer_Command] {
		case conf.COMMAND_CREATE_VERTEX, conf.COMMAND_CHANGE_VERTEX:
			{
				jsonObj := units[conf.InstructionPointer_JsonObj]
				Vertex := db_model.Vertex{}
				err := json.Unmarshal([]byte(jsonObj), &Vertex)
				if err != nil {
					log.CtxError(context.Background(), "lsm_tree init wrong, vertex unmarshal error: %+v", err)
					panic(err)
				}
				VertexTree.Insert(0, int64(Vertex.Idntifier), &Vertex)
			}
		case conf.COMMAND_REMOVE_VERTEX:
			{
				idStr := units[conf.InstructionPointer_Identifier]
				id, err := utils.ParseStringToInt64(idStr)
				if err != nil {
					log.CtxError(context.Background(), "lsm_tree init wrong, id parse error: %+v", err)
					panic(err)
				}
				VertexTree.Remove(0, id)
			}
		case conf.COMMAND_CREATE_EDGE, conf.COMMAND_CHANGE_EDGE:
			{
				jsonObj := units[conf.InstructionPointer_JsonObj]
				Edge := db_model.Edge{}
				err := json.Unmarshal([]byte(jsonObj), &Edge)
				if err != nil {
					log.CtxError(context.Background(), "lsm_tree init wrong, edge unmarshal error: %+v", err)
					panic(err)
				}
				EdgeTree.Insert(0, int64(Edge.Idntifier), &Edge)
			}
		case conf.COMMAND_REMOVE_EDGE:
			{
				idStr := units[conf.InstructionPointer_Identifier]
				id, err := utils.ParseStringToInt64(idStr)
				if err != nil {
					log.CtxError(context.Background(), "lsm_tree init wrong, id parse error: %+v", err)
					panic(err)
				}
				EdgeTree.Remove(0, id)
			}
		default:
			{
				err := errors.New("command no definition")
				log.CtxError(context.Background(), "lsm_tree init wrong, error: %+v", err)
				panic(err)
			}
		}
	}
	VertexTree.CleanStatus()
}