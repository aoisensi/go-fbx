package fbxf

type Node struct {
	Name       string
	Attributes []any
	Children   []*Node
}
