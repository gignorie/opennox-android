//go:build !windows
// +build !windows

package datapath

func registryStringKey(path, key string) (string, bool) {
	return "", false
}
