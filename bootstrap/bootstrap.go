package bootstrap

import (
	"github.com/Blackjack200/GracticeEssential/cmd"
	"github.com/Blackjack200/GracticeEssential/console"
	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/Blackjack200/GracticeEssential/util"
	df "github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func signalHandler(log *slog.Logger, callback func()) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func(fn func()) {
		<-c
		if err := server.Global().Close(); err != nil {
			log.Error("error shutting down server: %v", err)
		}
		if fn != nil {
			fn()
		}
	}(callback)
}

func Bootstrap(log *slog.Logger, cfgFunc func(config *df.Config), playerFunc func(*player.Player), end func(), s chat.Subscriber) (startFunc func()) {
	if err := server.SetupFunc(log, cfgFunc); err != nil {
		panic(err)
	}
	chat.Global.Subscribe(s)
	cmd.Setup()
	c := console.Setup(log)
	c.Run()
	signalHandler(log, func() {
		c.Stop()
		server.Stop()
	})
	startFunc = func() {
		server.Start()
		server.Loop(playerFunc, end)
	}
	return startFunc
}

func NewLogger() *slog.Logger {
	log := slog.Default()
	return log
}

func Default(log *slog.Logger, cfgFunc func(config *df.Config), playerFunc func(*player.Player), end func()) (startFunc func()) {
	return Bootstrap(log, cfgFunc, playerFunc, end, &util.LoggerSubscriber{Logger: log, Uuid: util.SelectNotNil[uuid.UUID](uuid.NewRandom())})
}
