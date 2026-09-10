//go:build !script

package script

// SetRuntime sets a global runtime instance.
func SetRuntime(g Game) {
	global = g
}
