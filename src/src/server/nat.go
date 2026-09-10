package server

import (
	"context"

	"github.com/opennox/libs/env"
	"github.com/opennox/nat"

	noxflags "github.com/opennox/opennox/v1/common/flags"
)

type natService struct {
	stop func()
}

func (s *Server) startNAT() error {
	if !s.UseNAT || !noxflags.HasGame(noxflags.GameOnline) || env.IsE2E() {
		return nil
	}
	port, hport := s.ServerPort(), s.HTTPPort()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop, _ := nat.Map(ctx, []nat.Port{
			{Port: port, Proto: "udp", Desc: "Nox game port"},
			{Port: hport, Proto: "tcp", Desc: "Nox HTTP port"},
		})
		if stop != nil {
			<-ctx.Done()
			stop()
		}
	}()
	s.nat.stop = cancel
	return nil
}

func (s *Server) stopNAT() {
	if s.nat.stop != nil {
		s.nat.stop()
		s.nat.stop = nil
	}
}
