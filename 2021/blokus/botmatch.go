package blokus

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

func sortedPlayerNames(players map[string]Player) (playerNames []string) {
	playerNames = make([]string, 0, len(players))
	for playerName := range players {
		playerNames = append(playerNames, playerName)
	}
	sort.Strings(playerNames)
	return
}

type botMatchCommandLine struct {
	RepeatGames uint
	Player1     Player
	Player2     Player
	Player1Name string
	Player2Name string
}

func (c *botMatchCommandLine) Parse(players map[string]Player) {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.UintVar(&c.RepeatGames, "n", 2, "repeat games between players (with alternating player order)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s <player 1> <player 2> [flags]\n\nAvailable players:\n", os.Args[0])
		for _, playerName := range sortedPlayerNames(players) {
			fmt.Fprintf(os.Stderr, "  - %s\n", playerName)
		}
		fmt.Fprintln(os.Stderr, "\nFlags:")
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

// BotMatchMain provides the main function for a bot match command line tool that can be customized by
// setting the available players. This way you can compile it to contain multiple variants of your player
// implementation, that you then can let play against each other.
// You can invoke the executable like this: <name of executable> <name of first player> <name of second player>
// By default, two games will be played. You can change this by adding the desired number of games with -n <number>
// The first player position will alternate between the two players.
// Results will be printed to stdout.
func BotMatchMain(players map[string]Player) {
	var cl botMatchCommandLine
	cl.Parse(players)
	result := RunRepeatedGames(cl.Player1, cl.Player2, cl.Player1Name, cl.Player2Name, cl.RepeatGames, os.Stdout)
	result.Print(os.Stdout)
	os.Stdout.Sync()
}
