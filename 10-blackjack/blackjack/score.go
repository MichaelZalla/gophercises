package blackjack

import "github.com/MichaelZalla/gophercises/09-deck/deck"

func hasAce(p *player) bool {
	return (p.Hand[0].Rank == deck.Ace || p.Hand[1].Rank == deck.Ace)
}

func getScore(p *player) int {

	soft := false

Tally:

	score := 0

	for _, c := range p.Hand {
		score += getValue(c, soft)
	}

	if score > maxScore {

		if !soft {
			soft = true
			goto Tally
		}

		return -1

	}

	return score

}

func getValue(c deck.Card, soft bool) int {

	switch c.Rank {
	case deck.Ace:
		if soft {
			return 1
		}
		return 11
	case deck.Jack:
		fallthrough
	case deck.Queen:
		fallthrough
	case deck.King:
		return 10
	default:
		return int(c.Rank) + 1
	}

}
