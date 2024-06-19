package service

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/entity"
)

type Card struct {
	len  int
	deck []entity.Card
}

func NewCard() *Card {
	card := Card{}

	suits := []string{"Cups", "Pentacles", "Swords", "Wands"}
	values := []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Page", "Knight", "Queen", "King"}

	for _, suit := range suits {
		for _, value := range values {
			card.deck = append(card.deck, entity.Card{Suit: suit, Value: value})
		}
	}

	for i := 1; i <= 22; i++ {
		card.deck = append(card.deck, entity.Card{Suit: "Major Arcana", Value: strconv.Itoa(i)})
	}

	card.len = len(card.deck)

	return &card
}

func (c *Card) Shuffle(ctx context.Context, n int) []entity.Card {
	if n <= 0 {
		n = 1
	}
	if n > c.len {
		n = c.len
	}

	uniqI := make([]int, c.len)
	for i := 0; i < c.len; i++ {
		uniqI[i] = i
	}

	shuffled := make([]entity.Card, n)

	for i := 0; i < n; i++ {
		j := rand.Intn(len(uniqI))
		shuffled[i] = c.deck[uniqI[j]]
		uniqI[j] = uniqI[len(uniqI)-1]
		uniqI = uniqI[:len(uniqI)-1]
	}

	return shuffled
}
