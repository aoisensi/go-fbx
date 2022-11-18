package fbx75

import "github.com/aoisensi/go-fbx/pkg/fbx"

type Connections struct {
	Cs []C
}

func (c *Connections) Node() *fbx.Node {
	if len(c.Cs) == 0 {
		return &fbx.Node{Name: "Connections"}
	}
	node := &fbx.Node{
		Name:     "Connections",
		Children: make([]*fbx.Node, len(c.Cs)),
	}
	for i := range node.Children {
		node.Children[i] = &fbx.Node{
			Name:       "C",
			Attributes: []any{"OO", c.Cs[i][0], c.Cs[i][1]},
		}
	}
	return node
}

type C [2]int64

func (fbx *FBX) Connect(o1, o2 Object) {
	c := C{o1.Base().ID, o2.Base().ID}
	fbx.Connections.Cs = append(fbx.Connections.Cs, c)
}
