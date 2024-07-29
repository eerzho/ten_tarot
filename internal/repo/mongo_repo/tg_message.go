package mongo_repo

import (
	"context"
	"time"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/eerzho/ten_tarot/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TGMessageTable = "tg_messages"
)

type (
	TGMessage struct {
		*mongo.Mongo
	}
)

func NewTGMessage(mg *mongo.Mongo) *TGMessage {
	return &TGMessage{mg}
}

func (t *TGMessage) CountByChatIDFromTime(ctx context.Context, chatID string, fromTime time.Time) (int, error) {
	const op = "mongo_repo.TGMessage.CountByChatIDFromTime"
	logger.Debug(
		op,
		logger.Any("chatID", chatID),
		logger.Any("fromTime", fromTime),
	)

	filter := bson.M{
		"chat_id": chatID,
		"created_at": bson.M{
			"$gte": fromTime.Format(time.DateTime),
		},
	}

	count, err := t.DB.
		Collection(TGMessageTable).
		CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (t *TGMessage) GetList(ctx context.Context, chatID string, page, count int) ([]model.TGMessage, error) {
	const op = "mongo_repo.TGMessage.GetList"
	logger.Debug(
		op,
		logger.Any("chatID", chatID),
		logger.Any("page", page),
		logger.Any("count", count),
	)

	opts := options.Find()
	if page > 0 && count > 0 {
		opts.SetSkip(int64((page - 1) * count))
		opts.SetLimit(int64(count))
	}

	filter := t.applyFilter(chatID)
	cursor, err := t.DB.
		Collection(TGMessageTable).
		Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			logger.OPError(op, closeErr)
		}
	}()

	var messages []model.TGMessage
	for cursor.Next(ctx) {
		var message model.TGMessage
		if err = cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (t *TGMessage) GetListCount(ctx context.Context, chatID string) (int, error) {
	const op = "mongo_repo.TGMessage.GetListCount"
	logger.Debug(op, logger.Any("chatID", chatID))

	filter := t.applyFilter(chatID)
	count, err := t.DB.
		Collection(TGMessageTable).
		CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (t *TGMessage) Create(ctx context.Context, message *model.TGMessage) error {
	const op = "mongo_repo.TGMessage.Create"
	logger.Debug(op, logger.Any("message", message))

	message.ID = primitive.NewObjectID().Hex()
	message.CreatedAt = time.Now().Format(time.DateTime)

	result, err := t.DB.
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

func (t *TGMessage) applyFilter(chatID string) bson.D {
	filter := bson.D{}
	if chatID != "" {
		filter = append(
			filter,
			bson.E{Key: "chat_id", Value: chatID},
		)
	}

	return filter
}
