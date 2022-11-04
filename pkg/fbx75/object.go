package fbx75

import "github.com/aoisensi/go-fbx/pkg/fbx"

type Objects struct {
	Objects []Object
}

func NewObjects() *Objects {
	return &Objects{Objects: make([]Object, 0, 16)}
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

type Geometry struct {
	ID                    int64
	Name                  string
	Vertices              []float64
	PolygonVertexIndex    []int32
	LayerElementUV        *LayerElementUV
	LayerElementSmoothing *LayerElementSmoothing
	LayerElementMaterial  *LayerElementMaterial
	Layer                 *Layer
}

func (s *Geometry) ObjectNode() *fbx.Node {
	node := &fbx.Node{
		Name:       "Geometry",
		Attributes: []any{s.ID, s.Name + "::Geometry", "Mesh"},
		Children: []*fbx.Node{
			{
				Name:       "GeometryVersion",
				Attributes: []any{int32(124)},
			},
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
	if s.LayerElementUV != nil {
		node.AddChild(s.LayerElementUV.Node())
	}
	if s.LayerElementSmoothing != nil {
		node.AddChild(s.LayerElementSmoothing.Node())
	}
	if s.LayerElementMaterial != nil {
		node.AddChild(s.LayerElementMaterial.Node())
	}
	if s.Layer != nil {
		node.AddChild(s.Layer.Node())
	}
	return node
}

type Model struct {
	ID   int64
	Name string
}

func (s *Model) ObjectNode() *fbx.Node {
	return &fbx.Node{
		Name:       "Model",
		Attributes: []any{s.ID, s.Name + "::Model", "Mesh"},
		Children: []*fbx.Node{
			Version(232),
			{
				Name:       "Shading",
				Attributes: []any{true},
			},
		},
	}
}

type Material struct {
	ID   int64
	Name string
}

func (s *Material) ObjectNode() *fbx.Node {
	return &fbx.Node{
		Name:       "Material",
		Attributes: []any{s.ID, s.Name + "::Material", ""},
		Children: []*fbx.Node{
			Version(102),
			{
				Name:       "ShadingModel",
				Attributes: []any{"Phong"},
			},
			{
				Name:       "MultiLayer",
				Attributes: []any{int32(0)},
			},
		},
	}
}
