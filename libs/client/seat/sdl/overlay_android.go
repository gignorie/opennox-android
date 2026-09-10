//go:build android

package sdl

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// drawAndroidOverlay renders virtual controls on top of the game frame.
// All elements are drawn with thin contours, neutral light-gray color (R: 200, G: 200, B: 200),
// and alpha not higher than 80-100 (90/255).
// If the overlay is hidden via TOGGLE, it immediately returns and does not draw.
func drawAndroidOverlay(r *sdl.Renderer, w, h int32) {
	if isMenuMode() || !isAndroidOverlayVisible() {
		return
	}

	_ = r.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	fw := float32(w)
	fh := float32(h)
	scale := fh / 1080.0
	if scale < 1.0 {
		scale = 1.0
	}

	// ── 1. Draw Virtual Buttons ──────────────────────────────────
	buttons := getButtons(fw, fh)
	for _, btn := range buttons {
		cx := int32(btn.cx)
		cy := int32(btn.cy)
		rad := int32(btn.radius)

		pressed := isButtonPressed(btn.id)
		if pressed {
			_ = r.SetDrawColor(240, 240, 240, 160)
			drawCircleOutline(r, cx, cy, rad, 28)
			drawCircleOutline(r, cx, cy, rad-2, 28)
		} else {
			_ = r.SetDrawColor(200, 200, 200, 90)
			drawCircleOutline(r, cx, cy, rad, 28)
		}

		drawButtonGlyph(r, btn.id, cx, cy, btn.radius/18.0)
	}

	// ── 2. Draw Floating Virtual Joystick (when active) ───────────
	active, baseX, baseY, curX, curY := getJoyRenderState()
	if active {
		bx := int32(baseX)
		by := int32(baseY)

		dx := curX - baseX
		dy := curY - baseY
		dist := float32(math.Hypot(float64(dx), float64(dy)))

		// Limit visual stick knob distance from base to max 80*scale px
		maxKnob := float32(80.0) * scale
		kx := curX
		ky := curY
		if dist > maxKnob && dist > 0 {
			kx = baseX + (dx/dist)*maxKnob
			ky = baseY + (dy/dist)*maxKnob
		}

		baseR := int32(math.Round(float64(75.0 * scale)))
		knobR := int32(math.Round(float64(32.0 * scale)))

		// Outer base circle
		_ = r.SetDrawColor(200, 200, 200, 80)
		drawCircleOutline(r, bx, by, baseR, 32)

		// Connecting line between base center and knob
		_ = r.SetDrawColor(200, 200, 200, 70)
		_ = r.DrawLine(bx, by, int32(kx), int32(ky))

		// Inner stick knob
		_ = r.SetDrawColor(220, 220, 220, 100)
		drawCircleOutline(r, int32(kx), int32(ky), knobR, 24)
	}
}

func drawCircleOutline(r *sdl.Renderer, cx, cy, radius int32, segments int) {
	if segments < 8 {
		segments = 8
	}
	if segments > 32 {
		segments = 32
	}
	var pts [33]sdl.Point
	step := (2.0 * math.Pi) / float64(segments)
	for i := 0; i <= segments; i++ {
		angle := float64(i) * step
		pts[i] = sdl.Point{
			X: cx + int32(math.Round(float64(radius)*math.Cos(angle))),
			Y: cy + int32(math.Round(float64(radius)*math.Sin(angle))),
		}
	}
	_ = r.DrawLines(pts[:segments+1])
}

