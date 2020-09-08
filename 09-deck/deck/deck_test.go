package deck

import "fmt"

func ExampleNew() {

	deck := New()

	fmt.Println(len(deck))
	fmt.Println(deck[0])
	fmt.Println(deck[len(deck)-1])

	// Output:
	// 52
	// Ace of Spades
	// King of Hearts

}

func ExampleWithStandard() {

	deck := New()

	fmt.Println(len(deck))
	fmt.Println(deck[0])
	fmt.Println(deck[len(deck)-1])

	// Output:
	// 52
	// Ace of Spades
	// King of Hearts

}

func ExampleWithoutCards() {

	faceCards := New(
		WithStandard(),
		WithoutCards(func(c Card) bool {
			return c.Rank >= Two && c.Rank <= Ten
		}))

	fmt.Println(len(faceCards))
	fmt.Println(faceCards[0])
	fmt.Println(faceCards[1])
	fmt.Println(faceCards[len(faceCards)-2])
	fmt.Println(faceCards[len(faceCards)-1])

	// Output:
	// 16
	// Ace of Spades
	// Jack of Spades
	// Queen of Hearts
	// King of Hearts

}

func ExampleWithJokers() {

	deck := New(WithJokers(10))

	fmt.Println(len(deck))

	fmt.Println(deck[0])
	fmt.Println(deck[0].Suit)
	// fmt.Println(deck[0].Rank)

	fmt.Println(deck[1])
	fmt.Println(deck[1].Suit)
	// fmt.Println(deck[1].Rank)

	fmt.Println(deck[2])
	fmt.Println(deck[2].Suit)
	// fmt.Println(deck[2].Rank)

	fmt.Println(deck[len(deck)-1])

	// Output:
	// 10
	// Joker
	// Joker
	// Joker
	// Joker
	// Joker
	// Joker
	// Joker

}

func ExampleConcat() {

	otherDecks := [][]Card{{}, {}}

	otherDecks[0] = New(WithStandard())
	otherDecks[1] = New(WithJokers(2))

	deck := New(Concat(otherDecks...))

	fmt.Println(len(deck) == len(otherDecks[0])+len(otherDecks[1]))
	fmt.Println(deck[0])
	fmt.Println(deck[len(deck)-2])
	fmt.Println(deck[len(deck)-1])

	// Output:
	// true
	// Ace of Spades
	// Joker
	// Joker

}

func ExampleSorted() {

	normalDeck := New(WithStandard())

	reverseDeck := New(Concat(normalDeck), Sorted(func(deck []Card) LessFn {
		return func(i, j int) bool {
			return deck[i].order() > deck[j].order()
		}
	}))

	fmt.Println(len(reverseDeck) == len(normalDeck))

	fmt.Println(reverseDeck[0])
	fmt.Println(reverseDeck[1])

	fmt.Println(reverseDeck[len(reverseDeck)-2])
	fmt.Println(reverseDeck[len(reverseDeck)-1])

	// Output:
	// true
	// King of Hearts
	// Queen of Hearts
	// Two of Spades
	// Ace of Spades

}

func ExampleShuffled() {

	normalDeck := New(WithStandard())

	shuffledDeck := New(Concat(normalDeck), Shuffled())

	fmt.Println(len(shuffledDeck) == len(normalDeck))

	// fmt.Println(shuffled[0])
	// fmt.Println(shuffled[len(shuffled)-1])

	// Output:
	// true

}
