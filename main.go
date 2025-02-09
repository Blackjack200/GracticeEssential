package main

import (
	"github.com/Blackjack200/GracticeEssential/bootstrap"
	"github.com/Blackjack200/GracticeEssential/mhandler"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"log/slog"
)

type myChatHandler struct {
	unreg func()
}

func (m myChatHandler) HandleChat(ctx *event.Context[*player.Player], message *string) {
	ctx.Val().Message("You can break block now")
	m.unreg()
}

type myBlockBreakHandler struct{}

func (m myBlockBreakHandler) HandleBlockBreak(ctx *event.Context[*player.Player], pos cube.Pos, drops *[]item.Stack, xp *int) {
	ctx.Cancel()
}

type myQuitHandler struct{}

func (m myQuitHandler) HandleQuit(p *player.Player) {
	slog.Info("Quited")
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
