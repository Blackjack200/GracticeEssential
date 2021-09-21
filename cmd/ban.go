package cmd

import (
	"sort"
	"strings"

	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/Blackjack200/GracticeEssential/server"
	"github.com/df-mc/dragonfly/server/cmd"
)

type Ban struct {
	Target string
}

func (b Ban) Run(src cmd.Source, o *cmd.Output) {
	defer o.Messages()
	if permission.OpEntry().Has(src.Name()) {
		if b.Target == "" {
			o.Error("Command argument error")
			return
		}
		if t, found := server.Global().PlayerByName(b.Target); found {
			t.Disconnect("Banned by admin")
		}
		permission.BanEntry().Add(b.Target)
		o.Printf("Banned player %v", b.Target)
	} else {
		o.Error("You are not operator")
	}
}

type Unban struct {
	Target string
}

func (b Unban) Run(src cmd.Source, o *cmd.Output) {
	if permission.OpEntry().Has(src.Name()) {
		if b.Target == "" {
			o.Error("Command argument error")
			return
		}
		permission.BanEntry().Delete(b.Target)
		o.Printf("Unbanned player %v", b.Target)
	} else {
		o.Error("You are not operator")
	}
}

type BanList struct {
}

func (BanList) Run(src cmd.Source, o *cmd.Output) {
	if permission.OpEntry().Has(src.Name()) {
		arr := permission.BanEntry().GetAll()
		sort.Strings(arr)
		o.Printf("There are %v total banned players:", len(arr))
		o.Print(strings.Join(arr, ", "))
	} else {
		o.Error("You are not operator")
	}
}
