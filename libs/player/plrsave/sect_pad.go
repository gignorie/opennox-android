package plrsave

func init() {
	RegisterSection(&Padding{})
}

type Padding struct{}

func (*Padding) ID() uint32 {
	return 11
}

func (s *Padding) SectName() string {
	return "PAD_DATA"
}

func (s *Padding) UnmarshalBinary(data []byte) error {
	return nil
}
