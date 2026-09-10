package spell

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	_ json.Marshaler   = Phoneme(0)
	_ json.Unmarshaler = (*Phoneme)(nil)
	_ yaml.Marshaler   = Phoneme(0)
	_ yaml.Unmarshaler = (*Phoneme)(nil)
)

const (
	PhonKA  = Phoneme(0) // upper-left
	PhonUN  = Phoneme(1) // up
	PhonIN  = Phoneme(2) // upper-right
	PhonET  = Phoneme(3) // left
	PhonEnd = Phoneme(4)
	PhonCHA = Phoneme(5) // right
	PhonRO  = Phoneme(6) // lower-left
	PhonZO  = Phoneme(7) // down
	PhonDO  = Phoneme(8) // lower-right
	PhonMax = 9
)

type Phoneme byte

func (p Phoneme) Valid() bool {
	return p >= 0 && p < PhonMax
}

func (p Phoneme) String() string {
	switch p {
	case PhonKA:
		return "ka"
	case PhonUN:
		return "un"
	case PhonIN:
		return "in"
	case PhonET:
		return "et"
	case PhonEnd:
		return "!"
	case PhonCHA:
		return "cha"
	case PhonRO:
		return "ro"
	case PhonZO:
		return "zo"
	case PhonDO:
		return "do"
	default:
		return fmt.Sprintf("Phoneme(%d)", int(p))
	}
}

func (p Phoneme) MarshalJSON() ([]byte, error) {
	if p.Valid() {
		return json.Marshal(p.String())
	}
	return json.Marshal(int(p))
}

func (p Phoneme) MarshalYAML() (interface{}, error) {
	if p.Valid() {
		return p.String(), nil
	}
	return int(p), nil
}

func (p *Phoneme) parseText(s string) error {
	switch strings.ToLower(s) {
	case "ka":
		*p = PhonKA
	case "un":
		*p = PhonUN
	case "in":
		*p = PhonIN
	case "et":
		*p = PhonET
	case "!", "":
		*p = PhonEnd
	case "cha":
		*p = PhonCHA
	case "ro":
		*p = PhonRO
	case "zo":
		*p = PhonZO
	case "do":
		*p = PhonDO
	default:
		return fmt.Errorf("unknown spell phoneme: %q", s)
	}
	return nil
}

func (p *Phoneme) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		var v int
		if err2 := json.Unmarshal(data, &v); err2 != nil {
			return err
		}
		*p = Phoneme(v)
		return nil
	}
	return p.parseText(s)
}

func (p *Phoneme) UnmarshalYAML(n *yaml.Node) error {
	var v uint8
	if err := n.Decode(&v); err == nil {
		*p = Phoneme(v)
		return nil
	}
	var s string
	if err := n.Decode(&s); err != nil {
		return err
	}
	return p.parseText(s)
}
