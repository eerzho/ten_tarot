package service

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/model"
)

type (
	Card interface {
		Shuffle(ctx context.Context, n int) ([]model.Card, error)
	}
	card struct {
		len  int
		deck []model.Card
	}
)

func NewCard() Card {
	c := card{}

	suits := []string{"Cups", "Pentacles", "Swords", "Wands"}
	values := []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Page", "Knight", "Queen", "King"}

	for _, suit := range suits {
		for _, value := range values {
			c.deck = append(c.deck, model.Card{Suit: suit, Value: value})
		}
	}

	for i := 1; i <= 22; i++ {
		c.deck = append(c.deck, model.Card{Suit: "Major Arcana", Value: strconv.Itoa(i)})
	}

	c.len = len(c.deck)

	return &c
}

func (c *card) Shuffle(ctx context.Context, n int) ([]model.Card, error) {
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

	shuffled := make([]model.Card, n)

	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			return shuffled[:i], ctx.Err()
		default:
			j := rand.Intn(len(uniqI))
			shuffled[i] = c.deck[uniqI[j]]
			uniqI[j] = uniqI[len(uniqI)-1]
			uniqI = uniqI[:len(uniqI)-1]
		}
	}

	return shuffled, nil
}
