//go:build android

package sdl

import (
	"math"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/opennox/libs/client/seat"
)

func init() {
	touchEventHandler = func(win *Window, ev *sdl.TouchFingerEvent) {
		win.processTouchEvent(ev)
	}
	SetAndroidOverlayVisible = setAndroidOverlayVisible
	ResetAndroidTouches = func() {
		androidCtrl.mu.Lock()
		defer androidCtrl.mu.Unlock()
		androidCtrl.reset(nil)
	}
	filterTouchMouse = func() bool {
		return true
	}
}

func (ctrl *touchCtrl) reset(win *Window) {
	if ctrl.worldTouch.active {
		if ctrl.worldTouch.holdingLMB && win != nil {
			win.sendMouseButton(seat.MouseButtonLeft, false)
		}
		ctrl.worldTouch = worldTouchState{}
	}
	if ctrl.joy.active {
		if (ctrl.joy.rmbHeld || ctrl.joy.moving) && win != nil {
			win.sendMouseButton(seat.MouseButtonRight, false)
		}
		ctrl.joy = joystickState{}
	}
	for i := range ctrl.btnTouches {
		if ctrl.btnTouches[i].active {
			if ctrl.btnTouches[i].btnID == btnAttack && win != nil {
				win.sendMouseButton(seat.MouseButtonLeft, false)
			} else if ctrl.btnTouches[i].scancode != 0 && win != nil {
				win.pushKey(ctrl.btnTouches[i].scancode, false)
			}
			ctrl.btnTouches[i] = buttonTouch{}
		}
	}
}

func setAndroidOverlayVisible(visible bool) {
	androidCtrl.mu.Lock()
	defer androidCtrl.mu.Unlock()
	androidCtrl.overlayVisible = visible
}

const (
	deadzonePx = float32(20.0) // 20 pixels deadzone as requested
)

type virtualButtonID int

const (
	btnNone virtualButtonID = iota
	// Top Service Bar
	btnEsc
	btnMap
	btnSpl
	btnInv
	btnKbd
	// Combat Block
	btnJump
	btnAttack
	btnSlot1
	btnSlot2
	btnSlot3
	btnSlot4
	btnSlot5
	// Potion Quick-Access Buttons
	btnPotionZ
	btnPotionX
	btnPotionC
)

type virtualButton struct {
	id       virtualButtonID
	scancode sdl.Scancode
	cx       float32
	cy       float32
	radius   float32
}

func getButtons(winW, winH float32) []virtualButton {
	scale := winH / 1080.0
	if scale < 1.0 {
		scale = 1.0 // Ensure button radius is at least 48 physical pixels
	}

	// In OpenNox, the Health and Mana indicators widget is at the bottom-right of the 4:3 game area
	gw := winH * 4.0 / 3.0
	if gw > winW {
		gw = winW
	}
	gameRight := (winW + gw) / 2.0
	widgetW := 91.0 * (gw / 640.0)
	widgetX := gameRight - widgetW

	potRadius := 28.0 * scale
	potY := winH - 45.0 * scale
	potC_X := widgetX - 35.0 * scale
	potX_X := widgetX - 95.0 * scale
	potZ_X := widgetX - 155.0 * scale

	return []virtualButton{
		// Top Service Bar (top-right corner)
		// Pause button (Escape) is the top-right button!
		{id: btnEsc, scancode: sdl.SCANCODE_ESCAPE, cx: winW - 65*scale, cy: 60*scale, radius: 48*scale},
		{id: btnMap, scancode: sdl.SCANCODE_TAB, cx: winW - 190*scale, cy: 60*scale, radius: 48*scale},
		{id: btnSpl, scancode: sdl.SCANCODE_B, cx: winW - 315*scale, cy: 60*scale, radius: 48*scale},
		{id: btnInv, scancode: sdl.SCANCODE_I, cx: winW - 440*scale, cy: 60*scale, radius: 48*scale},
		{id: btnKbd, scancode: 0, cx: winW - 565*scale, cy: 60*scale, radius: 48*scale},

		// Combat Block (bottom-right quarter)
		// Big Jump button
		{id: btnJump, scancode: sdl.SCANCODE_SPACE, cx: winW - 95*scale, cy: winH - 95*scale, radius: 62*scale},
		// Auxiliary Attack button
		{id: btnAttack, scancode: 0, cx: winW - 235*scale, cy: winH - 95*scale, radius: 54*scale},

		// Semicircle of 5 quick slot buttons
		{id: btnSlot1, scancode: sdl.SCANCODE_1, cx: winW - 350*scale, cy: winH - 160*scale, radius: 48*scale},
		{id: btnSlot2, scancode: sdl.SCANCODE_2, cx: winW - 310*scale, cy: winH - 280*scale, radius: 48*scale},
		{id: btnSlot3, scancode: sdl.SCANCODE_3, cx: winW - 210*scale, cy: winH - 360*scale, radius: 48*scale},
		{id: btnSlot4, scancode: sdl.SCANCODE_4, cx: winW - 95*scale, cy: winH - 370*scale, radius: 48*scale},
		{id: btnSlot5, scancode: sdl.SCANCODE_5, cx: winW - 45*scale, cy: winH - 260*scale, radius: 48*scale},

		// Quick potion buttons (Z: Cure Poison, X: Health, C: Mana) near health/mana indicators
		{id: btnPotionZ, scancode: sdl.SCANCODE_Z, cx: potZ_X, cy: potY, radius: potRadius},
		{id: btnPotionX, scancode: sdl.SCANCODE_X, cx: potX_X, cy: potY, radius: potRadius},
		{id: btnPotionC, scancode: sdl.SCANCODE_C, cx: potC_X, cy: potY, radius: potRadius},
	}
}

