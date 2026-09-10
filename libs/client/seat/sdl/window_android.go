//go:build android

package sdl

import (
	"fmt"
	"image"
	"log/slog"
	"math"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/opennox/libs/client/seat"
	noxlog "github.com/opennox/libs/log"
	"github.com/opennox/libs/noximage"
)

type rendererBackend struct {
	win      *Window
	renderer *sdl.Renderer
}

func (b *rendererBackend) Close() error {
	if b.renderer != nil {
		_ = b.renderer.Destroy()
		b.renderer = nil
	}
	return nil
}

func (b *rendererBackend) NewSurface(sz image.Point, filter bool) seat.Surface {
	if filter {
		sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	} else {
		sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
	}
	tex, err := b.renderer.CreateTexture(sdl.PIXELFORMAT_RGB555, sdl.TEXTUREACCESS_STREAMING, int32(sz.X), int32(sz.Y))
	if err != nil {
		b.win.log.Error("cannot create texture", "err", err)
	}
	return &Surface{renderer: b.renderer, sz: sz, tex: tex}
}

func (b *rendererBackend) Clear() {
	if b.renderer != nil {
		_ = b.renderer.SetDrawColor(0, 0, 0, 255)
		_ = b.renderer.Clear()
	}
}

func (b *rendererBackend) ScreenSize() image.Point {
	if b.renderer != nil {
		w, h := b.renderer.GetLogicalSize()
		if w > 0 && h > 0 {
			return image.Point{X: int(w), Y: int(h)}
		}
	}
	w, h := b.win.LogicalSize()
	return image.Point{X: int(w), Y: int(h)}
}

func (b *rendererBackend) ResizeScreen(sz image.Point) {
	b.win.prevSz = sz
}

func (b *rendererBackend) SetGamma(v float32) {
	_ = b.win.win.SetBrightness(v)
}

func (b *rendererBackend) Present() {
	if b.renderer != nil {
		// Draw virtual controls overlay on top of game frame in physical window coordinates
		winW, winH := b.win.PhysicalSize()
		if winW > 0 && winH > 0 && isAndroidOverlayVisible() {
			_ = b.renderer.SetViewport(nil)
			_ = b.renderer.SetScale(1.0, 1.0)
			drawAndroidOverlay(b.renderer, winW, winH)
			logW, logH := b.win.LogicalSize()
			if logW > 0 && logH > 0 {
				_ = b.renderer.SetLogicalSize(logW, logH)
			}
		}
		b.renderer.Present()
	}
}

func New(log *slog.Logger, title string, sz image.Point) (*Window, error) {
	log = noxlog.WithSystem(log, "sdl")
	sdl.SetHint(sdl.HINT_WINDOWS_DPI_AWARENESS, "permonitorv2")
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_TIMER); err != nil {
		return nil, fmt.Errorf("SDL Initialization failed: %w", err)
	}
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	win, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(sz.X), int32(sz.Y),
		sdl.WINDOW_FULLSCREEN|sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		sdl.Quit()
		return nil, fmt.Errorf("SDL Window creation failed: %w", err)
	}
	renderer, err := sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Warn("cannot create accelerated renderer, trying software", "err", err)
		renderer, err = sdl.CreateRenderer(win, -1, sdl.RENDERER_SOFTWARE)
		if err != nil {
			_ = win.Destroy()
			sdl.Quit()
			return nil, fmt.Errorf("SDL Renderer creation failed: %w", err)
		}
	}

	screenW, screenH := win.GetSize()
	screenSz := image.Point{X: int(screenW), Y: int(screenH)}
	log.Info("android physical screen size", "screen", screenSz)

	const logicalH = 380
	logicalW := int32(logicalH)
	if screenH > 0 {
		logicalW = int32(float32(logicalH) * (float32(screenW) / float32(screenH)))
		// Round down to multiple of 4 for alignment
		logicalW = logicalW &^ 3
		if logicalW < int32(logicalH) {
			logicalW = int32(logicalH)
		}
	}
	logicalSz := image.Point{X: int(logicalW), Y: logicalH}
	log.Info("android widescreen logical size", "screen", screenSz, "logical", logicalSz)
	_ = renderer.SetLogicalSize(logicalW, logicalH)

	h := &Window{
		log:    log,
		win:    win,
		prevSz: logicalSz,
		mode:   seat.Fullscreen,
	}
	h.b = &rendererBackend{win: h, renderer: renderer}
	return h, nil
}

