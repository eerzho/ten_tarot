package mongo_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/eerzho/event_manager/pkg/mongo"
	"github.com/eerzho/ten_tarot/internal/entity"
	"github.com/eerzho/ten_tarot/internal/failure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoD "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TgUserTable = "tg_users"

type TGUser struct {
	*mongo.Mongo
}

func NewTGUser(mg *mongo.Mongo) *TGUser {
	return &TGUser{mg}
}

func (t *TGUser) Create(ctx context.Context, user *entity.TGUser) error {
	user.ID = primitive.NewObjectID().Hex()

	result, err := t.DB.Collection(TgUserTable).InsertOne(ctx, user)
	if err != nil {
		return err
	}

	if _, ok := result.InsertedID.(string); !ok {
		return failure.ErrInvalidDocument
	}

	return nil
}

func (t *TGUser) ExistsByChatID(ctx context.Context, chatID string) (bool, error) {
	filter := bson.D{{"chat_id", chatID}}

	var user entity.TGUser
	err := t.DB.Collection(TgUserTable).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongoD.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (t *TGUser) Count(ctx context.Context, chatID, username string) (int, error) {
	filter := t.applyFilter(chatID, username)

	count, err := t.DB.Collection(TgUserTable).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (t *TGUser) ByChatID(ctx context.Context, chatID string) (*entity.TGUser, error) {
	filter := bson.D{{"chat_id", chatID}}

	var user entity.TGUser
	err := t.DB.Collection(TgUserTable).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongoD.ErrNoDocuments) {
			return nil, failure.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (t *TGUser) List(ctx context.Context, username, chatID string, page, count int) ([]entity.TGUser, error) {
	var users []entity.TGUser

	opts := options.Find()
	if page > 0 && count > 0 {
		opts.SetSkip(int64((page - 1) * count))
		opts.SetLimit(int64(count))
	}

	filter := t.applyFilter(chatID, username)
	cursor, err := t.DB.Collection(TgUserTable).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	for cursor.Next(ctx) {
		var user entity.TGUser
		if err = cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (t *TGUser) applyFilter(chatID, username string) bson.D {
	filter := bson.D{}
	if username != "" {
		filter = append(filter, bson.E{Key: "username", Value: bson.D{bson.E{Key: "$regex", Value: fmt.Sprintf("^%s", username)}}})
	}
	if chatID != "" {
		filter = append(filter, bson.E{Key: "chat_id", Value: chatID})
	}

	return filter
}
