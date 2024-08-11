package model

type SupportRequest struct {
	ID        string `bson:"_id" json:"id"`
	ChatID    string `bson:"chat_id" json:"chat_id"`
	Question  string `bson:"question" json:"question"`
	CreatedAt string `bson:"created_at" json:"created_at"`
}
