package plrsave

import (
	"errors"
	"fmt"
)

func init() {
	RegisterSection(&FileInfo{})
}

const FileInfoID = 1

type FileInfo struct {
	Vers  uint16
	Val0  uint32
	Path  string
	Str2  string
	Time  Time
	Val4  [3]byte
	Val5  [3]byte
	Val6  [3]byte
	Val7  [3]byte
	Val8  [3]byte
	Val9  byte
	Val10 byte
	Val11 byte
	Val12 byte
	Val13 byte
	Name  string
	Val15 byte
	Val16 byte
	Val17 byte
	Map   string
	Val19 byte
}

func (*FileInfo) ID() uint32 {
	return FileInfoID
}

func (s *FileInfo) SectName() string {
	return "File Info Data"
}

func (s *FileInfo) UnmarshalBinary(data []byte) error {
	*s = FileInfo{}
	s.Vers, data = readU16(data)
	if s.Vers > 12 {
		return fmt.Errorf("unsupported version of file info: %v", s.Vers)
	}
	s.Val0, data = readU32(data)

	s.Path, data = readString16(data)
	s.Str2, data = readString8(data)

	s.Time.Year, data = readU16(data)
	s.Time.Month, data = readU16(data)
	s.Time.DayOfWeek, data = readU16(data)
	s.Time.Day, data = readU16(data)
	s.Time.Hour, data = readU16(data)
	s.Time.Minute, data = readU16(data)
	s.Time.Second, data = readU16(data)
	s.Time.Milliseconds, data = readU16(data)

	s.Val4, data = readU24(data)
	s.Val5, data = readU24(data)
	s.Val6, data = readU24(data)
	s.Val7, data = readU24(data)
	s.Val8, data = readU24(data)

	s.Val9, data = readU8(data)
	s.Val10, data = readU8(data)
	s.Val11, data = readU8(data)
	s.Val12, data = readU8(data)
	s.Val13, data = readU8(data)

	s.Name, data = readWString8(data)

	s.Val15, data = readU8(data)
	s.Val16, data = readU8(data)
	s.Val17, data = readU8(data)

	if s.Vers >= 11 {
		s.Map, data = readString8(data)
	}
	if s.Vers >= 12 {
		s.Val19, data = readU8(data)
	}
	if len(data) > 0 {
		return errors.New("unexpected trailing data in file info")
	}
	return nil
}
