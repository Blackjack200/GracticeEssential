package cmd

import (
	"github.com/df-mc/dragonfly/server/world"
	"runtime"

	"github.com/df-mc/dragonfly/server/cmd"
)

type GC struct{}

func (GC) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if AllowImpl(src) {
		a, b := gc()
		o.Printf("Allocated Memory freed: %v MB", (b.Sys-a.Sys)/1024/1024)
	} else {
		o.Error("You are not operator")
	}
}

func (GC) Allow(s cmd.Source) bool {
	return AllowImpl(s)
}

func gc() (runtime.MemStats, runtime.MemStats) {
	var m runtime.MemStats
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	return m, m2
}
