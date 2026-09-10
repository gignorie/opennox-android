package server

import (
	"unsafe"

	"github.com/opennox/libs/player"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

var (
	_ = [1]struct{}{}[72-unsafe.Sizeof(SavePlayerInfo{})]
	_ = [1]struct{}{}[20-unsafe.Offsetof(SavePlayerInfo{}.NameBuf)]
)

type SavePlayerInfo struct {
	SkinColor      types.RGB  // 12184 (+0)
	HairColor      types.RGB  // 12187 (+3)
	MustacheColor  types.RGB  // 12190 (+6)
	GoateeColor    types.RGB  // 12193 (+9)
	BeardColor     types.RGB  // 12196 (+12)
	PantsColor     byte       // 12199 (+15)
	Shirt1Color    byte       // 12200 (+16)
	Shirt2Color    byte       // 12201 (+17)
	Shoes1Color    byte       // 12202 (+18)
	Shoes2Color    byte       // 12203 (+19)
	NameBuf        [25]uint16 // 12204 (+20)
	PlayerClassVal byte       // 12254 (+70)
	IsFemaleVal    byte       // 12255 (+71)
}

func (p *SavePlayerInfo) Name() string {
	return alloc.GoString16S(p.NameBuf[:])
}

func (p *SavePlayerInfo) SetName(v string) {
	alloc.StrCopyZero16(p.NameBuf[:], v)
}

func (p *SavePlayerInfo) PlayerClass() player.Class {
	if p == nil {
		return 0
	}
	return player.Class(p.PlayerClassVal)
}

var (
	_ = [1]struct{}{}[1277-unsafe.Offsetof(SaveGameInfo{}.Stage)]
	_ = [1]struct{}{}[1280-unsafe.Sizeof(SaveGameInfo{})] // FIXME: should be 1278
)

type SaveGameInfo struct {
	Flags      uint32         // 10980 (+0)
	PathBuf    [1024]byte     // 10984 (+4)
	Field1028  [128]byte      // 12008 (+1028)
	MapNameBuf [32]byte       // 12136 (+1156)
	Timestamp  SystemTime     // 12168 (+1188)
	Player     SavePlayerInfo // 12184 (+1204)
	Field1276  byte           // 12256 (+1276)
	Stage      byte           // 12257 (+1277)
}

func (s *SaveGameInfo) Path() string {
	return alloc.GoStringS(s.PathBuf[:])
}
