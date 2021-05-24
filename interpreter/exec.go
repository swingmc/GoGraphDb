package interpreter

import (
	"GoGraphDb/conf"
	"GoGraphDb/log"
	"context"
	"errors"
)

func (i *Interpreter) ExecuteSentence(subject string, verb string, object string) (int32, error){
	if i.transaction == nil{
		return 0, errors.New("no executing transaction, can't execute")
	}
	ctx := context.WithValue(context.Background(), "version_id", i.transaction.Version)
	log.CtxInfo(ctx, "execute sentence: sub: %+v, verb: %+v, object: %+v", subject, verb, object)
	command, err := i.judgeCommand(subject, verb, object)
	if err != nil {
		log.CtxError(context.Background(), "command error: %+v, sub: %+v, verb: %+v, object: %+v", err, subject, verb, object)
	}
	function, ok := instructionMap[command]
	if !ok {
		return 0, errors.New("command not exist!")
	}
	/*
	if IsWriteCommand(command) && i.transaction.ReadOnly {
		return 0, errors.New("ReadOnlyTransaction exec write command")
	}
	*/
	err = function(i.transaction, subject, verb, object)
	//err = i.exec(command, subject, verb, object)
	if err != nil {
		log.CtxError(ctx, "transaction: %+v exec wrong, error: %+V", i.transaction.Version, err)
		return 0,err
	}
	return 0, nil
}

func (i *Interpreter) judgeCommand(subject string, verb string, object string) (int32, error){
	if len(subject) == 0 || len(verb) == 0 || len(object) == 0 {
		return 0, errors.New("has nil word")
	}
	commad, ok := conf.VerbMap[verb]
	if ok {
		return commad, nil
	}
	if i.transaction.IsVertex(subject) && i.transaction.IsVertex(object) {
		return conf.GRAMMER_CREATE_EDGE, nil
	}
	return 0, errors.New("no match command!")
}
/*
func (i *Interpreter) exec(command int32, subject string, verb string, object string) error{
	switch command {
	case conf.GRAMMER_BIND_VERTEX:
		{
			id, err := utils.ParseStringToInt64(object)
			if err != nil {
				return err
			}
			i.transaction.VertexBind[subject] = id
		}
	case conf.GRAMMER_CREATE_VERTEX:
		{
			typeId, err := db_schema.SchemaInstance.VertexType(object)
			if err != nil {
				return err
			}
			id, err := memory_cache.CreateVertex(i.transaction.Version)
			if err != nil {
				return err
			}
			memory_cache.ModifyVertex(i.transaction.Version, id, &typeId, nil)
		}
	case conf.GRAMMER_SET_VERTEX_TYPE:
		{
			id, ok := i.transaction.VertexBind[subject]
			if !ok {
				return errors.New(fmt.Sprintf("set Vertex %+v type error, not exist",id))
			}
			vertexType, err := utils.ParseStringToInt32(object)
			if err != nil {
				return err
			}
			memory_cache.ModifyVertex(i.transaction.Version, id, &vertexType, nil)
		}
	case conf.GRAMMER_SET_VERTEX_PROPERTY:
		{
			id, ok := i.transaction.VertexBind[subject]
			if !ok {
				return errors.New(fmt.Sprintf("set Vertex %+v properties error, not exist",id))
			}
			properties := map[string]string{
				verb: object,
			}
			memory_cache.ModifyVertex(i.transaction.Version, id, nil, properties)
		}
	case conf.GRAMMER_BIND_EDGE:
	case conf.GRAMMER_CREATE_EDGE:
	case conf.GRAMMER_SET_EDGE_TYPE:
	case conf.GRAMMER_SET_EDGE_PROPERTY:
	}
	return nil
}
*/

//判断写命令
func IsWriteCommand(n int32) bool{
	return n > 10000
}
