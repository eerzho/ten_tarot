package mongo_repo

import (
	"bot/internal/failure"
	"bot/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"time"
)

const (
	TGMessageTable = "tg_messages"
)

type (
	TGMessage struct {
		lg  *slog.Logger
		mng *mongo.Database
	}
)

func NewTGMessage(lg *slog.Logger, mng *mongo.Database) *TGMessage {
	return &TGMessage{
		lg:  lg,
		mng: mng,
	}
}

func (t *TGMessage) CountByChatIDFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error) {
	const op = "mongo_repo.TGMessage.CountByChatIDFromTime"
	t.lg.Debug(
		op,
		slog.String("chatID", chatID),
		slog.Any("fromTime", fromTime),
	)

	filter := bson.M{
		"chat_id": chatID,
		"created_at": bson.M{
			"$gte": fromTime.Format(time.DateTime),
		},
	}

	count, err := t.mng.
		Collection(TGMessageTable).
		CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (t *TGMessage) Create(ctx context.Context, message *model.TGMessage) error {
	const op = "mongo_repo.TGMessage.Create"
	t.lg.Debug(op, slog.Any("message", message))

	message.ID = primitive.NewObjectID().Hex()
	message.CreatedAt = time.Now().Format(time.DateTime)

	result, err := t.mng.
		Collection(TGMessageTable).
		InsertOne(ctx, message)
	if err != nil {
		return err
	}

	if _, ok := result.InsertedID.(string); !ok {
		return failure.ErrInvalidDocument
	}

	return nil
}
