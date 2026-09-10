package discover

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net/netip"
	"os"
	"path/filepath"
	"strings"

	"github.com/opennox/lobby"

	"github.com/opennox/libs/common"
)

const (
	StaticFile = "game_ip.txt"
)

func init() {
	RegisterBackend("static", func(ctx context.Context, log *slog.Logger, out chan<- Server) error {
		list, err := staticIPs(log, StaticFile)
		if len(list) == 0 {
			return err
		}
		for _, s := range list {
			s.Priority = -1
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- s:
			}
		}
		return nil
	})
}

func staticIPs(log *slog.Logger, path string) ([]Server, error) {
	name := filepath.Base(path)
	log = log.With("path", name)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		log.Warn("no file")
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", name, err)
	}
	defer f.Close()

	var (
		out  []Server
		last error
	)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		addr, err := netip.ParseAddrPort(line)
		if err != nil {
			ip, err := netip.ParseAddr(line)
			if err != nil {
				last = fmt.Errorf("cannot parse server IP: %q: %v", line, err)
				log.Error("cannot parse", "err", last)
				continue
			}
			addr = netip.AddrPortFrom(ip, common.GamePort)
		}
		log.Info("result", "addr", addr)
		out = append(out, Server{
			Source:   name,
			Priority: priorityStatic,
			IP:       addr.Addr(),
			Game: lobby.Game{
				Address: addr.Addr().String(),
				Port:    int(addr.Port()),
			},
		})
	}
	if err := sc.Err(); err != nil {
		last = fmt.Errorf("error reading %s: %w", name, err)
		log.Error("read error", "err", last)
	}
	return out, last
}
