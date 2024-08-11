package service

import (
	"bot/internal/model"
	"context"
	"log/slog"
	"math/rand"
	"strconv"
)

type (
	Deck struct {
		lg         *slog.Logger
		totalCards int
		cards      []model.Card
	}
)

func NewDeck(lg *slog.Logger) *Deck {
	deck := Deck{lg: lg}

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

func (d *Deck) Shuffle(ctx context.Context, count int) ([]model.Card, error) {
	const op = "service.Deck.Shuffle"
	d.lg.Debug(op, slog.Int("count", count))

	if count <= 0 {
		count = 1
	}
	if count > d.totalCards {
		count = d.totalCards
	}

	indices := make([]int, d.totalCards)
	for i := 0; i < d.totalCards; i++ {
		indices[i] = i
	}

	shuffledDeck := make([]model.Card, count)

	for i := 0; i < count; i++ {
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
