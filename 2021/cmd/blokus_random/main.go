package main

import (
	"github.com/hschendel/sc/2021/blokus"
	"github.com/hschendel/sc/2021/blokus/example_players"
)

func main() {
	player := new(example_players.RandomPlayer)
	blokus.ClientMain(player)
}