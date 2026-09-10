package plrsave

func init() {
	RegisterSection(&FieldGuide{})
}

type FieldGuide struct {
	Data []byte
}

func (*FieldGuide) ID() uint32 {
	return 8
}

func (s *FieldGuide) SectName() string {
	return "FieldGuide Data"
}

func (s *FieldGuide) UnmarshalBinary(data []byte) error {
	s.Data = data
	return nil
}
