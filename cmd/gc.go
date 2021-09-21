package cmd

import (
	"runtime"

	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/df-mc/dragonfly/server/cmd"
)

type GC struct{}

func (GC) Run(src cmd.Source, o *cmd.Output) {
	if permission.OpEntry().Has(src.Name()) {
		a, b := gc()
		o.Printf("Allocated Memory freed: %v MB", (b.Sys-a.Sys)/1024/1024)
	} else {
		o.Error("You are not operator")
	}
}

func gc() (runtime.MemStats, runtime.MemStats) {
	var m runtime.MemStats
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	return m, m2
}
