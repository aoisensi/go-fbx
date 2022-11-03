package fbx

import "golang.org/x/exp/slices"

type FBX struct {
	Version int
	Nodes   []*Node
}

func (fbx *FBX) AddNode(child *Node) *FBX {
	fbx.Nodes = append(fbx.Nodes, child)
	return fbx
}

func (fbx *FBX) AddNodes(nodes ...*Node) *FBX {
	fbx.Nodes = append(fbx.Nodes, nodes...)
	return fbx
}

func (fbx *FBX) FindNode(name string) *Node {
	for _, node := range fbx.Nodes {
		if node.Name == name {
			return node
		}
	}
	return nil
}

func (fbx *FBX) RemoveNode(name string) *Node {
	for i, node := range fbx.Nodes {
		if node.Name == name {
			fbx.Nodes = slices.Delete(fbx.Nodes, i, i+1)
			return node
		}
	}
	return nil
}
