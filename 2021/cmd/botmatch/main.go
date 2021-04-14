package main

import (
	"flag"
	"fmt"
	"github.com/hschendel/sc/2021/blokus"
	"github.com/hschendel/sc/2021/blokus/example_players"
	"os"
	"sort"
)

var players = map[string]blokus.Player {
	"quick": new(example_players.QuickPlayer),
	"random": new(example_players.RandomPlayer),
	"restrict": new(example_players.RestrictingPlayer),
}

var playerNames []string

func init() {
	playerNames = make([]string, 0, len(players))
	for playerName := range players {
		playerNames = append(playerNames, playerName)
	}
	sort.Strings(playerNames)
}

type commandLine struct {
	RepeatGames uint
	Player1 blokus.Player
	Player2 blokus.Player
	Player1Name string
	Player2Name string
}

func (c *commandLine) Parse() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.UintVar(&c.RepeatGames, "r", 2, "repeat games between players (with alternating player order)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s <player 1> <player 2> [flags]\n\nAvailable players:\n", os.Args[0])
		for _, playerName := range playerNames {
			fmt.Fprintf(os.Stderr,"  - %s\n", playerName)
		}
		fmt.Fprintln(os.Stderr,"\nFlags:")
		fs.PrintDefaults()
		os.Stderr.Sync()
	}
	if len(os.Args) < 3 {
		fs.Usage()
		os.Exit(1)
	}
	c.Player1Name = os.Args[1]
	c.Player2Name = os.Args[2]
	c.Player1 = players[c.Player1Name]
	c.Player2 = players[c.Player2Name]
	if c.Player1 == nil || c.Player2 == nil {
		flag.Usage()
		os.Exit(1)
	}
	if err := fs.Parse(os.Args[3:]); err != nil {
		fmt.Fprintln(fs.Output(), err)
		os.Exit(1)
	}
}

func main() {
	var cl commandLine
	cl.Parse()
	result := blokus.RunRepeatedGames(cl.Player1, cl.Player2, cl.Player1Name, cl.Player2Name, cl.RepeatGames, os.Stdout)
	result.Print(os.Stdout)
	os.Stdout.Sync()
}
