package opennox

import (
	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/legacy"
)

func (s noxScriptNS) Unused5e(id string) int {
	str := s.s.Strings().GetStringInFile(strman.ID(id), "CScrFunc.c")
	return legacy.Sub_512E80(str)
}

func (s noxScriptNS) Unused74(arg1 int, arg2 int) {}

func (s noxScriptNS) Unknownb8(id int) bool {
	//TODO implement me
	panic("implement me")
}

func (s noxScriptNS) Unknownb9(id int) bool {
	//TODO implement me
	panic("implement me")
}

func (s noxScriptNS) Unknownc4() {
	noxClient.Nox_video_stopAllFades44E040()
}
