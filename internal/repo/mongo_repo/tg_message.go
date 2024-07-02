package mongo_repo

import (
	"context"
	"fmt"
	"time"

	"github.com/eerzho/event_manager/pkg/crypter"
	"github.com/eerzho/event_manager/pkg/mongo"
	"github.com/eerzho/ten_tarot/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TgMessageTable = "tg_messages"

type TGMessage struct {
	*mongo.Mongo
	c *crypter.Crypter
}

func NewTGMessage(m *mongo.Mongo, c *crypter.Crypter) *TGMessage {
	return &TGMessage{m, c}
}

func (t *TGMessage) CountByDay(ctx context.Context, chatID string) (int, error) {
	const op = "./internal/repo/mongo_repo/tg_message::CountMessagesByDay"

	startOfDay := time.Now().Truncate(24 * time.Hour)
	filter := bson.M{
		"chat_id": chatID,
		"created_at": bson.M{
			"$gte": startOfDay.Format(time.DateTime),
		},
	}

	count, err := t.DB.Collection(TgMessageTable).CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(count), nil
}

func (t *TGMessage) All(ctx context.Context, chatID string, page, count int) ([]entity.TGMessage, error) {
	const op = "./internal/repo/mongo_repo/tg_user::All"

	var messages []entity.TGMessage
	filter := bson.D{}
	if chatID != "" {
		filter = append(filter, bson.E{Key: "chat_id", Value: chatID})
	}

	opts := options.Find()
	if page == 0 {
		page = 1
	}
	if count == 0 {
		count = 10
	}
	opts.SetSkip(int64((page - 1) * count))
	opts.SetLimit(int64(count))

	cursor, err := t.DB.Collection(TgMessageTable).Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var message entity.TGMessage
		if err := cursor.Decode(&message); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		messages = append(messages, message)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, nil
}

func (t *TGMessage) Create(ctx context.Context, message *entity.TGMessage) error {
	const op = "./internal/repo/mongo_repo/tg_message::Create"

	decryptedText := message.Text
	decryptedAnswer := message.Answer
	defer func() {
		message.Text = decryptedText
		message.Answer = decryptedAnswer
	}()

	message.ID = primitive.NewObjectID().Hex()
	message.Text = t.c.Encrypt(decryptedText)
	message.Answer = t.c.Encrypt(decryptedAnswer)
	message.CreatedAt = time.Now().Format(time.DateTime)

	result, err := t.DB.Collection(TgMessageTable).InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, ok := result.InsertedID.(string); !ok {
		return fmt.Errorf("%s: document is nil", op)
	}

	return nil
}
