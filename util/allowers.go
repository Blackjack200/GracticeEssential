package util

import (
	"net"

	"github.com/df-mc/dragonfly/server"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

type serverAllower struct {
	parent  server.Allower
	current server.Allower
}

func (s *serverAllower) Allow(addr net.Addr, i login.IdentityData, c login.ClientData) (string, bool) {
	curt := s.current
	if str, ok := curt.Allow(addr, i, c); !ok {
		return str, ok
	}
	curt = s.parent
	for curt != nil {
		if str, ok := curt.Allow(addr, i, c); !ok {
			return str, ok
		}
		if n, ok := curt.(*serverAllower); ok {
			curt = n
		} else {
			break
		}
	}
	return "", true
}

func ServerAllower(parent, current server.Allower) server.Allower {
	return &serverAllower{parent: parent, current: current}
}
