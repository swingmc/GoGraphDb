package interpreter

import (
	"GoGraphDb/conf"
	"GoGraphDb/db/db_schema"
	"GoGraphDb/memory_cache"
	"GoGraphDb/transaction"
	"GoGraphDb/utils"
	"errors"
	"fmt"
)

var (
	instructionMap = map[int32]func(*transaction.Transaction, string, string, string)error{}
)

func init() {
	instructionMap[conf.GRAMMER_BIND_VERTEX] = bindVertex
	instructionMap[conf.GRAMMER_CREATE_VERTEX] = createVertex
	instructionMap[conf.GRAMMER_SET_VERTEX_TYPE] = setVertexType
	instructionMap[conf.GRAMMER_SET_VERTEX_PROPERTY] = setVertexProperty
}

func bindVertex(t *transaction.Transaction, subject string, verb string, object string)error{
	id, err := utils.ParseStringToInt64(object)
	if err != nil {
		return err
	}
	t.VertexBind[subject] = id
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
	memory_cache.ModifyVertex(t.Version, id, &typeId, nil)
	return nil
}

func setVertexType(t *transaction.Transaction, subject string, verb string, object string)error{
	id, ok := t.VertexBind[subject]
	if !ok {
		return errors.New(fmt.Sprintf("set Vertex %+v type error, not exist",id))
	}
	vertexType, err := utils.ParseStringToInt32(object)
	if err != nil {
		return err
	}
	memory_cache.ModifyVertex(t.Version, id, &vertexType, nil)
	return nil
}

func setVertexProperty(t *transaction.Transaction, subject string, verb string, object string)error{
	id, ok := t.VertexBind[subject]
	if !ok {
		return errors.New(fmt.Sprintf("set Vertex %+v properties error, not exist",id))
	}
	properties := map[string]string{
		verb: object,
	}
	memory_cache.ModifyVertex(t.Version, id, nil, properties)
	return nil
}