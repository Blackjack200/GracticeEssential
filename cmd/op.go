package cmd

import (
	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/df-mc/dragonfly/server/cmd"
)

type Op struct {
	Target string
}

func (b Op) Run(src cmd.Source, o *cmd.Output) {
	defer o.Messages()
	if permission.IsOperator(src.Name()) {
		if b.Target == "" {
			o.Error("Command argument error")
			return
		}
		if t, found := server.Global().PlayerByName(b.Target); found {
			op := &cmd.Output{}
			op.Print("You have been opped")
			t.SendCommandOutput(op)
		}
		permission.SetOperator(b.Target)
		o.Printf("Opped: %v", b.Target)
	} else {
		o.Error("You are not operator")
	}
}

type DeOp struct {
	Target string
}

func (b DeOp) Run(src cmd.Source, o *cmd.Output) {
	if permission.IsOperator(src.Name()) {
		if b.Target == "" {
			o.Error("Command argument error")
			return
		}
		permission.RemoveOperator(b.Target)
		o.Printf("De-opped: %v", b.Target)
	} else {
		o.Error("You are not operator")
	}
}
