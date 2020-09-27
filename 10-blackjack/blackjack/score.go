package blackjack

import "github.com/MichaelZalla/gophercises/09-deck/deck"

func hasAce(p *player) bool {
	return (p.Hand[0].Rank == deck.Ace || p.Hand[1].Rank == deck.Ace)
}

func min(a, b int) int {

	if a <= b {
		return a
	}

	return b

}

func getMinScore(p *player) int {

	score := 0

	for _, c := range p.Hand {
		score += min(int(c.Rank), 10)
	}

	return score

}

func getScore(p *player) int {

	score := getMinScore(p)

	if score <= 11 {

		// We can count any Aces as 11 to increase our score, but, if we have 2
		// Aces, counting both as 11 would put us over the max score limit;
		// therefore, double-Aces would be counted as 1+11=12

		for _, c := range p.Hand {
			if c.Rank == deck.Ace {
				score += 10
			}
		}

	}

	if score > maxScore {
		return bustScore
	}

	return score

}
