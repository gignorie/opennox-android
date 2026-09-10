package noxcolor

import "image/color"

var (
	_ Color16 = RGB565(0)
)

// ToRGB565Color converts color.Color to RGB565.
func ToRGB565Color(c color.Color) RGB565 {
	switch c := c.(type) {
	case RGB565:
		return c
	case Color:
		cl := c.ColorNRGBA()
		return RGB565Color(cl.R, cl.G, cl.B)
	case color.NRGBA:
		return RGB565Color(c.R, c.G, c.B)
	case color.Gray:
		return RGB565Color(c.Y, c.Y, c.Y)
	case color.Gray16:
		v := byte(c.Y >> 8)
		return RGB565Color(v, v, v)
	}
	cl := nrgbaModel(c)
	return RGB565Color(cl.R, cl.G, cl.B)
}

// RGB565Color converts RGBA color to RGB565.
func RGB565Color(r, g, b byte) RGB565 {
	return RGB565((uint16(b&0xf8) >> 3) | (uint16(g&0xfc) << 3) | (uint16(r&0xf8) << 8))
}

// RGB565 stores RGB color in 16 bits (565).
type RGB565 uint16

// Color16 implements Color16.
func (c RGB565) Color16() uint16 {
	return uint16(c)
}

// Color32 implements Color16.
func (c RGB565) Color32() uint32 {
	v := uint32(c)
	return v | v<<16
}

// ColorNRGBA implements Color.
func (c RGB565) ColorNRGBA() (v color.NRGBA) {
	v.R = byte((c >> 8) & 0xf8)
	v.G = byte((c >> 3) & 0xfc)
	v.B = byte((c << 3) & 0xf8)
	v.A = 0xff
	return
}

// ColorRGBA implements Color.
func (c RGB565) ColorRGBA() (v color.RGBA) {
	v.R = byte((c >> 8) & 0xf8)
	v.G = byte((c >> 3) & 0xfc)
	v.B = byte((c << 3) & 0xf8)
	v.A = 0xff
	return
}

// RGBA implements color.Color.
func (c RGB565) RGBA() (r, g, b, a uint32) {
	return c.ColorNRGBA().RGBA()
}
