package noxcolor

import "image/color"

var (
	_ Color16 = RGBA4444(0)
)

// ToRGBA4444Color converts color.Color to RGBA4444.
func ToRGBA4444Color(c color.Color) RGBA4444 {
	switch c := c.(type) {
	case RGBA4444:
		return c
	case Color:
		cl := c.ColorNRGBA()
		return RGBA4444Color(cl.R, cl.G, cl.B, cl.A)
	case color.NRGBA:
		return RGBA4444Color(c.R, c.G, c.B, c.A)
	case color.Gray:
		return RGBA4444Color(c.Y, c.Y, c.Y, 0xff)
	case color.Gray16:
		v := byte(c.Y >> 8)
		return RGBA4444Color(v, v, v, 0xff)
	case color.Alpha:
		if c.A == 0 {
			return RGBA4444Color(0, 0, 0, c.A)
		}
		return RGBA4444Color(0xff, 0xff, 0xff, c.A)
	case color.Alpha16:
		if c.A == 0 {
			return RGBA4444Color(0, 0, 0, byte(c.A>>8))
		}
		return RGBA4444Color(0xff, 0xff, 0xff, byte(c.A>>8))
	}
	cl := nrgbaModel(c)
	return RGBA4444Color(cl.R, cl.G, cl.B, cl.A)
}

// RGBA4444Color converts RGBA color to RGBA4444.
func RGBA4444Color(r, g, b, a byte) RGBA4444 {
	return RGBA4444((uint16(a&0xf0) >> 4) | (uint16(b&0xf0) << 0) | (uint16(g&0xf0) << 4) | (uint16(r&0xf0) << 8))
}

// RGB4444Color converts RGB color to RGBA4444.
func RGB4444Color(r, g, b byte) RGBA4444 {
	return RGBA4444Color(r, g, b, 0xff)
}

// RGBA4444 stores RGBA color in 16 bits (4444).
type RGBA4444 uint16

// Color16 implements Color16.
func (c RGBA4444) Color16() uint16 {
	return uint16(c)
}

// Color32 implements Color16.
func (c RGBA4444) Color32() uint32 {
	v := uint32(c)
	return v | v<<16
}

// ColorNRGBA implements Color.
func (c RGBA4444) ColorNRGBA() (v color.NRGBA) {
	v.R = byte((c >> 8) & 0xf0)
	v.G = byte((c >> 4) & 0xf0)
	v.B = byte((c >> 0) & 0xf0)
	v.A = byte((c << 4) & 0xf0)
	return
}

// ColorRGBA implements Color.
func (c RGBA4444) ColorRGBA() (v color.RGBA) {
	r, g, b, a := c.ColorNRGBA().RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// RGBA implements color.Color.
func (c RGBA4444) RGBA() (r, g, b, a uint32) {
	return c.ColorNRGBA().RGBA()
}
