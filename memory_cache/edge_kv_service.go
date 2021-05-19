package memory_cache

import (
	"GoGraphDb/conf"
	"GoGraphDb/db/db_model"
	"GoGraphDb/log"
	"GoGraphDb/utils"
	"context"
	"errors"
)

func GetEdge(versionId int64, id int64) *db_model.Edge {
	return EdgeTree.Get(versionId, id)
}

func CreateEdge(versionId int64) (int64,error){
	id := utils.GenId()
	if !checkWriteEdgeParam(versionId, id){
		log.CtxWarn(context.Background(), "CreateEdge Param Error, id: %+v, type: %+v", id, versionId)
		return 0,errors.New("CreateEdge Param Error")
	}
	Edge := db_model.Edge{
		Idntifier: conf.EdgeIdentifier(id),
	}
	err := EdgeTree.Insert(versionId, id, &Edge)
	if err != nil {
		log.CtxError(context.Background(), "create Edge error: %+v, id: %+v", err, id)
		return 0,err
	}
	log.CtxInfo(context.Background(), "create Edge id: %+v", id)
	return id, nil
}

func RemoveEdge(versionId int64, id int64) *db_model.Edge {
	return EdgeTree.Remove(versionId, id)
}

func ModifyEdge(versionId int64, id int64, EdgeType *int32, start *int64, end *int64, properties map[string]string) error{
	if !checkWriteEdgeParam(versionId, id){
		log.CtxWarn(context.Background(), "ModifyEdge Param Error, id: %+v, version_id: %+v", id, versionId)
		return errors.New("ModifyEdge Param Error")
	}
	oldEdge := EdgeTree.Get(versionId, id)
	if oldEdge == nil {
		return errors.New("Edge to be modified not exist")
	}
	newEdge := db_model.Edge{
		Idntifier:  oldEdge.Idntifier,
		EdgeType:   oldEdge.EdgeType,
		Start:      oldEdge.Start,
		End:        oldEdge.End,
		Properties: utils.MapClone(oldEdge.Properties),
	}
	if EdgeType != nil {
		newEdge.EdgeType = *EdgeType
	}
	if start != nil {
		newEdge.Start = conf.VertexIdentifier(*start)
	}
	if end != nil {
		newEdge.End = conf.VertexIdentifier(*end)
	}
	if properties != nil {
		if newEdge.Properties == nil {
			newEdge.Properties = properties
		}else {
			for key, value := range properties {
				newEdge.Properties[key] = value
			}
		}
	}
	err := EdgeTree.Insert(versionId, id, &newEdge)
	if err != nil {
		log.CtxError(context.Background(),"Modify Edge wrong, error: %+v", err)
		return err
	}
	return nil
}

func checkWriteEdgeParam(versionId int64, id int64) bool{
	if id <= 0 || versionId <= 0 {
		return false
	}
	return true
}