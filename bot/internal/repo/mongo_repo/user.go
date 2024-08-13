package mongo_repo

import (
	"bot/internal/def"
	"bot/internal/model"
	"context"
	"errors"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const UserTable = "users"

type User struct {
	lg  *slog.Logger
	mng *mongo.Database
}

func NewUser(lg *slog.Logger, mng *mongo.Database) *User {
	return &User{
		lg:  lg,
		mng: mng,
	}
}

func (u *User) ExistsByChatID(ctx context.Context, chatID string) bool {
	const op = "mongo_repo.User.ExistsByChatID"
	u.lg.Debug(op, slog.String("chatID", chatID))

	filter := bson.D{{Key: "chat_id", Value: chatID}}
	err := u.mng.
		Collection(UserTable).
		FindOne(ctx, filter).
		Err()

	return err == nil
}

func (u *User) GetByChatID(ctx context.Context, chatID string) (*model.User, error) {
	const op = "mongo_repo.User.GetByChatID"
	u.lg.Debug(op, slog.String("chatID", chatID))

	var user model.User
	filter := bson.D{{Key: "chat_id", Value: chatID}}

	err := u.mng.
		Collection(UserTable).
		FindOne(ctx, filter).
		Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, def.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	const op = "mongo_repo.User.Create"
	u.lg.Debug(op, slog.Any("user", user))

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := u.mng.
		Collection(UserTable).
		InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Update(ctx context.Context, user *model.User) error {
	const op = "mongo_repo.User.Update"
	u.lg.Debug(op, slog.Any("user", user))

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	result, err := u.mng.
		Collection(UserTable).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return def.ErrNotFound
	}

	return nil
}
