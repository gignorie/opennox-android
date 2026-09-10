package binenc

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

func allZeros(data []byte) bool {
	for _, v := range data {
		if v != 0 {
			return false
		}
	}
	return true
}

func CLen(data []byte) int {
	if i := bytes.IndexByte(data, 0); i >= 0 {
		return i
	}
	return len(data)
}

func CString(data []byte) string {
	data = data[:CLen(data)]
	return string(data)
}

func CStringSet0(data []byte, s string) int {
	zeros(data)
	n := copy(data, s)
	i := n
	if i >= len(data) {
		i = len(data) - 1
	}
	data[i] = 0
	return n
}

func CStringSet(data []byte, s string) int {
	zeros(data)
	return copy(data, s)
}

func CLen16(data []byte) int {
	for i := 0; i < len(data)/2; i++ {
		v := binary.LittleEndian.Uint16(data[2*i:])
		if v == 0 {
			return i
		}
	}
	return len(data)
}

func CString16(data []byte) string {
	data16 := make([]uint16, len(data)/2)
	for i := range data16 {
		v := binary.LittleEndian.Uint16(data[2*i:])
		if v == 0 {
			data16 = data16[:i]
			break
		}
		data16[i] = v
	}
	return string(utf16.Decode(data16))
}

func CStringSet16(data []byte, s string) int {
	zeros(data)
	data16 := utf16.Encode([]rune(s))
	for i, v := range data16 {
		binary.LittleEndian.PutUint16(data[2*i:], v)
	}
	return 2 * len(data16)
}

func zeros(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

type String struct {
	Value string
	Junk  []byte
}

func (s *String) Encode(data []byte) {
	i := CStringSet(data, s.Value)
	if len(s.Junk) != 0 {
		copy(data[i+1:], s.Junk)
	}
}

func (s *String) Decode(data []byte) {
	s.Value = CString(data)
	s.Junk = nil
	if i := len(s.Value); !allZeros(data[i+1:]) {
		s.Junk = make([]byte, len(data)-i-1)
		copy(s.Junk, data[i+1:])
	}
}
