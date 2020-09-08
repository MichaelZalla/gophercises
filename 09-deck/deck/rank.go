package deck

// Rank represents the value of a card within its suit.
type Rank uint8

// Enumeration of ranks in order of appearance in a newly minted deck.
//go:generate stringer -type=Rank
const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const minRank = Ace
const maxRank = King
