package util

import (
	"encoding/binary"
)

func Unit16Tobyte(in uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, in)
	return b
}

func Unit32Tobyte(in uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, in)
	return b
}

func byteToInt16(in []byte) int16 {
	return 1
}
