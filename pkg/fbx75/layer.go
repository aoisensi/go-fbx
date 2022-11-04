package fbx75

import "github.com/aoisensi/go-fbx/pkg/fbx"

type Layer struct {
	LayerElements []*LayerElement
}

func (s *Layer) Node() *fbx.Node {
	node := &fbx.Node{
		Name:       "Layer",
		Attributes: []any{int32(0)},
		Children:   []*fbx.Node{Version(101)},
	}
	for _, e := range s.LayerElements {
		node.AddChild(e.Node())
	}
	return node
}

type LayerElement struct {
	Type       string
	TypedIndex int32
}

func (s *LayerElement) Node() *fbx.Node {
	return &fbx.Node{
		Name: "LayerElement",
		Children: []*fbx.Node{
			{Name: "Type", Attributes: []any{s.Type}},
			{Name: "TypedIndex", Attributes: []any{s.TypedIndex}},
		},
	}
}

type LayerElementUV struct {
	UV []float64
}

func (s *LayerElementUV) Node() *fbx.Node {
	return &fbx.Node{
		Name:       "LayerElementUV",
		Attributes: []any{int32(0)},
		Children: []*fbx.Node{
			Version(101),
			Name("UVMap"),
			{
				Name:       "MappingInformationType",
				Attributes: []any{"ByPolygonVertex"},
			},
			{
				Name:       "ReferenceInformationType",
				Attributes: []any{"Direct"},
			},
			{
				Name:       "UV",
				Attributes: []any{s.UV},
			},
		},
	}
}

type LayerElementSmoothing struct {
	Smoothing []int32
}

func (s *LayerElementSmoothing) Node() *fbx.Node {
	return &fbx.Node{
		Name:       "LayerElementSmoothing",
		Attributes: []any{int32(0)},
		Children: []*fbx.Node{
			Version(102),
			Name(""),
			{
				Name:       "MappingInformationType",
				Attributes: []any{"ByPolygon"},
			},
			{
				Name:       "ReferenceInformationType",
				Attributes: []any{"Direct"},
			},
			{
				Name:       "Smoothing",
				Attributes: []any{s.Smoothing},
			},
		},
	}
}

type LayerElementMaterial struct {
}

func (s *LayerElementMaterial) Node() *fbx.Node {
	return &fbx.Node{
		Name:       "LayerElementMaterial",
		Attributes: []any{int32(0)},
		Children: []*fbx.Node{
			Version(101),
			Name(""),
			{
				Name:       "MappingInformationType",
				Attributes: []any{"AllSame"},
			},
			{
				Name:       "ReferenceInformationType",
				Attributes: []any{"IndexToDirect"},
			},
			{
				Name:       "Materials",
				Attributes: []any{[]int32{0}},
			},
		},
	}
}
