package mongo_repo

import (
	"bot/internal/failure"
	"bot/internal/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"time"
)

const (
	TGUserTable = "tg_users"
)

type (
	TGUser struct {
		lg  *slog.Logger
		mng *mongo.Database
	}
)

func NewTGUser(lg *slog.Logger, mng *mongo.Database) *TGUser {
	return &TGUser{
		lg:  lg,
		mng: mng,
	}
}

func (t *TGUser) ExistsByChatID(ctx context.Context, chatID string) bool {
	const op = "mongo_repo.TGUser.ExistsByChatID"
	t.lg.Debug(op, slog.String("chatID", chatID))

	filter := bson.D{{Key: "chat_id", Value: chatID}}
	err := t.mng.
		Collection(TGUserTable).
		FindOne(ctx, filter).
		Err()

	return err == nil
}

func (t *TGUser) GetByChatID(ctx context.Context, chatID string) (*model.TGUser, error) {
	const op = "mongo_repo.TGUser.GetByChatID"
	t.lg.Debug(op, slog.String("chatID", chatID))

	var user model.TGUser
	filter := bson.D{{Key: "chat_id", Value: chatID}}

	err := t.mng.
		Collection(TGUserTable).
		FindOne(ctx, filter).
		Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, failure.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (t *TGUser) Create(ctx context.Context, user *model.TGUser) error {
	const op = "mongo_repo.TGUser.Create"
	t.lg.Debug(op, slog.Any("user", user))

	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = time.Now().Format(time.DateTime)

	result, err := t.mng.
		Collection(TGUserTable).
		InsertOne(ctx, user)
	if err != nil {
		return err
	}

	if _, ok := result.InsertedID.(string); !ok {
		return failure.ErrInvalidDocument
	}

	return nil
}

func (t *TGUser) Update(ctx context.Context, user *model.TGUser) error {
	const op = "mongo_repo.TGUser.Update"
	t.lg.Debug(op, slog.Any("user", user))

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	result, err := t.mng.
		Collection(TGUserTable).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return failure.ErrNotFound
	}

	return nil
}
