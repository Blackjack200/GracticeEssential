package cmd

import (
	"github.com/Blackjack200/GracticeEssential/permission"
	"github.com/df-mc/dragonfly/server/cmd"
)

var AllowImpl = func(s cmd.Source) bool {
	if t, ok := s.(cmd.NamedTarget); ok {
		return permission.OpEntry().Has(t.Name())
	}
	return false
}

func Setup() {
	cmd.Register(cmd.New("help", "Provides help/list of commands.", []string{"?"}, Help{}))

	cmd.Register(cmd.New("version", "Gets the version of this server in use.", []string{"ver", "about"}, Version{}))
	cmd.Register(cmd.New("status", "Reads back the server's performance.", []string{"stat"}, Status{}))
	cmd.Register(cmd.New("list", "Lists all online players", nil, List{}))
	cmd.Register(cmd.New("gc", "Fires garbage collection tasks.", nil, GC{}))
	cmd.Register(cmd.New("stop", "Stops the server.", nil, Stop{}))

	cmd.Register(cmd.New("op", "Grants operator status to a player.", nil, Op{}))
	cmd.Register(cmd.New("deop", "Revokes operator status from a player.", nil, DeOp{}))

	cmd.Register(cmd.New("banlist", "View all players banned from this server", nil, BanList{}))
	cmd.Register(cmd.New("ban", "Adds player to banlist.", nil, Ban{}))
	cmd.Register(cmd.New("unban", "Removes player from banlist.", nil, Unban{}))
	cmd.Register(cmd.New("kick", "Kicks a player from the server.", nil, Kick{}))

	cmd.Register(cmd.New("difficulty", "Sets the game difficulty", nil, Difficulty{}))
	cmd.Register(cmd.New("defaultgamemode", "Sets the default game mode.", nil, DefaultGameMode{}))
	cmd.Register(cmd.New("gamemode", "Sets your game mode.", []string{"gm"}, GameMode{}))

	cmd.Register(cmd.New("setworldspawn", "Sets the world spawn.", nil, SetWorldSpawn{}))
}
