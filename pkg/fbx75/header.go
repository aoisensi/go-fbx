package fbx75

import (
	"time"

	"github.com/aoisensi/go-fbx/pkg/fbx"
)

type FBXHeaderExtension struct {
	CreationTimeStamp time.Time
	Creator           string
}

func NewFBXHeaderExtension() *FBXHeaderExtension {
	return &FBXHeaderExtension{
		CreationTimeStamp: time.Now(),
		Creator:           "github.com/aoisensi/go-fbx",
	}
}

func (s *FBXHeaderExtension) Node() *fbx.Node {
	node := &fbx.Node{
		Name: "FBXHeaderExtension",
		Children: []*fbx.Node{
			{
				Name:       "FBXHeaderVersion",
				Attributes: []any{int32(1003)},
			},
			{
				Name:       "FBXVersion",
				Attributes: []any{int32(7500)},
			},
			{
				Name:       "EncryptionType",
				Attributes: []any{int32(0)},
			},
			{
				Name: "CreationTimeStamp",
				Children: []*fbx.Node{
					Version(1000),
					{Name: "Year", Attributes: []any{int32(s.CreationTimeStamp.Year())}},
					{Name: "Month", Attributes: []any{int32(int(s.CreationTimeStamp.Month()))}},
					{Name: "Day", Attributes: []any{int32(s.CreationTimeStamp.Day())}},
					{Name: "Hour", Attributes: []any{int32(s.CreationTimeStamp.Hour())}},
					{Name: "Minute", Attributes: []any{int32(s.CreationTimeStamp.Minute())}},
					{Name: "Second", Attributes: []any{int32(s.CreationTimeStamp.Second())}},
					{Name: "Millisecond", Attributes: []any{int32(0)}},
				},
			},
			{
				Name:       "Creator",
				Attributes: []any{s.Creator},
			},
		},
	}
	return node
}
