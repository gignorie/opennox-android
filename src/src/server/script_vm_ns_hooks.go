package server

import (
	"github.com/opennox/libs/script"
	ns4 "github.com/opennox/noxscript/ns/v4"
)

func (s NoxScriptNS) OnFrame(fnc ns4.FrameFunc) {
	s.s.OnScriptFrame(fnc)
}
func (s NoxScriptNS) OnMapEvent(ev ns4.MapEvent, fnc ns4.MapEventFunc) {
	switch ev {
	case ns4.MapInitialize:
		s.s.OnMapEvent(script.EventMapInitialize, fnc)
	case ns4.MapEntry:
		s.s.OnMapEvent(script.EventMapEntry, fnc)
	case ns4.MapExit:
		s.s.OnMapEvent(script.EventMapExit, fnc)
	case ns4.MapShutdown:
		s.s.OnMapEvent(script.EventMapShutdown, fnc)
	}
}
