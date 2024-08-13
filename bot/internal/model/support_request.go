package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SupportRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChatID    string             `bson:"chat_id" json:"chat_id"`
	Question  string             `bson:"question" json:"question"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}
