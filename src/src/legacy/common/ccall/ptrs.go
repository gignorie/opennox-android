package ccall

import (
	"unsafe"
)

func NewFuncs[
	F any,
](wrap func(cfnc unsafe.Pointer) F) *Funcs[F] {
	return &Funcs[F]{
		wrap:  wrap,
		byPtr: make(map[unsafe.Pointer]F),
	}
}

type Funcs[F any] struct {
	wrap  func(p unsafe.Pointer) F
	byPtr map[unsafe.Pointer]F
}

func (r *Funcs[F]) Register(cfnc unsafe.Pointer, fnc F) {
	r.byPtr[cfnc] = fnc
}

func (r *Funcs[F]) Get(cfnc unsafe.Pointer) F {
	if cfnc == nil {
		var zero F
		return zero
	}
	fnc, ok := r.byPtr[cfnc]
	if ok {
		return fnc
	}
	return r.wrap(cfnc)
}
