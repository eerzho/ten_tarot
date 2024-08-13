package mongo_repo

import (
	"bot/internal/model"
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const MessageTable = "messages"

type Message struct {
	lg  *slog.Logger
	mng *mongo.Database
}

func NewMessage(lg *slog.Logger, mng *mongo.Database) *Message {
	return &Message{
		lg:  lg,
		mng: mng,
	}
}

func (m *Message) CountByChatIDAndFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error) {
	const op = "mongo_repo.Message.CountByChatIDFromTime"
	m.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.Any("fromTime", fromTime),
	)

	filter := bson.M{
		"chat_id": chatID,
		"created_at": bson.M{
			"$gte": fromTime,
		},
	}

	count, err := m.mng.
		Collection(MessageTable).
		CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (m *Message) Create(ctx context.Context, message *model.Message) error {
	const op = "mongo_repo.Message.Create"
	m.lg.Debug(op, slog.Any("message", message))

	message.ID = primitive.NewObjectID()
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	_, err := m.mng.
		Collection(MessageTable).
		InsertOne(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
