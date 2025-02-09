package cmd

import (
	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

type Stop struct{}

func (Stop) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	out := &cmd.Output{}
	out.Print("Stopping the server")
	for p := range server.Global().Players(nil) {
		p.SendCommandOutput(out)
	}
	go server.Global().Close()
}

func (Stop) Allow(s cmd.Source) bool {
	return AllowImpl(s)
}
