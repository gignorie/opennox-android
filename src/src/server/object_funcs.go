package server

import (
	"unsafe"

	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

var (
	createFuncs    = make(map[string]unsafe.Pointer)
	initFuncs      = make(map[string]objectDefFunc)
	initParseFuncs = map[string]ObjectParseFunc{
		"SkullInit":     objectDirectionInitParse,
		"DirectionInit": objectDirectionInitParse,
		"BreakInit": func(objt *ObjectType, _ []string) error {
			objt.Field9 = 2
			return nil
		},
	}
	updateFuncs       = make(map[string]objectDefFunc)
	updateParseFuncs  = make(map[string]ObjectParseFunc)
	collideFuncs      = make(map[string]objectDefFunc)
	collideParseFuncs = make(map[string]ObjectParseFunc)
	useFuncs          = make(map[string]objectDefFunc)
	useParseFuncs     = make(map[string]ObjectParseFunc)
	damageFuncs       = make(map[string]unsafe.Pointer)
	damageSoundFuncs  = make(map[string]unsafe.Pointer)
	deathFuncs        = make(map[string]objectDefFunc)
	deathParseFuncs   = make(map[string]ObjectParseFunc)
	dropFuncs         = make(map[string]unsafe.Pointer)
	dropParseFuncs    = map[string]ObjectParseFunc{
		"AudEventDrop": func(t *ObjectType, args []string) error {
			if len(args) != 0 {
				if snd := sound.ByName(args[0]); snd != 0 {
					t.s.dropSoundTable[t.ind] = snd
				}
			}
			return nil
		},
	}
	pickupFuncs      = make(map[string]unsafe.Pointer)
	pickupParseFuncs = map[string]ObjectParseFunc{
		"AudEventPickup": func(t *ObjectType, args []string) error {
			if len(args) != 0 {
				if snd := sound.ByName(args[0]); snd != 0 {
					t.s.pickupSoundTable[t.ind] = snd
				}
			}
			return nil
		},
	}
	xferFuncs = make(map[string]unsafe.Pointer)
)

var (
	DefaultDamage      unsafe.Pointer
	DefaultDamageSound unsafe.Pointer
	DefaultXfer        unsafe.Pointer
)

func init() {
	RegisterObjectCreate("NoCreate", nil)
	RegisterObjectCreate("PlayerCreate", nil)

	RegisterObjectInit("NoInit", nil, 0)
	RegisterObjectInit("ProjectileInit", nil, 0)
	RegisterObjectInit("GruntInit", nil, 0)
	RegisterObjectInit("SkeletonInit", nil, 0)
	RegisterObjectInit("TowerInit", nil, 0)
	RegisterObjectInit("RewardMarkerInit", nil, 220)
	RegisterObjectInit("AnkhInit", nil, 5124)

	RegisterObjectUpdate("NoUpdate", nil, 0)
	RegisterObjectUpdate("HomingProjectileUpdate", nil, 4)
	RegisterObjectUpdate("SpikeBlockUpdate", nil, 2200)
	RegisterObjectUpdate("TowerUpdate", nil, 8)
	RegisterObjectUpdate("WeaponArmorUpdate", nil, 8)
	RegisterObjectUpdate("DamageRoundoffUpdate", nil, 8) // used in demo instead of WeaponArmorUpdate

	RegisterObjectCollide("NoCollide", nil, 0)

	RegisterObjectUseC("AmmoUse", nil, 3)
	RegisterObjectUseC("BowUse", nil, 1)
}

type ObjectParseFunc func(objt *ObjectType, args []string) error

type objectDefFunc struct {
	Func     unsafe.Pointer
	DataSize uintptr
}

func RegisterObjectCreate(name string, fnc unsafe.Pointer) {
	if _, ok := createFuncs[name]; ok {
		panic("already registered")
	}
	createFuncs[name] = fnc
}

func RegisterObjectInit(name string, fnc unsafe.Pointer, sz uintptr) {
	if _, ok := initFuncs[name]; ok {
		panic("already registered")
	}
	initFuncs[name] = objectDefFunc{Func: fnc, DataSize: sz}
}

func RegisterObjectUpdate(name string, fnc unsafe.Pointer, sz uintptr) {
	if _, ok := updateFuncs[name]; ok {
		panic("already registered")
	}
	updateFuncs[name] = objectDefFunc{Func: fnc, DataSize: sz}
}

func RegisterObjectUpdateParse(name string, fnc ObjectParseFunc) {
	if _, ok := updateParseFuncs[name]; ok {
		panic("already registered")
	}
	updateParseFuncs[name] = fnc
}

func RegisterObjectCollide(name string, fnc unsafe.Pointer, sz uintptr) {
	if _, ok := collideFuncs[name]; ok {
		panic("already registered")
	}
	collideFuncs[name] = objectDefFunc{Func: fnc, DataSize: sz}
}

func RegisterObjectCollideParse(name string, fnc ObjectParseFunc) {
	if _, ok := collideParseFuncs[name]; ok {
		panic("already registered")
	}
	collideParseFuncs[name] = fnc
}

type UseFuncPtr struct {
	Ptr unsafe.Pointer
}

func (p UseFuncPtr) Get() UseFunc {
	if p.Ptr == nil {
		return nil
	}
	return objUse.Get(p.Ptr)
}

type UseFunc func(obj, obj2 *Object) bool

var objUse = ccall.NewFuncs(func(cfnc unsafe.Pointer) UseFunc {
	return func(obj, obj2 *Object) bool {
		return ccall.CallIntPtr2(cfnc, obj.CObj(), obj2.CObj()) != 0
	}
})

func RegisterObjectUseC(name string, cfnc unsafe.Pointer, sz uintptr) {
	if _, ok := useFuncs[name]; ok {
		panic("already registered")
	}
	useFuncs[name] = objectDefFunc{Func: cfnc, DataSize: sz}
}

func RegisterObjectUse(name string, cfnc unsafe.Pointer, fnc UseFunc, sz uintptr) {
	if _, ok := useFuncs[name]; ok {
		panic("already registered")
	}
	useFuncs[name] = objectDefFunc{Func: cfnc, DataSize: sz}
	objUse.Register(cfnc, fnc)
}

func RegisterObjectUseParse(name string, fnc ObjectParseFunc) {
	if _, ok := useParseFuncs[name]; ok {
		panic("already registered")
	}
	useParseFuncs[name] = fnc
}

func RegisterObjectDamage(name string, fnc unsafe.Pointer) {
	if _, ok := damageFuncs[name]; ok {
		panic("already registered")
	}
	damageFuncs[name] = fnc
}

func RegisterObjectDamageSound(name string, fnc unsafe.Pointer) {
	if _, ok := damageSoundFuncs[name]; ok {
		panic("already registered")
	}
	damageSoundFuncs[name] = fnc
}

func RegisterObjectDeath(name string, fnc unsafe.Pointer, sz uintptr) {
	if _, ok := deathFuncs[name]; ok {
		panic("already registered")
	}
	deathFuncs[name] = objectDefFunc{Func: fnc, DataSize: sz}
}

func RegisterObjectDeathParse(name string, fnc ObjectParseFunc) {
	if _, ok := deathParseFuncs[name]; ok {
		panic("already registered")
	}
	deathParseFuncs[name] = fnc
}

type DropFuncPtr struct {
	Ptr unsafe.Pointer
}

func (p DropFuncPtr) Get() DropFunc {
	if p.Ptr == nil {
		return nil
	}
	return objDrop.Get(p.Ptr)
}

type DropFunc func(obj, obj2 *Object, pos types.Pointf) bool

var objDrop = ccall.NewFuncs(func(cfnc unsafe.Pointer) DropFunc {
	return func(obj, obj2 *Object, pos types.Pointf) bool {
		cpos, free := alloc.New(types.Pointf{})
		defer free()
		*cpos = pos

		return ccall.CallIntPtr3(cfnc, obj.CObj(), obj2.CObj(), unsafe.Pointer(cpos)) != 0
	}
})

func RegisterObjectDropC(name string, cfnc unsafe.Pointer) {
	if _, ok := dropFuncs[name]; ok {
		panic("already registered")
	}
	dropFuncs[name] = cfnc
}

func RegisterObjectDrop(name string, cfnc unsafe.Pointer, fnc DropFunc) {
	if _, ok := dropFuncs[name]; ok {
		panic("already registered")
	}
	dropFuncs[name] = cfnc
	objDrop.Register(cfnc, fnc)
}

type PickupFuncPtr struct {
	Ptr unsafe.Pointer
}

func (p PickupFuncPtr) Get() PickupFunc {
	if p.Ptr == nil {
		return nil
	}
	return objPickup.Get(p.Ptr)
}

type PickupFunc func(who, it *Object, a3, a4 int) bool

var objPickup = ccall.NewFuncs(func(cfnc unsafe.Pointer) PickupFunc {
	return func(who, it *Object, a3, a4 int) bool {
		return ccall.CallIntUPtr4(cfnc, uintptr(who.CObj()), uintptr(it.CObj()), uintptr(a3), uintptr(a4)) != 0
	}
})

func RegisterObjectPickupC(name string, cfnc unsafe.Pointer) {
	if _, ok := pickupFuncs[name]; ok {
		panic("already registered")
	}
	pickupFuncs[name] = cfnc
}

func RegisterObjectPickup(name string, cfnc unsafe.Pointer, fnc PickupFunc) {
	if _, ok := pickupFuncs[name]; ok {
		panic("already registered")
	}
	pickupFuncs[name] = cfnc
	objPickup.Register(cfnc, fnc)
}

func RegisterObjectXfer(name string, fnc unsafe.Pointer) {
	if _, ok := xferFuncs[name]; ok {
		panic("already registered")
	}
	xferFuncs[name] = fnc
}
