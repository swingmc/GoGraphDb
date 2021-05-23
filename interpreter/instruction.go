package interpreter

import (
	"GoGraphDb/conf"
	"GoGraphDb/db/db_model"
	"GoGraphDb/db/db_schema"
	"GoGraphDb/memory_cache"
	"GoGraphDb/transaction"
	"GoGraphDb/utils"
	"errors"
	"fmt"
	"strings"
)

var (
	instructionMap = map[int32]func(*transaction.Transaction, string, string, string)error{}
)

func init() {
	instructionMap[conf.GRAMMER_BIND_VERTEX] = bindVertex
	instructionMap[conf.GRAMMER_CREATE_VERTEX] = createVertex
	instructionMap[conf.GRAMMER_SET_VERTEX_TYPE] = setVertexType
	instructionMap[conf.GRAMMER_SET_VERTEX_PROPERTY] = setVertexProperty
	instructionMap[conf.GRAMMER_BIND_EDGE] = bindEdge
	instructionMap[conf.GRAMMER_CREATE_EDGE] = createEdge
	instructionMap[conf.GRAMMER_SET_EDGE_TYPE] = setEdgeType
	instructionMap[conf.GRAMMER_SET_EDGE_PROPERTY] = setEdgeProperty
}

func bindVertex(t *transaction.Transaction, subject string, verb string, object string)error{
	id, err := utils.ParseStringToInt64(object)
	if err != nil {
		return err
	}
	t.VertexSetBind[subject] = map[int64]*db_model.Vertex{id : memory_cache.VertexTree.Get(t.Version, id)}
	return nil
}

func createVertex(t *transaction.Transaction, subject string, verb string, object string)error{
	typeId, err := db_schema.SchemaInstance.VertexType(object)
	if err != nil {
		return err
	}
	id, err := memory_cache.CreateVertex(t.Version)
	if err != nil {
		return err
	}
	err = memory_cache.ModifyVertex(t.Version, id, &typeId, nil)
	if err != nil {
		return err
	}
	t.VertexSetBind[subject] = map[int64]*db_model.Vertex{id : memory_cache.VertexTree.Get(t.Version, id)}
	return nil
}

func setVertexType(t *transaction.Transaction, subject string, verb string, object string)error{
	vs, ok := t.VertexSetBind[subject]
	if !ok {
		return errors.New(fmt.Sprintf("set Vertex %+v type error, not exist", vs))
	}
	vertexType, err := utils.ParseStringToInt32(object)
	if err != nil {
		return err
	}
	for id, _ := range vs {
		err := memory_cache.ModifyVertex(t.Version, id, &vertexType, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func setVertexProperty(t *transaction.Transaction, subject string, verb string, object string)error{
	vs, ok := t.VertexSetBind[subject]
	if !ok {
		return errors.New(fmt.Sprintf("set Vertex %+v properties error, not exist", vs))
	}
	property := strings.Split(object, ":")
	if len(property) != 2{
		return errors.New("Property format not match")
	}
	for id, _ := range vs {
		err := memory_cache.ModifyVertex(t.Version, id, nil, map[string]string{property[0]:property[1]})
		if err != nil {
			return err
		}
	}
	return nil
}

func bindEdge(t *transaction.Transaction, subject string, verb string, object string)error{
	id, err := utils.ParseStringToInt64(object)
	if err != nil {
		return err
	}
	t.EdgeSetBind[subject] = map[int64]*db_model.Edge{id : memory_cache.EdgeTree.Get(t.Version, id)}
	return nil
}

func createEdge(t *transaction.Transaction, subject string, verb string, object string)error{
	startVertexSet, ok := t.VertexSetBind[subject]
	if !ok {
		return errors.New(fmt.Sprintf("create edge error, startVertex not exist"))
	}
	endVertexSet, ok := t.VertexSetBind[subject]
	if !ok {
		return errors.New(fmt.Sprintf("create edge error, endVertex not exist"))
	}
	edgeSet := map[int64]*db_model.Edge{}
	for start, _ := range startVertexSet{
		for end, _ := range endVertexSet{
			id, err := memory_cache.CreateEdge(t.Version)
			if err != nil{
				return err
			}
			err = memory_cache.ModifyEdge(t.Version, id, nil, &start, &end, nil)
			if err != nil{
				return err
			}
			edgeSet[id] = memory_cache.EdgeTree.Get(t.Version, id)
		}
	}
	t.EdgeSetBind[verb] = edgeSet
	return nil
}

func setEdgeType(t *transaction.Transaction, subject string, verb string, object string)error {
	es, ok := t.EdgeSetBind[subject]
	if !ok {
		return errors.New(fmt.Sprintf("set Edge %+v type error, not exist", es))
	}
	edgeType, err := utils.ParseStringToInt32(object)
	if err != nil {
		return err
	}
	for id, _ := range es {
		err = memory_cache.ModifyEdge(t.Version, id, &edgeType, nil, nil, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func setEdgeProperty(t *transaction.Transaction, subject string, verb string, object string)error{
	es, ok := t.VertexSetBind[subject]
	if !ok {
		return errors.New(fmt.Sprintf("set Edge %+v properties error, not exist", es))
	}
	property := strings.Split(object, ":")
	if len(property) != 2{
		return errors.New("Property format not match")
	}
	for id, _ := range es {
		err := memory_cache.ModifyEdge(t.Version, id, nil,nil,nil, map[string]string{property[0]:property[1]})
		if err != nil {
			return err
		}
	}
	return nil
}