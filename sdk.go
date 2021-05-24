package main

import (
	"GoGraphDb/db/db_model"
	"GoGraphDb/interpreter"
	"GoGraphDb/log"
	"GoGraphDb/utils"
	"context"
	"errors"
	"fmt"
)

type InterpreterWrapper struct{
	i interpreter.Interpreter
}

type VertexWrapper struct{
	Idntifier  int64
	VertexType int32
	OutE       map[int64]bool
	InE        map[int64]bool
	Properties map[string]string
}

type EdgeWrapper struct{
	Idntifier  int64
	EdgeType   int32
	Start      int64
	End        int64
	Properties map[string]string
}

func NewInterpreter() *InterpreterWrapper {
	return &InterpreterWrapper{}
}

func wrapVertex(v *db_model.Vertex) *VertexWrapper{
	oE := map[int64]bool{}
	iE := map[int64]bool{}
	for id, _ := range v.OutE{
		oE[int64(id)] = true
	}
	for id, _ := range v.OutE{
		iE[int64(id)] = true
	}
	return &VertexWrapper{
		Idntifier: int64(v.Idntifier),
		VertexType: v.VertexType,
		OutE: oE,
		InE: iE,
		Properties: utils.MapClone(v.Properties),
	}
}

func wrapEdge(e *db_model.Edge) *EdgeWrapper{
	return &EdgeWrapper{
		Idntifier:  int64(e.Idntifier),
		EdgeType:   e.EdgeType,
		Start:      int64(e.Start),
		End:        int64(e.End),
		Properties: utils.MapClone(e.Properties),
	}
}

func (w *InterpreterWrapper) BeginTransaction() error{
	return w.i.BeginTransaction()
}

func (w *InterpreterWrapper) EndTransaction() error{
	return w.i.EndTransaction()
}

func (w *InterpreterWrapper) RollbackTransaction() error{
	return w.i.RollbackTransaction()
}

func (w *InterpreterWrapper) RawSql(subject string, verb string, object string) (int32, error){
	return w.i.RawSql(subject, verb, object)
}

func (w *InterpreterWrapper) GetData() (interface{},error){
	obj, err := w.i.GetData()
	if err != nil {
		log.CtxError(context.Background(),"GetData error: %+v", obj)
		return nil, errors.New(fmt.Sprintf("GetData error: %+v", obj))
	}
	vertexMap, ok := obj.(map[int64]*db_model.Vertex)
	if ok {
		wrapperMap := map[int64]*VertexWrapper{}
		for id, vertex := range vertexMap{
			wrapperMap[id] = wrapVertex(vertex)
		}
		return wrapperMap, nil
	}
	edgeMap, ok := obj.(map[int64]*db_model.Edge)
	if ok {
		wrapperMap := map[int64]*EdgeWrapper{}
		for id, edge := range edgeMap{
			wrapperMap[id] = wrapEdge(edge)
		}
		return wrapperMap, nil
	}
	return nil, errors.New(fmt.Sprintf("data struct class error: %+v", obj))
}