package script

import (
	"log/slog"
	"reflect"
	"sort"

	"golang.org/x/exp/maps"
)

// EventType is a type of script events.
type EventType string

const (
	EventMapInitialize = EventType("MapInitialize")
	EventMapEntry      = EventType("MapEntry")
	EventMapExit       = EventType("MapExit")
	EventMapShutdown   = EventType("MapShutdown")
	EventPlayerDeath   = EventType("PlayerDeath")
)

var vms = make(map[string]VMRuntime)

// RegisterVM registers a new script VM runtime.
func RegisterVM(r VMRuntime) {
	if r.Name == "" {
		panic("name must be set")
	}
	if r.NewMap == nil {
		panic("new map function must be set")
	}
	if _, ok := vms[r.Name]; ok {
		panic("already registered")
	}
	vms[r.Name] = r
}

// VMRuntimes returns all registered VM runtimes.
func VMRuntimes() []VMRuntime {
	keys := maps.Keys(vms)
	sort.Strings(keys)
	var out []VMRuntime
	for _, k := range keys {
		out = append(out, vms[k])
	}
	return out
}

// VMRuntime is a type for registering new script runtime implementations.
type VMRuntime struct {
	// Name is a short name for the VM runtime. Must be unique.
	Name string
	// Title is a human-friendly name for the VM runtime.
	Title string
	// NewMap creates a new script VM for map scripts.
	// If there's no scripts for the map, function may return nil, nil.
	NewMap func(log *slog.Logger, g Game, maps string, name string) (VM, error)
}

// VM is an interface for a virtual machine running the script engine.
type VM interface {
	// Exec executes text as a script source code.
	Exec(s string) (reflect.Value, error)
	// ExecFile executes a script from the file or directory.
	ExecFile(path string) error
	// OnFrame must be called when a new game frame is complete.
	OnFrame()
	// OnEvent is called when a certain scrip event happens.
	OnEvent(typ EventType)
	// GetSymbol tries to find a function or variable with a given name and type.
	// If symbol is found, but type doesn't match, it returns an error.
	// If symbol is not found, it returns reflect.Value{}, false, nil.
	GetSymbol(name string, typ reflect.Type) (reflect.Value, bool, error)
	// Close the VM and free resources.
	Close() error
}

// GetVMSymbol is a generic helper for VM.GetSymbol.
// It is suitable for functions mostly, since it returns value directly, not a pointer.
func GetVMSymbol[T any](vm VM, name string) (T, error) {
	var zero T
	rv, ok, err := vm.GetSymbol(name, reflect.TypeOf(zero))
	if err != nil || !ok {
		return zero, err
	}
	return rv.Interface().(T), nil
}

// GetVMSymbolPtr is a generic helper for VM.GetSymbol.
// It is similar to GetVMSymbol, but returns a pointer to the value. Useful for variables.
func GetVMSymbolPtr[T any](vm VM, name string) (*T, error) {
	var zero T
	rv, ok, err := vm.GetSymbol(name, reflect.TypeOf(zero))
	if err != nil || !ok {
		return nil, err
	}
	return rv.Addr().Interface().(*T), nil
}
