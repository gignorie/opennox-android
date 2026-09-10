package plrsave

func init() {
	RegisterSection(&Enchant{})
}

type Enchant struct {
	Data []byte
}

func (*Enchant) ID() uint32 {
	return 6
}

func (s *Enchant) SectName() string {
	return "Enchantment Data"
}

func (s *Enchant) UnmarshalBinary(data []byte) error {
	s.Data = data
	return nil
}
