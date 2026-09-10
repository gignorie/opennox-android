package server

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/opennox/libs/spell"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func init() {
	RegisterObjectUseParse("ConsumeUse", parseUseConsume)
	RegisterObjectUseParse("ConsumeConfuseUse", parseUseConsume)
	RegisterObjectUseParse("CastUse", parseUseCast)
	RegisterObjectUseParse("EnchantUse", parseUseEnchant)
	RegisterObjectUseParse("PotionUse", parseUsePotion)
}

type UseData interface {
	UseDataPtr() unsafe.Pointer
}

func useDataPtrAs[T any, P interface {
	*T
	UseData
}](p UseDataPtr) P {
	if alloc.IsDead(p.Ptr) {
		panic("object already deleted")
	}
	return (*T)(p.Ptr)
}

type UseDataPtr struct {
	Ptr unsafe.Pointer
}

func (p *UseDataPtr) Free() {
	if p.Ptr != nil {
		alloc.FreePtr(p.Ptr)
		p.Ptr = nil
	}
}

func (p *UseDataPtr) SetPtr(ptr unsafe.Pointer) {
	if ptr == nil {
		p.Free()
		return
	}
	if p.Ptr != ptr {
		p.Free()
	}
	p.Ptr = ptr
}

func (p *UseDataPtr) Set(data UseData) {
	if data == nil {
		p.Free()
		return
	}
	ptr := data.UseDataPtr()
	if p.Ptr != ptr {
		p.Free()
	}
	p.Ptr = ptr
}

func (p UseDataPtr) AsCast() *CastUseData {
	return useDataPtrAs[CastUseData](p)
}

func (p UseDataPtr) AsPotion() *PotionUseData {
	return useDataPtrAs[PotionUseData](p)
}

func (p UseDataPtr) AsConsume() *ConsumeUseData {
	return useDataPtrAs[ConsumeUseData](p)
}

func (p UseDataPtr) AsEnchant() *EnchantUseData {
	return useDataPtrAs[EnchantUseData](p)
}

func (p UseDataPtr) AsSpellReward() *SpellRewardUseData {
	return useDataPtrAs[SpellRewardUseData](p)
}

func (p UseDataPtr) AsAbilityReward() *AbilityRewardUseData {
	return useDataPtrAs[AbilityRewardUseData](p)
}

func (p UseDataPtr) AsFieldGuide() *FieldGuideUseData {
	return useDataPtrAs[FieldGuideUseData](p)
}

type ConsumeUseData struct {
	Value int32
}

func (d *ConsumeUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

type PotionUseData struct {
	Value int32
}

func (d *PotionUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

type CastUseData struct {
	Spell int32
}

func (d *CastUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

type EnchantUseData struct {
	Enchant int32
	Dur     int32
}

func (d *EnchantUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

type SpellRewardUseData struct {
	Spell byte
}

func (d *SpellRewardUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

type AbilityRewardUseData struct {
	Ability byte
}

func (d *AbilityRewardUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

type FieldGuideUseData struct {
	CreatureBuf [64]byte
}

func (d *FieldGuideUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

func (d *FieldGuideUseData) Creature() string {
	return alloc.GoStringS(d.CreatureBuf[:])
}

func (d *FieldGuideUseData) SetCreature(name string) {
	alloc.StrCopyZero(d.CreatureBuf[:], name)
}

func parseUseConsume(objt *ObjectType, args []string) error {
	use := objt.UseData.AsConsume()
	var s string
	if len(args) != 0 {
		s = args[0]
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	use.Value = int32(v)
	return nil
}

func parseUsePotion(objt *ObjectType, args []string) error {
	use := objt.UseData.AsPotion()
	var s string
	if len(args) != 0 {
		s = args[0]
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	use.Value = int32(v)
	return nil
}

func parseUseCast(objt *ObjectType, args []string) error {
	use := objt.UseData.AsCast()
	var s string
	if len(args) != 0 {
		s = args[0]
	}
	ind := spell.ParseID(s)
	if ind <= 0 {
		return fmt.Errorf("unsupported spell: %q", s)
	}
	use.Spell = int32(ind)
	return nil
}

func parseUseEnchant(objt *ObjectType, args []string) error {
	use := objt.UseData.AsEnchant()
	var s1, s2 string
	if len(args) >= 1 {
		s1 = args[0]
	}
	if len(args) >= 2 {
		s2 = args[1]
	}
	id, ok := ParseEnchant(s1)
	if !ok {
		return fmt.Errorf("unsupported enchant: %q", s1)
	}
	use.Enchant = int32(id)

	dur, err := strconv.Atoi(s2)
	if err != nil {
		return err
	}
	use.Dur = int32(dur)
	return nil
}
