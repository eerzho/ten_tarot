package service

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
)

type (
	Deck struct {
		totalCards int
		cards      []model.Card
	}
)

func NewDeck() *Deck {
	deck := Deck{}

	suits := []string{"Cups", "Pentacles", "Swords", "Wands"}
	values := []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Page", "Knight", "Queen", "King"}

	for _, suit := range suits {
		for _, value := range values {
			deck.cards = append(
				deck.cards,
				model.Card{Suit: suit, Value: value},
			)
		}
	}

	for i := 1; i <= 22; i++ {
		deck.cards = append(
			deck.cards,
			model.Card{Suit: "Major Arcana", Value: strconv.Itoa(i)},
		)
	}

	deck.totalCards = len(deck.cards)

	return &deck
}

func (d *Deck) Shuffle(ctx context.Context, numCards int) ([]model.Card, error) {
	const op = "service.Deck.Shuffle"
	logger.Debug(op, logger.Any("numCards", numCards))

	if numCards <= 0 {
		numCards = 1
	}
	if numCards > d.totalCards {
		numCards = d.totalCards
	}

	indices := make([]int, d.totalCards)
	for i := 0; i < d.totalCards; i++ {
		indices[i] = i
	}

	shuffledDeck := make([]model.Card, numCards)

	for i := 0; i < numCards; i++ {
		select {
		case <-ctx.Done():
			return shuffledDeck[:i], ctx.Err()
		default:
			j := rand.Intn(len(indices))
			shuffledDeck[i] = d.cards[indices[j]]
			indices[j] = indices[len(indices)-1]
			indices = indices[:len(indices)-1]
		}
	}

	return shuffledDeck, nil
}
