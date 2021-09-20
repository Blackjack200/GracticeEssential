package cmd

import (
	"github.com/blackjack200/gracticerssential/permission"
	"github.com/blackjack200/gracticerssential/server"
	"github.com/df-mc/dragonfly/server/cmd"
)

type Ban struct {
	Target string
}

func (b Ban) Run(src cmd.Source, o *cmd.Output) {
	defer o.Messages()
	if permission.IsOperator(src.Name()) {
		if b.Target == "" {
			o.Error("Command argument error")
			return
		}
		if t, found := server.Global().PlayerByName(b.Target); found {
			t.Disconnect("Banned by admin")
		}
		permission.SetBanned(b.Target)
		o.Printf("Banned player %v", b.Target)
	} else {
		o.Error("You are not operator")
	}
}

type Unban struct {
	Target string
}

func (b Unban) Run(src cmd.Source, o *cmd.Output) {
	if permission.IsOperator(src.Name()) {
		if b.Target == "" {
			o.Error("Command argument error")
			return
		}
		permission.Unban(b.Target)
		o.Printf("Unbanned player %v", b.Target)
	} else {
		o.Error("You are not operator")
	}
}
