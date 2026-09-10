//go:build windows
// +build windows

package ifs

func Normalize(path string) string {
	return path
}

func Denormalize(path string) string {
	return path
}
