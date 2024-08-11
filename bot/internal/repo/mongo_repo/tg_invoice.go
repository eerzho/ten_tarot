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
	TGInvoiceTable = "tg_invoice"
)

type (
	TGInvoice struct {
		lg  *slog.Logger
		mng *mongo.Database
	}
)

func NewTGInvoice(lg *slog.Logger, mng *mongo.Database) *TGInvoice {
	return &TGInvoice{
		lg:  lg,
		mng: mng,
	}
}

func (t TGInvoice) GetByID(ctx context.Context, id string) (*model.TGInvoice, error) {
	const op = "mongo_repo.TGInvoice.GetByID"
	t.lg.Debug(op, slog.String("id", id))

	var invoice model.TGInvoice
	filter := bson.M{"_id": id}

	err := t.mng.
		Collection(TGInvoiceTable).
		FindOne(ctx, filter).
		Decode(&invoice)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, failure.ErrNotFound
		}
		return nil, err
	}

	return &invoice, nil
}

func (t TGInvoice) Create(ctx context.Context, invoice *model.TGInvoice) error {
	const op = "mongo_repo.TGInvoice.Create"
	t.lg.Debug(op, slog.Any("invoice", invoice))

	invoice.ID = primitive.NewObjectID().Hex()
	invoice.CreatedAt = time.Now().Format(time.DateTime)

	result, err := t.mng.
		Collection(TGInvoiceTable).
		InsertOne(ctx, invoice)
	if err != nil {
		return err
	}

	if _, ok := result.InsertedID.(string); !ok {
		return failure.ErrInvalidDocument
	}

	return nil
}

func (t TGInvoice) Update(ctx context.Context, invoice *model.TGInvoice) error {
	const op = "mongo_repo.TGInvoice.Update"
	t.lg.Debug(op, slog.Any("invoice", invoice))

	filter := bson.M{"_id": invoice.ID}
	update := bson.M{"$set": invoice}

	result, err := t.mng.
		Collection(TGInvoiceTable).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return failure.ErrNotFound
	}

	return nil
}
