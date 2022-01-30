package main

import (
	"github.com/Blackjack200/GracticeEssential/cmd"
	"github.com/Blackjack200/GracticeEssential/console"
	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/Blackjack200/GracticeEssential/util"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	if err := server.Setup(log); err != nil {
		logrus.Fatal(err)
	}

	chat.Global.Subscribe(&util.LoggerSubscriber{Logger: log})
	cmd.Setup()
	c := console.Setup(log)
	c.Run()
	server.CloseOnProgramEnd(log, func() {
		c.Stop()
	})
	if err := server.Global().Start(); err != nil {
		logrus.Fatal(err)
	}

	server.Global().Allow(permission.BanEntry().ServerAllower("You are banned", false))
	server.Loop(func(p *player.Player) {
	}, func() {
	})
}
