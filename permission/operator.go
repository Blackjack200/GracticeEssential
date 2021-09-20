package permission

import (
	"github.com/blackjack200/gracticerssential/util"
	"path/filepath"
	"strings"
)

var _opFile = filepath.Join(util.WorkingPath, "ops.txt")
var _opList = loadOperators()

func IsOperator(n string) bool {
	for _, l := range _opList {
		if l == n {
			return true
		}
	}
	return false
}

func SetOperator(n string) {
	if !IsOperator(n) {
		_mu.Lock()
		_opList = append(_opList, n)
		writeOperators()
		_mu.Unlock()
	}
}

func RemoveOperator(n string) {
	if IsOperator(n) {
		_mu.Lock()
		var a []string
		for _, l := range _opList {
			if l != n {
				a = append(a, l)
			}
		}
		_opList = a
		writeOperators()
		_mu.Unlock()
	}
}

func writeOperators() {
	util.WriteFile(_opFile, []byte(strings.Join(_opList, "\n")))
}

func loadOperators() []string {
	if !util.FileExist(_opFile) {
		util.WriteFile(_opFile, nil)
	}
	return strings.Split(string(util.ReadFile(_opFile)), "\n")
}
