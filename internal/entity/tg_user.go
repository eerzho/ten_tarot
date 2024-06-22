package entity

type TGUser struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	ChatID    string `bson:"chat_id" json:"chat_id"`
	Username  string `bson:"username" json:"username"`
	CreatedAt string `bson:"created_at" json:"created_at"`
}
