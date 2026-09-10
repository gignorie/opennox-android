package legacy

/*
#include "GAME3_3.h"
#include "GAME4_3.h"
int nox_objectDropAudEvent_4EE2F0(nox_object_t* a1, nox_object_t* a2, float2* a3);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_objectDropAudEvent_4EE2F0 server.DropFunc
)

func init() {
	server.RegisterObjectDropC("DefaultDrop", C.nox_xxx_dropDefault_4ED290)
	server.RegisterObjectDropC("ArmorDrop", C.nox_xxx_dropArmor_53EB70)
	server.RegisterObjectDropC("WeaponDrop", C.nox_xxx_dropWeapon_53AB10)
	server.RegisterObjectDropC("TreasureDrop", C.nox_xxx_dropTreasure_4ED710)
	server.RegisterObjectDropC("GlyphDrop", C.nox_GlyphDrop_4ED500)
	server.RegisterObjectDropC("PotionDrop", C.sub_4EDDE0)
	server.RegisterObjectDropC("TrapDrop", C.nox_xxx_dropTrap_4ED580)
	server.RegisterObjectDropC("FoodDrop", C.nox_xxx_dropFood_4EDE50)
	server.RegisterObjectDropC("CrownDrop", C.nox_xxx_dropCrown_4ED5E0)
	server.RegisterObjectDrop("AudEventDrop", C.nox_objectDropAudEvent_4EE2F0, func(obj, obj2 *server.Object, pos types.Pointf) bool {
		return Nox_objectDropAudEvent_4EE2F0(obj, obj2, pos)
	})
	server.RegisterObjectDropC("AnkhTradableDrop", C.nox_xxx_dropAnkhTradable_4EE370)
}

//export nox_objectDropAudEvent_4EE2F0
func nox_objectDropAudEvent_4EE2F0(cobj1 *nox_object_t, cobj2 *nox_object_t, a3 *C.float2) int {
	pt := readUnalignedPointf(unsafe.Pointer(a3))
	return bool2int(Nox_objectDropAudEvent_4EE2F0(asObjectS(cobj1), asObjectS(cobj2), pt))
}

func Nox_xxx_dropDefault_4ED290(obj1 *server.Object, obj2 *server.Object, a3 *types.Pointf) int {
	return int(C.nox_xxx_dropDefault_4ED290(asObjectC(obj1), asObjectC(obj2), (*C.float2)(unsafe.Pointer(a3))))
}
