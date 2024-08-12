package model

import (
	"bot/internal/def"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	State         def.UserState      `bson:"state" json:"state"`
	ChatID        string             `bson:"chat_id" json:"chat_id"`
	Username      string             `bson:"username" json:"username"`
	QuestionCount int                `bson:"question_count" json:"question_count"`
	CreatedAt     time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}
