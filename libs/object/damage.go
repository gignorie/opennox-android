package object

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var DamageTypeNames = []string{
	"BLADE", "FLAME", "CRUSH",
	"IMPALE", "DRAIN", "POISON",
	"DISPEL_UNDEAD", "EXPLOSION", "BITE",
	"ELECTRIC", "CLAW", "IMPACT",
	"LAVA", "DEATH_MAGIC", "PLASMA",
	"MANA_BOMB", "ZAP_RAY", "AIRBORNE_ELECTRIC",
}

var goDamageTypeNames = []string{
	"DamageBlade",
	"DamageFlame",
	"DamageCrush",
	"DamageImpale",
	"DamageDrain",
	"DamagePoison",
	"DamageDispelUndead",
	"DamageExplosion",
	"DamageBite",
	"DamageElectric",
	"DamageClaw",
	"DamageImpact",
	"DamageLava",
	"DamageDeathMagic",
	"DamagePlasma",
	"DamageManaBomb",
	"DamageZapRay",
	"DamageAirborneElectric",
}

func ParseDamageType(name string) (DamageType, error) {
	s := strings.ToUpper(name)
	s = strings.TrimPrefix(s, "DAMAGE_")
	if s == "TRUE" {
		return DamageTrue, nil
	}
	for i, v := range DamageTypeNames {
		if v == s {
			return DamageType(i), nil
		}
	}
	return 0, fmt.Errorf("invalid damage name: %q", name)
}

type DamageType int32

const (
	DamageTrue = DamageType(iota - 1)
	DamageBlade
	DamageFlame
	DamageCrush
	DamageImpale
	DamageDrain
	DamagePoison
	DamageDispelUndead
	DamageExplosion
	DamageBite
	DamageElectric
	DamageClaw
	DamageImpact
	DamageLava
	DamageDeathMagic
	DamagePlasma
	DamageManaBomb
	DamageZapRay
	DamageAirborneElectric
)

func (v DamageType) String() string {
	if v == DamageTrue {
		return "DAMAGE_TRUE"
	}
	if int(v) < len(DamageTypeNames) {
		return "DAMAGE_" + DamageTypeNames[v]
	}
	return "DAMAGE_" + strconv.Itoa(int(v))
}

func (v DamageType) GoString() string {
	if v == DamageTrue {
		return "DamageTrue"
	}
	if int(v) < len(goDamageTypeNames) {
		return goDamageTypeNames[v]
	}
	return "Damage(" + strconv.Itoa(int(v)) + ")"
}

func (v DamageType) MarshalJSON() ([]byte, error) {
	if v == DamageTrue {
		return json.Marshal("DAMAGE_TRUE")
	}
	if int(v) < len(DamageTypeNames) {
		return json.Marshal("DAMAGE_" + DamageTypeNames[v])
	}
	return json.Marshal(int(v))
}
