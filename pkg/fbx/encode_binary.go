package fbx

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"reflect"
	"strings"
)

var (
	footer1 = []byte{0xfa, 0xbc, 0xab, 0x09, 0xd0, 0xc8, 0xd4, 0x66, 0xb1, 0x76, 0xfb, 0x83, 0x1c, 0xf7, 0x26, 0x7e}
	footer2 = make([]byte, 4)
	footer3 = make([]byte, 120)
	footer4 = []byte{0xf8, 0x5a, 0x8c, 0x6a, 0xde, 0xf5, 0xd9, 0x7e, 0xec, 0xe9, 0x0c, 0xe3, 0x75, 0x8f, 0x29, 0x0b}
)

var (
	dataTrue  byte = 'Y'
	dataFalse byte = 'T'
)

type BinaryEncoder struct {
	big bool
	w   io.Writer
	buf *bytes.Buffer
}

func NewBinaryEncoder(w io.Writer) *BinaryEncoder {
	return &BinaryEncoder{w: w}
}

func (e *BinaryEncoder) Encode(fbxf *FBX) error {
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
	if err := e.write(make([]byte, 16-e.buf.Len()%16)); err != nil {
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
		default:
			if err := e.write(a); err != nil {
				return err
			}
		case []bool, []int32, []int64, []float32, []float64:
			l := reflect.ValueOf(a).Len()
			if err := e.write(uint32(l)); err != nil {
				return err
			}
			data := new(bytes.Buffer)
			if ab, ok := a.([]bool); ok {
				for _, c := range ab {
					if err := binary.Write(data, le, dataBool(c)); err != nil {
						return err
					}
				}
			} else {
				if err := binary.Write(data, le, a); err != nil {
					return err
				}
			}
			encode := uint32(0)
			if data.Len() >= 128 {
				cdata := new(bytes.Buffer)
				cw := zlib.NewWriter(cdata)
				if _, err := io.Copy(cw, data); err != nil {
					cw.Close()
					return err
				}
				cw.Close()
				data = cdata
				encode = 1
			}
			if err := e.write(encode); err != nil {
				return err
			}
			if err := e.write(uint32(data.Len())); err != nil {
				return err
			}
			if err := e.write(data.Bytes()); err != nil {
				return err
			}
		case []byte:
			if err := e.write(uint32(len(a))); err != nil {
				return err
			}
			if err := e.write(a); err != nil {
				return err
			}
		case string:
			a = strings.ReplaceAll(a, "::", "\x00\x01")
			if err := e.write(uint32(len([]byte(a)))); err != nil {
				return err
			}
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

func dataBool(c bool) byte {
	if c {
		return dataTrue
	} else {
		return dataFalse
	}
}

func (e *BinaryEncoder) write(data any) error {
	if s, ok := data.(bool); ok {
		return e.write(dataBool(s))
	}
	if s, ok := data.([]bool); ok {
		for _, c := range s {
			if err := e.write(c); err != nil {
				return err
			}
		}
		return nil
	}
	if s, ok := data.(string); ok {
		return binary.Write(e.buf, le, []byte(s))
	}
	return binary.Write(e.buf, le, data)
}
