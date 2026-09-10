package netstr

import (
	"errors"
	"net"
	"net/netip"

	"github.com/opennox/libs/noxnet"
	"github.com/opennox/libs/noxnet/netmsg"
)

type LobbyWaitOptions struct {
	OnResult       func(addr netip.AddrPort, data []byte)
	OnPassRequired func()
	OnPing         func(addr netip.AddrPort, data []byte)
	OnConnectErr   func(errcode noxnet.ConnectError) bool
	OnJoinOK       func()
	OnJoinFail     func()
}

func WaitForLobbyResults(conn net.PacketConn, srvAddr netip.Addr, flag RecvFlags, opts LobbyWaitOptions) (int, error) {
	argp := 0
	if flag.Has(RecvCanRead) {
		var err error
		argp, err = canReadConn(DebugSockets, Log, conn)
		if err != nil {
			return 0, err
		} else if argp == 0 {
			// TODO: is it an error at all?
			return 0, errors.New("nothing to read")
		}
	} else {
		argp = 1
	}
	buf := make([]byte, 256)
	for {
		buf = buf[:cap(buf)]
		n, from, err := readFrom(DebugSockets, Log, conn, buf)
		if err != nil {
			return 0, err
		}
		buf = buf[:n]
		op := netmsg.Op(buf[2])
		if op < netmsg.MSG_CLIENT_ACCEPT {
			if op == netmsg.MSG_SERVER_INFO || srvAddr == from.Addr() {
				switch op {
				case netmsg.MSG_SERVER_INFO:
					if from.Addr().IsValid() {
						opts.OnResult(from, buf)
					}
				case netmsg.MSG_PASSWORD_REQUIRED:
					opts.OnPassRequired()
				case netmsg.MSG_SERVER_PING:
					opts.OnPing(from, buf)
				case netmsg.MSG_SERVER_ERROR:
					if !opts.OnConnectErr(noxnet.ConnectError(buf[3])) {
						break
					}
				case netmsg.MSG_SERVER_JOIN_OK:
					opts.OnJoinOK()
				case netmsg.MSG_SERVER_JOIN_FAIL:
					opts.OnJoinFail()
				}
			}
		}
		if flag&RecvCanRead == 0 || (flag&RecvJustOne) != 0 {
			return n, nil
		}
		argp, err = canReadConn(DebugSockets, Log, conn)
		if err != nil {
			return n, err
		} else if argp == 0 {
			return n, nil
		}
	}
}
