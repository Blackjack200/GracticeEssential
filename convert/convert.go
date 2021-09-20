package convert

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/world"
)

func ParseGameMode(v string) (world.GameMode, error) {
	switch v {
	case "0", "s", "survival":
		return world.GameModeSurvival{}, nil
	case "1", "c", "creative":
		return world.GameModeCreative{}, nil
	case "2", "a", "adventure":
		return world.GameModeAdventure{}, nil
	case "3", "spectator":
		return world.GameModeSpectator{}, nil
	}
	return nil, fmt.Errorf("unknown %v", v)
}

func DumpGameMode(g world.GameMode) (string, error) {
	switch g.(type) {
	case world.GameModeSurvival:
		return "survival", nil
	case world.GameModeCreative:
		return "creative", nil
	case world.GameModeAdventure:
		return "adventure", nil
	case world.GameModeSpectator:
		return "spectator", nil
	default:
		return "", fmt.Errorf("unknown %T", g)
	}
}
