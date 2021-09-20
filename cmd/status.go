package cmd

import (
	"github.com/blackjack200/gracticerssential/permission"
	"github.com/blackjack200/gracticerssential/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"runtime"
)

type Status struct{}

func (Status) Run(src cmd.Source, o *cmd.Output) {
	if permission.IsOperator(src.Name()) {
		stat := getMemStats()
		o.Printf("Uptime: %v", server.Global().Uptime().String())
		o.Printf("Goroutine Count: %v", runtime.NumGoroutine())
		o.Printf("Allocated Memory: %dMB", stat.Sys/1024/1024)
		o.Printf("Virtual Memory: %dMB", stat.HeapSys/1024/1024)
		o.Printf("Stack Memory: %dMB", stat.StackSys/1024/1024)
		o.Printf("Heap Object: %d", (stat.Mallocs-stat.Frees)/1024/1024)
		o.Printf("GC cycles: %d", stat.NumGC)
	} else {
		o.Error("You are not operator")
	}
}

func getMemStats() runtime.MemStats {
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	return m2
}
