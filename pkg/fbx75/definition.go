package fbx75

import "github.com/aoisensi/go-fbx/pkg/fbx"

type Definitions struct {
}

func (s *Definitions) Node() *fbx.Node {
	return &fbx.Node{
		Name: "Definitions",
		Children: []*fbx.Node{
			Count(3),
			{
				Name:       "ObjectType",
				Attributes: []any{"Model"},
			},
			{
				Name:       "ObjectType",
				Attributes: []any{"Geometry"},
			},
			{
				Name:       "ObjectType",
				Attributes: []any{"Material"},
			},
		},
	}
}
