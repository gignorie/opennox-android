//go:build !windows
// +build !windows

package ifs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

func TestNormalize(t *testing.T) {
	dir, err := os.MkdirTemp("", "nox_fs_")
	must.NoError(t, err)
	defer os.RemoveAll(dir)

	dir1 := filepath.Join(dir, "AbC", "def")
	err = os.MkdirAll(dir1, 0755)
	must.NoError(t, err)
	file1 := filepath.Join(dir1, "File.txt")

	err = os.WriteFile(file1, []byte("data"), 0644)
	must.NoError(t, err)

	must.EqOp(t,
		file1,
		Normalize(strings.Join([]string{dir, "abc", "Def", "FILE.TXT"}, "\\")),
	)
	must.EqOp(t,
		filepath.Join(dir1, "NotExistent"),
		Normalize(strings.Join([]string{dir, "ABC", "DeF", "NotExistent"}, "\\")),
	)
}
