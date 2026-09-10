package player

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	_ json.Marshaler   = Class(0)
	_ json.Unmarshaler = (*Class)(nil)
)

const (
	Warrior  = Class(0)
	Wizard   = Class(1)
	Conjurer = Class(2)
)

type Class byte

func (c Class) String() string {
	switch c {
	case Warrior:
		return "warrior"
	case Wizard:
		return "wizard"
	case Conjurer:
		return "conjurer"
	default:
		return "unknown"
	}
}

func (c Class) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Class) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	switch s {
	case "":
		*c = 0
	case "warrior":
		*c = Warrior
	case "wizard":
		*c = Wizard
	case "conjurer":
		*c = Conjurer
	default:
		return fmt.Errorf("invalid class: %q", s)
	}
	return nil
}
