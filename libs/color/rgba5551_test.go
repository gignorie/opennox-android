package noxcolor

import (
	"image/color"
	"testing"

	"github.com/shoenig/test/must"
)

func TestRGBA5551(t *testing.T) {
	for i := 0; i < 0xFFFF; i++ {
		c := RGBA5551(i)
		cl := c.ColorNRGBA()
		c2 := RGBA5551Color(cl.R, cl.G, cl.B, cl.A)
		must.EqOp(t, c, c2, must.Sprintf("0x%x", int(c)))
	}
}

func TestRGBA5551Builtin(t *testing.T) {
	m := ModelRGBA5551
	for _, c := range []struct {
		name  string
		exp   color.Color
		model color.Model
		exp16 Color16
	}{
		{name: "transparent", exp: color.Transparent, model: color.Alpha16Model, exp16: TransparentRGBA5551},
		{name: "opaque", exp: color.Opaque, model: color.Alpha16Model, exp16: RGB5551Color(0xff, 0xff, 0xff)},
		{name: "black", exp: color.Black, model: color.Gray16Model, exp16: RGB5551Color(0, 0, 0)},
		{name: "white", exp: color.White, model: color.Gray16Model, exp16: RGB5551Color(0xff, 0xff, 0xff)},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := m.Convert16(c.exp)
			must.EqOp(t, c.exp16, got)
			must.EqOp(t, c.exp, c.model.Convert(got))
		})
	}
}
