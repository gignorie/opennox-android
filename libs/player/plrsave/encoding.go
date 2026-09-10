package plrsave

import (
	"encoding/binary"
	"unicode/utf16"
)

func readU8(data []byte) (byte, []byte) {
	return data[0], data[1:]
}

func readU16(data []byte) (uint16, []byte) {
	return binary.LittleEndian.Uint16(data), data[2:]
}

func readU24(data []byte) ([3]byte, []byte) {
	var v [3]byte
	copy(v[:], data)
	return v, data[3:]
}

func readU32(data []byte) (uint32, []byte) {
	return binary.LittleEndian.Uint32(data), data[4:]
}

func readString8(data []byte) (string, []byte) {
	sz, data := readU8(data)
	str := data[:sz]
	data = data[sz:]
	return string(str), data
}

func readString16(data []byte) (string, []byte) {
	sz, data := readU16(data)
	str := data[:sz]
	data = data[sz:]
	return string(str), data
}

func readWString(data []byte, sz int) (string, []byte) {
	sbuf := make([]uint16, sz)
	for i := range sbuf {
		sbuf[i] = binary.LittleEndian.Uint16(data)
		data = data[2:]
	}
	return string(utf16.Decode(sbuf)), data
}

func readWString8(data []byte) (string, []byte) {
	sz := int(data[0])
	data = data[1:]
	return readWString(data, sz)
}

func readWString16(data []byte) (string, []byte) {
	sz := int(binary.LittleEndian.Uint16(data))
	data = data[2:]
	return readWString(data, sz)
}
