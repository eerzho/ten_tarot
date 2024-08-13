package mongo_repo

import (
	"bot/internal/model"
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const SPTable = "support_requests"

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

func (sp *SupportRequest) Create(ctx context.Context, supportRequest *model.SupportRequest) error {
	const op = "mongo_repo.SupportRequest.Create"
	sp.lg.Debug(op, slog.Any("supportRequest", supportRequest))

	supportRequest.ID = primitive.NewObjectID()
	supportRequest.CreatedAt = time.Now()
	supportRequest.UpdatedAt = time.Now()

	_, err := sp.mng.Collection(SPTable).InsertOne(ctx, supportRequest)
	if err != nil {
		return err
	}

	return nil
}