// LogicalSize returns the logical pixel dimensions of the window/screen.
func (win *Window) LogicalSize() (int32, int32) {
	if win == nil {
		return 0, 0
	}
	if rb, ok := win.b.(*rendererBackend); ok && rb.renderer != nil {
		w, h := rb.renderer.GetLogicalSize()
		if w > 0 && h > 0 {
			return w, h
		}
	}
	if win.win != nil {
		w, h := win.win.GetSize()
		return w, h
	}
	return int32(win.prevSz.X), int32(win.prevSz.Y)
}

// PhysicalSize returns the physical pixel dimensions of the window/screen.
func (win *Window) PhysicalSize() (int32, int32) {
	if win == nil {
		return 0, 0
	}
	if win.win != nil {
		w, h := win.win.GetSize()
		if w > 0 && h > 0 {
			return w, h
		}
	}
	if rb, ok := win.b.(*rendererBackend); ok && rb.renderer != nil {
		w, h, err := rb.renderer.GetOutputSize()
		if err == nil && w > 0 && h > 0 {
			return w, h
		}
	}
	return int32(win.prevSz.X), int32(win.prevSz.Y)
}

// WindowToLogical returns logical pixel coordinates for physical window coordinates.
func (win *Window) WindowToLogical(pixelX, pixelY int32) (int32, int32, bool) {
	if win == nil {
		return pixelX, pixelY, true
	}
	physW, physH := win.PhysicalSize()
	logW, logH := win.LogicalSize()
	if physW > 0 && physH > 0 && logW > 0 && logH > 0 {
		lx := int32(math.Round(float64(pixelX) * float64(logW) / float64(physW)))
		ly := int32(math.Round(float64(pixelY) * float64(logH) / float64(physH)))
		return clampInt32(lx, 0, logW-1), clampInt32(ly, 0, logH-1), true
	}
	return pixelX, pixelY, true
}

func clampInt32(v, lo, hi int32) int32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// TouchToLogical converts normalized touch coordinates (0.0 to 1.0) into logical pixel coordinates.
func (win *Window) TouchToLogical(touchX, touchY float32) (int32, int32) {
	logW, logH := win.LogicalSize()
	pixelX := int32(math.Round(float64(touchX * float32(logW))))
	pixelY := int32(math.Round(float64(touchY * float32(logH))))
	return pixelX, pixelY
}

type Surface struct {
	renderer *sdl.Renderer
	sz       image.Point
	tex      *sdl.Texture
}

func (s *Surface) Size() image.Point {
	return s.sz
}

func (s *Surface) Update(data *noximage.Image16) {
	if s.sz != data.Size() {
		panic("invalid image size")
	}
	if s.tex != nil && len(data.Pix) > 0 {
		raw := unsafe.Slice((*byte)(unsafe.Pointer(&data.Pix[0])), len(data.Pix)*2)
		_ = s.tex.Update(nil, raw, s.sz.X*2)
	}
}

func (s *Surface) Draw(vp image.Rectangle) {
	if s.tex == nil || s.renderer == nil {
		return
	}
	dst := sdl.Rect{
		X: int32(vp.Min.X),
		Y: int32(vp.Min.Y),
		W: int32(vp.Dx()),
		H: int32(vp.Dy()),
	}
	_ = s.renderer.Copy(s.tex, nil, &dst)
}

func (s *Surface) Destroy() {
	if s.tex != nil {
		_ = s.tex.Destroy()
		s.tex = nil
	}
	s.renderer = nil
}
