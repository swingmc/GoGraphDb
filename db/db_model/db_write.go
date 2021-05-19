package db_model

import (
	"GoGraphDb/conf"
	"GoGraphDb/db"
	"GoGraphDb/log"
	_ "GoGraphDb/log"
	"GoGraphDb/utils"
	"context"
	"encoding/json"
	"errors"
	"strconv"
)


func (v *Vertex) Flush(versionId int64, changed int32) error{
	bytes, err := json.Marshal(v)
	if err != nil {
		log.CtxError(context.Background(),"vertex flush error: %+v", err)
		return err
	}
	switch changed {
	case conf.Modify_Nochange:
		return nil
	case conf.Modify_Create:
		db.DbWriter.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_CREATE_VERTEX + conf.Splitor + strconv.FormatInt(int64(v.Idntifier), 10) + conf.Splitor + string(bytes) + "\n")
	case conf.Modify_Changed:
		db.DbWriter.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_CHANGE_VERTEX + conf.Splitor + strconv.FormatInt(int64(v.Idntifier), 10) + conf.Splitor + string(bytes) + "\n")
	case conf.Modify_Removed:
		db.DbWriter.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_REMOVE_VERTEX + conf.Splitor +
			strconv.FormatInt(int64(v.Idntifier), 10) +
			"\n")
	default:
		{
			err = errors.New("vertex 'changed' status no defination: " + utils.ParseInt64ToString(int64(changed)))
			log.CtxError(context.Background(),"vertex flush error: %+v", err)
			return err
		}
	}
	return nil
}

func (v *Vertex) FlushAsUndoBase(versionId int64, changed int32) error{
	bytes, err := json.Marshal(v)
	if err != nil {
		log.CtxError(context.Background(),"vertex flush_undo error: %+v", err)
		return err
	}
	switch changed {
	case conf.Modify_Create:
		log.UndoLogFile.WriteString(
			utils.LogTimeByVersion(versionId) +
				conf.Splitor + conf.COMMAND_CREATE_VERTEX + conf.Splitor +
				strconv.FormatInt(int64(v.Idntifier), 10) +
				conf.Splitor + string(bytes) + "\n")
	case conf.Modify_Changed:
		log.UndoLogFile.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_CHANGE_VERTEX + conf.Splitor + strconv.FormatInt(int64(v.Idntifier), 10) + conf.Splitor + string(bytes) + "\n")
	case conf.Modify_Removed:
		log.UndoLogFile.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_REMOVE_VERTEX + conf.Splitor + strconv.FormatInt(int64(v.Idntifier), 10) + conf.Splitor + string(bytes) + "\n")
	default:
		{
			err = errors.New("vertex flush'changed' status no defination: " + utils.ParseInt64ToString(int64(changed)))
			log.CtxError(context.Background(),"vertex undo flush error: %+v", err)
			return err
		}
	}

	return nil
}

func (v *Edge) Flush(versionId int64, changed int32) error{
	bytes, err := json.Marshal(v)
	if err != nil {
		log.CtxError(context.Background(),"Edgeflush error: %+v", err)
		return err
	}
	switch changed {
	case conf.Modify_Nochange:
		return nil
	case conf.Modify_Create:
		db.DbWriter.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_CREATE_EDGE+ conf.Splitor + strconv.FormatInt(int64(v.Idntifier), 10) + conf.Splitor + string(bytes) + "\n")
	case conf.Modify_Changed:
		db.DbWriter.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_CHANGE_EDGE+ conf.Splitor + strconv.FormatInt(int64(v.Idntifier), 10) + conf.Splitor + string(bytes) + "\n")
	case conf.Modify_Removed:
		db.DbWriter.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_REMOVE_EDGE+ conf.Splitor + strconv.FormatInt(int64(v.Idntifier), 10) + "\n")
	default:
		{
			err = errors.New("Edge'changed' status no defination: " + utils.ParseInt64ToString(int64(changed)))
			log.CtxError(context.Background(),"Edgeflush error: %+v", err)
			return err
		}
	}
	return nil
}

func (v *Edge) FlushAsUndoBase(versionId int64, changed int32) error{
	bytes, err := json.Marshal(v)
	if err != nil {
		log.CtxError(context.Background(),"Edgeflush_undo error: %+v", err)
		return err
	}
	switch changed {
	case conf.Modify_Create:
		log.UndoLogFile.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_CREATE_EDGE+ conf.Splitor + strconv.FormatInt(int64(v.Idntifier), 10) + conf.Splitor + string(bytes) + "\n")
	case conf.Modify_Changed:
		log.UndoLogFile.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_CHANGE_EDGE+ conf.Splitor + strconv.FormatInt(int64(v.Idntifier), 10) + conf.Splitor + string(bytes) + "\n")
	case conf.Modify_Removed:
		log.UndoLogFile.WriteString(utils.LogTimeByVersion(versionId) + conf.Splitor + conf.COMMAND_REMOVE_EDGE+ conf.Splitor + strconv.FormatInt(int64(v.Idntifier), 10) + conf.Splitor + string(bytes) + "\n")
	default:
		{
			err = errors.New("Edge flush 'changed' status no defination: " + utils.ParseInt64ToString(int64(changed)))
			log.CtxError(context.Background(),"Edgeundo flush error: %+v", err)
			return err
		}
	}
	return nil
}