type buttonTouch struct {
	active   bool
	fingerID sdl.FingerID
	btnID    virtualButtonID
	scancode sdl.Scancode
}

type joystickState struct {
	active    bool
	fingerID  sdl.FingerID
	baseX     float32
	baseY     float32
	curX      float32
	curY      float32
	moving    bool
	rmbHeld   bool
	running   bool
	startTime time.Time
}

type worldTouchState struct {
	active     bool
	fingerID   sdl.FingerID
	startX     float32
	startY     float32
	curX       float32
	curY       float32
	startTime  time.Time
	holdingLMB bool
}

type touchCtrl struct {
	mu             sync.Mutex
	overlayVisible bool

	btnTouches [8]buttonTouch
	joy        joystickState
	worldTouch worldTouchState

	lastDirX float32
	lastDirY float32
}

var androidCtrl = &touchCtrl{
	overlayVisible: true,
	lastDirX:       0,
	lastDirY:       1,
}

func isMenuMode() bool {
	if CheckInMenu != nil {
		return CheckInMenu()
	}
	return false
}

func isGuiOrInventory() bool {
	if isMenuMode() {
		return true
	}
	if CheckInInventory != nil && CheckInInventory() {
		return true
	}
	return false
}

func isAndroidOverlayVisible() bool {
	androidCtrl.mu.Lock()
	defer androidCtrl.mu.Unlock()
	return androidCtrl.overlayVisible
}

func getJoyRenderState() (active bool, baseX, baseY, curX, curY float32) {
	androidCtrl.mu.Lock()
	defer androidCtrl.mu.Unlock()
	return androidCtrl.joy.active, androidCtrl.joy.baseX, androidCtrl.joy.baseY, androidCtrl.joy.curX, androidCtrl.joy.curY
}

func isButtonPressed(id virtualButtonID) bool {
	androidCtrl.mu.Lock()
	defer androidCtrl.mu.Unlock()
	for i := range androidCtrl.btnTouches {
		if androidCtrl.btnTouches[i].active && androidCtrl.btnTouches[i].btnID == id {
			return true
		}
	}
	return false
}

func (ctrl *touchCtrl) hitButton(pixelX, pixelY, winW, winH float32) (virtualButtonID, sdl.Scancode) {
	for _, btn := range getButtons(winW, winH) {
		dx := pixelX - btn.cx
		dy := pixelY - btn.cy
		if (dx*dx + dy*dy) <= (btn.radius * btn.radius * 1.5) {
			return btn.id, btn.scancode
		}
	}
	return btnNone, 0
}

