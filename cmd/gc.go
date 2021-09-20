package cmd

import (
	"github.com/blackjack200/gracticerssential/permission"
	"github.com/df-mc/dragonfly/server/cmd"
	"runtime"
)

type GC struct{}

func (GC) Run(src cmd.Source, o *cmd.Output) {
	if permission.IsOperator(src.Name()) {
		a, b := gc()
		o.Printf("Allocated Memory freed: %d MB", (b.Sys-a.Sys)/1024/1024)
		o.Printf("Virtual Memory freed: %d MB", (b.HeapSys-a.HeapSys)/1024/1024)
	} else {
		o.Error("You are not operator")
	}
}

func gc() (after runtime.MemStats, before runtime.MemStats) {
	var m runtime.MemStats
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	return before, m2
}
