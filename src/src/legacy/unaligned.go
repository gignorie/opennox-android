package legacy

import (
	"unsafe"

	"github.com/opennox/libs/types"
)

// readUnalignedPointf copies 8 bytes from an unaligned pointer safely into a Go types.Pointf.
// Uses byte-level copies to avoid SIGBUS/BUS_ADRALN on ARM architectures (VFP vldr requires 4-byte alignment).
func readUnalignedPointf(p unsafe.Pointer) types.Pointf {
	if p == nil {
		return types.Pointf{}
	}
	var pt types.Pointf
	for i := 0; i < 8; i++ {
		*(*byte)(unsafe.Add(unsafe.Pointer(&pt), i)) = *(*byte)(unsafe.Add(p, i))
	}
	return pt
}

// readUnalignedPointfPtr returns an aligned pointer to a copy of the unaligned types.Pointf.
func readUnalignedPointfPtr(p unsafe.Pointer) *types.Pointf {
	if p == nil {
		return nil
	}
	pt := readUnalignedPointf(p)
	return &pt
}

// writeUnalignedPointf safely writes a types.Pointf into an unaligned destination buffer byte-by-byte.
func writeUnalignedPointf(dst unsafe.Pointer, pt types.Pointf) {
	if dst == nil {
		return
	}
	for i := 0; i < 8; i++ {
		*(*byte)(unsafe.Add(dst, i)) = *(*byte)(unsafe.Add(unsafe.Pointer(&pt), i))
	}
}

// readUnalignedRectf copies 16 bytes from an unaligned pointer safely into a Go types.Rectf.
func readUnalignedRectf(p unsafe.Pointer) types.Rectf {
	if p == nil {
		return types.Rectf{}
	}
	var r types.Rectf
	for i := 0; i < 16; i++ {
		*(*byte)(unsafe.Add(unsafe.Pointer(&r), i)) = *(*byte)(unsafe.Add(p, i))
	}
	return r
}

// readUnalignedFloat32 copies 4 bytes from an unaligned pointer safely into a Go float32.
func readUnalignedFloat32(p unsafe.Pointer) float32 {
	if p == nil {
		return 0
	}
	var f float32
	for i := 0; i < 4; i++ {
		*(*byte)(unsafe.Add(unsafe.Pointer(&f), i)) = *(*byte)(unsafe.Add(p, i))
	}
	return f
}
