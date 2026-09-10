package stdlib

import "reflect"

var Symbols = make(map[string]map[string]reflect.Value)

// This block ensures that we compile the stdlib bindings.
// It may fail if the compiler version doesn't match the one used to generate bindings.
//
// To fix this, we update stdlib package files from the stdlib package in yaegi.
// We do not import stdlib from yaegi because it will compile ALL stdlib packages, even if we whitelist only a few.

var (
	_ _context_Context
	_ _fmt_Formatter
	_ _image_Image
	_ _image_color_Color
	_ _io_Closer
	_ _math_rand_Source
	_ _sort_Interface
)
