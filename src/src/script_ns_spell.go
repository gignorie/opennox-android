package opennox

import (
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
	ns4 "github.com/opennox/noxscript/ns/v4"
	nab "github.com/opennox/noxscript/ns/v4/abil"
	"github.com/opennox/noxscript/ns/v4/enchant"
	nsp "github.com/opennox/noxscript/ns/v4/spell"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func (s noxScriptNS) CastSpellLvl(name nsp.Spell, lvl int, source, target ns4.Positioner) {
	sp := spell.ParseID(string(name))
	if !sp.Valid() {
		return
	}
	if source == nil || target == nil {
		return
	}
	srcH, _ := source.(server.Obj)
	targH, _ := target.(server.Obj)
	src := server.ToObject(srcH)
	if src != nil && src.Flags().HasAny(object.FlagDestroyed|object.FlagDead) {
		return
	}
	targPos := target.Pos()
	if src == nil {
		asObjectS(nox_xxx_imagCasterUnit_1569664).SetPos(source.Pos())
		src = nox_xxx_imagCasterUnit_1569664
	}
	src.Direction1 = server.DirFromVec(targPos.Sub(src.Pos()))
	// TODO: pass spell level
	s.s.castSpellBy(sp, lvl, src, targH, targPos)
}

func (s noxScriptNS) CastSpell(name nsp.Spell, source, target ns4.Positioner) {
	s.CastSpellLvl(name, -1, source, target)
}

func (s noxScriptNS) NewTrap(pos ns4.Positioner, spells []ns4.TrapSpell) ns4.Obj {
	if pos == nil {
		return nil
	}
	trap := s.s.NewObjectByTypeID("Glyph")
	if trap == nil {
		return nil
	}
	p := pos.Pos()
	s.s.CreateObjectAt(trap, nil, p)

	idata := trap.InitDataGlyph()
	*idata = server.GlyphInitData{
		SpellsCnt: 0,
		SpellArg: server.SpellAcceptArg{
			Pos: p,
		},
	}
	for _, sp := range spells {
		spl := spell.ParseID(string(sp.Spell))
		if !spl.Valid() {
			continue
		}
		idata.Spells[idata.SpellsCnt] = uint32(spl)
		idata.SpellsCnt++
	}
	return s.toObj(trap)
}

func newUseItem[T comparable, P interface {
	*T
	server.UseData
}](s noxScriptNS, typ string, pos ns4.Positioner, cfnc unsafe.Pointer, noPickup bool, data T) ns4.Obj {
	obj := s.createObject(typ, pos)
	if obj == nil {
		return nil
	}
	use, _ := alloc.New[T](data)
	*use = data
	var puse P = use
	obj.SetUse(cfnc, puse)
	if noPickup {
		obj.SetPickup(nil)
	}
	return s.toObj(obj)
}

func (s noxScriptNS) newPotion(typ string, pos ns4.Positioner, val int) ns4.Obj {
	return newUseItem(s, typ, pos, legacy.Get_nox_xxx_usePotion_53EF70(), false, server.PotionUseData{
		Value: int32(val),
	})
}

func (s noxScriptNS) NewHealthPotion(pos ns4.Positioner, health int) ns4.Obj {
	return s.newPotion("RedPotion", pos, health)
}

func (s noxScriptNS) NewManaPotion(pos ns4.Positioner, mana int) ns4.Obj {
	return s.newPotion("BluePotion", pos, mana)
}

func (s noxScriptNS) NewSpellBook(pos ns4.Positioner, name nsp.Spell) ns4.Obj {
	sp := spell.ParseID(string(name))
	if !sp.Valid() {
		return s.CreateObject("CommonSpellBook", pos)
	}
	typ := "CommonSpellBook"
	flags := s.s.Spells.Flags(sp)
	if flags.Has(things.SpellClassWizard) {
		typ = "WizardSpellBook"
	} else if flags.Has(things.SpellClassConjurer) {
		typ = "ConjurerSpellBook"
	}
	return newUseItem(s, typ, pos, legacy.Get_nox_xxx_useSpellReward_53F9E0(), false, server.SpellRewardUseData{
		Spell: byte(sp),
	})
}

func (s noxScriptNS) NewAbilityBook(pos ns4.Positioner, name nab.Ability) ns4.Obj {
	a := s.s.abilities.nox_xxx_abilityNameToN_424D80(string(name))
	if !a.Valid() {
		return s.CreateObject("AbilityBook", pos)
	}
	return newUseItem(s, "AbilityBook", pos, legacy.Get_nox_xxx_useAbilityReward_53FAE0(), false, server.AbilityRewardUseData{
		Ability: byte(a),
	})
}

func (s noxScriptNS) NewEnchantUseItem(typ ns4.ObjTypeName, pos ns4.Positioner, enc enchant.Enchant, dur ns4.Duration) ns4.Obj {
	if typ == "" {
		typ = "YellowPotion"
	}
	id, ok := server.ParseEnchant(string(enc))
	if !ok {
		return s.CreateObject(typ, pos)
	}
	return newUseItem(s, typ, pos, legacy.Get_nox_xxx_useEnchant_53ED60(), true, server.EnchantUseData{
		Enchant: int32(id),
		Dur:     int32(s.s.AsFrames(dur)),
	})
}

func (s noxScriptNS) NewSpellUseItem(typ ns4.ObjTypeName, pos ns4.Positioner, name nsp.Spell, lvl int) ns4.Obj {
	if typ == "" {
		typ = "SpellBook"
	}
	sp := spell.ParseID(string(name))
	if !sp.Valid() {
		return s.CreateObject(typ, pos)
	}
	return newUseItem(s, typ, pos, legacy.Get_nox_xxx_useCast_53ED90(), true, server.CastUseData{
		Spell: int32(sp),
		// TODO: pass spell level in ExtData
	})
}

func (s noxScriptNS) NewFieldGuide(pos ns4.Positioner, creature string) ns4.Obj {
	var data server.FieldGuideUseData
	data.SetCreature(creature)
	return newUseItem(s, "FieldGuide", pos, legacy.Get_sub_53F930(), false, data)
}

func (obj nsObj) AwardSpell(name nsp.Spell) bool {
	sp := spell.ParseID(string(name))
	if !sp.Valid() {
		return false
	}
	return legacy.Nox_xxx_spellGrantToPlayer_4FB550(obj.SObj(), sp, 1, 0, 0) != 0
}

func (obj nsObj) Enchant(enc enchant.Enchant, dt ns4.Duration) {
	id, ok := server.ParseEnchant(string(enc))
	if !ok {
		return
	}
	s := obj.Server()
	obj.ApplyEnchant(id, s.AsFrames(dt), 5)
}

func (obj nsObj) EnchantOff(enc enchant.Enchant) {
	e, ok := server.ParseEnchant(string(enc))
	if !ok {
		return
	}
	obj.DisableEnchant(e)
}

func (obj nsObj) TrapSpells(spells ...nsp.Spell) {
	out := make([]spell.ID, 0, len(spells))
	for _, name := range spells {
		sp := spell.ParseID(string(name))
		if !sp.Valid() {
			continue
		}
		out = append(out, sp)
	}
	obj.Object.SetTrapSpells(out...)
}

func (obj nsObj) TrapSpellsAdv(spells []ns4.TrapSpell) {
	out := make([]spell.ID, 0, len(spells))
	for _, s := range spells {
		sp := spell.ParseID(string(s.Spell))
		if !sp.Valid() {
			continue
		}
		out = append(out, sp)
	}
	obj.Object.SetTrapSpells(out...)
}

func (g nsObjGroup) AwardSpell(sp nsp.Spell) {
	spl := spell.ParseID(string(sp))
	if !spl.Valid() {
		return
	}
	g.EachObject(true, func(it ns4.Obj) bool {
		// TODO: why it.AwardSpell is different?
		val := 0
		obj := server.ToObject(it.(server.Obj))
		if noxflags.HasGame(noxflags.GameModeCoop) && it.Class().Has(object.ClassPlayer) && obj.UpdateDataPlayer().Player.SpellLvl[spl] == 0 {
			val = 1
		}
		legacy.Nox_xxx_spellGrantToPlayer_4FB550(obj, spl, 1, val, 0)
		return true
	})
}
