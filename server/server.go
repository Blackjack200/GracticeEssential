package server

import (
	"fmt"
	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/Blackjack200/GracticeEssential/util"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/pelletier/go-toml"
	"log/slog"
	"os"
	"time"
)

var _global *server.Server
var _startDate time.Time

func Global() *server.Server {
	return _global
}

func SetupFunc(l *slog.Logger, cfgFunc func(*server.Config)) error {
	util.PanicFunc(func(v interface{}) {
		panic(v)
	})
	if cfg, err := readConfig(); err != nil {
		return err
	} else {
		cfg := util.SelectNotNil[server.Config](cfg.Config(l))
		cfg.Allower = permission.BanEntry().ServerAllower("You are banned", false)
		if cfgFunc != nil {
			cfgFunc(&cfg)
		}
		_global = cfg.New()
	}
	return nil
}

func Start() {
	Global().Listen()
	_startDate = time.Now()
}

var Stop = func() {}

func readConfig() (server.UserConfig, error) {
	c := server.DefaultConfig()
	if !util.FileExist("config.toml") {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}

func Loop(h func(p *player.Player), end func()) {
	for p := range Global().Accept() {
		h(p)
		if end != nil {
			end()
		}
	}
}

func Uptime() time.Duration {
	return time.Now().Sub(_startDate)
}
