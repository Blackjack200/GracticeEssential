package cmd

import (
	"github.com/df-mc/dragonfly/server/world"
	"sort"
	"strings"

	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/df-mc/dragonfly/server/cmd"
)

type List struct{}

func (List) Run(_ cmd.Source, o *cmd.Output, tx *world.Tx) {
	players := server.Global().Players(nil)
	var names []string
	i := 0
	for p := range players {
		names[i] = p.Name()
		i++
	}
	sort.Strings(names)
	o.Printf("There are %v/%v players online:", len(names), server.Global().MaxPlayerCount())
	o.Print(strings.Join(names, ", "))
}
