package mongo_repo

import (
	"context"
	"errors"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

const (
	TGInvoiceTable = "tg_invoice"
)

type (
	TGInvoice struct {
		*mongo.Mongo
	}
)

func NewTGInvoice(mg *mongo.Mongo) *TGInvoice {
	return &TGInvoice{mg}
}

func (t TGInvoice) Create(ctx context.Context, invoice *model.TGInvoice) error {
	invoice.ID = primitive.NewObjectID().Hex()

	result, err := t.DB.Collection(TGInvoiceTable).InsertOne(ctx, invoice)
	if err != nil {
		return err
	}

	if _, ok := result.InsertedID.(string); !ok {
		return failure.ErrInvalidDocument
	}

	return nil
}

func (t TGInvoice) Update(ctx context.Context, invoice *model.TGInvoice) error {
	filter := bson.M{"_id": invoice.ID}
	update := bson.M{"$set": invoice}

	result, err := t.DB.Collection(TGInvoiceTable).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return failure.ErrNotFound
	}

	return nil
}

func (t TGInvoice) GetByID(ctx context.Context, id string) (*model.TGInvoice, error) {
	var invoice model.TGInvoice
	filter := bson.M{"_id": id}

	if err := t.DB.Collection(TGInvoiceTable).FindOne(ctx, filter).Decode(&invoice); err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocuments) {
			return nil, failure.ErrNotFound
		}
		return nil, err
	}

	return &invoice, nil
}
