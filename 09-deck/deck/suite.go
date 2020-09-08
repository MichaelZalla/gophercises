package deck

// Suit associated with a Card.
type Suit int

// Enumeration of Suits, ordered as they appear in a standard 52-card deck.
//go:generate stringer -type=Suit
const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
	Joker
)

var suits = [...]Suit{Spades, Diamonds, Clubs, Hearts}
