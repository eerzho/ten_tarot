package mongo_repo

import (
	"bot/internal/failure"
	"bot/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"time"
)

const SupportRequestTable = "support_requests"

type SupportRequest struct {
	lg  *slog.Logger
	mng *mongo.Database
}

func NewSupportRequest(lg *slog.Logger, mng *mongo.Database) *SupportRequest {
	return &SupportRequest{
		lg:  lg,
		mng: mng,
	}
}

func (s *SupportRequest) Create(ctx context.Context, supportRequest *model.SupportRequest) error {
	const op = "mongo_repo.SupportRequest.Create"
	s.lg.Debug(op, slog.Any("supportRequest", supportRequest))

	supportRequest.ID = primitive.NewObjectID().Hex()
	supportRequest.CreatedAt = time.Now().Format(time.DateTime)

	result, err := s.mng.Collection(SupportRequestTable).InsertOne(ctx, supportRequest)
	if err != nil {
		return err
	}

	if _, ok := result.InsertedID.(string); !ok {
		return failure.ErrInvalidDocument
	}

	return nil
}
