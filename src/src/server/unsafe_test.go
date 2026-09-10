package server

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestFuncAddr(t *testing.T) {
	fnc1 := func() int {
		return 1
	}
	fnc2 := func() int {
		return 2
	}
	must.NotEq(t, interfaceAddr(fnc1), interfaceAddr(fnc2))
	must.Eq(t, interfaceAddr(fnc1), interfaceAddr(fnc1))
	must.Eq(t, interfaceAddr(fnc2), interfaceAddr(fnc2))
}
