package discover

import (
	"context"
	"log/slog"
	"math/rand"
	"net"
	"time"

	"github.com/opennox/libs/common"
)

const (
	lanTimeout = 400 * time.Millisecond
)

func init() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	const backend = "lan"
	RegisterBackend(backend, func(ctx context.Context, log *slog.Logger, out chan<- Server) error {
		if deadline, ok := ctx.Deadline(); ok {
			if dt := time.Until(deadline); dt > lanTimeout {
				var cancel func()
				ctx, cancel = context.WithTimeout(ctx, lanTimeout)
				defer cancel()
			}
		}
		// we don't care about our own port here, because we don't need to port-forward for LAN
		l, err := net.ListenUDP("udp4", &net.UDPAddr{IP: nil, Port: 0})
		if err != nil {
			return err
		}
		defer l.Close()

		token := rnd.Uint32()
		if _, err := l.WriteTo(encodeGameDiscovery(token), &net.UDPAddr{
			IP:   net.IPv4(255, 255, 255, 255), // multicast
			Port: common.GamePort,
		}); err != nil {
			return err
		}
		start := time.Now()
		errc := make(chan error, 1)
		go func() {
			defer close(errc)
			buf := make([]byte, 256)
			for {
				buf = buf[:cap(buf)]
				n, addr, err := l.ReadFromUDPAddrPort(buf)
				if err != nil {
					errc <- err
					return
				}
				buf = buf[:n]
				m := decodeGameInfo(buf)
				if m == nil || m.Token != token {
					continue
				}
				if !addr.Addr().Is4() {
					log.Warn("skipping IPv6 response", "addr", addr)
					continue
				}
				if g := convGameInfo(addr, m, buf); g != nil {
					select {
					case <-ctx.Done():
						return
					case out <- Server{
						IP:       addr.Addr(),
						Source:   backend,
						Priority: priorityLAN,
						Ping:     time.Since(start),
						NoPing:   true, // already pinged
						Game:     *g,
					}:
					}
				}
			}
		}()
		select {
		case <-ctx.Done():
			_ = l.Close()
			<-errc
			return ctx.Err()
		case err := <-errc:
			return err
		}
	})
}