func drawButtonGlyph(r *sdl.Renderer, id virtualButtonID, cx, cy int32, s float32) {
	i := func(v float32) int32 {
		return int32(math.Round(float64(v * s)))
	}
	switch id {
	case btnEsc:
		// Pause bars "||"
		_ = r.DrawLine(cx-i(3), cy-i(6), cx-i(3), cy+i(6))
		_ = r.DrawLine(cx+i(3), cy-i(6), cx+i(3), cy+i(6))

	case btnMap:
		// "M"
		_ = r.DrawLines([]sdl.Point{
			{X: cx - i(5), Y: cy + i(6)},
			{X: cx - i(5), Y: cy - i(6)},
			{X: cx, Y: cy + i(1)},
			{X: cx + i(5), Y: cy - i(6)},
			{X: cx + i(5), Y: cy + i(6)},
		})

	case btnSpl:
		// "S"
		_ = r.DrawLines([]sdl.Point{
			{X: cx + i(4), Y: cy - i(6)},
			{X: cx - i(4), Y: cy - i(6)},
			{X: cx - i(4), Y: cy},
			{X: cx + i(4), Y: cy},
			{X: cx + i(4), Y: cy + i(6)},
			{X: cx - i(4), Y: cy + i(6)},
		})

	case btnInv:
		// "I"
		_ = r.DrawLine(cx, cy-i(6), cx, cy+i(6))
		_ = r.DrawLine(cx-i(3), cy-i(6), cx+i(3), cy-i(6))
		_ = r.DrawLine(cx-i(3), cy+i(6), cx+i(3), cy+i(6))

	case btnKbd:
		// Keyboard icon
		_ = r.DrawLines([]sdl.Point{
			{X: cx - i(7), Y: cy - i(5)},
			{X: cx + i(7), Y: cy - i(5)},
			{X: cx + i(7), Y: cy + i(5)},
			{X: cx - i(7), Y: cy + i(5)},
			{X: cx - i(7), Y: cy - i(5)},
		})
		_ = r.DrawLine(cx-i(4), cy+i(2), cx+i(4), cy+i(2))
		_ = r.DrawPoint(cx-i(4), cy-i(2))
		_ = r.DrawPoint(cx, cy-i(2))
		_ = r.DrawPoint(cx+i(4), cy-i(2))

	case btnJump:
		// Upward chevron "^"
		_ = r.DrawLines([]sdl.Point{
			{X: cx - i(10), Y: cy + i(3)},
			{X: cx, Y: cy - i(7)},
			{X: cx + i(10), Y: cy + i(3)},
		})
		_ = r.DrawLines([]sdl.Point{
			{X: cx - i(10), Y: cy + i(9)},
			{X: cx, Y: cy - i(1)},
			{X: cx + i(10), Y: cy + i(9)},
		})

	case btnAttack:
		// "A"
		_ = r.DrawLines([]sdl.Point{
			{X: cx - i(6), Y: cy + i(8)},
			{X: cx, Y: cy - i(8)},
			{X: cx + i(6), Y: cy + i(8)},
		})
		_ = r.DrawLine(cx-i(4), cy+i(2), cx+i(4), cy+i(2))

	case btnSlot1:
		_ = r.DrawLine(cx, cy-i(6), cx, cy+i(6))
		_ = r.DrawLine(cx-i(3), cy-i(3), cx, cy-i(6))
		_ = r.DrawLine(cx-i(3), cy+i(6), cx+i(3), cy+i(6))

	case btnSlot2:
		_ = r.DrawLines([]sdl.Point{
			{X: cx - i(4), Y: cy - i(6)},
			{X: cx + i(4), Y: cy - i(6)},
			{X: cx + i(4), Y: cy - i(1)},
			{X: cx - i(4), Y: cy + i(6)},
			{X: cx + i(4), Y: cy + i(6)},
		})

	case btnSlot3:
		_ = r.DrawLines([]sdl.Point{
			{X: cx - i(4), Y: cy - i(6)},
			{X: cx + i(4), Y: cy - i(6)},
			{X: cx, Y: cy},
			{X: cx + i(4), Y: cy + i(1)},
			{X: cx + i(4), Y: cy + i(6)},
			{X: cx - i(4), Y: cy + i(6)},
		})

	case btnSlot4:
		_ = r.DrawLine(cx+i(2), cy-i(6), cx+i(2), cy+i(6))
		_ = r.DrawLines([]sdl.Point{
			{X: cx - i(4), Y: cy - i(6)},
			{X: cx - i(4), Y: cy + i(1)},
			{X: cx + i(4), Y: cy + i(1)},
		})

	case btnSlot5:
		_ = r.DrawLines([]sdl.Point{
			{X: cx + i(4), Y: cy - i(6)},
			{X: cx - i(4), Y: cy - i(6)},
			{X: cx - i(4), Y: cy},
			{X: cx + i(4), Y: cy},
			{X: cx + i(4), Y: cy + i(6)},
			{X: cx - i(4), Y: cy + i(6)},
		})

	case btnPotionZ:
		// "Z"
		_ = r.DrawLines([]sdl.Point{
			{X: cx - i(4), Y: cy - i(5)},
			{X: cx + i(4), Y: cy - i(5)},
			{X: cx - i(4), Y: cy + i(5)},
			{X: cx + i(4), Y: cy + i(5)},
		})

	case btnPotionX:
		// "X"
		_ = r.DrawLine(cx-i(4), cy-i(5), cx+i(4), cy+i(5))
		_ = r.DrawLine(cx+i(4), cy-i(5), cx-i(4), cy+i(5))

	case btnPotionC:
		// "C"
		_ = r.DrawLines([]sdl.Point{
			{X: cx + i(4), Y: cy - i(5)},
			{X: cx - i(4), Y: cy - i(5)},
			{X: cx - i(4), Y: cy + i(5)},
			{X: cx + i(4), Y: cy + i(5)},
		})
	}
}
