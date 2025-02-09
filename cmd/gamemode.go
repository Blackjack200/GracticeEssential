package cmd

import (
	"github.com/Blackjack200/GracticeEssential/convert"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type GameMode struct {
	GameMode string
}

func (g GameMode) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if AllowImpl(src) {
		if p, ok := src.(*player.Player); ok {
			mode, err := convert.ParseGameMode(g.GameMode)
			if err != nil {
				o.Error(err)
				return
			}
			p.SetGameMode(mode)
			o.Printf("Set game mode to %v", convert.MustString(convert.DumpGameMode(mode)))
		} else {
			o.Error("This command must use in game")
		}
	} else {
		o.Error("You are not operator")
	}
}

func (g GameMode) Allow(s cmd.Source) bool {
	return AllowImpl(s)
}
