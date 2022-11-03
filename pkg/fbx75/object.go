package fbx75

import "github.com/aoisensi/go-fbx/pkg/fbx"

type Objects struct {
	Objects []Object
}

func (s *Objects) Node() *fbx.Node {
	if s.Objects == nil {
		return &fbx.Node{Name: "Objects"}
	}
	node := &fbx.Node{
		Name:     "Objects",
		Children: make([]*fbx.Node, 0, len(s.Objects)),
	}
	for _, c := range s.Objects {
		node.Children = append(node.Children, c.ObjectNode())
	}
	return node
}

type Object interface {
	ObjectNode() *fbx.Node
}

type ObjectGeometry struct {
	ID                 int64
	Name               string
	Vertices           []float64
	PolygonVertexIndex []int32
}

func (s *ObjectGeometry) ObjectNode() *fbx.Node {
	return &fbx.Node{
		Name:       "Geometry",
		Attributes: []any{s.ID, s.Name + "::Geometry", "Mesh"},
		Children: []*fbx.Node{
			{
				Name:       "Vertices",
				Attributes: []any{s.Vertices},
			},
			{
				Name:       "PolygonVertexIndex",
				Attributes: []any{s.PolygonVertexIndex},
			},
		},
	}
}

type ObjectModel struct {
	ID   int64
	Name string
}

func (s *ObjectModel) ObjectNode() *fbx.Node {
	return &fbx.Node{
		Name:       "Model",
		Attributes: []any{s.ID, s.Name + "::Model", "Mesh"},
		Children:   []*fbx.Node{Version(232)},
	}
}
