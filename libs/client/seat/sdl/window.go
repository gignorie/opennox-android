package sdl

import (
	"image"
	"log/slog"
	"os"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/opennox/libs/client/seat"
	"github.com/opennox/libs/env"
)

var (
	debugGpad = os.Getenv("NOX_DEBUG_GPAD") == "true"
	// CheckInMenu is an optional callback to check if the game is currently in menu mode.
	CheckInMenu func() bool
	// CheckInInventory is an optional callback to check if the inventory or GUI dialog is active.
	CheckInInventory func() bool
	// SetAndroidOverlayVisible is an optional callback to toggle Android virtual controls visibility.
	SetAndroidOverlayVisible func(visible bool)
	// ResetAndroidTouches is an optional callback to reset Android touch state.
	ResetAndroidTouches func()
)

var _ seat.Seat = &Window{}

type backend interface {
	Close() error
	NewSurface(sz image.Point, filter bool) seat.Surface
	Clear()
	ScreenSize() image.Point
	ResizeScreen(sz image.Point)
	SetGamma(v float32)
	Present()
}

type Window struct {
	log      *slog.Logger
	win      *sdl.Window
	b        backend
	prevPos  image.Point
	prevSz   image.Point
	textInp  bool
	mode     seat.ScreenMode
	rel          bool
	mpos         image.Point
	onResize     []func(sz image.Point)
	onInput      []func(ev seat.InputEvent)
	pendingTapDown int
	pendingTapUp   int
}

func (win *Window) Close() error {
	if win.win == nil {
		return nil
	}
	if win.b != nil {
		_ = win.b.Close()
		win.b = nil
	}
	err := win.win.Destroy()
	win.win = nil
	win.onResize = nil
	win.onInput = nil
	sdl.Quit()
	return err
}

func (win *Window) NewSurface(sz image.Point, filter bool) seat.Surface {
	return win.b.NewSurface(sz, filter)
}

func (win *Window) Clear() {
	win.b.Clear()
}

func (win *Window) ScreenSize() image.Point {
	return win.b.ScreenSize()
}

func (win *Window) screenPos() image.Point {
	x, y := win.win.GetPosition()
	return image.Point{
		X: int(x), Y: int(y),
	}
}

func (win *Window) displayRect() sdl.Rect {
	disp, err := win.win.GetDisplayIndex()
	if err != nil {
		win.log.Warn("can't get display index", "err", err)
		return sdl.Rect{}
	}
	rect, err := sdl.GetDisplayBounds(disp)
	if err != nil {
		win.log.Warn("can't get display bounds", "err", err)
		return sdl.Rect{}
	}
	return rect
}

func (win *Window) ScreenMaxSize() image.Point {
	rect := win.displayRect()
	return image.Point{
		X: int(rect.W), Y: int(rect.H),
	}
}

func (win *Window) setSize(sz image.Point) {
	win.log.Info("window", "size", sz)
	win.win.SetSize(int32(sz.X), int32(sz.Y))
}

func (win *Window) center() {
	win.win.SetPosition(sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED)
}

func (win *Window) ResizeScreen(sz image.Point) {
	win.b.ResizeScreen(sz)
}

func (win *Window) setRelative(rel bool) {
	if win.rel == rel {
		return
	}
	win.rel = rel
	win.win.SetGrab(rel)
	sdl.SetRelativeMouseMode(rel)
}

func (win *Window) SetScreenMode(mode seat.ScreenMode) {
	if win.mode == mode {
		return
	}
	if win.mode == seat.Windowed {
		// preserve size and pos, so we can restore them later
		win.prevSz = win.ScreenSize()
		win.prevPos = win.screenPos()
	}
	switch mode {
	case seat.Windowed:
		win.win.SetFullscreen(0)
		win.win.SetResizable(true)
		win.win.SetBordered(true)
		win.setSize(win.prevSz)
		if win.prevPos != (image.Point{}) {
			win.win.SetPosition(int32(win.prevPos.X), int32(win.prevPos.Y))
		} else {
			win.center()
		}
		if env.IsDevMode() || env.IsE2E() {
			sdl.ShowCursor(sdl.ENABLE)
		} else {
			sdl.ShowCursor(sdl.DISABLE)
		}
		win.setRelative(false)
	case seat.Fullscreen:
		win.win.SetResizable(false)
		win.win.SetBordered(false)
		win.setSize(win.ScreenMaxSize())
		win.win.SetFullscreen(uint32(sdl.WINDOW_FULLSCREEN_DESKTOP))
		sdl.ShowCursor(sdl.DISABLE)
		win.setRelative(true)
	case seat.Borderless:
		win.win.SetFullscreen(0)
		win.win.SetResizable(false)
		win.win.SetBordered(true)
		win.setSize(win.ScreenMaxSize())
		win.center()
		sdl.ShowCursor(sdl.DISABLE)
		win.setRelative(false)
	}
	win.mode = mode
}

func (win *Window) SetGamma(v float32) {
	win.b.SetGamma(v)
}

func (win *Window) OnScreenResize(fnc func(sz image.Point)) {
	win.onResize = append(win.onResize, fnc)
}

func (win *Window) Present() {
	win.b.Present()
}
