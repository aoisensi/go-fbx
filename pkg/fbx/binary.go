package fbx

import "encoding/binary"

const binMagicSize = 23

const binHeader = "Kaydara FBX Binary  \x00\x1a\x00"

var le = binary.LittleEndian

type binNodeHeader struct {
	Term      uint64
	AttrsNum  uint64
	AttrsSize uint64
	NameSize  uint8
}
