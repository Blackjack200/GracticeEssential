package cmd

import (
	"github.com/Blackjack200/GracticeEssential/convert"
	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/df-mc/dragonfly/server/cmd"
)

type Difficulty struct {
	Diff string
}

func (d Difficulty) Run(src cmd.Source, o *cmd.Output) {
	if permission.OpEntry().Has(src.Name()) {
		if di, err := convert.ParseDifficulty(d.Diff); err != nil {
			o.Error(err)
		} else {
			server.Global().World().SetDifficulty(di)
			o.Printf("Set game difficulty to %v", convert.MustString(convert.DumpDifficulty(di)))
		}
	} else {
		o.Error("You are not operator")
	}
}
