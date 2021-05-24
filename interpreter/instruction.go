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
	instructionMap[conf.GRAMMER_REMOVE_VERTEX] = removeVertex
	instructionMap[conf.GRAMMER_FILTER_VERTEX_BY_TYPE] = filterVertexByType
	instructionMap[conf.GRAMMER_IN_EDGE] = inEdge
	instructionMap[conf.GRAMMER_OUT_EDGE] = OutEdge
	instructionMap[conf.GRAMMER_SHOW_VERTEX] = showVertex

	instructionMap[conf.GRAMMER_BIND_EDGE] = bindEdge
	instructionMap[conf.GRAMMER_CREATE_EDGE] = createEdge
	instructionMap[conf.GRAMMER_SET_EDGE_TYPE] = setEdgeType
	instructionMap[conf.GRAMMER_SET_EDGE_PROPERTY] = setEdgeProperty
	instructionMap[conf.GRAMMER_REMOVE_EDGE] = removeEdge
	instructionMap[conf.GRAMMER_FILTER_EDGE_BY_TYPE] = filterEdgeByType
	instructionMap[conf.GRAMMER_EDGE_START] = edgeStart
	instructionMap[conf.GRAMMER_EDGE_END] = edgeEnd
	instructionMap[conf.GRAMMER_SHOW_EDGE] = showEdge
}

func bindVertex(t *transaction.Transaction, subject string, verb string, object string)error{
	id, err := utils.ParseStringToInt64(object)
	if err != nil {
		return err
	}
	t.VertexWrite(subject, map[int64]*db_model.Vertex{id : memory_cache.VertexTree.Get(t.Version, id)})
	//t.vertexSetBind[subject] = map[int64]*db_model.Vertex{id : memory_cache.VertexTree.Get(t.Version, id)}
	return nil
}

func createVertex(t *transaction.Transaction, subject string, verb string, object string)error{
	typeId, err := db_schema.VertexType(object)
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
	t.VertexWrite(subject, map[int64]*db_model.Vertex{id : memory_cache.VertexTree.Get(t.Version, id)})
	//t.vertexSetBind[subject] = map[int64]*db_model.Vertex{id : memory_cache.VertexTree.Get(t.Version, id)}
	return nil
}

func setVertexType(t *transaction.Transaction, subject string, verb string, object string)error{
	vs, ok := t.VertexRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("set Vertex %+v type error, not exist", subject))
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
	vs, ok := t.VertexRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("set Vertex %+v properties error, not exist", subject))
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

func removeVertex(t *transaction.Transaction, subject string, verb string, object string)error{
	vs, ok := t.VertexRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("vertexSet not defined: %+v", subject))
	}
	for id, _ := range vs {
		err := memory_cache.RemoveVertex(t.Version, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func filterVertexByType(t *transaction.Transaction, subject string, verb string, object string)error{
	vs, ok := t.VertexRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("vertexSet not defined: %+v", subject))
	}
	vertexType, err := utils.ParseStringToInt32(object)
	if err != nil {
		return err
	}
	newVs := map[int64]*db_model.Vertex{}
	for id, vertex := range vs{
		if vertex.VertexType == vertexType{
			newVs[id] = vertex
		}
	}
	t.VertexWrite(subject, newVs)
	return nil
}

func inEdge(t *transaction.Transaction, subject string, verb string, object string)error{
	vs, ok := t.VertexRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("vertexSet not defined: %+v", subject))
	}
	edges := map[int64]*db_model.Edge{}
	for _, v := range vs{
		for id, _ := range v.InE{
			edge := memory_cache.EdgeTree.Get(t.Version, int64(id))
			if edge != nil{
				edges[int64(id)] = edge
			}
		}
	}
	t.EdgeWrite(object, edges)
	return nil
}

func OutEdge(t *transaction.Transaction, subject string, verb string, object string)error{
	vs, ok := t.VertexRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("vertexSet not defined: %+v", subject))
	}
	edges := map[int64]*db_model.Edge{}
	for _, v := range vs{
		for id, _ := range v.OutE{
			edge := memory_cache.EdgeTree.Get(t.Version, int64(id))
			if edge != nil{
				edges[int64(id)] = edge
			}
		}
	}
	t.EdgeWrite(object, edges)
	return nil
}

func showVertex(t *transaction.Transaction, subject string, verb string, object string)error{
	vs, ok := t.VertexRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("show Vertex %+v error, not exist", subject))
	}
	t.DataChan <- vs
	return nil
}


func bindEdge(t *transaction.Transaction, subject string, verb string, object string)error{
	id, err := utils.ParseStringToInt64(object)
	if err != nil {
		return err
	}
	t.EdgeWrite(subject, map[int64]*db_model.Edge{id : memory_cache.EdgeTree.Get(t.Version, id)})
	//t.edgeSetBind[subject] = map[int64]*db_model.Edge{id : memory_cache.EdgeTree.Get(t.Version, id)}
	return nil
}

func createEdge(t *transaction.Transaction, subject string, verb string, object string)error{
	startVertexSet, ok := t.VertexRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("create edge error, startVertex not exist: " + subject))
	}
	endVertexSet, ok := t.VertexRead(object)
	if !ok {
		return errors.New(fmt.Sprintf("create edge error, endVertex not exist: " + object))
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
	t.EdgeWrite(verb, edgeSet)
	return nil
}

func setEdgeType(t *transaction.Transaction, subject string, verb string, object string)error {
	es, ok := t.EdgeRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("set Edge %+v type error, not exist", subject))
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
	es, ok := t.EdgeRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("set Edge %+v properties error, not exist", subject))
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

func removeEdge(t *transaction.Transaction, subject string, verb string, object string)error{
	vs, ok := t.EdgeRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("edgeSet not defined: %+v", subject))
	}
	for id, _ := range vs {
		err := memory_cache.RemoveEdge(t.Version, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func filterEdgeByType(t *transaction.Transaction, subject string, verb string, object string)error{
	es, ok := t.EdgeRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("edgeSet not defined: %+v", subject))
	}
	edgeType, err := utils.ParseStringToInt32(object)
	if err != nil {
		return err
	}
	newEs := map[int64]*db_model.Edge{}
	for id, edge := range es {
		if edge.EdgeType == edgeType {
			newEs[id] = edge
		}
	}
	t.EdgeWrite(subject, newEs)
	return nil
}

func edgeStart(t *transaction.Transaction, subject string, verb string, object string)error{
	es, ok := t.EdgeRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("edgeSet not defined: %+v", subject))
	}
	vertexSet := map[int64]*db_model.Vertex{}
	for _,edge := range es{
		vertex := memory_cache.GetVertex(t.Version, int64(edge.Start))
		if vertex != nil {
			vertexSet[int64(edge.Start)] = vertex
		}
	}
	t.VertexWrite(object, vertexSet)
	return nil
}

func edgeEnd(t *transaction.Transaction, subject string, verb string, object string)error{
	es, ok := t.EdgeRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("edgeSet not defined: %+v", subject))
	}
	vertexSet := map[int64]*db_model.Vertex{}
	for _,edge := range es{
		vertex := memory_cache.GetVertex(t.Version, int64(edge.End))
		if vertex != nil {
			vertexSet[int64(edge.End)] = vertex
		}
	}
	t.VertexWrite(object, vertexSet)
	return nil
}


func showEdge(t *transaction.Transaction, subject string, verb string, object string)error{
	es, ok := t.EdgeRead(subject)
	if !ok {
		return errors.New(fmt.Sprintf("show Edge %+v error, not exist", subject))
	}
	t.DataChan <- es
	return nil
}
