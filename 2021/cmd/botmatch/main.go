package main

import (
	"github.com/hschendel/sc/2021/blokus"
	"github.com/hschendel/sc/2021/blokus/example_players"
)

var players = map[string]blokus.Player{
	"quick":    new(example_players.QuickPlayer),
	"random":   new(example_players.RandomPlayer),
	"restrict": new(example_players.RestrictingPlayer),
}

func main() {
	blokus.BotMatchMain(players)
}
