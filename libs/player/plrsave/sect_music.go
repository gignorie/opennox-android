package plrsave

func init() {
	RegisterSection(&Music{})
}

type Music struct {
	Data []byte
}

func (*Music) ID() uint32 {
	return 12
}

func (s *Music) SectName() string {
	return "Music Data"
}

func (s *Music) UnmarshalBinary(data []byte) error {
	s.Data = data
	return nil
}
