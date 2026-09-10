package opennox

import (
	"context"
	"errors"
	"log/slog"
	"time"

	noxlog "github.com/opennox/libs/log"
	"github.com/opennox/lobby"

	"github.com/opennox/opennox/v1/common/discover"
	noxflags "github.com/opennox/opennox/v1/common/flags"
)

func init() {
	configStrPtr("network.lobby.address", "NOX_LOBBY_ADDR", discover.LobbyServer, &discover.LobbyServer)
}

type LobbyServerFunc func(s *LobbyServerInfo)

func isCtxTimeout(err error) bool {
	return err != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded))
}

// discoverAndPingServers discovers game servers and sends them to the UI.
// It should run in a goroutine, and will communicate via discoverDone channel.
func discoverAndPingServers(ctx context.Context, log *slog.Logger) {
	log = noxlog.WithSystem(log, "discover")
	start := time.Now()
	list, err := discover.ListServersWith(ctx, log, lobbyBroadcast)
	if err != nil && !isCtxTimeout(err) {
		log.Error("discover failed", "err", err)
	}
	select {
	case discoverDone <- list:
		log.Info("done", "dt", time.Since(start))
	default:
		log.Info("discarded", "dt", time.Since(start))
	}
}

func gameModeToFlags(mode lobby.GameMode) uint16 {
	switch mode {
	case lobby.ModeKOTR:
		return uint16(noxflags.GameModeKOTR)
	case lobby.ModeCTF:
		return uint16(noxflags.GameModeCTF)
	case lobby.ModeFlagBall:
		return uint16(noxflags.GameModeFlagBall)
	case lobby.ModeChat:
		return uint16(noxflags.GameModeChat)
	case lobby.ModeArena:
		return uint16(noxflags.GameModeArena)
	case lobby.ModeElimination:
		return uint16(noxflags.GameModeElimination)
	case lobby.ModeQuest:
		return uint16(noxflags.GameModeQuest)
	case lobby.ModeCoop:
		return uint16(noxflags.GameModeCoop)
	}
	return 0
}
