package opennox

import (
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func nox_objectDropAudEvent_4EE2F0(obj1 *server.Object, obj2 *server.Object, a3 types.Pointf) bool {
	s := noxServer
	if obj1 == nil || obj2 == nil {
		return false
	}
	ok := legacy.Nox_xxx_dropDefault_4ED290(obj1, obj2, &a3) != 0
	if !ok {
		return false
	}
	if snd := s.DropSound(obj2.TypeInd); snd != 0 {
		s.Audio.EventObj(snd, obj1, 0, 0)
	}
	return true
}
