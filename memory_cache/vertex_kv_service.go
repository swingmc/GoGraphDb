package memory_cache

import (
	"GoGraphDb/conf"
	"GoGraphDb/db/db_model"
	"GoGraphDb/log"
	"GoGraphDb/utils"
	"context"
	"errors"
)

func GetVertex(versionId int64, id int64) *db_model.Vertex {
	return VertexTree.Get(versionId, id)
}

func CreateVertex(versionId int64) (int64,error){
	id := utils.GenId()
	if !checkWriteVertexParam(versionId, id){
		log.CtxWarn(context.Background(), "CreateVertex Param Error, id: %+v, type: %+v", id, versionId)
		return 0,errors.New("CreateVertex Param Error")
	}
	vertex := db_model.Vertex{
		Idntifier: conf.VertexIdentifier(id),
	}
	err := VertexTree.Insert(versionId, id, &vertex)
	if err != nil {
		log.CtxError(context.Background(), "create vertex error: %+v, id: %+v", err, id)
		return 0,err
	}
	log.CtxInfo(context.Background(), "create vertex id: %+v", id)
	return id, nil
}

func RemoveVertex(versionId int64, id int64) error{
	return VertexTree.Remove(versionId, id)
}

func ModifyVertex(versionId int64, id int64, vertexType *int32, properties map[string]string, inE map[int64]bool, outE map[int64]bool) error{
	if !checkWriteVertexParam(versionId, id){
		log.CtxWarn(context.Background(), "ModifyVertex Param Error, id: %+v, version_id: %+v", id, versionId)
		return errors.New("ModifyVertex Param Error")
	}
	oldVertex := VertexTree.Get(versionId, id)
	if oldVertex == nil {
		return errors.New("Vertex to be modified not exist")
	}
	newVertex := db_model.Vertex{
		Idntifier:  oldVertex.Idntifier,
		VertexType: oldVertex.VertexType,
		OutE:       utils.EdgeSetClone(oldVertex.OutE),
		InE:        utils.EdgeSetClone(oldVertex.InE),
		Properties: utils.MapClone(oldVertex.Properties),
	}
	if vertexType != nil {
		newVertex.VertexType = *vertexType
	}
	if properties != nil {
		if newVertex.Properties == nil {
			newVertex.Properties = properties
		}else {
			for key, value := range properties {
				newVertex.Properties[key] = value
			}
		}
	}
	if inE != nil {
		{
			for key, value := range inE {
				newVertex.InE[conf.EdgeIdentifier(key)] = value
			}
		}
	}
	if outE != nil {
		{
			for key, value := range outE {
				newVertex.OutE[conf.EdgeIdentifier(key)] = value
			}
		}
	}
	err := VertexTree.Insert(versionId, id, &newVertex)
	if err != nil {
		log.CtxError(context.Background(),"Modify Vertex wrong, error: %+v", err)
		return err
	}
	return nil
}

func checkWriteVertexParam(versionId int64, id int64) bool{
	if id <= 0 || versionId <= 0 {
		return false
	}
	return true
}

func FlushVertexTree() error{
	log.CtxInfo(context.Background(), "FlushVertexTree")
	err := VertexTree.Flush()
	if err != nil {
		log.CtxError(context.Background(), "FlushVertexTree error, error: %+v", err)
		return err
	}
	return nil
}