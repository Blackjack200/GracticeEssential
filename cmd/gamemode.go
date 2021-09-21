package cmd

import (
	"github.com/Blackjack200/GracticeEssential/convert"
	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type GameMode struct {
	GameMode string
}

func (g GameMode) Run(src cmd.Source, o *cmd.Output) {
	if permission.OpEntry().Has(src.Name()) {
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
