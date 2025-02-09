package cmd

import (
	"github.com/Blackjack200/GracticeEssential/convert"
	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

type Difficulty struct {
	Diff string
}

func (d Difficulty) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if di, err := convert.ParseDifficulty(d.Diff); err != nil {
		o.Error(err)
	} else {
		server.Global().World().SetDifficulty(di)
		o.Printf("Set game difficulty to %v", convert.MustString(convert.DumpDifficulty(di)))
	}
}

func (d Difficulty) Allow(s cmd.Source) bool {
	return AllowImpl(s)
}
