package fbx75

var (
	AxisXP = &Axis{0, +1}
	AxisYP = &Axis{1, +1}
	AxisZP = &Axis{2, +1}
	AxisXN = &Axis{0, -1}
	AxisYN = &Axis{1, -1}
	AxisZN = &Axis{2, -1}
)

type Axis struct {
	axis, sign int32
}
