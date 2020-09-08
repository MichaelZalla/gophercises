package blackjack

import "fmt"

// Player represents a participant in the Blackjack game. A Player may be the Dealer.
type Player struct {
	ID     int
	Score  int
	Dealer bool
	Hand   Hand
}

func (p Player) String() string {
	if p.Dealer {
		return "Dealer"
	}
	return fmt.Sprintf("Player %d", p.ID)
}

func showHand(p *Player) {
	fmt.Printf("\tScore: %d\n", getScore(p))
	fmt.Printf("\t%v\n", p.Hand)
}
