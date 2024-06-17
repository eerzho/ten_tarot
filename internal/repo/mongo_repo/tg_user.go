package mongo_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/eerzho/event_manager/internal/entity"
	"github.com/eerzho/event_manager/internal/failure"
	"github.com/eerzho/event_manager/pkg/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TgUserTable = "tg_users"

type TGUser struct {
	*mongo.Mongo
}

func NewTGUser(m *mongo.Mongo) *TGUser {
	return &TGUser{m}
}

func (t *TGUser) All(ctx context.Context, username, chatID string, page, count int) ([]entity.TGUser, error) {
	const op = "./internal/repo/mongo_repo/tg_user::All"

	var users []entity.TGUser

	filter := bson.D{}
	if username != "" {
		filter = append(filter, bson.E{Key: "username", Value: bson.D{bson.E{Key: "$regex", Value: fmt.Sprintf("^%s", username)}}})
	}
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

	cursor, err := t.DB.Collection(TgUserTable).Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	for cursor.Next(ctx) {
		var user entity.TGUser
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (t *TGUser) ByChatID(ctx context.Context, chatID string) (*entity.TGUser, error) {
	const op = "./internal/repo/mongo_repo/tg_user::ByChatID"

	var user entity.TGUser

	filter := bson.D{{"chat_id", chatID}}

	err := t.DB.Collection(TgUserTable).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, failure.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (t *TGUser) Create(ctx context.Context, user *entity.TGUser) error {
	const op = "./internal/repo/mongo_repo/tg_user::Create"

	user.ID = uuid.New().String()

	result, err := t.DB.Collection(TgUserTable).InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, ok := result.InsertedID.(string); !ok {
		return fmt.Errorf("%s: document is nil", op)
	}

	return nil
}
