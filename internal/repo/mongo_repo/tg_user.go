package mongo_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/eerzho/ten_tarot/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TGUserTable = "tg_users"
)

type (
	TGUser struct {
		*mongo.Mongo
	}
)

func NewTGUser(mg *mongo.Mongo) *TGUser {
	return &TGUser{mg}
}

func (t *TGUser) ExistsByChatID(ctx context.Context, chatID string) bool {
	const op = "mongo_repo.TGUser.ExistsByChatID"
	logger.Debug(op, logger.Any("chatID", chatID))

	filter := bson.D{{"chat_id", chatID}}
	err := t.DB.
		Collection(TGUserTable).
		FindOne(ctx, filter).
		Err()
	if err != nil {
		return false
	}

	return true
}

func (t *TGUser) GetByChatID(ctx context.Context, chatID string) (*model.TGUser, error) {
	const op = "mongo_repo.TGUser.GetByChatID"
	logger.Debug(op, logger.Any("chatID", chatID))

	var user model.TGUser
	filter := bson.D{{"chat_id", chatID}}

	err := t.DB.
		Collection(TGUserTable).
		FindOne(ctx, filter).
		Decode(&user)
	if err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocuments) {
			return nil, failure.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (t *TGUser) GetList(ctx context.Context, username, chatID string, page, count int) ([]model.TGUser, error) {
	const op = "mongo_repo.TGUser.GetList"
	logger.Debug(
		op,
		logger.Any("username", username),
		logger.Any("chatID", chatID),
		logger.Any("page", page),
		logger.Any("count", count),
	)

	opts := options.Find()
	if page > 0 && count > 0 {
		opts.SetSkip(int64((page - 1) * count))
		opts.SetLimit(int64(count))
	}

	filter := t.applyFilter(chatID, username)
	cursor, err := t.DB.
		Collection(TGUserTable).
		Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			logger.OPError(op, closeErr)
		}
	}()

	var users []model.TGUser
	for cursor.Next(ctx) {
		var user model.TGUser
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

func (t *TGUser) GetListCount(ctx context.Context, chatID, username string) (int, error) {
	const op = "mongo_repo.TGUser.GetListCount"
	logger.Debug(
		op,
		logger.Any("chatID", chatID),
		logger.Any("username", username),
	)

	filter := t.applyFilter(chatID, username)
	count, err := t.DB.
		Collection(TGUserTable).
		CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (t *TGUser) Create(ctx context.Context, user *model.TGUser) error {
	const op = "mongo_repo.TGUser.Create"
	logger.Debug(op, logger.Any("user", user))

	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = time.Now().Format(time.DateTime)

	result, err := t.DB.
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
	logger.Debug(op, logger.Any("user", user))

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	result, err := t.DB.
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

func (t *TGUser) applyFilter(chatID, username string) bson.D {
	filter := bson.D{}
	if username != "" {
		filter = append(
			filter,
			bson.E{
				Key: "username",
				Value: bson.D{
					bson.E{
						Key:   "$regex",
						Value: fmt.Sprintf("^%s", username),
					},
				},
			},
		)
	}
	if chatID != "" {
		filter = append(
			filter,
			bson.E{Key: "chat_id", Value: chatID},
		)
	}

	return filter
}
