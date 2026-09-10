package discover

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/netip"
	"time"

	"github.com/opennox/lobby"
	"github.com/opennox/xwis"
)

func init() {
	xwisRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	xwisName := fmt.Sprintf("jack%06x", xwisRand.Intn(0x1000000))
	const backend = "xwis"
	RegisterFallback("xwis", func(ctx context.Context, log *slog.Logger, out chan<- Server) error {
		log.Info("using name", "user", xwisName)
		cli, err := xwis.NewClient(ctx, xwisName, xwisName)
		if err != nil {
			return err
		}
		defer cli.Close()

		rooms, err := cli.ListRooms(ctx)
		if err != nil {
			return err
		}

		for _, r := range rooms {
			if r.Game == nil {
				continue
			}
			log := log.With("addr", r.Game.Addr, "name", r.Game.Name)
			ip, err := netip.ParseAddr(r.Game.Addr)
			if err != nil {
				log.Warn("cannot parse", "err", err)
				continue
			}
			log.Info("result")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- Server{
				IP:       ip,
				Source:   backend,
				Priority: priorityXWIS,
				NoPing:   true,
				Game:     *lobby.GameFromXWIS(r.Game),
			}:
			}
		}
		return nil
	})
}
