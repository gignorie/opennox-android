package noxcolor

import "image/color"

var _ color.Model = Model(0)

const (
	// ModelRGBA5551 stores RGBA color in 16 bits (5551).
	// Alpha channel is inverted to be compatible with ModelRGB555.
	ModelRGBA5551 = Model(iota)
	// ModelRGB555 stores RGB color in 15 bits (555).
	ModelRGB555
	// ModelRGB565 stores RGB color in 16 bits (565).
	ModelRGB565
	// ModelRGBA4444 stores RGBA color in 16 bits (4444).
	ModelRGBA4444
)

type Model int

func (m Model) Convert(c color.Color) color.Color {
	return m.Convert16(c)
}

func (m Model) Convert16(c color.Color) Color16 {
	switch m {
	case ModelRGBA5551:
		return ToRGBA5551Color(c)
	case ModelRGB555:
		return ToRGB555Color(c)
	case ModelRGB565:
		return ToRGB565Color(c)
	case ModelRGBA4444:
		return ToRGBA4444Color(c)
	default:
		panic("unsupported model")
	}
}

// RGB creates a Color16, according to this color model.
func (m Model) RGB(r, g, b byte) Color16 {
	return m.NRGBA(r, g, b, 0xff)
}

// NRGBA creates a Color16, according to this color model.
func (m Model) NRGBA(r, g, b, a byte) Color16 {
	switch m {
	case ModelRGBA5551:
		return RGBA5551Color(r, g, b, a)
	case ModelRGB555:
		return RGB555Color(r, g, b)
	case ModelRGB565:
		return RGB565Color(r, g, b)
	case ModelRGBA4444:
		return RGBA4444Color(r, g, b, a)
	default:
		panic("unsupported model")
	}
}

// FromUint32 unpacks color from uint32, according to this color model.
func (m Model) FromUint32(v uint32) Color16 {
	c := uint16(v >> 16)
	switch m {
	case ModelRGBA5551:
		return RGBA5551(c)
	case ModelRGB555:
		return RGB555(c)
	case ModelRGB565:
		return RGB565(c)
	case ModelRGBA4444:
		return RGBA4444(c)
	default:
		panic("unsupported model")
	}
}

func nrgbaModel(c color.Color) color.NRGBA {
	if c, ok := c.(color.NRGBA); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	if a == 0xffff {
		return color.NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 0xff}
	}
	if a == 0 {
		return color.NRGBA{0, 0, 0, 0}
	}
	r = (r * 0xffff) / a
	g = (g * 0xffff) / a
	b = (b * 0xffff) / a
	return color.NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}
