package fbx

import (
	"io"
)

type Decoder struct {
	r    io.Reader
	body []byte
	p    int
	fbxf *FBX
	big  bool
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Decode() (*FBX, error) {
	var err error
	d.body, err = io.ReadAll(d.r)
	if err != nil {
		return nil, err
	}

	d.fbxf = new(FBX)

	binVersion := d.readBinaryVersion()
	if binVersion == nil {
		// ASCII
		panic("not implemented ascii decode")
	} else {
		// Binary
		d.fbxf.Version = int(*binVersion)
		d.fbxf.Nodes = make([]*Node, 0, 64)
		for {
			node := d.readBinNode()
			if node == nil {
				break
			}
			d.fbxf.Nodes = append(d.fbxf.Nodes, node)
		}
		return d.fbxf, nil
	}
}

// if nil, data is not binary
func (d *Decoder) readBinaryVersion() *int32 {
	magic := d.readBinSeeker(binMagicSize)
	if string(magic) != binHeader {
		d.p -= binMagicSize
		return nil
	}
	v := d.readBinInt32()
	d.big = v >= 7500
	return &v
}
