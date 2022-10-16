package fbxf

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"math"
)

const binMagicSize = 23

var le = binary.LittleEndian

type binNodeHeader struct {
	Term      uint64
	AttrsNum  uint64
	AttrsSize uint64
	NameSize  uint8
}

func (d *Decoder) readBinNode() *Node {
	header := d.readBinNodeHeader()
	if header == nil {
		return nil
	}
	return d.readBinNodeBody(header)
}

func (d *Decoder) readBinNodeBody(header *binNodeHeader) *Node {
	node := new(Node)

	// Name
	node.Name = string(d.body[d.p : d.p+int(header.NameSize)])
	d.p += int(header.NameSize)

	// Attributes
	node.Attributes = make([]any, int(header.AttrsNum))
	for i := 0; i < int(header.AttrsNum); i++ {
		node.Attributes[i] = d.readBinAttribute()
	}

	for {
		if d.p >= int(header.Term) {
			return node
		}
		childHeader := d.readBinNodeHeader()
		if childHeader == nil {
			return node
		}
		if node.Children == nil {
			node.Children = make([]*Node, 0, 16)
		}
		child := d.readBinNodeBody(childHeader)
		node.Children = append(node.Children, child)
	}
}

func (d *Decoder) readBinAttribute() any {
	code := d.readBinByte()
	switch code {
	case 'C': // bool
		return d.readBinBool()
	case 'Y': // int16
		return d.readBinInt16()
	case 'I': // int32
		return d.readBinInt32()
	case 'L': // int64
		return d.readBinInt64()
	case 'F': // float32
		return d.readBinFloat32()
	case 'D': // float64
		return d.readBinFloat64()
	}
	// slice
	if code == 'b' ||
		code == 'i' || code == 'l' ||
		code == 'f' || code == 'd' {
		num := d.readBinInt32()
		enc := d.readBinInt32()
		size := d.readBinInt32()

		data := d.readBinSeeker(int(size))

		switch code {
		case 'b': // []bool
			slice := readAttrSlice[byte](data, num, enc)
			bools := make([]bool, len(data))
			for i, b := range slice {
				bools[i] = byteToBool(b)
			}
			return bools
		case 'i': // []int32
			return readAttrSlice[int32](data, num, enc)
		case 'l': // []int64
			return readAttrSlice[int64](data, num, enc)
		case 'f': // []float32
			return readAttrSlice[float32](data, num, enc)
		case 'd': // []float64
			return readAttrSlice[float64](data, num, enc)
		}
	}
	if code == 'S' || code == 'R' {
		size := d.readBinInt32()
		data := make([]byte, size)
		copy(data, d.readBinSeeker(int(size)))
		if code == 'S' {
			return string(data)
		}
		return data
	}
	panic("unknown type code detected")
}

// if null it is node term marker
func (d *Decoder) readBinNodeHeader() *binNodeHeader {
	header := new(binNodeHeader)

	header.Term = d.readBinSize()
	header.AttrsNum = d.readBinSize()
	header.AttrsSize = d.readBinSize()
	header.NameSize = d.readBinByte()
	if header.Term == 0 && header.AttrsNum == 0 &&
		header.AttrsSize == 0 && header.NameSize == 0 {
		return nil
	}
	return header
}

func (d *Decoder) readBinSize() uint64 {
	if d.big {
		return uint64(d.readBinInt64())
	} else {
		return uint64(uint32(d.readBinInt32()))
	}
}

func (d *Decoder) readBinFloat64() float64 {
	return math.Float64frombits(uint64(d.readBinInt64()))
}

func (d *Decoder) readBinFloat32() float32 {
	return math.Float32frombits(uint32(d.readBinInt32()))
}

func (d *Decoder) readBinInt64() int64 {
	return int64(le.Uint64(d.readBinSeeker(8)))
}

func (d *Decoder) readBinInt32() int32 {
	return int32(le.Uint32(d.readBinSeeker(4)))
}

func (d *Decoder) readBinInt16() int16 {
	return int16(le.Uint16(d.readBinSeeker(2)))
}

func (d *Decoder) readBinBool() bool {
	return byteToBool(d.readBinByte())
}

func (d *Decoder) readBinByte() byte {
	return d.readBinSeeker(1)[0]
}

func (d *Decoder) readBinSeeker(size int) []byte {
	s := d.body[d.p : d.p+size]
	d.p += size
	return s
}

func readAttrSlice[T any](data []byte, size int32, enc int32) []T {
	result := make([]T, size)
	var buf io.Reader
	buf = bytes.NewBuffer(data)
	if enc == 1 {
		buf, _ = zlib.NewReader(buf)
	} else if enc != 0 {
		panic("unknown encode type")
	}
	binary.Read(buf, le, result)
	return result
}

func byteToBool(d byte) bool {
	return d&1 == 1
}
