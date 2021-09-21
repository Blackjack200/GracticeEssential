package cmd

import (
	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type SetWorldSpawn struct{}

func (SetWorldSpawn) Run(src cmd.Source, o *cmd.Output) {
	if permission.OpEntry().Has(src.Name()) {
		if p, ok := src.(*player.Player); ok {
			s := cube.PosFromVec3(p.Position())
			server.Global().World().SetSpawn(s)
			o.Printf("Set the world spawn point to (%v, %v, %v)", s[0], s[1], s[2])
		} else {
			o.Error("This command must use in game")
		}
	} else {
		o.Error("You are not operator")
	}
}
