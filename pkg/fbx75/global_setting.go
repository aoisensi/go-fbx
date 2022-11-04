package fbx75

import "github.com/aoisensi/go-fbx/pkg/fbx"

type GlobalSettings struct {
	UpAxis    *Axis
	FrontAxis *Axis
	CoordAxis *Axis
}

func (s *GlobalSettings) Node() *fbx.Node {
	ps := &Properties70{
		Ps: make([]P, 0, 64),
	}
	axes := map[string]*Axis{
		"UpAxis":    s.UpAxis,
		"FrontAxis": s.FrontAxis,
		"CoordAxis": s.CoordAxis,
	}
	for name, axis := range axes {
		if axis == nil {
			continue
		}
		ps.Ps = append(
			ps.Ps,
			&PInt{name, axis.axis},
			&PInt{name + "Sign", axis.sign},
		)
	}
	return &fbx.Node{
		Name: "GlobalSettings",
		Children: []*fbx.Node{
			Version(1000),
			ps.Node(),
		},
	}
}
