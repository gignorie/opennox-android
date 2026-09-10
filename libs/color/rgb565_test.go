package noxcolor

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestRGB565(t *testing.T) {
	for i := 0; i < 0xFFFF; i++ {
		c := RGB565(i)
		cl := c.ColorNRGBA()
		c2 := RGB565Color(cl.R, cl.G, cl.B)
		must.EqOp(t, c, c2, must.Sprintf("0x%x", int(c)))
	}
}
