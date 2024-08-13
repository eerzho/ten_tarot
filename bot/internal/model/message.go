package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChatID       string             `bson:"chat_id" json:"chat_id"`
	BotAnswer    string             `bson:"bot_answer" json:"bot_answer"`
	UserQuestion string             `bson:"user_question" json:"user_question"`
	CreatedAt    time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}
