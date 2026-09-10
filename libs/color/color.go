package noxcolor

import "image/color"

// Color extends color.Color to also implement ColorNRGBA.
type Color interface {
	color.Color
	// ColorNRGBA returns a color.NRGBA for this color.
	ColorNRGBA() color.NRGBA
	// ColorRGBA returns a color.RGBA for this color.
	ColorRGBA() color.RGBA
}

// Color16 is an interface for colors that occupy 16 bits.
type Color16 interface {
	Color
	// Color16 returns color value as uint16.
	Color16() uint16
	// Color32 returns color value as uint32.
	Color32() uint32
}
