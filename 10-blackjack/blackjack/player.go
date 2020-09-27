package blackjack

import "fmt"

// player represents a participant in the Blackjack game. A player may be the Dealer.
type player struct {
	ID     int
	Score  int
	Dealer bool
	Hand   Hand
}

func (p player) String() string {
	if p.Dealer {
		return "Dealer"
	}
	return fmt.Sprintf("Player %d", p.ID)
}

func showHand(p *player) {
	fmt.Printf("\t(Score: %d) %v\n", getScore(p), p.Hand)
}
