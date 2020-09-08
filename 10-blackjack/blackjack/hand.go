package blackjack

import (
	"strings"

	"github.com/MichaelZalla/gophercises/09-deck/deck"
)

// Hand is a set of cards belonging to a Player
type Hand []deck.Card

func (h Hand) String() string {

	var sb strings.Builder

	sb.WriteString("{ ")

	for i, c := range h {
		sb.WriteString(c.String())
		if i < len(h)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString(" }")

	return sb.String()

}
