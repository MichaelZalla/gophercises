package deck

import "fmt"

func ExampleCard() {

	fmt.Println(Card{Suit: Hearts})

	fmt.Println(Card{Rank: Four})

	fmt.Println(Card{})

	fmt.Println(Card{Rank: Ace, Suit: Hearts})

	fmt.Println(Card{Rank: Two, Suit: Spades})

	fmt.Println(Card{Rank: Nine, Suit: Diamonds})

	fmt.Println(Card{Rank: Jack, Suit: Clubs})

	fmt.Println(Card{Suit: Joker})

	fmt.Println(Card{Rank: Seven, Suit: Joker})

	// Output:
	// Ace of Hearts
	// Four of Spades
	// Ace of Spades
	// Ace of Hearts
	// Two of Spades
	// Nine of Diamonds
	// Jack of Clubs
	// Joker
	// Joker

}
