package eval

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func allowedFile(mode fs.FileMode) bool {
	if mode&fs.ModeSymlink != 0 {
		return false
	} else if !mode.IsDir() && !mode.IsRegular() {
		return false
	}
	return true
}

type FS interface {
	fs.FS
	fs.StatFS
}

type File interface {
	fs.File
	fs.ReadDirFile
}

var _ FS = (*goFS)(nil)

func newModFS(log *slog.Logger, mapsDir, modsDir string) FS {
	mapsFS := &goFS{
		log:  log.With("vfs", "maps"),
		root: mapsDir,
	}
	if modsDir == "" {
		return mapsFS
	}
	return &overlayFS{Sub: []FS{
		mapsFS,
		&goFS{
			log:  log.With("vfs", "mods"),
			root: modsDir,
		},
	}}
}

// goFS is a modified FS tailored for serving OpenNox mod/map Go packages.
//
// It has the following assumptions/constrains:
// - File system is read-only.
// - Files outside root cannot be accessed.
// - Symlinks and other non-regular files are ignored.
// - All paths must start with /src/ - this is to make Yaegi happy (it assumes packages in GOPATH must be in src).
type goFS struct {
	log  *slog.Logger
	root string // folder mapped to prefix in virtual FS
}

type modFileInfo struct {
	m *goFS // to correctly print path
	fs.FileInfo
}

func (fi modFileInfo) Name() string {
	p := fi.m.virtPath(fi.FileInfo.Name())
	return p
}

type modDirEntry struct {
	m *goFS // to correctly print path
	fs.DirEntry
}

func (e modDirEntry) Name() string {
	p := e.m.virtPath(e.DirEntry.Name())
	return p
}

func (e modDirEntry) Info() (fs.FileInfo, error) {
	fi, err := e.DirEntry.Info()
	if err != nil {
		return nil, err
	}
	return modFileInfo{m: e.m, FileInfo: fi}, nil
}

type modFile struct {
	m *goFS // to correctly print path in Stat
	f File
}

func (f modFile) Read(data []byte) (int, error) {
	return f.f.Read(data)
}

func (f modFile) Close() error {
	return f.f.Close()
}

func (f modFile) Stat() (fs.FileInfo, error) {
	fi, err := f.f.Stat()
	if err != nil {
		return nil, err
	}
	return modFileInfo{m: f.m, FileInfo: fi}, nil
}

func (f modFile) ReadDir(n int) ([]fs.DirEntry, error) {
	list, err := f.f.ReadDir(n)
	if err != nil {
		return nil, err
	}
	out := make([]fs.DirEntry, 0, len(list))
	for i := 0; i < len(list); i++ {
		fi := list[i]
		if allowedFile(fi.Type()) {
			out = append(out, modDirEntry{m: f.m, DirEntry: fi})
		}
	}
	return out, nil
}

func (m *goFS) virtPath(p string) string {
	if !filepath.IsAbs(p) {
		p = filepath.Clean(p)
		return p
	}
	if !strings.HasPrefix(p, m.root) {
		return filepath.Base(p)
	}
	p = p[len(m.root):]
	p = filepath.Clean(p)
	p = filepath.Join(rootPref, p)
	return p
}

func (m *goFS) mapPath(p string) (string, bool) {
	if filepath.IsAbs(p) {
		if !strings.HasPrefix(p, rootPref) {
			return "", false
		}
		p = p[len(rootPref):]
	}
	p = filepath.Clean(p)
	p = filepath.Join(m.root, p)
	return p, true
}

func (m *goFS) Stat(orig string) (fs.FileInfo, error) {
	log := m.log.With("path", orig)
	p, ok := m.mapPath(orig)
	if !ok {
		log.Debug("stat", "err", fs.ErrNotExist)
		return nil, fs.ErrNotExist
	}
	log = log.With("real", p)
	fi, err := os.Stat(p)
	if err != nil {
		log.Debug("stat", "err", errors.Unwrap(err))
		return nil, err
	}
	log.Debug("stat")
	// do not expose other methods except the base ones, fix file name
	return modFileInfo{m: m, FileInfo: fi}, nil
}

func (m *goFS) Open(orig string) (fs.File, error) {
	log := m.log.With("path", orig)
	p, ok := m.mapPath(orig)
	if !ok {
		log.Debug("open", "err", fs.ErrNotExist)
		return nil, fs.ErrNotExist
	}
	log = log.With("real", p)
	f, err := os.Open(p)
	if err != nil {
		log.Debug("open", "err", err)
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		_ = f.Close()
		log.Debug("open", "err", err)
		return nil, err
	}
	if !allowedFile(fi.Mode()) {
		_ = f.Close()
		log.Debug("open", "err", fs.ErrNotExist)
		return nil, fs.ErrNotExist
	}
	log.Debug("open")
	// do not expose other methods except the base ones
	return modFile{m: m, f: f}, nil
}

var _ FS = (*overlayFS)(nil)

type overlayFS struct {
	Sub []FS // first match is used
}

func (l *overlayFS) Stat(name string) (fs.FileInfo, error) {
	for _, s := range l.Sub {
		f, err := fs.Stat(s, name)
		if !errors.Is(err, fs.ErrNotExist) {
			return f, err
		}
	}
	return nil, fs.ErrNotExist
}

func (l *overlayFS) Open(name string) (fs.File, error) {
	// First check if any of the FS has this dir/file.
	for _, s := range l.Sub {
		fi, err := s.Stat(name)
		if err == nil {
			if !allowedFile(fi.Mode()) {
				continue
			}
			if fi.IsDir() {
				// For directories, use overlay implementation.
				return &overlayDir{fs: l, path: name, fi: fi}, nil
			}
			// For files, open directly.
			return s.Open(name)
		} else if !errors.Is(err, fs.ErrNotExist) {
			// Report other errors directly.
			return nil, err
		}
	}
	return nil, fs.ErrNotExist
}

var _ File = (*overlayDir)(nil)

type overlayDir struct {
	fs     *overlayFS
	path   string
	fi     fs.FileInfo
	list   []fs.DirEntry
	loaded bool
}

func (d *overlayDir) Stat() (fs.FileInfo, error) {
	return d.fi, nil
}

func (d *overlayDir) Read(_ []byte) (int, error) {
	return 0, fs.ErrInvalid
}

func (d *overlayDir) Close() error {
	return nil
}

func (d *overlayDir) ReadDir(n int) ([]fs.DirEntry, error) {
	if !d.loaded {
		d.loaded = true

		var byName = make(map[string]struct{})
		for _, s := range d.fs.Sub {
			list, err := fs.ReadDir(s, d.path)
			if err != nil {
				continue
			}
			for _, fi := range list {
				name := fi.Name()
				if _, ok := byName[name]; ok {
					continue
				}
				byName[name] = struct{}{}
				d.list = append(d.list, fi)
			}
		}
		slices.SortFunc(d.list, func(a, b fs.DirEntry) int {
			return strings.Compare(a.Name(), b.Name())
		})
	}
	if n == 0 {
		return nil, nil
	}
	if n < 0 || n > len(d.list) {
		n = len(d.list)
	}
	list := d.list[:n]
	d.list = d.list[n:]
	return list, nil
}
