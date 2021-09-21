package permission

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/Blackjack200/GracticeEssential/util"
)

var _mu = sync.Mutex{}
var _banFile = filepath.Join(util.WorkingPath, "banned-players.txt")
var _banList = loadBanList()

func IsBanned(n string) bool {
	if n == "CONSOLE" {
		return true
	}
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
