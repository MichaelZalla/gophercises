package blackjack

import (
	"fmt"
	"log"

	"github.com/MichaelZalla/gophercises/09-deck/deck"
)

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

const dealerID = 0

const minPlayerID = 1

const maxScore = 21
const bustScore = -1

const dealerSoftHitLimit = 17

// NewGame creates a new game of Blackjack with the given options applied
func NewGame(options ...GameOptionFn) Game {

	// Make a Deck and a Dealer (who serves as Player 0)

	game := Game{
		Dealer: &player{
			ID:     dealerID,
			Dealer: true,
		},
		Players:      []*player{},
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

	fmt.Println()

	// Return the completed game

	return game

}

// Players adds the given number of players to a game (excludes the Dealer)
func Players(n int) func(game Game) Game {

	return func(game Game) Game {

		for i := minPlayerID; i < n+1; i++ {
			game.Players = append(game.Players, &player{
				ID:     i,
				Dealer: false,
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

	fmt.Printf("%s's hand: \n", game.Dealer)
	fmt.Printf("\t(Score: ???) { %s, **HIDDEN ** }\n", game.Dealer.Hand[0])
	fmt.Printf("\n")

	winners, topScore := getRoundWinners(game)

	// Rounds in our Game's History store the IDs of the round's winners

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
		top, bottom := game.Deck[0:n], game.Deck[n:]
		p.Hand = append(p.Hand, top...)
		game.Deck = bottom
	}

}

func getRoundWinners(game *Game) ([]*player, int) {

	var winners []*player

	// Check whether any players were dealt Blackjack (ignores the Dealer)

	for _, p := range game.Players {

		p.Score = getScore(p)

		if p.Dealer {
			continue
		}

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

	if sumPlayerScoreAfterRound < 0 {
		return []*player{game.Dealer}, game.Dealer.Score
	}

	// Dealer plays a turn

	topScore := playTurnAsDealer(game, game.Dealer)

	// At this point, every participant has played their turn; determine the top
	// score for the round

	for _, p := range game.Players {
		if p.Score > topScore {
			topScore = p.Score
		}
	}

	// Determine which participants scored the top score

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

	score := getScore(p)

	showHand(p)

	var choice string

	for score != bustScore && score != maxScore {

		fmt.Printf("\tWill you (h)it or (s)tand?\n\t")
		fmt.Scanln(&choice)

		switch choice {
		case "h":

			fmt.Println("\tYou have chosen to hit...")
			deal(game, p, 1)

			p.Score = getScore(p)

			score = p.Score

			showHand(p)

		case "s":

			fmt.Println("\tYou have chosen to stand.")

			return score

		default:
			fmt.Printf("\t'%s' is not a valid choice!\n", choice)
		}

	}

	if score == maxScore {
		fmt.Println("\tYou have Blackjack!")
	} else if score == bustScore {
		fmt.Println("\tYou went bust!")
	}

	return score

}

func playTurnAsDealer(game *Game, d *player) int {

	defer fmt.Println()

	fmt.Printf("%s's turn...\n", d)

	showHand(d)

	for d.Score != bustScore && d.Score < maxScore && (d.Score < dealerSoftHitLimit || getMinScore(d) < dealerSoftHitLimit) {

		fmt.Printf("\t%s hits.\n", d)

		deal(game, d, 1)

		showHand(d)

		d.Score = getScore(d)

	}

	fmt.Printf("\t%s stands.\n", d)

	return d.Score

}
