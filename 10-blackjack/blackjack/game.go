package blackjack

import (
	"fmt"
	"log"

	"github.com/MichaelZalla/gophercises/09-deck/deck"
)

// round stores information about a particular game round, including the winners
type round struct {
	Winners  []int
	TopScore int
}

// Game stores all information related to the current Blackjack game
type Game struct {
	Dealer       *player
	Players      []*player
	Deck         []deck.Card
	CurrentRound int
	Rounds       int
	History      []round
}

// GameOptionFn allows callers to customize a game (i.e., when calling NewGame)
type GameOptionFn func(game Game) Game

const maxScore = 21

const dealerSoftHitLimit = 17

// NewGame creates a new game of Blackjack
func NewGame(options ...GameOptionFn) Game {

	// Generate a Deck and a Dealer (who serves as Player 0)

	game := Game{
		Dealer: &player{
			ID:     0,
			Dealer: true,
		},
		Rounds:       1,
		CurrentRound: 1,
	}

	// Add Players, etc

	for _, option := range options {
		game = option(game)
	}

	// Ensures Dealer is always the last to be dealt

	game.Players = append(game.Players, game.Dealer)

	// Play out all of the rounds, keeping a history of the round winners

	for game.CurrentRound <= game.Rounds {
		playRound(&game)
		game.CurrentRound++
	}

	// Display game summary

	fmt.Printf("\n==========================================\n")
	fmt.Printf("\tGame Summary\n")
	fmt.Printf("==========================================\n\n")

	for i, r := range game.History {
		fmt.Printf("Winner IDs for Round %d: %v (top score: %d)\n", i+1, r.Winners, r.TopScore)
	}

	// Return the completed game

	return game

}

// Players adds the given number of players to a game (excludes the dealer)
func Players(n int) func(game Game) Game {

	return func(game Game) Game {

		// @NOTE(mzalla) Reserves 0 for Dealer's ID
		for i := 1; i < n+1; i++ {
			game.Players = append(game.Players, &player{
				ID: i,
			})
		}

		return game

	}

}

// Rounds sets a limit on the number of rounds to be played as part of the game
func Rounds(n int) func(game Game) Game {

	return func(game Game) Game {

		if n > 0 {
			game.Rounds = n
		}

		return game

	}

}

func playRound(game *Game) {

	// Reset the deck, and all player hands

	game.Deck = deck.New(
		deck.WithStandard(),
		deck.Shuffled())

	for _, p := range game.Players {
		p.Hand = []deck.Card{}
	}

	// Deal a hand to each player

	for i := 0; i < 2; i++ {
		for j := 0; j < len(game.Players); j++ {
			deal(game, game.Players[j], 1)
		}
	}

	// NewGame blocks until getRoundWinners returns

	fmt.Printf("\n==========================================\n")
	fmt.Printf("\tStarting round #%d\n", game.CurrentRound)
	fmt.Printf("==========================================\n\n")

	winners, topScore := getRoundWinners(game)

	var winnerIDs []int

	for _, p := range winners {
		winnerIDs = append(winnerIDs, p.ID)
	}

	// Record the round history

	game.History = append(game.History, round{
		Winners:  winnerIDs,
		TopScore: topScore,
	})

	// Print a summary of the round

	fmt.Printf("Round winners (top score: %d):\n", topScore)
	for _, p := range winners {
		fmt.Printf("\t%s\n", p)
	}

}

func deal(game *Game, p *player, n int) {

	if len(game.Deck) < n {
		log.Fatal("Called deal() on deck with too few cards!")
	}

	for i := 0; i < n; i++ {
		front, back := game.Deck[0:n], game.Deck[n:]
		p.Hand = append(p.Hand, front...)
		game.Deck = back
	}

}

func getRoundWinners(game *Game) ([]*player, int) {

	var winners []*player

	// Check whether any players were dealt Blackjack (ignores the Dealer)

	for _, p := range game.Players {

		p.Score = getScore(p)

		if int(p.Score) == maxScore {
			winners = append(winners, p)
		}

	}

	if len(winners) > 0 {
		return winners, maxScore
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
		return []*player{game.Dealer}, getScore(game.Dealer)
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

	return winners, topScore

}

func playTurnAsPlayer(game *Game, p *player) int {

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

func playTurnAsDealer(game *Game, d *player) int {

	fmt.Printf("%s's turn...\n", d)

	defer fmt.Println()
	defer showHand(d)

	score := getScore(d)

	for score > 0 && (score < dealerSoftHitLimit || getMinScore(d) < dealerSoftHitLimit) {

		fmt.Printf("\t%s hits.\n", d)

		deal(game, d, 1)

		showHand(d)

		score = getScore(d)

	}

	fmt.Printf("\t%s stands.\n", d)

	return score

}
