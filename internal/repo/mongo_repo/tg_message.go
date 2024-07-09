package mongo_repo

import (
	"context"
	"time"

	"github.com/eerzho/event_manager/pkg/crypter"
	"github.com/eerzho/event_manager/pkg/mongo"
	"github.com/eerzho/ten_tarot/internal/entity"
	"github.com/eerzho/ten_tarot/internal/failure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TgMessageTable = "tg_messages"

type TGMessage struct {
	*mongo.Mongo
	c *crypter.Crypter
}

func NewTGMessage(c *crypter.Crypter, mg *mongo.Mongo) *TGMessage {
	return &TGMessage{
		mg,
		c,
	}
}

func (t *TGMessage) Count(ctx context.Context, chatID string) (int, error) {
	filter := t.applyFilter(chatID)
	count, err := t.DB.Collection(TgMessageTable).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (t *TGMessage) Create(ctx context.Context, message *entity.TGMessage) error {
	message.ID = primitive.NewObjectID().Hex()

	originalText := message.Text
	originalAnswer := message.Answer
	defer func() {
		message.Text = originalText
		message.Answer = originalAnswer
	}()
	message.Text = t.c.Encrypt(originalText)
	message.Answer = t.c.Encrypt(originalAnswer)

	result, err := t.DB.Collection(TgMessageTable).InsertOne(ctx, message)
	if err != nil {
		return err
	}

	if _, ok := result.InsertedID.(string); !ok {
		return failure.ErrInvalidDocument
	}

	return nil
}

func (t *TGMessage) CountByTime(ctx context.Context, chatID string, st time.Time) (int, error) {
	filter := bson.M{
		"chat_id": chatID,
		"created_at": bson.M{
			"$gte": st.Format(time.DateTime),
		},
	}
	count, err := t.DB.Collection(TgMessageTable).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (t *TGMessage) List(ctx context.Context, chatID string, page, count int) ([]entity.TGMessage, error) {
	var messages []entity.TGMessage

	opts := options.Find()
	if page > 0 && count > 0 {
		opts.SetSkip(int64((page - 1) * count))
		opts.SetLimit(int64(count))
	}

	filter := t.applyFilter(chatID)
	cursor, err := t.DB.Collection(TgMessageTable).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	for cursor.Next(ctx) {
		var message entity.TGMessage
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

func (t *TGMessage) applyFilter(chatID string) bson.D {
	filter := bson.D{}
	if chatID != "" {
		filter = append(filter, bson.E{Key: "chat_id", Value: chatID})
	}

	return filter
}
