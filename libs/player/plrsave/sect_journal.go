package plrsave

func init() {
	RegisterSection(&Journal{})
}

type Journal struct {
	Data []byte
}

func (*Journal) ID() uint32 {
	return 9
}

func (s *Journal) SectName() string {
	return "Journal Data"
}

func (s *Journal) UnmarshalBinary(data []byte) error {
	s.Data = data
	return nil
}
