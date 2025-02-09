package util

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/text"
)

// LoggerSubscriber is an implementation of Subscriber that forwards messages sent to the logger
type LoggerSubscriber struct {
	Logger *slog.Logger
	Uuid   uuid.UUID
}

func (c *LoggerSubscriber) UUID() uuid.UUID {
	return c.Uuid
}

// Message ...
func (c *LoggerSubscriber) Message(a ...interface{}) {
	s := make([]string, len(a))
	for i, b := range a {
		s[i] = fmt.Sprint(b)
	}
	t := text.ANSI(strings.TrimSpace(strings.Join(s, " ")) + "§r")
	c.Logger.Info(t)
}
