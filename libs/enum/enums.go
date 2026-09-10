package enum

import (
	"encoding/json"
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

type number interface {
	~uint64 | ~uint32 | ~uint16 | ~uint8
}

type Enum[T number] interface {
	json.Marshaler
	Has(c2 T) bool
	HasAny(c2 T) bool
	Split() []T
	String() string
}

func bitSize[T number]() int {
	var zero T
	return bits.OnesCount64(uint64(^zero))
}

func SplitBits[T number](v T) (out []T) {
	for i := 0; i < bitSize[T](); i++ {
		v2 := T(1) << i
		if v&v2 != 0 {
			out = append(out, v2)
		}
	}
	return
}

func StringBits[T number](v T, names []string) string {
	var out []string
	for i := 0; i < bitSize[T](); i++ {
		v2 := T(1) << i
		if v&v2 != 0 {
			if i < len(names) {
				out = append(out, names[i])
			} else {
				out = append(out, "0x"+strconv.FormatUint(uint64(v2), 16))
			}
		}
	}
	return strings.Join(out, " | ")
}

func StringBitsRaw[T number](v T) string {
	var out []string
	for i := 0; i < bitSize[T](); i++ {
		v2 := T(1) << i
		if v&v2 != 0 {
			out = append(out, "0x"+strconv.FormatUint(uint64(v2), 16))
		}
	}
	return strings.Join(out, " | ")
}

func Parse[T number](ename string, s string, names []string) (T, error) {
	s = strings.ToUpper(s)
	if s == "" || s == "NULL" {
		return 0, nil
	}
	for i, v := range names {
		if s == v {
			return 1 << i, nil
		}
	}
	return 0, fmt.Errorf("invalid %s name: %q", ename, s)
}

func ParseSet[T number](ename string, s string, names []string) (T, error) {
	var (
		out  T
		last error
	)
	for _, w := range strings.Split(s, "+") {
		w = strings.TrimSpace(w)
		v, err := Parse[T](ename, w, names)
		if err != nil {
			last = err
		}
		out |= v
	}
	return out, last
}

func ParseMulti[T number](ename string, s string, lists [][]string) (T, error) {
	s = strings.ToUpper(s)
	if s == "" || s == "NULL" {
		return 0, nil
	}
	for _, names := range lists {
		for i, name := range names {
			if s == name {
				return 1 << i, nil
			}
		}
	}
	return 0, fmt.Errorf("invalid %s name: %q", ename, s)
}

func ParseSetMulti[T number](ename string, s string, lists [][]string) (T, error) {
	var (
		out  T
		last error
	)
	for _, w := range strings.Split(s, "+") {
		w = strings.TrimSpace(w)
		v, err := ParseMulti[T](ename, w, lists)
		if err != nil {
			last = err
		}
		out |= v
	}
	return out, last
}

func MarshalJSONArray[E interface{ Split() []T }, T interface{ String() string }](v E) ([]byte, error) {
	var arr []string
	for _, s := range v.Split() {
		arr = append(arr, s.String())
	}
	return json.Marshal(arr)
}
