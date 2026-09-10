//go:build !server

package opennox

import (
	"image"
	"unsafe"

	"github.com/spf13/viper"
	"github.com/tawesoft/golib/v2/dialog"

	"github.com/opennox/libs/client/seat/sdl"
	"github.com/opennox/libs/env"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/input"
	"github.com/opennox/opennox/v1/client/render"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/version"
	"github.com/opennox/opennox/v1/legacy"
)

func init() {
	viper.SetDefault(configVideoFiltering, true)
	viper.SetDefault(configVideoStretch, false)
}

func isMultiplayer() bool {
	return nox_client_isConnected() || noxflags.HasGame(noxflags.GameOnline) || noxflags.HasGame(noxflags.GameClient)
}

// isGameHUDVisible checks if the in-game HUD (health/mana globes, skill icons, weapon bar)
// is currently visible on screen. When HUD is visible, the touch controls overlay MUST be visible.
func (c *Client) isGameHUDVisible() bool {
	if c == nil {
		return false
	}
	// 1. If in the main menu shell (Start Game, Multiplayer, Options, Character Select), HUD is not active.
	// flag_815132 is set to 1 by nox_xxx_wndLoadMainBG_4A2210 / nox_xxx_cliWaitForJoinData_43BFE0
	// and cleared to 0 by initGameSession435CC0 when the game session starts.
	if nox_client_gui_flag_815132 != 0 {
		return false
	}
	// 2. If watching cutscene/movies, HUD is not active
	if c.GameGetStateCode() == client.StateMovies {
		return false
	}
	// 3. If in Esc Quit / Pause menu, HUD is blocked by modal menu
	if legacy.Nox_gui_xxx_check_446360() == 1 {
		return false
	}
	// 4. If NPC conversation dialog is open, HUD is hidden/inactive so user can tap dialogue choices
	if legacy.Sub_47A260() == 1 {
		return false
	}
	if root := legacy.Get_dword_5d4594_1123524(); root != nil && !root.GetFlags().IsHidden() {
		return false
	}
	// 5. If Trader / Shop window is open, HUD is inactive so user can tap shop items.
	// In Nox, Sub_478030() returns dword_5d4594_1098624 (1 when shop is open, 0 when closed).
	// Sub_479590() returned dword_5d4594_1098628 (tab state, set to 1 by sub_478F80 on exit, causing permanent HUD hide).
	if legacy.Sub_478030() != 0 {
		return false
	}
	// 6. If a blocking modal dialog box (OK/Cancel, Quit Confirm) is currently open and not hidden
	if dia := nox_gui_curDialog_830224; dia != nil && !dia.GetFlags().IsHidden() {
		return false
	}

	// 7. Authoritative OpenNox flag: HUD (HP, Mana, Skills, Weapon) is being rendered
	if nox_client_getRenderGUI() != 0 {
		return true
	}

	// 8. Direct check: Health/Mana meter window (dword_5d4594_1090276) is created and visible
	if meterPtr := memmap.Uint32(0x5D4594, 1090276); meterPtr != 0 {
		win := legacy.AsWindowP(unsafe.Pointer(uintptr(meterPtr)))
		if win != nil && !win.GetFlags().IsHidden() {
			return true
		}
	}

	// 9. Fallback: if player unit exists or play state is active (3)
	if c.ClientPlayerUnit() != nil || gameGetPlayState() == 3 || isMultiplayer() {
		return true
	}

	return false
}

func (c *Client) initSeat(sz image.Point) error {
	if sdl.SetAndroidOverlayVisible != nil {
		sdl.SetAndroidOverlayVisible(true)
	}
	var lastMenuState bool = true
	sdl.CheckInMenu = func() bool {
		inMenu := c == nil || !c.isGameHUDVisible()
		if inMenu != lastMenuState {
			lastMenuState = inMenu
			c.Log.Info("touch overlay mode changed", "inMenu", inMenu, "renderGUI", nox_client_getRenderGUI(), "mainMenuFlag", nox_client_gui_flag_815132)
			if sdl.ResetAndroidTouches != nil {
				sdl.ResetAndroidTouches()
			}
		}
		return inMenu
	}
	sdl.CheckInInventory = func() bool {
		if c == nil {
			return false
		}
		// 1. Inventory open (sliding open or fully open)
		if b := memmap.Uint8(0x5D4594, 1049868); b == 1 || b == 2 {
			return true
		}
		// 2. Dragging item
		if c.dragndropGetItem() != nil {
			return true
		}
		// 3. GUI cursor active
		if nox_xxx_guiCursor_477600() != 0 {
			return true
		}
		// 4. Trader / Shop dialog
		if legacy.Sub_478030() != 0 {
			return true
		}
		return false
	}
	sst, err := sdl.New(c.Log, "OpenNox "+version.ClientVersion(), sz)
	if err != nil {
		return err
	}
	c.Seat = sst
	if sz := c.Seat.ScreenSize(); sz.X > 0 && sz.Y > 0 {
		c.videoSetGameMode(sz)
	}
	if env.IsE2E() {
		c.Seat = e2eWrapSeat(c.Seat)
	}
	c.Win, err = render.New(c.Seat)
	if err != nil {
		_ = c.Seat.Close()
		return err
	}

	inp := input.New(c.Log, c.Seat, false, c.Strings().Lang())
	c.Inp = inp

	inp.OnQuit(mainloopStop)
	inp.OnToggleFullScreen(c.Win.ToggleWindowMode)
	inp.OnKeyPress(gameexOnKeyboardPress)
	inp.OnMouseWheel(func(dv int) {
		// mix event handler is triggered only for wheel events
		call_OnLibraryNotice_265(dv)
	})
	inp.OnInputString(func(str string) {
		for _, r := range str {
			if r == '\n' || r == '\r' {
				legacy.NoxInputOnChar(10)
				continue
			}
			var c uint16
			if r >= 0x0410 && r <= 0x044F {
				// Convert Cyrillic Unicode 'А'..'я' to Windows-1251 (0xC0..0xFF) for legacy Nox
				c = uint16(0xC0 + (r - 0x0410))
			} else if r == 0x0401 { // 'Ё'
				c = 0xA8
			} else if r == 0x0451 { // 'ё'
				c = 0xB8
			} else {
				c = uint16(r)
			}
			legacy.NoxInputOnChar(c)
		}
	})

	c.Win.OnViewResize(inp.SetWinSize)
	OnPixBufferResize(inp.SetDrawWinSize)

	c.Win.SetFiltering(viper.GetBool(configVideoFiltering))
	c.Win.SetStretched(viper.GetBool(configVideoStretch))
	if err != nil {
		return err
	}
	sst.SetGamma(getGamma())
	return nil
}

func (c *Client) freeSeat() {
	if c.Seat != nil {
		c.Seat.Close()
		c.Seat = nil
	}
}

func errorMessage(format string, args ...any) {
	dialog.Error(format, args...)
}
