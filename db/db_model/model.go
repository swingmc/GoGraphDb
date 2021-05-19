package db_model

import "GoGraphDb/conf"

type Edge struct {
	Idntifier  conf.EdgeIdentifier
	EdgeType   int32
	Start      conf.VertexIdentifier
	End        conf.VertexIdentifier
	Properties map[string]string
}

type Vertex struct {
	Idntifier  conf.VertexIdentifier
	VertexType int32
	OutE       map[conf.EdgeIdentifier]bool
	InE        map[conf.EdgeIdentifier]bool
	Properties map[string]string
}


func NewEdge(id int64, start int64, end int64) *Edge {
	return &Edge{
		Idntifier:  conf.EdgeIdentifier(id),
		Start: conf.VertexIdentifier(start),
		End: conf.VertexIdentifier(end),
		Properties: map[string]string{},
	}
}

func NewVertex(id int64, vertexType int32) *Vertex {
	return &Vertex{
		Idntifier:  conf.VertexIdentifier(id),
		VertexType: vertexType,
		OutE:       map[conf.EdgeIdentifier]bool{},
		InE:        map[conf.EdgeIdentifier]bool{},
		Properties: map[string]string{},
	}
}
