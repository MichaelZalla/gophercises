package deck

import (
	"fmt"
)

// Card defines a given card in a deck
type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {

	if c.Suit == Joker {
		return c.Suit.String()
	}

	return fmt.Sprintf("%s of %s", c.Rank.String(), c.Suit.String())

}
