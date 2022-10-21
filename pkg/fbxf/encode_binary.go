package fbxf

import (
	"bytes"
	"encoding/binary"
	"io"
)

var footer1 = make([]byte, 16)
var footer2 = make([]byte, 4)
var footer3 = make([]byte, 120)
var footer4 = []byte{0xf8, 0x5a, 0x8c, 0x6a, 0xde, 0xf5, 0xd9, 0x7e, 0xec, 0xe9, 0x0c, 0xe3, 0x75, 0x8f, 0x29, 0x0b}

type BinaryEncoder struct {
	big bool
	w   io.Writer
	buf *bytes.Buffer
}

func NewBinaryEncoder(w io.Writer) *BinaryEncoder {
	return &BinaryEncoder{w: w}
}

func (e *BinaryEncoder) Encode(fbxf *FBXF) error {
	e.buf = bytes.NewBuffer(make([]byte, 0, 1024*1024*16))
	e.big = fbxf.Version >= 7500
	if err := e.write(binHeader); err != nil {
		return err
	}
	if err := e.write(int32(fbxf.Version)); err != nil {
		return err
	}
	for _, node := range fbxf.Nodes {
		if err := e.writeNode(node); err != nil {
			return err
		}
	}
	if err := e.writeNodeHeader(nil); err != nil {
		return err
	}
	if err := e.write(footer1); err != nil {
		return err
	}
	// padding
	if err := e.write(make([]byte, e.buf.Len()%16)); err != nil {
		return err
	}
	if err := e.write(footer2); err != nil {
		return err
	}
	if err := e.write(int32(fbxf.Version)); err != nil {
		return err
	}
	if err := e.write(footer3); err != nil {
		return err
	}
	if err := e.write(footer4); err != nil {
		return err
	}
	_, err := e.buf.WriteTo(e.w)
	return err
}

func (e *BinaryEncoder) writeNode(node *Node) error {
	termPos := e.buf.Len()
	if err := e.writeSize(0); err != nil {
		return err
	}
	if err := e.writeSize(uint64(len(node.Attributes))); err != nil {
		return err
	}
	attrsSizePos := e.buf.Len()
	if err := e.writeSize(0); err != nil {
		return err
	}
	if err := e.write(uint8(len(node.Name))); err != nil {
		return err
	}
	if err := e.write([]byte(node.Name)); err != nil {
		return err
	}
	attrsStartPos := e.buf.Len()
	for _, a := range node.Attributes {
		var code byte
		switch a.(type) {
		case bool:
			code = 'C'
		case int16:
			code = 'Y'
		case int32:
			code = 'I'
		case int64:
			code = 'L'
		case float32:
			code = 'F'
		case float64:
			code = 'D'
		case []bool:
			code = 'b'
		case []int32:
			code = 'i'
		case []int64:
			code = 'l'
		case []float32:
			code = 'f'
		case []float64:
			code = 'd'
		case []byte:
			code = 'R'
		case string:
			code = 'S'
		default:
			panic("invalied attribute type detected.")
		}
		if err := e.write(code); err != nil {
			return err
		}
		switch a := a.(type) {
		case bool:
			if err := e.writeBool(a); err != nil {
				return err
			}
		case []bool:
			for _, c := range a {
				if err := e.writeBool(c); err != nil {
					return err
				}
			}
		default:
			if err := e.write(a); err != nil {
				return err
			}
		}
	}
	attrsSize := e.buf.Len() - attrsStartPos
	for _, childNode := range node.Children {
		if err := e.writeNode(childNode); err != nil {
			return err
		}
	}
	if len(node.Attributes) == 0 || len(node.Children) > 0 {
		if err := e.writeNodeHeader(nil); err != nil {
			return err
		}
	}
	term := e.buf.Len()
	e.writeSizeAt(uint64(term), termPos)
	e.writeSizeAt(uint64(attrsSize), attrsSizePos)
	return nil
}

func (e *BinaryEncoder) writeNodeHeader(header *binNodeHeader) error {
	if header == nil {
		header = new(binNodeHeader)
	}
	if err := e.writeSize(header.Term); err != nil {
		return err
	}
	if err := e.writeSize(header.AttrsNum); err != nil {
		return err
	}
	if err := e.writeSize(header.AttrsSize); err != nil {
		return err
	}
	return e.write(header.NameSize)
}

func (e *BinaryEncoder) writeBool(c bool) error {
	if c {
		return e.write('Y')
	} else {
		return e.write('T')
	}
}

func (e *BinaryEncoder) writeSizeAt(size uint64, pos int) {
	buf := e.buf.Bytes()
	i := 4
	if e.big {
		i = 8
	}
	for ; i > 0; i-- {
		buf[pos] = byte(size)
		size >>= 8
		pos++
	}
}

func (e *BinaryEncoder) writeSize(size uint64) error {
	if e.big {
		return e.write(size)
	} else {
		return e.write(uint32(size))
	}
}

func (e *BinaryEncoder) write(data any) error {
	if s, ok := data.(string); ok {
		return binary.Write(e.buf, le, []byte(s))
	}
	return binary.Write(e.buf, le, data)
}
