package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/Blackjack200/GracticeEssential/util"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

var _global *server.Server
var _log *logrus.Logger

func Global() *server.Server {
	return _global
}

func Setup(l *logrus.Logger) error {
	_log = l
	if cfg, err := readConfig(); err != nil {
		return err
	} else {
		_global = server.New(&cfg, l)
	}
	return nil
}

func readConfig() (server.Config, error) {
	c := server.DefaultConfig()
	if !util.FileExist("config.toml") {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}

func Loop(h func(p *player.Player), end func()) {
	for {
		if p, err := Global().Accept(); err != nil {
			break
		} else {
			h(p)
		}
		end()
	}
}

func CloseOnProgramEnd(log *logrus.Logger, f func()) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func(fn func()) {
		<-c
		if err := Global().Close(); err != nil {
			log.Errorf("error shutting down server: %v", err)
		}
		fn()
	}(f)
}
