package permission

import (
	"github.com/blackjack200/gracticerssential/util"
	"path/filepath"
	"strings"
	"sync"
)

var _mu = sync.Mutex{}
var _banFile = filepath.Join(util.WorkingPath, "banned-players.txt")
var _banList = loadBanList()

func IsBanned(n string) bool {
	for _, l := range _banList {
		if l == n {
			return true
		}
	}
	return false
}

func SetBanned(n string) {
	if !IsBanned(n) {
		_mu.Lock()
		_banList = append(_banList, n)
		writeBanList()
		_mu.Unlock()
	}
}

func Unban(n string) {
	if IsBanned(n) {
		_mu.Lock()
		var a []string
		for _, l := range _banList {
			if l != n {
				a = append(a, l)
			}
		}
		_banList = a
		writeBanList()
		_mu.Unlock()
	}
}

func writeBanList() {
	util.WriteFile(_banFile, []byte(strings.Join(_banList, "\n")))
}

func loadBanList() []string {
	if !util.FileExist(_banFile) {
		util.WriteFile(_banFile, nil)
	}
	return strings.Split(string(util.ReadFile(_banFile)), "\n")
}
