package permission

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"net"
)

func (e *Entry) ServerAllower(msg string, detectHas bool) server.Allower {
	return entryServerAllower{
		e:         e,
		detectHas: detectHas,
		msg:       msg,
	}
}

func (e *Entry) CmdAllower(detectHas bool) cmd.Allower {
	return entryCmdAllower{
		e:         e,
		detectHas: detectHas,
	}
}

type entryServerAllower struct {
	e         *Entry
	detectHas bool
	msg       string
}

func (e entryServerAllower) Allow(_ net.Addr, d login.IdentityData, _ login.ClientData) (string, bool) {
	if e.detectHas {
		return e.msg, e.e.Has(d.DisplayName)
	}
	return e.msg, !e.e.Has(d.DisplayName)
}

type entryCmdAllower struct {
	e         *Entry
	detectHas bool
}

func (e entryCmdAllower) Allow(s cmd.Source) bool {
	if e.detectHas {
		return e.e.Has(s.Name())
	}
	return !e.e.Has(s.Name())
}
