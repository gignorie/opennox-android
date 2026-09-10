package opennox

import (
	"github.com/opennox/libs/spell"

	"github.com/opennox/opennox/v1/server"
)

func castInvulnerability(spellID spell.ID, _, _, _ *server.Object, args *server.SpellAcceptArg, lvl int) int {
	return castBuffSpell(spellID, server.ENCHANT_INVULNERABLE, lvl, args.Obj, spellBuffConf{
		DurOpt: "InvulnerabilityEnchantDuration", Defensive: true,
	})
}
