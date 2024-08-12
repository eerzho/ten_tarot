package model

import (
	"bot/internal/def"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type          def.InvoiceType    `bson:"type" json:"type"`
	ChatID        string             `bson:"chat_id" json:"chat_id"`
	ChargeID      string             `bson:"charge_id" json:"charge_id,omitempty"`
	StarsCount    int                `bson:"stars_count" json:"stars_count"`
	QuestionCount int                `bson:"question_count" json:"question_count"`
	CreatedAt     time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}
