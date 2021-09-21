package main

import (
	"github.com/Blackjack200/GracticeEssential/cmd"
	"github.com/Blackjack200/GracticeEssential/console"
	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	if err := server.Setup(log); err != nil {
		logrus.Fatal(err)
	}
	cmd.Setup()
	c := console.Setup(log)
	c.Run()
	server.CloseOnProgramEnd(log, func() {
		c.Stop()
	})
	if err := server.Global().Start(); err != nil {
		logrus.Fatal(err)
	}
	server.Loop(func(p *player.Player) {
		if permission.BanEntry().Has(p.Name()) {
			p.Disconnect("You are banned")
		}
	}, func() {
	})
}
