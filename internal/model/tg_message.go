package model

type TGMessage struct {
	ID           string `bson:"_id" json:"id"`
	ChatID       string `bson:"chat_id" json:"chat_id"`
	BotAnswer    string `bson:"bot_answer" json:"bot_answer"`
	UserQuestion string `bson:"user_question" json:"user_question"`
	CreatedAt    string `bson:"created_at" json:"created_at"`
}
