package main

import "strconv"

type suit int

const (
	Spades suit = iota
	Hearts
	Clubs
	Diamonds
)

type Card struct {
	name  string
	suit  suit
	value int
}

func newCard(suit suit, value int) Card {
	var letter = "S"
	if suit == 1 {
		letter = "H"
	}
	if suit == 2 {
		letter = "C"
	}
	if suit == 3 {
		letter = "D"
	}

	return Card{
		name:  letter + strconv.Itoa(value),
		suit:  suit,
		value: value,
	}
}

func newDeck() [52]Card {
	var newDeck = [52]Card{}
	for i := 0; i < 4; i++ {
		for j := 2; j <= 14; j++ {
			newDeck[(i*13)+j-2] = newCard(suit(i), j)
		}
	}
	return newDeck
}
