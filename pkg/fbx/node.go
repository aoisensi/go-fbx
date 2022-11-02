package fbx

type Node struct {
	Name       string
	Attributes []any
	Children   []*Node
}
