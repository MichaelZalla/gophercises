package deck

import (
	"math/rand"
	"sort"
)

// Filter allows the caller to create a new deck from an exisitng deck.
type Filter = func([]Card) []Card

// New generates an empty deck and applies all filters to the deck. If
// no filters are passed, New defaults to using the WithStandard filter.
func New(filters ...Filter) []Card {

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
func WithStandard() Filter {

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

// WithoutCards takes a deck and returns a new deck in which any matching cards
// in cardsToRemove have been removed.
func WithoutCards(cardsToRemove []Card) Filter {

	return func(deck []Card) []Card {

		filtered := []Card{}

		for i := 0; i < len(deck); i++ {
			keep := true
			for _, c := range cardsToRemove {
				if deck[i].Suit == c.Suit && deck[i].Rank == c.Rank {
					keep = false
				}
			}
			if keep {
				filtered = append(filtered, deck[i])
			}
		}

		return filtered

	}

}

// WithJokers combines a deck with some number of new Joker cards.
func WithJokers(n int) Filter {

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

// WithDecks combines a deck with a set of existing decks.
func WithDecks(decks ...[]Card) Filter {

	return func(deck []Card) []Card {

		for _, otherDeck := range decks {
			deck = append(deck, otherDeck...)
		}

		return deck

	}

}

// Sorted takes a deck and performs an in-place sort specified by sortFn. If no
// sortFn is passed, Sorted() will use a default sort function that sorts a deck
// by Rank and Suit. Sorted returns the deck that is passed to it.
func Sorted(getLessFn LessFnGetter) Filter {

	if getLessFn == nil {
		getLessFn = defaultLessFnGetter
	}

	return func(deck []Card) []Card {

		sort.Slice(deck, getLessFn(deck))

		return deck

	}

}

// Shuffled will shuffle a deck randomly. Shuffle returns the deck that is
// passed to it. Shuffle does not call rand.Seed().
func Shuffled() Filter {

	return func(deck []Card) []Card {

		rand.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})

		return deck

	}

}
