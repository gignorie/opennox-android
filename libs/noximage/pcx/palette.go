package pcx

import (
	"errors"
	"image/color"
	"io"
)

func DecodePalette(r io.Reader) (color.Palette, error) {
	var buf [7]byte
	_, err := io.ReadFull(r, buf[:7])
	if err != nil {
		return nil, err
	}
	if string(buf[:7]) != "PALETTE" {
		return nil, errors.New("not a palette file")
	}
	pal := make(color.Palette, 256)
	for i := range pal {
		_, err = io.ReadFull(r, buf[:3])
		if err != nil {
			return nil, err
		}
		pal[i] = color.NRGBA{R: buf[0], G: buf[1], B: buf[2], A: 0xff}
	}
	return pal, nil
}
