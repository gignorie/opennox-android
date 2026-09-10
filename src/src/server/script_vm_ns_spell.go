package server

import (
	"strings"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/types"
	ns4 "github.com/opennox/noxscript/ns/v4"
	"github.com/opennox/noxscript/ns/v4/effect"
)

var nsFxNames = make(map[string]netmsg.Op)

func init() {
	for fx := netmsg.MSG_FX_PARTICLEFX; fx <= netmsg.MSG_FX_MANA_BOMB_CANCEL; fx++ {
		nsFxNames[fx.String()] = fx
	}
}

func (s NoxScriptNS) Effect(effect effect.Effect, p1, p2 ns4.Positioner) {
	var pos, pos2 types.Pointf
	if p1 != nil {
		pos = p1.Pos()
	}
	if p2 != nil {
		pos2 = p2.Pos()
	}
	name := "MSG_FX_" + strings.ToUpper(string(effect))
	switch fx := nsFxNames[name]; fx {
	case netmsg.MSG_FX_BLUE_SPARKS,
		netmsg.MSG_FX_YELLOW_SPARKS,
		netmsg.MSG_FX_CYAN_SPARKS,
		netmsg.MSG_FX_VIOLET_SPARKS,
		netmsg.MSG_FX_EXPLOSION,
		netmsg.MSG_FX_LESSER_EXPLOSION,
		netmsg.MSG_FX_COUNTERSPELL_EXPLOSION,
		netmsg.MSG_FX_THIN_EXPLOSION,
		netmsg.MSG_FX_TELEPORT,
		netmsg.MSG_FX_SMOKE_BLAST,
		netmsg.MSG_FX_DAMAGE_POOF,
		netmsg.MSG_FX_RICOCHET,
		netmsg.MSG_FX_WHITE_FLASH,
		netmsg.MSG_FX_TURN_UNDEAD,
		netmsg.MSG_FX_MANA_BOMB_CANCEL:
		s.s.Nox_xxx_netSendPointFx_522FF0(fx, pos)
	case netmsg.MSG_FX_LIGHTNING,
		netmsg.MSG_FX_ENERGY_BOLT,
		netmsg.MSG_FX_PLASMA,
		netmsg.MSG_FX_DRAIN_MANA,
		netmsg.MSG_FX_CHARM,
		netmsg.MSG_FX_GREATER_HEAL,
		netmsg.MSG_FX_DEATH_RAY,
		netmsg.MSG_FX_SENTRY_RAY:
		s.s.Nox_xxx_netSendRayFx_5232F0(fx, pos.Point(), pos2.Point())
	case netmsg.MSG_FX_SPARK_EXPLOSION:
		s.s.Nox_xxx_netSparkExplosionFx_5231B0(pos, byte(pos2.X))
	case netmsg.MSG_FX_JIGGLE:
		s.s.Nox_xxx_earthquakeSend_4D9110(pos, int(pos2.X))
	case netmsg.MSG_FX_GREEN_BOLT:
		s.s.Nox_xxx_netSendFxGreenBolt_523790(pos.Point(), pos2.Point(), 30)
	case netmsg.MSG_FX_VAMPIRISM:
		s.s.Nox_xxx_netSendVampFx_523270(fx, pos.Point(), pos2.Point(), 100)
	}
}
