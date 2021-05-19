package utils

import (
	"GoGraphDb/conf"
	"context"
	"encoding/json"
	"strconv"
)

func ParseInt64ToString(n int64) string{
	return strconv.FormatInt(n ,10)
}

func ParseStringToInt64(s string) (int64,error){
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0,err
	}
	return num,nil
}

func ParseStringToInt32(s string) (int32,error){
	num, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0,err
	}
	return int32(num),nil
}

func CtxToString(ctx context.Context) string{
	 bytes, err := (json.Marshal(ctx))
	 if err != nil {
	 	return err.Error()
	 }
	 return string(bytes)
}

func MapClone(properties map[string]string) map[string]string {
	clone := make(map[string]string)
	for k, v := range properties {
		clone[k] = v
	}
	return clone
}

func EdgeSetClone(properties map[conf.EdgeIdentifier]bool) map[conf.EdgeIdentifier]bool {
	clone := make(map[conf.EdgeIdentifier]bool)
	for k, v := range properties {
		clone[k] = v
	}
	return clone
}

func VertexSetClone(properties map[conf.VertexIdentifier]bool) map[conf.VertexIdentifier]bool {
	clone := make(map[conf.VertexIdentifier]bool)
	for k, v := range properties {
		clone[k] = v
	}
	return clone
}
