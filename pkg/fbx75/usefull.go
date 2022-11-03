package fbx75

import "github.com/aoisensi/go-fbx/pkg/fbx"

func Count(c int) *fbx.Node {
	return &fbx.Node{
		Name:       "Count",
		Attributes: []any{int32(c)},
	}
}

func Version(v int) *fbx.Node {
	return &fbx.Node{
		Name:       "Version",
		Attributes: []any{int32(v)},
	}
}
