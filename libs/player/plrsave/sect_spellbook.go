package plrsave

func init() {
	RegisterSection(&SpellBook{})
}

type SpellBook struct {
	Data []byte
}

func (*SpellBook) ID() uint32 {
	return 5
}

func (s *SpellBook) SectName() string {
	return "Spellbook Data"
}

func (s *SpellBook) UnmarshalBinary(data []byte) error {
	s.Data = data
	return nil
}
