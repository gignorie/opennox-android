package plrsave

func init() {
	RegisterSection(&GUIData{})
}

type GUIData struct {
	Data []byte
}

func (*GUIData) ID() uint32 {
	return 7
}

func (s *GUIData) SectName() string {
	return "GUI Data"
}

func (s *GUIData) UnmarshalBinary(data []byte) error {
	s.Data = data
	return nil
}
