package deck

import "sort"

type by func(i, j *Card) bool

func (byFn by) sort(deck []Card) {
	cs := &cardSorter{
		deck: deck,
		byFn: byFn,
	}
	sort.Sort(cs)
}

type cardSorter struct {
	deck []Card
	byFn by
}

func (cs cardSorter) Len() int {
	return len(cs.deck)
}

func (cs cardSorter) Swap(i, j int) {
	cs.deck[i], cs.deck[j] = cs.deck[j], cs.deck[i]
}

func (cs cardSorter) Less(i, j int) bool {
	return cs.byFn(&cs.deck[i], &cs.deck[j])
}

func order(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

func defaultSort(i, j *Card) bool {

	return order(*i) < order(*j)

}
