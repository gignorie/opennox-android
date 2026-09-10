//go:build !android

package sdl

import (
	"fmt"
	"image"
	"log/slog"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/opennox/libs/client/seat"
	"github.com/opennox/libs/client/seat/opengl"
	noxlog "github.com/opennox/libs/log"
)

type glBackend struct {
	win *Window
	gl  opengl.Window
}

func (b *glBackend) Close() error {
	b.gl.Close()
	return nil
}

func (b *glBackend) NewSurface(sz image.Point, filter bool) seat.Surface {
	return b.gl.NewSurface(sz, filter)
}

func (b *glBackend) Clear() {
	b.gl.Clear()
}

func (b *glBackend) ScreenSize() image.Point {
	w, h := b.win.win.GLGetDrawableSize()
	return image.Point{
		X: int(w), Y: int(h),
	}
}

func (b *glBackend) ResizeScreen(sz image.Point) {
	if b.win.mode != seat.Windowed {
		return
	}
	b.win.setSize(sz)
	b.win.prevSz = sz
}

func (b *glBackend) SetGamma(v float32) {
	b.gl.SetGamma(v)
}

func (b *glBackend) Present() {
	b.win.win.GLSwap()
}

func New(log *slog.Logger, title string, sz image.Point) (*Window, error) {
	log = noxlog.WithSystem(log, "sdl")
	sdl.SetHint(sdl.HINT_WINDOWS_DPI_AWARENESS, "permonitorv2")
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_TIMER); err != nil {
		return nil, fmt.Errorf("SDL Initialization failed: %w", err)
	}
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL version: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL version: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_CORE, 1); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL core: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG, 1); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL forward: %w", err)
	}
	if err := sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1); err != nil {
		return nil, fmt.Errorf("cannot set OpenGL attribute: %w", err)
	}

	win, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(sz.X), int32(sz.Y),
		sdl.WINDOW_RESIZABLE|sdl.WINDOW_OPENGL|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		sdl.Quit()
		return nil, fmt.Errorf("SDL Window creation failed: %w", err)
	}
	h := &Window{
		log:    log,
		win:    win,
		prevSz: sz,
	}
	h.SetScreenMode(seat.Windowed)

	gtx, err := win.GLCreateContext()
	if err != nil {
		_ = win.Destroy()
		sdl.Quit()
		return nil, fmt.Errorf("OpenGL creation failed: %w", err)
	}
	err = win.GLMakeCurrent(gtx)
	if err != nil {
		_ = win.Destroy()
		sdl.Quit()
		return nil, fmt.Errorf("OpenGL bind failed: %w", err)
	}
	sdl.GLSetSwapInterval(0)
	glb := &glBackend{win: h}
	if err := glb.gl.Init(log); err != nil {
		_ = win.Destroy()
		sdl.Quit()
		return nil, err
	}
	h.b = glb
	return h, nil
}
