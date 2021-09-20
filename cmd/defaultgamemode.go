package cmd

import (
	"github.com/Blackjack200/GracticeEssential/convert"
	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/df-mc/dragonfly/server/cmd"
)

type DefaultGameMode struct {
	GameMode string
}

func (d DefaultGameMode) Run(src cmd.Source, o *cmd.Output) {
	if permission.IsOperator(src.Name()) {
		mode, err := convert.ParseGameMode(d.GameMode)
		if err != nil {
			o.Error(err)
			return
		}
		str, _ := convert.DumpGameMode(mode)
		server.Global().World().SetDefaultGameMode(mode)
		o.Printf("Set default game mode to %v", str)
	} else {
		o.Error("You are not operator")
	}
}
