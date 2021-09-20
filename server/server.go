package server

import (
	"fmt"
	"github.com/blackjack200/gracticerssential/util"
	"github.com/df-mc/dragonfly/server"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

var _global *server.Server

func Global() *server.Server {
	return _global
}

func Setup(l *logrus.Logger) error {
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
