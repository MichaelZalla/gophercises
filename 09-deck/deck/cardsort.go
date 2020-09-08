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

func defaultSort(i, j *Card) bool {

	if i.Suit < j.Suit {
		return true
	}

	if i.Suit > j.Suit {
		return false
	}

	if i.Rank < j.Rank {
		return true
	}

	if i.Rank > j.Rank {
		return false
	}

	return true

}
