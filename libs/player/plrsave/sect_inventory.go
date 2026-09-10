package plrsave

func init() {
	RegisterSection(&Inventory{})
}

type Inventory struct {
	Data []byte
}

func (*Inventory) ID() uint32 {
	return 4
}

func (s *Inventory) SectName() string {
	return "Inventory Data"
}

func (s *Inventory) UnmarshalBinary(data []byte) error {
	s.Data = data
	return nil
}
