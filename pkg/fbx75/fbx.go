package fbx75

import (
	"io"

	"github.com/aoisensi/go-fbx/pkg/fbx"
)

var fileID = []byte{0x28, 0xb3, 0x2a, 0xeb, 0xb6, 0x24, 0xcc, 0xc2, 0xbf, 0xc8, 0xb0, 0x2a, 0xa9, 0x2b, 0xfc, 0xf1}
var creationTime = "1970-01-01 10:00:00:000"

type FBX struct {
	Documents   *Documents
	Objects     *Objects
	Connections *Connections
}

func NewFBX() *FBX {
	return &FBX{
		Documents:   NewDocuments(),
		Objects:     &Objects{},
		Connections: &Connections{},
	}
}

func (f *FBX) WriteTo(w io.Writer) (int64, error) {
	fbxr := &fbx.FBX{
		Version: 7500,
		Nodes:   []*fbx.Node{},
	}
	fbxr.AddNodes(
		&fbx.Node{Name: "FileId", Attributes: []any{fileID}},
		&fbx.Node{Name: "CreationTime", Attributes: []any{creationTime}},
	)
	if f.Documents != nil {
		fbxr.Nodes = append(fbxr.Nodes, f.Documents.Node())
	}
	fbxr.Nodes = append(fbxr.Nodes, &fbx.Node{Name: "References"})
	fbxr.Nodes = append(fbxr.Nodes, &fbx.Node{
		Name:     "Definitions",
		Children: []*fbx.Node{Version(100), Count(0)},
	})
	if f.Objects != nil {
		fbxr.Nodes = append(fbxr.Nodes, f.Objects.Node())
	}
	if f.Connections != nil {
		fbxr.Nodes = append(fbxr.Nodes, f.Connections.Node())
	}

	return fbx.NewBinaryEncoder(w).Encode(fbxr)
}
