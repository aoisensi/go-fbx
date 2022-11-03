package fbx

import "golang.org/x/exp/slices"

type Node struct {
	Name       string
	Attributes []any
	Children   []*Node
}

func (node *Node) AddChild(child *Node) *Node {
	node.Children = append(node.Children, child)
	return node
}

func (node *Node) AddChildren(child ...*Node) *Node {
	node.Children = append(node.Children, child...)
	return node
}

func (node *Node) FindFirstChild(name string) *Node {
	for _, c := range node.Children {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (node *Node) FindAllChildren(name string) []*Node {
	children := make([]*Node, 0)
	for _, c := range node.Children {
		if c.Name == name {
			children = append(children, c)
		}
	}
	return children
}

func (node *Node) RemoveChild(name string) *Node {
	for i, c := range node.Children {
		if c.Name == name {
			node.Children = slices.Delete(node.Children, i, i+1)
			return c
		}
	}
	return nil
}
