package db_schema

import (
	"GoGraphDb/conf"
	"GoGraphDb/log"
	"GoGraphDb/utils"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Schema struct{
	VertexTypeMap map[string]int32
	EdgeTypeMap map[string]int32
}

var (
	schemaFilePath = conf.ProjectRootPath + conf.DataSchemaPath

	SchemaInstance = Schema{
		VertexTypeMap: map[string]int32{},
		EdgeTypeMap:   map[string]int32{},
	}
)

func init() {
	schemaFile, err := os.OpenFile(schemaFilePath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	schemaFlag := conf.SchemaFile_VertexFlag
	reader := bufio.NewReader(schemaFile)
	for {
		row, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.CtxError(context.Background(), "SchemaInstance init wrong, read data file error: %+v", err)
				panic(err)
			}
			break
		}
		//去除换行符
		row = strings.ReplaceAll(row, "\r", "")
		row = strings.ReplaceAll(row, "\n", "")
		if len(row) == 0 {
			log.CtxWarn(context.Background(), "SchemaInstance file has empty line")
			continue
		}
		if strings.HasPrefix(row ,conf.SchemaFile_VertexFlag) {
			schemaFlag = conf.SchemaFile_VertexFlag
			continue
		}
		if row == conf.SchemaFile_EdgeFlag {
			schemaFlag = conf.SchemaFile_EdgeFlag
			continue
		}

		fmt.Print(row)
		fmt.Println(row == conf.SchemaFile_VertexFlag)
		fmt.Println(strings.HasPrefix(row ,conf.SchemaFile_VertexFlag))
		fmt.Println(strings.Compare(row, conf.SchemaFile_VertexFlag))
		units := strings.Split(row, conf.Splitor)
		if len(units) != 2 {
			err = errors.New(fmt.Sprintf("sentence grammer error, sentence: %+v", units))
			log.CtxError(context.Background(), "SchemaInstance init wrong, error: %+v", err)
			panic(err)
		}
		if schemaFlag == conf.SchemaFile_VertexFlag {
			typeId, err := utils.ParseStringToInt32(units[1])
			if err != nil {
				log.CtxError(context.Background(), "SchemaInstance init wrong, error: %+v", err)
				panic(err)
			}
			SchemaInstance.VertexTypeMap[units[0]] = typeId
		} else {
			typeId, err := utils.ParseStringToInt32(units[1])
			if err != nil {
				log.CtxError(context.Background(), "SchemaInstance init wrong, error: %+v", err)
				panic(err)
			}
			SchemaInstance.EdgeTypeMap[units[0]] = typeId
		}
	}
}

func (s *Schema) VertexType(name string) (int32, error){
	typeId, ok := s.VertexTypeMap[name]
	if !ok {
		return 0, errors.New("no such vertex type: " + name)
	}
	return typeId, nil
}

func (s *Schema) EdgeType(name string) int32{
	typeId, ok := s.EdgeTypeMap[name]
	if !ok {
		return 0
	}
	return typeId
}
