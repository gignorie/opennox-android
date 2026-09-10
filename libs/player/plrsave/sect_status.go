package plrsave

import (
	"errors"
	"fmt"
)

func init() {
	RegisterSection(&Status{})
}

type Status struct {
	Vers      uint16
	Val1      byte
	MaxHealth uint16
	MaxMana   uint16
	Val4      uint16
	Val5      uint16
	Val6      byte
	Val7      byte
	Val8      uint16
	Val9      uint32
	Val10     uint16
}

func (*Status) ID() uint32 {
	return 3
}

func (s *Status) SectName() string {
	return "Status Data"
}

func (s *Status) UnmarshalBinary(data []byte) error {
	*s = Status{}
	s.Vers, data = readU16(data)
	if s.Vers > 2 {
		return fmt.Errorf("unsupported version of status data: %v", s.Vers)
	}
	s.Val1, data = readU8(data)
	if s.Val1 != 0 {
		s.MaxHealth, data = readU16(data)
		s.MaxMana, data = readU16(data)
		s.Val4, data = readU16(data)
		s.Val5, data = readU16(data)
		s.Val6, data = readU8(data)
		s.Val7, data = readU8(data)
		s.Val8, data = readU16(data)
		s.Val9, data = readU32(data)
		if s.Vers >= 2 {
			s.Val10, data = readU16(data)
		}
	}
	if len(data) > 0 {
		return errors.New("unexpected trailing data in status")
	}
	return nil
}
