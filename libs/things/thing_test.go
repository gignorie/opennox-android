package things

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestThingFixAttrs(t *testing.T) {
	arr := fixThingAttrs("MASS = 6  DESTROY = DefaultDestroy")
	must.Eq(t, []string{"MASS = 6", "DESTROY = DefaultDestroy"}, arr)
}
