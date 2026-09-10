package legacy

/*
#include "defs.h"
#include "GAME3_3.h"
#include "GAME4_3.h"
int nox_xxx_pickupGold_4F3A60_obj_pickup(int a1, int a2, int a3);
int nox_objectPickupAudEvent_4F3D50(nox_object_t* a1, nox_object_t* a2, int a3);
*/
import "C"
import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/player"

	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_xxx_pickupDefault_4F31E0         server.PickupFunc
	Nox_objectPickupAudEvent_4F3D50      server.PickupFunc
	Nox_xxx_pickupPotion_4F37D0          server.PickupFunc
	Nox_xxx_playerClassCanUseItem_57B3D0 func(item *server.Object, cl player.Class) bool
	Sub_57B370                           func(cl object.Class, sub object.SubClass, typ int) byte
)

func init() {
	server.RegisterObjectPickup("DefaultPickup", C.nox_xxx_pickupDefault_4F31E0, func(who, it *server.Object, a3, a4 int) bool {
		return Nox_xxx_pickupDefault_4F31E0(who, it, a3, a4)
	})
	server.RegisterObjectPickupC("FoodPickup", C.nox_xxx_pickupFood_4F3350)
	server.RegisterObjectPickupC("UsePickup", C.nox_xxx_pickupUse_4F34D0)
	server.RegisterObjectPickupC("ArmorPickup", C.nox_xxx_pickupArmor_53E7F0)
	server.RegisterObjectPickupC("WeaponPickup", C.sub_53A720)
	server.RegisterObjectPickupC("OblivionPickup", C.nox_xxx_sendMsgOblivionPickup_53A9C0)
	server.RegisterObjectPickupC("TreasurePickup", C.nox_xxx_pickupTreasure_4F3580)
	server.RegisterObjectPickupC("TrapPickup", C.nox_xxx_pickupTrap_4F3510)
	server.RegisterObjectPickup("PotionPickup", C.nox_xxx_pickupPotion_4F37D0, func(who, it *server.Object, a3, a4 int) bool {
		return Nox_xxx_pickupPotion_4F37D0(who, it, a3, a4)
	})
	server.RegisterObjectPickupC("GoldPickup", C.nox_xxx_pickupGold_4F3A60_obj_pickup)
	server.RegisterObjectPickupC("AmmoPickup", C.nox_xxx_pickupAmmo_4F3B00)
	server.RegisterObjectPickupC("SpellBookPickup", C.nox_xxx_pickupSpellbook_4F3C60)
	server.RegisterObjectPickupC("AbilityBookPickup", C.nox_xxx_pickupAbilitybook_4F3CE0)
	server.RegisterObjectPickupC("CrownPickup", C.sub_4F3400)
	server.RegisterObjectPickup("AudEventPickup", C.nox_objectPickupAudEvent_4F3D50, func(who, it *server.Object, a3, a4 int) bool {
		return Nox_objectPickupAudEvent_4F3D50(who, it, a3, a4)
	})
	server.RegisterObjectPickupC("AnkhTradablePickup", C.sub_4F3DD0)
}

//export nox_xxx_pickupDefault_4F31E0
func nox_xxx_pickupDefault_4F31E0(cobj1 *nox_object_t, cobj2 *nox_object_t, a3 int) int {
	return bool2int(Nox_xxx_pickupDefault_4F31E0(asObjectS(cobj1), asObjectS(cobj2), a3, 0))
}

//export nox_objectPickupAudEvent_4F3D50
func nox_objectPickupAudEvent_4F3D50(cobj1 *nox_object_t, cobj2 *nox_object_t, a3 int) int {
	return bool2int(Nox_objectPickupAudEvent_4F3D50(asObjectS(cobj1), asObjectS(cobj2), a3, 0))
}

//export nox_xxx_pickupPotion_4F37D0
func nox_xxx_pickupPotion_4F37D0(cobj1 *nox_object_t, cobj2 *nox_object_t, a3 int) int {
	return bool2int(Nox_xxx_pickupPotion_4F37D0(asObjectS(cobj1), asObjectS(cobj2), a3, 0))
}

//export sub_57B370
func sub_57B370(cl, sub, typ int32) byte {
	return Sub_57B370(object.Class(cl), object.SubClass(sub), int(typ))
}

//export sub_419E10
func sub_419E10(u *nox_object_t, a2 int32) {
	GetServer().S().Players.SetXxx(asObjectS(u), a2)
}

//export sub_419E60
func sub_419E60(u *nox_object_t) int {
	return bool2int(GetServer().S().Players.CheckXxx(asObjectS(u)))
}

//export sub_419EA0
func sub_419EA0() int {
	return bool2int(GetServer().S().Players.AnyXxx())
}

//export nox_xxx_playerClassCanUseItem_57B3D0
func nox_xxx_playerClassCanUseItem_57B3D0(item *nox_object_t, cl int8) int {
	return bool2int(Nox_xxx_playerClassCanUseItem_57B3D0(asObjectS(item), player.Class(cl)))
}
