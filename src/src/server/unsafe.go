package server

import "unsafe"

// interfaceAddr extracts an address of an object stored in the interface value.
func interfaceAddr(v any) uintptr {
	ifc := *(*[2]uintptr)(unsafe.Pointer(&v))
	return ifc[1]
}
