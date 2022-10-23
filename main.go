package main

import (
	"github.com/Blackjack200/GracticeEssential/bootstrap"
	"github.com/Blackjack200/GracticeEssential/mhandler"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sirupsen/logrus"
)

type myChatHandler struct {
	unreg func()
}

func (m myChatHandler) HandleChat(ctx *event.Context, msg *string) {
	*msg = "You can break block now"
	m.unreg()
}

type myBlockBreakHandler struct{}

func (myBlockBreakHandler) HandleBlockBreak(ctx *event.Context, _ cube.Pos, _ *[]item.Stack) {
	ctx.Cancel()
}

type myQuitHandler struct{}

func (myQuitHandler) HandleQuit() {
	logrus.Info("Quited")
}

func main() {
	log := bootstrap.NewLogger()
	bootstrap.Default(log, nil, func(p *player.Player) {
		h := mhandler.New()
		p.Handle(h)
		unreg := h.Register(myBlockBreakHandler{})
		h.Register(myChatHandler{unreg: unreg})
		h.Register(myQuitHandler{})
	}, nil)()
}
