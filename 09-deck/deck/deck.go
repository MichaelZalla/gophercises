package deck

import (
	"math/rand"
	"sort"
	"time"
)

// FilterFn allows the caller to create a new deck from an exisitng deck.
type FilterFn = func([]Card) []Card

// New generates an empty deck and applies all filters to the deck. If
// no filters are passed, New defaults to using the WithStandard filter.
func New(filters ...FilterFn) []Card {

	d := []Card{}

	if len(filters) == 0 {
		filters = append(filters, WithStandard())
	}

	for _, o := range filters {
		d = o(d)
	}

	return d

}

// WithStandard combines a deck with a standard 52-card deck.
func WithStandard() FilterFn {

	return func(deck []Card) []Card {

		for _, suit := range suits {
			for rank := minRank; rank <= maxRank; rank++ {
				deck = append(deck, Card{
					Suit: suit,
					Rank: rank,
				})
			}

		}

		return deck

	}

}

// Difference takes a filter function and returns a new deck in which any
// original cards matching the filter are removed.
func Difference(filterFn func(c Card) bool) FilterFn {

	return func(deck []Card) []Card {

		filtered := []Card{}

		for _, card := range deck {
			if !filterFn(card) {
				filtered = append(filtered, card)
			}
		}

		return filtered

	}

}

// WithJokers combines a deck with some number of new Joker cards.
func WithJokers(n int) FilterFn {

	return func(deck []Card) []Card {

		for i := 0; i < n; i++ {
			deck = append(deck, Card{
				Suit: Joker,
				Rank: Rank(i),
			})
		}

		return deck

	}

}

// Concat combines a deck with a set of existing decks.
func Concat(decks ...[]Card) FilterFn {

	return func(deck []Card) []Card {

		for _, otherDeck := range decks {
			deck = append(deck, otherDeck...)
		}

		return deck

	}

}

// WithCopies multiplies a deck n times.
func WithCopies(n int) FilterFn {

	return func(deck []Card) []Card {

		var copies []Card

		for i := 0; i < n; i++ {
			copies = append(copies, deck...)
		}

		return copies

	}

}

// Sorted takes a deck and performs an in-place sort specified by sortFn. If no
// sortFn is passed, Sorted() will use a default sort function that sorts a deck
// by Rank and Suit. Sorted returns the deck that is passed to it.
func Sorted(getLessFn LessFnGetter) FilterFn {

	if getLessFn == nil {
		getLessFn = defaultLessFnGetter
	}

	return func(deck []Card) []Card {

		sort.Slice(deck, getLessFn(deck))

		return deck

	}

}

var shuffleRandSource = rand.New(rand.NewSource(time.Now().Unix()))

// Shuffled will shuffle a deck randomly. Shuffle returns the deck that is
// passed to it. Shuffle uses its own rand.Source.
func Shuffled() FilterFn {

	return func(deck []Card) []Card {

		shuffled := make([]Card, len(deck))

		p := shuffleRandSource.Perm(len(deck))

		for i, j := range p {
			shuffled[i] = deck[j]
		}

		return shuffled

	}

}
