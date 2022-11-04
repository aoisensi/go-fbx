package fbx75

import "github.com/aoisensi/go-fbx/pkg/fbx"

type Properties70 struct {
	Ps []P
}

func (p *Properties70) Node() *fbx.Node {
	node := &fbx.Node{
		Name:     "Properties70",
		Children: make([]*fbx.Node, len(p.Ps)),
	}
	for i, p := range p.Ps {
		node.Children[i] = p.PNode()
	}
	return node
}

type P interface {
	PNode() *fbx.Node
}

type PInt struct {
	Name  string
	Value int32
}

func (p *PInt) PNode() *fbx.Node {
	return &fbx.Node{
		Name: "P",
		Attributes: []any{
			p.Name,
			"int",
			"Integer",
			"",
			p.Value,
		},
	}
}
