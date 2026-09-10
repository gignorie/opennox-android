package opennox

import (
	"github.com/opennox/libs/client/keybind"

	"github.com/opennox/opennox/v1/common/memmap"
)

var (
	keyBinding *keybind.Binding
)

var noxMouseSelectOpt = []string{
	"Left",
	"Right",
	"Middle",
	"Wheel",
}

func inputInitMouse() {
	// indicates that mouse is present so cursor should be drawn
	*memmap.PtrUint32(0x5D4594, 1193108) = 1
}
