package model

import "bot/internal/constant"

type TGUser struct {
	ID            string             `bson:"_id,omitempty" json:"id"`
	State         constant.UserState `bson:"state" json:"state"`
	ChatID        string             `bson:"chat_id" json:"chat_id"`
	Username      string             `bson:"username" json:"username"`
	QuestionCount int                `bson:"question_count" json:"question_count"`
	CreatedAt     string             `bson:"created_at" json:"created_at"`
}
