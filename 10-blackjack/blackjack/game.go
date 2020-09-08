package blackjack

import (
	"fmt"
	"log"

	"github.com/MichaelZalla/gophercises/09-deck/deck"
)

// Game stores all information related to the current Blackjack game
type Game struct {
	Dealer  *Player
	Players []*Player
	Deck    []deck.Card
}

// GameOptionFn allows callers to customize a game (i.e., when calling NewGame)
type GameOptionFn func(game Game) Game

const maxScore = 21

const dealerSoftHitLimit = 17

// NewGame creates a new game of Blackjack
func NewGame(options ...GameOptionFn) (*Game, []*Player) {

	// Generate a Deck and a Dealer (who serves as Player 0)

	game := Game{
		Deck: deck.New(
			deck.WithStandard(),
			deck.Shuffled()),
		Dealer: &Player{
			ID:     0,
			Dealer: true,
		},
	}

	// Add Players, etc

	for _, option := range options {
		game = option(game)
	}

	// Ensures Dealer is always the last to be dealt

	game.Players = append(game.Players, game.Dealer)

	// Deal a hand to each player

	for i := 0; i < 2; i++ {
		for j := 0; j < len(game.Players); j++ {
			deal(&game, game.Players[j], 1)
		}
	}

	// NewGame blocks until getWinners returns

	return &game, getWinners(&game)

}

// Players adds the given number of players to a game (excludes the dealer)
func Players(n int) func(game Game) Game {

	return func(game Game) Game {

		// @NOTE(mzalla) Reserves 0 for Dealer's ID
		for i := 1; i < n+1; i++ {
			game.Players = append(game.Players, &Player{
				ID: i,
			})
		}

		return game

	}

}

func deal(game *Game, player *Player, n int) {

	if len(game.Deck) < n {
		log.Fatal("Called deal() on deck with too few cards!")
	}

	for i := 0; i < n; i++ {
		front, back := game.Deck[0:n], game.Deck[n:]
		player.Hand = append(player.Hand, front...)
		game.Deck = back
	}

}

func getWinners(game *Game) []*Player {

	var winners []*Player

	// Check whether any players were dealt Blackjack (ignores the Dealer)

	for _, p := range game.Players {

		p.Score = getScore(p)

		if int(p.Score) == maxScore {
			winners = append(winners, p)
		}

	}

	if len(winners) > 0 {
		return winners
	}

	// Otherwise, each player (besides the Dealer) plays a turn

	sumPlayerScoreAfterRound := 0

	for _, p := range game.Players {
		if p.Dealer {
			continue
		}
		sumPlayerScoreAfterRound += playTurnAsPlayer(game, p)
	}

	// If all players went bust that round, Dealer wins!

	if sumPlayerScoreAfterRound == 0 {
		return []*Player{game.Dealer}
	}

	// Dealer plays a turn

	topScore := playTurnAsDealer(game, game.Dealer)

	// Determine the top score for the round

	for _, p := range game.Players {
		p.Score = getScore(p)
		if p.Score > topScore {
			topScore = p.Score
		}
	}

	// Determine which Player(s) took the top score

	for _, p := range game.Players {
		if p.Score == topScore {
			winners = append(winners, p)
		}
	}

	return winners

}

func playTurnAsPlayer(game *Game, p *Player) int {

	fmt.Printf("%s's turn...\n", p)

	defer fmt.Println()
	defer showHand(p)

	score := getScore(p)

	var choice string

	for score > 0 && score != maxScore {

		showHand(p)

		fmt.Printf("\tWill you hit (h), or stand (s)?\n\t")
		fmt.Scanln(&choice)

		switch choice {
		case "h":

			fmt.Println("\tYou have chosen to hit...")
			deal(game, p, 1)

			p.Score = getScore(p)

			score = p.Score

		case "s":

			fmt.Println("\tYou have chosen to stand.")

			return score

		default:
			fmt.Printf("\t'%s' is not a valid choice!\n", choice)
		}

	}

	return score

}

func playTurnAsDealer(game *Game, d *Player) int {

	fmt.Printf("%s's turn...\n", d)

	defer fmt.Println()
	defer showHand(d)

	score := getScore(d)

	for score > 0 && (score < dealerSoftHitLimit || score == dealerSoftHitLimit && hasAce(d)) {

		fmt.Printf("\t%s hits.\n", d)

		deal(game, d, 1)

		showHand(d)

		score = getScore(d)

	}

	fmt.Printf("\t%s stands.\n", d)

	return score

}
