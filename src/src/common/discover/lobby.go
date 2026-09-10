package discover

import (
	"context"
	"log/slog"
	"net/netip"

	"github.com/opennox/lobby"

	"github.com/opennox/opennox/v1/internal/version"
)

var (
	// LobbyServer is an address of the Nox lobby server used for discovery.
	LobbyServer = "http://nox.nwca.xyz:8088"
)

func newLobbyClient() *lobby.Client {
	cli := lobby.NewClient(LobbyServer)
	cli.SetUserAgent("OpenNox/" + version.Version())
	return cli
}

func init() {
	const backend = "lobby"
	RegisterBackend(backend, func(ctx context.Context, log *slog.Logger, out chan<- Server) error {
		cli := newLobbyClient()
		rooms, err := cli.ListGames(ctx)
		if err != nil {
			return err
		}
		for _, g := range rooms {
			log := log.With("addr", g.Address, "name", g.Name)
			ip, err := netip.ParseAddr(g.Address)
			if err != nil {
				log.Error("cannot parse address", "err", err)
				continue
			}
			g := g.Game
			log.Info("result")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- Server{
				IP:       ip,
				Source:   backend,
				Priority: priorityLobby,
				Game:     g,
			}:
			}
		}
		return nil
	})
}
