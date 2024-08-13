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

const InvoiceTable = "invoices"

type Invoice struct {
	lg  *slog.Logger
	mng *mongo.Database
}

func NewInvoice(lg *slog.Logger, mng *mongo.Database) *Invoice {
	return &Invoice{
		lg:  lg,
		mng: mng,
	}
}

func (i *Invoice) GetByID(ctx context.Context, id string) (*model.Invoice, error) {
	const op = "mongo_repo.Invoice.GetByID"
	i.lg.Debug(op, slog.String("id", id))

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		i.lg.Warn(op, slog.String("error", err.Error()))
		return nil, def.ErrNotFound
	}

	var invoice model.Invoice
	filter := bson.M{"_id": objectID}

	err = i.mng.
		Collection(InvoiceTable).
		FindOne(ctx, filter).
		Decode(&invoice)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, def.ErrNotFound
		}
		return nil, err
	}

	return &invoice, nil
}

func (i *Invoice) Create(ctx context.Context, invoice *model.Invoice) error {
	const op = "mongo_repo.Invoice.Create"
	i.lg.Debug(op, slog.Any("invoice", invoice))

	invoice.ID = primitive.NewObjectID()
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	_, err := i.mng.
		Collection(InvoiceTable).
		InsertOne(ctx, invoice)
	if err != nil {
		return err
	}

	return nil
}

func (i *Invoice) Update(ctx context.Context, invoice *model.Invoice) error {
	const op = "mongo_repo.Invoice.Update"
	i.lg.Debug(op, slog.Any("invoice", invoice))

	invoice.UpdatedAt = time.Now()

	filter := bson.M{"_id": invoice.ID}
	update := bson.M{"$set": invoice}
	result, err := i.mng.
		Collection(InvoiceTable).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return def.ErrNotFound
	}

	return nil
}
