package noxcolor

import "image/color"

var (
	_ Color16 = RGBA5551(0)
)

const (
	TransparentRGBA5551   = RGBA5551(0x8000)
	Transparent32RGBA5551 = 0x80000000
)

// ToRGBA5551Color converts color.Color to RGBA5551.
func ToRGBA5551Color(c color.Color) RGBA5551 {
	switch c := c.(type) {
	case RGBA5551:
		return c
	case RGB555:
		return RGBA5551(c & 0x7fff) // alpha channel is inverted (0 = opaque)
	case Color:
		cl := c.ColorNRGBA()
		return RGBA5551Color(cl.R, cl.G, cl.B, cl.A)
	case color.NRGBA:
		return RGBA5551Color(c.R, c.G, c.B, c.A)
	case color.Gray:
		return RGBA5551Color(c.Y, c.Y, c.Y, 0xff)
	case color.Gray16:
		v := byte(c.Y >> 8)
		return RGBA5551Color(v, v, v, 0xff)
	case color.Alpha:
		if c.A == 0 {
			return RGBA5551Color(0, 0, 0, c.A)
		}
		return RGBA5551Color(0xff, 0xff, 0xff, c.A)
	case color.Alpha16:
		if c.A == 0 {
			return RGBA5551Color(0, 0, 0, byte(c.A>>8))
		}
		return RGBA5551Color(0xff, 0xff, 0xff, byte(c.A>>8))
	}
	cl := nrgbaModel(c)
	return RGBA5551Color(cl.R, cl.G, cl.B, cl.A)
}

// RGBA5551Color converts RGBA color to RGBA5551.
func RGBA5551Color(r, g, b, a byte) RGBA5551 {
	if a >= 128 {
		return RGBA5551((uint16(b&0xf8) >> 3) | (uint16(g&0xf8) << 2) | (uint16(r&0xf8) << 7))
	}
	return RGBA5551((uint16(b&0xf8)>>3)|(uint16(g&0xf8)<<2)|(uint16(r&0xf8)<<7)) | TransparentRGBA5551
}

// RGB5551Color converts RGB color to RGBA5551.
func RGB5551Color(r, g, b byte) RGBA5551 {
	return RGBA5551((uint16(b&0xf8) >> 3) | (uint16(g&0xf8) << 2) | (uint16(r&0xf8) << 7))
}

// RGBA5551 stores RGBA color in 16 bits (5551).
// Alpha channel is inverted to be compatible with RGB555.
type RGBA5551 uint16

// Color16 implements Color16.
func (c RGBA5551) Color16() uint16 {
	return uint16(c)
}

// Color32 implements Color16.
func (c RGBA5551) Color32() uint32 {
	if c == TransparentRGBA5551 {
		return Transparent32RGBA5551
	}
	v := uint32(c)
	return v | v<<16
}

// ColorNRGBA implements Color.
func (c RGBA5551) ColorNRGBA() (v color.NRGBA) {
	v.R = byte((c >> 7) & 0xf8)
	v.G = byte((c >> 2) & 0xf8)
	v.B = byte((c << 3) & 0xf8)
	v.A = 0xff
	if v.R == 0xf8 {
		v.R = 0xff
	}
	if v.G == 0xf8 {
		v.G = 0xff
	}
	if v.B == 0xf8 {
		v.B = 0xff
	}
	if c>>15 != 0 {
		v.A = 0
	}
	return
}

// ColorRGBA implements Color.
func (c RGBA5551) ColorRGBA() (v color.RGBA) {
	if c>>15 != 0 {
		return color.RGBA{}
	}
	v.R = byte((c >> 7) & 0xf8)
	v.G = byte((c >> 2) & 0xf8)
	v.B = byte((c << 3) & 0xf8)
	if v.R == 0xf8 {
		v.R = 0xff
	}
	if v.G == 0xf8 {
		v.G = 0xff
	}
	if v.B == 0xf8 {
		v.B = 0xff
	}
	v.A = 0xff
	return
}

// RGBA implements color.Color.
func (c RGBA5551) RGBA() (r, g, b, a uint32) {
	return c.ColorNRGBA().RGBA()
}
