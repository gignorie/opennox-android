package opennox

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/opennox/libs/console"
	"github.com/opennox/libs/datapath"
	"github.com/opennox/libs/ifs"
	"github.com/opennox/libs/maps"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

var noxMapsIgnoreMode = false

var noxCmdSetMaps = &console.Command{
	Token: "maps",
	Help:  "set map-related variables",
	Flags: console.Server | console.Cheat,
}

func init() {
	noxConsole.Register(&console.Command{
		Token:  "load",
		HelpID: "loadhelp",
		Flags:  console.Server | console.Cheat,
		Func:   noxCmdLoad,
	})
	noxCmdSet.Register(noxCmdSetMaps)
	noxCmdSetMaps.Register(&console.Command{
		Token: "allow.all",
		Help:  "ignore all map mode checks when using load command",
		Flags: console.Server | console.Cheat,
		Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			return noxCmdSetBool(c, tokens, func(allow bool) {
				noxMapsIgnoreMode = allow
				if allow {
					c.Print(console.ColorLightYellow, "any map can be loaded now")
					c.Print(console.ColorLightRed, "the game and the menu may act weirdly in this mode")
				} else {
					c.Print(console.ColorLightYellow, "maps loading will check game mode")
				}
			})
		},
	})
}

func noxCmdLoad(ctx context.Context, c *console.Console, tokens []string) bool {
	if len(tokens) != 1 {
		return false
	}
	if err := noxLoadMap(tokens[0], mapLoadOptions{
		Force: noxMapsIgnoreMode,
	}); err != nil {
		c.Print(console.ColorRed, err.Error())
	}
	return true
}

type mapLoadOptions struct {
	Force bool
}

func noxLoadMap(name string, opts mapLoadOptions) error {
	ignoreMode := opts.Force
	if len(name) != 0 && name[0] != '#' {
		if err := nox_common_checkMapFile(name); err != nil {
			return fmt.Errorf("Error checking map file %q: %w", name, err)
		}
	}
	mode := nox_xxx_mapGetTypeMB_4CFFA0(memmap.PtrOff(0x973F18, 2408))
	if noxflags.HasGame(noxflags.GameOnline) {
		if !ignoreMode && (mode == 0 || mode.Has(noxflags.GameModeCoopTeam)) {
			return errors.New("Switching maps to Solo is not allowed")
		}
		if noxflags.HasGame(noxflags.GameModeChat) {
			if mode.Has(noxflags.GameModeCTF|noxflags.GameModeFlagBall) && noxServer.Teams.Count() != 2 {
				legacy.Nox_xxx_wndGuiTeamCreate_4185B0()
			}
		} else if !ignoreMode && !noxflags.GetGame().Has(mode) {
			v6 := strMan.GetStringInFile("NoMapLoadNewMode", "parsecmd.c")
			nox_xxx_printCentered_445490(v6)
			return errors.New("Swithing to a different mode is not allowed")
		}
	} else {
		if !ignoreMode && !mode.Has(noxflags.GameModeCoopTeam) {
			return errors.New("Switching to non-Solo maps is not allowed")
		}
	}
	fname := name
	if !strings.HasSuffix(fname, maps.Ext) {
		fname = name + maps.Ext
	} else {
		fname = name
	}
	if fname == "" || fname[0] != '#' {
		if _, err := ifs.Stat(datapath.Maps(name, fname)); err != nil {
			str := strMan.GetStringInFile("CannotAccessMap", "parsecmd.c")
			nox_xxx_printCentered_445490(fmt.Sprintf(str, fname))
			return fmt.Errorf("Cannot access map %q", fname)
		}
	}
	legacy.Nox_xxx_mapLoadOrSaveMB_4DCC70(1)
	noxServer.SwitchMap(fname)
	legacy.Sub_41D650()
	str := strMan.GetStringInFile("maploaded", "parsecmd.c")
	str = strings.ReplaceAll(str, "%S", "%s")
	nox_xxx_printCentered_445490(fmt.Sprintf(str, fname))
	return nil
}
