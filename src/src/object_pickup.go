package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/player"
	"github.com/opennox/libs/spell"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

var (
	cheatAllowAll = false
)

func nox_xxx_inventoryServPlace_4F36F0(obj *server.Object, it *server.Object, a3 int, a4 int) bool {
	s := noxServer
	if obj == nil || it == nil {
		return false
	}
	if obj.CarryCapacity == 0 {
		return false
	}
	if it.Flags().Has(object.FlagDestroyed) {
		return false
	}
	if obj.Flags().Has(object.FlagDead) {
		return false
	}
	if !s.Types.ByInd(int(it.TypeInd)).Allowed() {
		return false
	}
	if !obj.Class().HasAny(object.MaskUnits) {
		return false
	}
	if !it.CallPickup(obj, a3, a4) {
		return false
	}
	if it.Flags().Has(object.FlagNoCollide) {
		it.ObjFlags &^= object.FlagNoCollide
		if it.Collide != nil {
			legacy.Sub_5117F0(it)
		}
	}
	if it.ScriptPickup.Func != -1 {
		s.noxScript.ScriptCallback(&it.ScriptPickup, obj, it, server.NoxEventInventoryPlace)
		it.ScriptPickup.Func = -1
	}
	return true
}

func nox_xxx_pickupDefault_4F31E0(obj, item *server.Object, a3, a4 int) bool {
	s := noxServer
	if !noxflags.HasGame(noxflags.GameModeQuest) && item.TeamPtr().Has() && !obj.TeamPtr().SameAs(item.TeamPtr()) {
		if tm := s.Teams.ByID(item.TeamVal.ID); tm != nil {
			if obj.Class().Has(object.ClassPlayer) {
				ud := obj.UpdateDataPlayer()
				s.NetInformTextMsg(ud.Player.PlayerIndex(), 16, int(tm.ColorInd))
			}
			return false
		}
	}
	if item.InvHolder != nil {
		return false
	}
	if obj.CarryCapacity == 0 {
		return false
	}
	weight := 0
	for it := obj.InvFirstItem; it != nil; it = it.InvNextItem {
		weight += int(it.Weight)
	}
	if int(item.Weight) > int(obj.CarryCapacity)*2-weight {
		s.NetPriMsgToPlayer(obj, "pickup.c:CarryingTooMuch", 0)
		s.Log.Info("can't pickup; carrying too much", "obj", obj.String(), "item", item.String())
		return false
	}
	if item.Class().Has(object.ClassFood) {
		cnt := obj.CountInventoryWithType(int(item.TypeInd))
		maxCnt := 3
		if noxflags.HasGame(noxflags.GameModeQuest | noxflags.GameModeCoop) {
			maxCnt = 9
		}
		if cnt >= maxCnt {
			s.NetPriMsgToPlayer(obj, "pickup.c:MaxSameItem", 0)
			s.Log.Info("can't pickup; carrying too many items of this kind", "obj", obj.String(), "item", item.String())
			return false
		}
	}
	s.ObjectDeleteLast(item)
	legacy.Nox_xxx_inventoryPutImpl_4F3070(obj, item, a3)
	return true
}

func nox_objectPickupAudEvent_4F3D50(obj1 *server.Object, obj2 *server.Object, a3, a4 int) bool {
	s := noxServer
	if obj1 == nil || obj2 == nil {
		return false
	}
	if !nox_xxx_pickupDefault_4F31E0(obj1, obj2, a3, a4) {
		return false
	}
	if snd := s.PickupSound(obj2.TypeInd); snd != 0 {
		s.Audio.EventObj(snd, obj1, 0, 0)
	}
	return true
}

func sub_57B370(cl object.Class, sub object.SubClass, typ int) byte {
	s := noxServer
	if cl.HasAny(object.ClassWeapon | object.ClassWand) {
		m := s.Modif.Nox_xxx_getProjectileClassById413250(typ)
		if m == nil {
			return 0
		}
		return m.Classes62
	}
	if cl.Has(object.ClassArmor) {
		m := s.Modif.Nox_xxx_equipClothFindDefByTT413270(typ)
		if m == nil {
			return 0
		}
		return m.Classes62
	}
	if cl.Has(object.ClassFood) {
		return byte(^(uint32(sub) >> 5) | 0xFE)
	}
	return 0xFF
}

func nox_xxx_playerClassCanUseItem_57B3D0(item *server.Object, cl player.Class) bool {
	if cheatAllowAll {
		return true
	}
	return ((byte(1) << cl) & sub_57B370(item.Class(), item.SubClass(), int(item.TypeInd))) != 0
}

func nox_xxx_pickupPotion_4F37D0(obj *server.Object, potion *server.Object, a3, a4 int) bool {
	s := noxServer
	if noxflags.HasGame(0x2000) && !noxflags.HasGame(4096) && obj.Class().Has(object.ClassPlayer) && !nox_xxx_playerClassCanUseItem_57B3D0(potion, obj.UpdateDataPlayer().Player.PlayerClass()) {
		s.NetPriMsgToPlayer(obj, "pickup.c:ObjectEquipClassFail", 0)
		s.Audio.EventObj(sound.SoundNoCanDo, obj, 2, obj.NetCode)
		return false
	}
	if !s.Players.CheckXxx(obj) {
		use := potion.UseDataPotion()
		consumed := false
		if use != nil && potion.SubClass().AsFood().Has(object.FoodHealthPotion) && obj.HealthData != nil {
			dhp := int(use.Value)
			if obj.Class().Has(object.ClassPlayer) {
				ud := obj.UpdateDataPlayer()
				if mult := s.Players.ClassStatsMult(ud.Player.PlayerClass()); mult != nil {
					dhp = int(float64(dhp) * float64(mult.Health))
				}
			}
			if dhp+int(obj.HealthData.Cur) < int(obj.HealthData.Max) {
				legacy.Nox_xxx_unitAdjustHP_4EE460(obj, dhp)
				s.Audio.EventObj(sound.SoundRestoreHealth, obj, 0, 0)
				consumed = true
			}
		}
		if use != nil && potion.SubClass().AsFood().Has(object.FoodManaPotion) && obj.Class().Has(object.ClassPlayer) {
			ud := obj.UpdateDataPlayer()
			dmp := int(use.Value)
			if mult := s.Players.ClassStatsMult(ud.Player.PlayerClass()); mult != nil {
				dmp = int(float64(dmp) * float64(mult.Mana))
			}
			if dmp+int(ud.ManaCur) < int(ud.ManaMax) {
				legacy.Nox_xxx_playerManaAdd_4EEB80(obj, dmp)
				s.Audio.EventObj(sound.SoundRestoreMana, obj, 0, 0)
				consumed = true
			}
		}
		if potion.SubClass().AsFood().Has(object.FoodCurePoisonPotion) && obj.Class().Has(object.ClassPlayer) && int32(obj.Poison540) != 0 {
			legacy.Nox_xxx_removePoison_4EE9D0(obj)
			aud := s.Spells.DefByInd(spell.SPELL_CURE_POISON).GetOnSound()
			s.Audio.EventObj(aud, obj, 0, 0)
			s.DelayedDelete(potion)
			return true
		}
		if consumed {
			s.DelayedDelete(potion)
			return true
		}
	}
	legacy.Nox_xxx_decay_5116F0(potion)
	ok := nox_xxx_pickupDefault_4F31E0(obj, potion, a3, a4)
	if ok {
		s.Audio.EventObj(sound.SoundPotionPickup, obj, 0, 0)
	}
	return ok
}
