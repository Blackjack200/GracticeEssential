package cmd

import (
	"github.com/df-mc/dragonfly/server/world"
	"runtime"

	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type Version struct{}

func (Version) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	o.Printf("This server is running %v", "dragonfly")
	o.Printf("Server version: %v", server.Version())
	o.Printf("Compatible Minecraft version: %v (protocol version: %v)", protocol.CurrentVersion, protocol.CurrentProtocol)
	o.Printf("Golang version: %v", runtime.Version())
	o.Printf("Compiler: %v", runtime.Compiler)
	o.Printf("ARCH/GOODS: %v/%v", runtime.GOARCH, runtime.GOOS)
}
