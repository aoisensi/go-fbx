package fbx75

import (
	"math/rand"

	"github.com/aoisensi/go-fbx/pkg/fbx"
)

type Documents struct {
	Documents []*Document
}

func NewDocuments() *Documents {
	return &Documents{
		Documents: []*Document{
			{
				ID:   rand.Int63(),
				Type: "Scene",
				Name: "Scene",
			},
		},
	}
}

func (s *Documents) Node() *fbx.Node {
	l := len(s.Documents)
	if l == 0 {
		return &fbx.Node{
			Name:     "Documents",
			Children: []*fbx.Node{Count(0)},
		}
	}
	node := &fbx.Node{
		Name:     "Documents",
		Children: make([]*fbx.Node, l+1),
	}
	node.Children[0] = Count(l)
	for i, doc := range s.Documents {
		node.Children[i+1] = &fbx.Node{
			Name:       "Document",
			Attributes: []any{doc.ID, doc.Name, doc.Type},
		}
	}
	return node
}

type Document struct {
	ID   int64
	Name string
	Type string
}
