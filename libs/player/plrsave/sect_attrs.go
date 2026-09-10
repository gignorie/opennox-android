package plrsave

import (
	"errors"
	"fmt"
)

func init() {
	RegisterSection(&Attributes{})
}

type Attributes struct {
	Vers    uint16
	Flags   uint32
	Name    string
	Val50   uint32
	Val54   uint32
	Val58   uint32
	Val62   uint32
	Val66   byte
	Val67   byte
	Val68   [3]byte
	Val71   [3]byte
	Val74   [3]byte
	Val77   [3]byte
	Val80   [3]byte
	Val83   byte
	Val84   byte
	Val85   byte
	Val86   byte
	Val87   byte
	Val88   byte
	Val320  uint32
	Unused1 [9]uint32
	Val4696 uint32
}

func (*Attributes) ID() uint32 {
	return 2
}

func (s *Attributes) SectName() string {
	return "Attrib Data"
}

func (s *Attributes) UnmarshalBinary(data []byte) error {
	*s = Attributes{}
	s.Vers, data = readU16(data)
	if s.Vers > 5 {
		return fmt.Errorf("unsupported version of attributes data: %v", s.Vers)
	}
	if s.Vers >= 5 {
		s.Flags, data = readU32(data)
	}
	s.Name, data = readWString8(data)

	s.Val50, data = readU32(data)
	s.Val54, data = readU32(data)
	s.Val58, data = readU32(data)
	s.Val62, data = readU32(data)
	s.Val66, data = readU8(data)
	s.Val67, data = readU8(data)
	s.Val68, data = readU24(data)
	s.Val71, data = readU24(data)
	s.Val74, data = readU24(data)
	s.Val77, data = readU24(data)
	s.Val80, data = readU24(data)
	if s.Vers >= 2 {
		s.Val83, data = readU8(data)
		s.Val84, data = readU8(data)
		s.Val85, data = readU8(data)
		s.Val86, data = readU8(data)
		s.Val87, data = readU8(data)
	}
	s.Val88, data = readU8(data)
	if s.Vers >= 3 {
		s.Val320, data = readU32(data)
		if s.Vers == 3 {
			for i := range s.Unused1 {
				s.Unused1[i], data = readU32(data)
			}
		}
	}
	if s.Vers >= 4 {
		s.Val4696, data = readU32(data)
	}
	if len(data) > 0 {
		return errors.New("unexpected trailing data in attributes")
	}
	return nil
}
