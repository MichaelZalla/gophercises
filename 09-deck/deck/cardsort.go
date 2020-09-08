package deck

// type by func(i, j *Card) bool

// func (byFn by) sort(deck []Card) {
// 	cs := &cardSorter{
// 		deck: deck,
// 		byFn: byFn,
// 	}
// 	sort.Sort(cs)
// }

// type cardSorter struct {
// 	deck []Card
// 	byFn by
// }

// func (cs cardSorter) Len() int {
// 	return len(cs.deck)
// }

// func (cs cardSorter) Swap(i, j int) {
// 	cs.deck[i], cs.deck[j] = cs.deck[j], cs.deck[i]
// }

// func (cs cardSorter) Less(i, j int) bool {
// 	return cs.byFn(&cs.deck[i], &cs.deck[j])
// }

func (c Card) order() int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

// LessFn is a comparator function for comparing two cards during sorting
type LessFn func(i, j int) bool

// LessFnGetter accepts a deck and returns a corresponding comparator function
type LessFnGetter func(deck []Card) LessFn

func defaultLessFnGetter(deck []Card) LessFn {
	return func(i, j int) bool {
		return deck[i].order() < deck[j].order()
	}
}
