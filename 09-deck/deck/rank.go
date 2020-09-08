package deck

// Rank represents the value of a card within its suit.
type Rank int

// Enumeration of ranks in order of appearance in a newly minted deck.
//go:generate stringer -type=Rank
const (
	Ace Rank = iota
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
