package noxcolor

import "image/color"

var (
	_ Color16 = RGB555(0)
)

// ToRGB555Color converts color.Color to RGB555.
func ToRGB555Color(c color.Color) RGB555 {
	switch c := c.(type) {
	case RGB555:
		return c
	case RGBA5551:
		return RGB555(c & 0x7fff)
	case Color:
		cl := c.ColorNRGBA()
		return RGB555Color(cl.R, cl.G, cl.B)
	case color.NRGBA:
		return RGB555Color(c.R, c.G, c.B)
	case color.Gray:
		return RGB555Color(c.Y, c.Y, c.Y)
	case color.Gray16:
		v := byte(c.Y >> 8)
		return RGB555Color(v, v, v)
	}
	cl := nrgbaModel(c)
	return RGB555Color(cl.R, cl.G, cl.B)
}

// RGB555Color converts RGB color to RGB555.
func RGB555Color(r, g, b byte) RGB555 {
	return RGB555((uint16(b&0xf8) >> 3) | (uint16(g&0xf8) << 2) | (uint16(r&0xf8) << 7))
}

// RGB555 stores RGB color in 15 bits (555).
type RGB555 uint16

// Color16 implements Color16.
func (c RGB555) Color16() uint16 {
	return uint16(c)
}

// Color32 implements Color16.
func (c RGB555) Color32() uint32 {
	v := uint32(c)
	return v | v<<16
}

// ColorNRGBA implements Color.
func (c RGB555) ColorNRGBA() (v color.NRGBA) {
	v.R = byte((c >> 7) & 0xf8)
	v.G = byte((c >> 2) & 0xf8)
	v.B = byte((c << 3) & 0xf8)
	v.A = 0xff
	return
}

// ColorRGBA implements Color.
func (c RGB555) ColorRGBA() (v color.RGBA) {
	v.R = byte((c >> 7) & 0xf8)
	v.G = byte((c >> 2) & 0xf8)
	v.B = byte((c << 3) & 0xf8)
	v.A = 0xff
	return
}

// RGBA implements color.Color.
func (c RGB555) RGBA() (r, g, b, a uint32) {
	return c.ColorNRGBA().RGBA()
}
