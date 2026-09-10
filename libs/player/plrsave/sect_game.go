package plrsave

func init() {
	RegisterSection(&GameData{})
}

type GameData struct {
	Data []byte
}

func (*GameData) ID() uint32 {
	return 10
}

func (s *GameData) SectName() string {
	return "Game Data"
}

func (s *GameData) UnmarshalBinary(data []byte) error {
	s.Data = data
	return nil
}
