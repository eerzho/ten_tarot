package mongo_repo

import (
	"context"
	"time"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/eerzho/ten_tarot/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const SupportRequestTable = "support_requests"

type SupportRequest struct {
	*mongo.Mongo
}

func NewSupportRequest(mg *mongo.Mongo) *SupportRequest {
	return &SupportRequest{mg}
}

func (s *SupportRequest) Create(ctx context.Context, supportRequest *model.SupportRequest) error {
	const op = "mongo_repo.SupportRequest.Create"
	logger.Debug(op, logger.Any("supportRequest", supportRequest))

	supportRequest.ID = primitive.NewObjectID().Hex()
	supportRequest.CreatedAt = time.Now().Format(time.DateTime)

	result, err := s.DB.Collection(SupportRequestTable).InsertOne(ctx, supportRequest)
	if err != nil {
		return err
	}

	if _, ok := result.InsertedID.(string); !ok {
		return failure.ErrInvalidDocument
	}

	return nil
}