func (ctrl *touchCtrl) addButtonTouch(id sdl.FingerID, btnID virtualButtonID, scancode sdl.Scancode) {
	for i := range ctrl.btnTouches {
		if !ctrl.btnTouches[i].active {
			ctrl.btnTouches[i] = buttonTouch{
				active:   true,
				fingerID: id,
				btnID:    btnID,
				scancode: scancode,
			}
			return
		}
	}
}

func (ctrl *touchCtrl) removeButtonTouch(id sdl.FingerID) *buttonTouch {
	for i := range ctrl.btnTouches {
		if ctrl.btnTouches[i].active && ctrl.btnTouches[i].fingerID == id {
			res := ctrl.btnTouches[i]
			ctrl.btnTouches[i] = buttonTouch{}
			return &res
		}
	}
	return nil
}

func (ctrl *touchCtrl) isButtonFinger(id sdl.FingerID) bool {
	for i := range ctrl.btnTouches {
		if ctrl.btnTouches[i].active && ctrl.btnTouches[i].fingerID == id {
			return true
		}
	}
	return false
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// processTouchEvent is the main entry-point for touch events on Android.
func (win *Window) processTouchEvent(ev *sdl.TouchFingerEvent) {
	ctrl := androidCtrl
	ctrl.mu.Lock()
	defer ctrl.mu.Unlock()

	// ── STAGE 1: Physical Window Pixel Translation ───────────────
	winW, winH := win.PhysicalSize()
	logW, logH := win.LogicalSize()
	if winW <= 0 || winH <= 0 {
		return
	}
	if logW <= 0 || logH <= 0 {
		logW, logH = winW, winH
	}

	pixelX := clampInt(int(math.Round(float64(ev.X*float32(winW)))), 0, int(winW)-1)
	pixelY := clampInt(int(math.Round(float64(ev.Y*float32(winH)))), 0, int(winH)-1)
	curX := float32(pixelX)
	curY := float32(pixelY)

	logPixelX := clampInt(int(math.Round(float64(ev.X*float32(logW)))), 0, int(logW)-1)
	logPixelY := clampInt(int(math.Round(float64(ev.Y*float32(logH)))), 0, int(logH)-1)
	logX := int32(logPixelX)
	logY := int32(logPixelY)

	// If in menu mode (intro, main menu, pause menu, modal dialogs, NPC conversation):
	// Direct mouse interaction with Left Mouse Button hold/drag/release in logical pixels.
	if isMenuMode() {
		switch ev.Type {
		case sdl.FINGERDOWN:
			if ctrl.worldTouch.active && ctrl.worldTouch.holdingLMB {
				win.sendMouseButton(seat.MouseButtonLeft, false)
			}
			ctrl.worldTouch = worldTouchState{
				active:     true,
				fingerID:   ev.FingerID,
				startX:     curX,
				startY:     curY,
				curX:       curX,
				curY:       curY,
				startTime:  time.Now(),
				holdingLMB: true,
			}
			win.sendMouseMotion(logX, logY)
			win.sendMouseButton(seat.MouseButtonLeft, true)
		case sdl.FINGERMOTION:
			ctrl.worldTouch.curX = curX
			ctrl.worldTouch.curY = curY
			win.sendMouseMotion(logX, logY)
		case sdl.FINGERUP:
			win.sendMouseMotion(logX, logY)
			if ctrl.worldTouch.active && ctrl.worldTouch.holdingLMB {
				win.sendMouseButton(seat.MouseButtonLeft, false)
			}
			ctrl.worldTouch = worldTouchState{}
		}
		return
	}

	// Active Gameplay Mode:
	switch ev.Type {
	case sdl.FINGERDOWN:
		// ── STAGE 2: Overlay UI Hit Testing (in physical pixelX, pixelY) ──
		// 1. Check virtual buttons
		btnID, scancode := ctrl.hitButton(curX, curY, float32(winW), float32(winH))
		if btnID != btnNone {
			if btnID == btnAttack {
				// Attack button: send MouseButtonLeft DOWN
				// Target coordinates: screen center + offset in last stick direction (in logical pixels)
				dirX := ctrl.lastDirX
				dirY := ctrl.lastDirY
				if dirX == 0 && dirY == 0 {
					dirY = 1.0
				}
				centerX := float32(logW) / 2.0
				centerY := float32(logH) / 2.0
				targetX := int(centerX + dirX*100.0)
				targetY := int(centerY + dirY*100.0)
				targetX = clampInt(targetX, 0, int(logW)-1)
				targetY = clampInt(targetY, 0, int(logH)-1)

				win.sendMouseMotion(int32(targetX), int32(targetY))
				win.sendMouseButton(seat.MouseButtonLeft, true)

				ctrl.addButtonTouch(ev.FingerID, btnID, 0)
				return // CONSUMED
			}

			if btnID == btnKbd {
				if sdl.IsScreenKeyboardShown(win.win) {
					win.SetTextInput(false)
				} else {
					win.SetTextInput(true)
				}
				ctrl.addButtonTouch(ev.FingerID, btnID, 0)
				return // CONSUMED
			}

			// Other buttons (Esc, Jump, Inv, Map, Spl, Slots 1-5, Potion Z, X, C)
			win.pushKey(scancode, true)
			ctrl.addButtonTouch(ev.FingerID, btnID, scancode)
			return // CONSUMED
		}

		// 2. Virtual Joystick strictly in bottom-left corner:
		// pixelX < winW * 0.28 && pixelY > winH * 0.45
		if curX < float32(winW)*0.28 && curY > float32(winH)*0.45 {
			if ctrl.joy.active && ctrl.joy.fingerID != ev.FingerID {
				if ctrl.joy.rmbHeld || ctrl.joy.moving {
					win.sendMouseButton(seat.MouseButtonRight, false)
				}
				ctrl.joy = joystickState{}
			}
			if !ctrl.joy.active {
				ctrl.joy = joystickState{
					active:    true,
					fingerID:  ev.FingerID,
					baseX:     curX,
					baseY:     curY,
					curX:      curX,
					curY:      curY,
					moving:    false,
					rmbHeld:   false,
					running:   false,
					startTime: time.Now(),
				}
				return // CONSUMED
			}
		}

		// ── STAGE 3: World / Game GUI (Tap-to-Interact & Drag-and-Drop) ──
		if ctrl.worldTouch.active && ctrl.worldTouch.holdingLMB {
			win.sendMouseButton(seat.MouseButtonLeft, false)
		}
		inGui := isGuiOrInventory()
		ctrl.worldTouch = worldTouchState{
			active:     true,
			fingerID:   ev.FingerID,
			startX:     curX,
			startY:     curY,
			curX:       curX,
			curY:       curY,
			startTime:  time.Now(),
			holdingLMB: inGui,
		}
		win.sendMouseMotion(logX, logY)
		if inGui {
			win.sendMouseButton(seat.MouseButtonLeft, true)
		}

	case sdl.FINGERMOTION:
		// 1. Button touch ignores drag
		if ctrl.isButtonFinger(ev.FingerID) {
			return // CONSUMED
		}

		// 2. Virtual Joystick movement
		if ctrl.joy.active && ctrl.joy.fingerID == ev.FingerID {
			ctrl.joy.curX = curX
			ctrl.joy.curY = curY
			dx := curX - ctrl.joy.baseX
			dy := curY - ctrl.joy.baseY
			dist := float32(math.Hypot(float64(dx), float64(dy)))

			scale := float32(winH) / 1080.0
			if scale < 1.0 {
				scale = 1.0
			}
			joyDeadzone := float32(12.0) * scale

			if dist > joyDeadzone {
				if !ctrl.joy.moving {
					ctrl.joy.moving = true
					ctrl.joy.rmbHeld = true
					win.sendMouseButton(seat.MouseButtonRight, true)
				}
				normX := dx / dist
				normY := dy / dist
				ctrl.lastDirX = normX
				ctrl.lastDirY = normY

				// Walk vs Run threshold with hysteresis to eliminate jitter:
				// Base radius is 75*scale px.
				// Mild deflection (<32*scale px) is walking; strong deflection (>=42*scale px) is running.
				walkDistPx := float32(32.0) * scale
				runDistPx := float32(42.0) * scale
				if ctrl.joy.running {
					if dist < walkDistPx {
						ctrl.joy.running = false
					}
				} else {
					if dist >= runDistPx {
						ctrl.joy.running = true
					}
				}

				// Cursor distance from screen center:
				// OpenNox engine treats world distance <= 100 as walk, > 100 as run.
				cursorOffset := float32(50.0) // guaranteed walk
				if ctrl.joy.running {
					cursorOffset = float32(180.0) // guaranteed run
				}

				centerX := float32(logW) / 2.0
				centerY := float32(logH) / 2.0
				targetX := int(centerX + normX*cursorOffset)
				targetY := int(centerY + normY*cursorOffset)
				targetX = clampInt(targetX, 0, int(logW)-1)
				targetY = clampInt(targetY, 0, int(logH)-1)
				win.sendMouseMotion(int32(targetX), int32(targetY))
			} else {
				if ctrl.joy.moving {
					ctrl.joy.moving = false
					ctrl.joy.running = false
				}
			}
			return // CONSUMED
		}

		// 3. Stage 3 World Touch / Inventory motion
		if ctrl.worldTouch.active {
			ctrl.worldTouch.curX = curX
			ctrl.worldTouch.curY = curY

			// If finger was held or dragged in world mode, start holding LMB
			if !ctrl.worldTouch.holdingLMB {
				dx := curX - ctrl.worldTouch.startX
				dy := curY - ctrl.worldTouch.startY
				dist := float32(math.Hypot(float64(dx), float64(dy)))
				scale := float32(winH) / 1080.0
				if scale < 1.0 {
					scale = 1.0
				}
				maxTapDist := float32(25.0) * scale
				if dist > maxTapDist || time.Since(ctrl.worldTouch.startTime) > 250*time.Millisecond {
					ctrl.worldTouch.holdingLMB = true
					win.sendMouseButton(seat.MouseButtonLeft, true)
				}
			}

			win.sendMouseMotion(logX, logY)
		}

	case sdl.FINGERUP:
		// 1. Button release
		if bt := ctrl.removeButtonTouch(ev.FingerID); bt != nil {
			if bt.btnID == btnAttack {
				win.sendMouseButton(seat.MouseButtonLeft, false)
				return // CONSUMED
			}
			if bt.btnID == btnKbd {
				return // CONSUMED
			}
			if bt.scancode != 0 {
				win.pushKey(bt.scancode, false)
			}
			return // CONSUMED
		}

		// 2. Virtual Joystick release
		if ctrl.joy.active && ctrl.joy.fingerID == ev.FingerID {
			if ctrl.joy.rmbHeld || ctrl.joy.moving {
				win.sendMouseButton(seat.MouseButtonRight, false)
			} else {
				// Short tap in bottom-left corner without movement (<25px, <450ms): Tap-to-Interact
				dx := curX - ctrl.joy.baseX
				dy := curY - ctrl.joy.baseY
				dist := float32(math.Hypot(float64(dx), float64(dy)))
				dt := time.Since(ctrl.joy.startTime)
				scale := float32(winH) / 1080.0
				if scale < 1.0 {
					scale = 1.0
				}
				maxTapDist := float32(25.0) * scale
				if dist < maxTapDist && dt < 450*time.Millisecond {
					tapX, tapY, _ := win.WindowToLogical(int32(ctrl.joy.baseX), int32(ctrl.joy.baseY))
					win.triggerTapClick(tapX, tapY)
				}
			}
			ctrl.joy = joystickState{}
			return // CONSUMED
		}

		// 3. Stage 3 World Touch release (Tap-to-Interact or LMB release)
		if ctrl.worldTouch.active {
			dx := curX - ctrl.worldTouch.startX
			dy := curY - ctrl.worldTouch.startY
			dist := float32(math.Hypot(float64(dx), float64(dy)))
			dt := time.Since(ctrl.worldTouch.startTime)
			startX := ctrl.worldTouch.startX
			startY := ctrl.worldTouch.startY
			wasHolding := ctrl.worldTouch.holdingLMB

			ctrl.worldTouch = worldTouchState{}

			if wasHolding {
				win.sendMouseMotion(logX, logY)
				win.sendMouseButton(seat.MouseButtonLeft, false)
			} else {
				scale := float32(winH) / 1080.0
				if scale < 1.0 {
					scale = 1.0
				}
				maxTapDist := float32(25.0) * scale
				if dist < maxTapDist && dt < 450*time.Millisecond {
					tapX, tapY, _ := win.WindowToLogical(int32(startX), int32(startY))
					win.triggerTapClick(tapX, tapY)
				}
			}
		}
	}
}
